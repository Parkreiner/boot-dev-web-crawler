package crawler

import (
	"fmt"
	"slices"
)

type crawlerReportEntry struct {
	url       string
	frequency int
}

func (c *CrawlerConfig) CrawlerReport() string {
	totalUrls := len(c.pages)
	serialized := make([]crawlerReportEntry, totalUrls)
	for url, freq := range c.pages {
		serialized = append(serialized, crawlerReportEntry{
			url:       url,
			frequency: freq,
		})
	}

	slices.SortStableFunc(serialized, func(a, b crawlerReportEntry) int {
		return b.frequency - a.frequency
	})

	report := fmt.Sprintf("Config has %d URLs :\n", totalUrls)
	for _, entry := range serialized {
		report += fmt.Sprintf("- %d - %s", entry.frequency, entry.url)
	}

	return report
}
