package crawler

import (
	"fmt"
	"net/url"
)

func (c *crawlerConfig) CrawlAllPages() map[string]int {
	c.crawled = true

	c.wg.Add(1)
	go c.crawlPage(c.baseUrl.String())
	c.wg.Wait()

	return c.pages
}

func (c *crawlerConfig) addPageVisit(normalizedUrl string) (isFirst bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, visited := c.pages[normalizedUrl]; visited {
		c.pages[normalizedUrl]++
		return false
	}

	c.pages[normalizedUrl] = 1
	return true
}

func (c *crawlerConfig) PageCount() int {
	c.mu.Lock()
	pageCount := len(c.pages)
	c.mu.Unlock()

	return pageCount
}

func (c *crawlerConfig) crawlPage(rawCurrentURL string) {
	c.concurrencyControl <- struct{}{}
	defer func() {
		<-c.concurrencyControl
		c.wg.Done()
	}()

	if c.PageCount() > c.maxPages {
		return
	}

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
