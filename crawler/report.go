package crawler

import (
	"errors"
	"fmt"
	"slices"
)

// Generates a stringified representation of
func (c *crawlerConfig) CrawlerReport() (string, error) {
	if !c.crawled {
		return "", errors.New("cannot generate a report for a crawler that has not crawled yet")
	}

	totalUrls := len(c.pages)
	if totalUrls == 0 {
		return "No URLs to report", nil
	}

	type crawlerReportEntry struct {
		url       string
		frequency int
	}

	serialized := make([]crawlerReportEntry, 0, totalUrls)
	for url, freq := range c.pages {
		serialized = append(serialized, crawlerReportEntry{
			url:       url,
			frequency: freq,
		})
	}

	slices.SortStableFunc(
		serialized,
		func(a crawlerReportEntry, b crawlerReportEntry) int {
			return b.frequency - a.frequency
		},
	)

	report := fmt.Sprintf("config has %d URLs:\n", totalUrls)
	for _, entry := range serialized {
		report += fmt.Sprintf("- %d: %s\n", entry.frequency, entry.url)
	}

	return report, nil
}
