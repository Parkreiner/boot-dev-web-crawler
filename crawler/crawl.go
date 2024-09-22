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

func (c *CrawlerConfig) addPageVisit(normalizedUrl string) (isFirst bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, visited := c.pages[normalizedUrl]; visited {
		c.pages[normalizedUrl]++
		return false
	}

	c.pages[normalizedUrl] = 1
	return true
}

func (c *CrawlerConfig) crawlPage(rawCurrentURL string) {
	c.concurrencyControl <- struct{}{}
	defer func() {
		<-c.concurrencyControl
		c.wg.Done()
	}()

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - crawlPage: couldn't parse URL '%s': %v\n", rawCurrentURL, err)
		return
	}

	// skip other websites
	if currentURL.Hostname() != c.baseUrl.Hostname() {
		return
	}

	normalizedUrl, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - normalizedURL: %v", err)
		return
	}

	isFirstVisit := c.addPageVisit(normalizedUrl)
	if !isFirstVisit {
		return
	}

	fmt.Printf("crawling %s\n", rawCurrentURL)

	htmlBody, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - getHTML: %v", err)
		return
	}

	nextURLs, err := getURLsFromHTML(htmlBody, c.baseUrl)
	if err != nil {
		fmt.Printf("Error - getURLsFromHTML: %v", err)
		return
	}

	for _, nextURL := range nextURLs {
		c.wg.Add(1)
		go c.crawlPage(nextURL)
	}
}
