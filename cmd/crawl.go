package cmd

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

// Urls : All urls would be stored in the map present here and it also contains a mutex for locking the resource
type Urls struct {
	Access      sync.RWMutex
	UrlsScraped map[string]bool
}

// Add : Add to the UrlsScraped map
func (URLS *Urls) Read(url string) bool {
	URLS.Access.RLock()
	defer URLS.Access.RUnlock()
	if !URLS.UrlsScraped[url] {
		return false
	}
	return true
}

func (URLS *Urls) Write(url string) bool {
	if !URLS.Read(url) {
		URLS.Access.Lock()

		URLS.UrlsScraped[url] = true
		URLS.Access.Unlock()
		return true
	}
	return false
}

// SetUrlsMap : Setter function to set the map
func (URLS *Urls) SetUrlsMap() {
	URLS.UrlsScraped = make(map[string]bool)
}

// URLS : instance of Urls
var URLS = Urls{}

// sema : This channel will act as a semaphore to help allow a certain number of running goroutines(concurrent) at a time
var sema chan bool

// CrawlWebsite : Crawl a given website
func CrawlWebsite(str string, concurrent int) error {

	if str != "" && concurrent >= 2 {
		var wg sync.WaitGroup
		sema = make(chan bool, concurrent)
		defer close(sema)
		defer wg.Wait()

		// Set Map in URLS
		URLS.SetUrlsMap()

		site := CheckStringInitial(str)
		fmt.Println("Crawling " + site + " ...")

		wg.Add(1)
		sema <- true
		go CrawlURL(str, site, &wg)
	} else if concurrent < 2 {
		fmt.Println("The minimum number of goroutines needed to run this application is 2. \n\t Please give 2 or a higher number of goroutines")
	} else {
		fmt.Println("Please give a proper argument for the website you want to scrape")
	}

	return nil

}

func getIndent(depth int) string {
	return strings.Repeat("| ", depth)
}

func getHTTPGETRequest(str string) *http.Request {

	/*
		#### Custom Http Get Request so we can time it out at the correct moment
		Was implemented with contecxt before
	*/
	request, err := http.NewRequest(http.MethodGet, str, nil)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Error occurred during HttpGet %v : %v \n", str, err))
	}
	// ####

	return request
}

func getHTTPClient() *http.Client {
	/*
		#### Custom Client for correct Timing Out of the Client and Dial
	*/
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: time.Duration(time.Second * 5),
		}).Dial,
		TLSHandshakeTimeout: time.Duration(time.Second * 5),
	}
	var netClient = &http.Client{
		Timeout:   time.Duration(time.Second * 5),
		Transport: netTransport,
	}
	// ####
	return netClient
}

// CrawlURL : Crawl a given URL
func CrawlURL(str string, site string, wg *sync.WaitGroup) error {

	//defer func() { <-sema }()
	defer wg.Done()

	if !strings.Contains(str, "http") {
		str = "http://" + str
	}

	httpClient := getHTTPClient()
	request := getHTTPGETRequest(str)
	response, err := httpClient.Do(request)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Error occurred during http.Get : %v \n", err))
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		os.Stderr.WriteString(fmt.Sprintf("Incorrect HTTP status obtained : %v \n", response.Status))
	}
	// ####

	htmlTokenizer := html.NewTokenizer(response.Body)

	for {
		tagType := htmlTokenizer.Next()
		switch {
		case tagType == html.ErrorToken:
			return errors.New("Error at CrawlWebsite : Encountered ErrorTag (End of scraping)")
		case tagType == html.StartTagToken: // Check if encountered tag is starting tag or ending tag
			token := htmlTokenizer.Token()
			if token.Data == "a" { // Upon encountering <a>
				for _, a := range token.Attr {
					if a.Key == "href" {
						scrapedURL := a.Val
						scrapedURL = CheckScrapedHref(scrapedURL, site)
						if scrapedURL != "" {
							if URLS.Write(scrapedURL) {
								sema <- true
								wg.Add(2)
								fmt.Printf("%v %v \n", getIndent(strings.Count(a.Val, "/")), a.Val)
								go CrawlURL(scrapedURL, site, wg)
								go func() { defer wg.Done(); time.Sleep(3 * time.Second); <-sema }()
							}
						}
					}
				}
			}
		}
	}

	return nil
}

// CheckStringInitial : Initial check for the string which has been given
func CheckStringInitial(str string) string {
	r, _ := regexp.Compile("^https?://.*")

	if r.MatchString(str) {
		return str
	}
	str = "http://" + str

	return str
}

/*
CheckScrapedHref : Check the href and change it accordingly
based on whether it starts with a / or the name of the site itself
*/
func CheckScrapedHref(scrapedURL string, site string) string {

	ret := ""
	/*
		pass Arguments for the scraped strings
		Contains
			http://
			https://
			if it is not present in the array of urls

		changes to be made if the above is not applicable
			if the name starts with / add the site in front of the name

	*/
	rslash, _ := regexp.Compile("^/.*")
	slashMatch := rslash.MatchString(scrapedURL)
	slashMatchURL := ""
	if slashMatch {
		//fmt.Println("__________SLASH MATCH___________")
		if strings.Contains(scrapedURL, "//") {
			//fmt.Println("__________DOUBLE FORWARD MATCHED___________")
			//slashMatchURL = "http://" + strings.TrimLeft(scrapedURL, "/")
			//fmt.Println(slashMatchURL)
		} else {
			slashMatchURL = strings.TrimRight(site, "/") + scrapedURL
			//fmt.Println(slashMatchURL)
		}
		ret = slashMatchURL
	}

	return ret
}
