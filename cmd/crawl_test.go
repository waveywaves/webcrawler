package cmd

import (
	"testing"
)

func TestCrawlWebsite(t *testing.T) {
	var urls = []string{
		"google.com",
		"webscraper.io",
	}
	for _, u := range urls {
		err := CrawlWebsite(u)
