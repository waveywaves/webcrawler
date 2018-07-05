package cmd

import (
	"net/http"
	"testing"
)

var tt = []struct {
	name                 string
	website              string
	concurrentGoroutines int
}{
	{"No arguments and 0 goroutines", "", 0},
	{"Wrong argument and 1 goroutine", "webscraper", 0},
	{"Website and minimum goroutines", "http://www.guimp.com/", 2},
	{"Website and lot of goroutines", "webscraper.io", 300},
}

func TestCrawlWebsite(t *testing.T) {
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := CrawlWebsite(tc.website, tc.concurrentGoroutines)
			if err != nil {
				t.Fatalf("\n Error during testing %v ", err)
			}
		})
	}
}

var ttReq = []struct {
	name    string
	website string
}{
	{"Guimp.Com", "http://www.guimp.com/"},
	{"Webscraper.io", "http://webscraper.io"},
}

func TestGetHTTPGETRequest(t *testing.T) {
	for _, tc := range ttReq {
		t.Run(tc.name, func(t *testing.T) {
			var req = getHTTPGETRequest(tc.website)
			httpClient := getHTTPClient()

			res, err := httpClient.Do(req)
			if err != nil {
				t.Fatalf("Could not read response from site %v : %v \n", tc.website, err)
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				t.Errorf("Expected Status OK : got %v", res.Status)
			}
		})
	}
}
