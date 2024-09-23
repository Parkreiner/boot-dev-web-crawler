package main

import (
	"fmt"
	"os"
	"strconv"
	crawler "web_crawler/blogCrawler"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(args) == 1 {
		fmt.Println("Missing crawler configuration settings (maxConcurrency and maxPages)")
		os.Exit(1)
	}

	if len(args) == 2 {
		fmt.Println("Missing crawler setting maxPages")
		os.Exit(1)
	}

	if len(args) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	rawMaxConcurrency := args[1]
	maxConcurrency, err := strconv.Atoi(rawMaxConcurrency)
	if err != nil {
		fmt.Printf("%s is not an integer (maxConcurrency)", rawMaxConcurrency)
		os.Exit(1)
	}

	rawMaxPages := args[2]
	maxPages, err := strconv.Atoi(rawMaxPages)
	if err != nil {
		fmt.Printf("%s is not an integer (maxPages)", rawMaxPages)
		os.Exit(1)
	}

	baseUrl := args[0]
	cfg, err := crawler.Configure(baseUrl, maxConcurrency, maxPages)

	if err != nil {
		fmt.Printf("configuration error - %v", err)
		return
	}

	fmt.Println("starting crawl of: " + baseUrl)
	pages := cfg.CrawlAllPages()
	fmt.Println("Crawling complete")

	report, err := cfg.CrawlerReport()
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	fmt.Println(report)
	crawler.PrintReport(pages, baseUrl)
}
