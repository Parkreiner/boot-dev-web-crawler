package crawler

import (
	"net/url"
	"sync"
)

type CrawlerConfig struct {
	pages              map[string]int // Key is URLs, ints are frequency
	baseUrl            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{} // Will be buffered
	wg                 *sync.WaitGroup
}

func Configure(rawBaseUrl string, maxConcurrency int) (CrawlerConfig, error) {
	baseUrl, err := url.Parse(rawBaseUrl)
	if err != nil {
		return CrawlerConfig{}, err
	}

	newConfig := CrawlerConfig{
		pages:              map[string]int{},
		baseUrl:            baseUrl,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
	}

	return newConfig, nil
}
