package cmd

import (
	"testing"
)

func TestCrawlWebsite(t *testing.T) {
	var urls = []string{
		"webscraper.io",
	}
	for _, u := range urls {
		err := CrawlWebsite(u)
		if err != nil {
			t.Fail()
		}
	}
}
