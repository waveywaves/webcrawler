package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var concurrentGoroutines int

// RootCmd : Root of the cmdline tool
var RootCmd = &cobra.Command{
	Use:   "webcrawler [url]",
	Short: "Crawl a website to get a sitemap",
	Long:  "Crawl a website to get a sitemap",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Please pass the website you want to crawl")
		}
		return nil
	},
	Run: RunRootCmd,
}

// RunRootCmd : The command which will
func RunRootCmd(cmd *cobra.Command, args []string) {
	err := CrawlWebsite(args[0], concurrentGoroutines)
	if err != nil {
		os.Stderr.WriteString("Error at RunRootCmd : " + err.Error() + " \n")
	}
}

// Execute :
func Execute() {

	RootCmd.PersistentFlags().IntVar(&concurrentGoroutines, "n", 5, "number of concurrent Goroutines")

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
