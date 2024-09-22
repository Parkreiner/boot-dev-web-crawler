package main

import (
	"fmt"
	"normalize_url/crawler"
	"os"
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

	const maxConcurrency = 3
	baseUrl := args[0]
	cfg, err := crawler.Configure(baseUrl, maxConcurrency)

	if err != nil {
		fmt.Printf("configuration error - %v", err)
		return
	}

	fmt.Println("starting crawl of: " + baseUrl)
	cfg.CrawlAllPages()
	fmt.Println("Crawling complete")

	fmt.Println(cfg.CrawlerReport())
}
