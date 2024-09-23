package main

import (
	"fmt"
	"os"
	"web_crawler/crawler"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	const maxConcurrency = 10
	baseUrl := args[0]
	cfg, err := crawler.Configure(baseUrl, maxConcurrency)

	if err != nil {
		fmt.Printf("configuration error - %v", err)
		return
	}

	fmt.Println("starting crawl of: " + baseUrl)
	cfg.CrawlAllPages()
	fmt.Println("Crawling complete")

	report, err := cfg.CrawlerReport()
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	fmt.Println(report)
}
