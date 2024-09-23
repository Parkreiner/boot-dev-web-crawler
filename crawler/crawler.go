package crawler

import (
	"net/url"
	"sync"
)

type crawlerConfig struct {
	crawled            bool
	pages              map[string]int // Key is URLs, ints are frequency
	baseUrl            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{} // Will be buffered
	wg                 *sync.WaitGroup
}

func Configure(rawBaseUrl string, maxConcurrency int) (crawlerConfig, error) {
	baseUrl, err := url.Parse(rawBaseUrl)
	if err != nil {
		return crawlerConfig{}, err
	}

	newConfig := crawlerConfig{
		crawled:            false,
		pages:              map[string]int{},
		baseUrl:            baseUrl,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
	}

	return newConfig, nil
}
