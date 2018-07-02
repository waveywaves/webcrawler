package cmd

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

// Urls : All urls would be stored in the map present here and it also contains a mutex for locking mutex
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

		fmt.Println(url)
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

// CrawlWebsite : Crawl a given website
func CrawlWebsite(str string) error {

	var wg sync.WaitGroup
	URLS.SetUrlsMap()

	site := CheckStringInitial(str)
	fmt.Println("Crawling " + site + " ...")

	wg.Add(1)
	go CrawlURL(str, site, &wg)
	wg.Wait()

	return nil
}

// CrawlURL : Crawl a given URL
func CrawlURL(str string, site string, wg *sync.WaitGroup) error {
	defer wg.Done()
	if !strings.Contains(str, "http") {
		str = "http://" + str
	}

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: time.Second * 200,
		}).Dial,
		TLSHandshakeTimeout: time.Second * 200,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 200,
		Transport: netTransport,
	}
	response, err := netClient.Get(str)
	if err != nil {
		log.Fatalf("Error occurred during http.Get : %v \n", err)
		return err
	}
	defer response.Body.Close()

	// Cannot assign struct field for map

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
								wg.Add(1)
								go CrawlURL(scrapedURL, site, wg)
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
			slashMatchURL = "http://" + strings.TrimLeft(scrapedURL, "/")
			//fmt.Println(slashMatchURL)
		} else {
			slashMatchURL = strings.TrimRight(site, "/") + scrapedURL
			//fmt.Println(slashMatchURL)
		}
		ret = slashMatchURL
	} else {
		checkSite, _ := regexp.Compile("^http.*\\.(.*)\\..*")
		match := checkSite.FindStringSubmatch(scrapedURL)
		if len(match) == 2 && strings.Contains(site, match[1]) {
			//fmt.Println("_________SITE MATCH__________")
			//fmt.Println(scrapedURL)
			ret = scrapedURL
		}
	}

	return ret
}
