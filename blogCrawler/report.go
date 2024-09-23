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

// This version of the function is required for the Boot.dev challenge
func PrintReport(pages map[string]int, baseUrl string) {
	fmt.Println("=============================")
	fmt.Printf("  REPORT for %s\n", baseUrl)
	fmt.Println("=============================")

	type crawlerReportEntry struct {
		url       string
		frequency int
	}

	serialized := make([]crawlerReportEntry, 0, len(pages))
	for url, freq := range pages {
		serialized = append(serialized, crawlerReportEntry{
			url:       url,
			frequency: freq,
		})
	}

	slices.SortStableFunc(
		serialized,
		func(a crawlerReportEntry, b crawlerReportEntry) int {
			freqDelta := b.frequency - a.frequency
			if freqDelta != 0 {
				return freqDelta
			}

			if a.url < b.url {
				return -1
			}

			if a.url > b.url {
				return 1
			}

			return 0
		},
	)

	for _, entry := range serialized {
		fmt.Printf("Found %d internal links to %s\n", entry.frequency, entry.url)
	}
}
