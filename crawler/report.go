package crawler

import (
	"fmt"
	"sort"
)

type crawlerEntry struct {
	url       string
	frequency int
}

type crawlerEntryList []crawlerEntry

func (l crawlerEntryList) Len() int {
	return len(l)
}
func (l crawlerEntryList) Less(i int, j int) bool {
	return (l)[i].frequency < (l)[j].frequency
}
func (l crawlerEntryList) Swap(i int, j int) {
	l[i], l[j] = l[j], l[i]
}

func (c *CrawlerConfig) CrawlerReport() string {
	totalUrls := len(c.pages)
	serialized := make(crawlerEntryList, totalUrls)
	for url, freq := range c.pages {
		serialized = append(serialized, crawlerEntry{
			url:       url,
			frequency: freq,
		})
	}

	sort.Sort(sort.Reverse(serialized))
	report := fmt.Sprintf("Config has %d URLs :\n", totalUrls)
	for _, entry := range serialized {
		report += fmt.Sprintf("- %d - %s", entry.frequency, entry.url)
	}

	return report
}
