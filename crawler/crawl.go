package crawler

import (
	"fmt"
	"net/url"
)

func (c *CrawlerConfig) CrawlAllPages() {
	c.wg.Add(1)
	go c.crawlPage(c.baseUrl.Hostname())
	c.wg.Wait()
}

func (c *CrawlerConfig) crawlPage(rawCurrentURL string) {
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - crawlPage: couldn't parse URL '%s': %v\n", rawCurrentURL, err)
		return
	}

	// skip other websites
	if currentURL.Hostname() != c.baseUrl.Hostname() {
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - normalizedURL: %v", err)
		return
	}

	// increment if visited
	if _, visited := c.pages[normalizedURL]; visited {
		c.pages[normalizedURL]++
		return
	}

	// mark as visited
	c.pages[normalizedURL] = 1

	fmt.Printf("crawling %s\n", rawCurrentURL)

	htmlBody, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - getHTML: %v", err)
		return
	}

	nextURLs, err := getURLsFromHTML(htmlBody, c.baseUrl.Hostname())
	if err != nil {
		fmt.Printf("Error - getURLsFromHTML: %v", err)
		return
	}

	for _, nextURL := range nextURLs {
		c.crawlPage(nextURL)
	}
}
