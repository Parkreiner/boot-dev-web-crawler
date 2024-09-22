package crawler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func normalizeURL(input string) (string, error) {
	parsed, err := url.Parse(input)
	if err != nil {
		return "", err
	}

	normalized := parsed.Scheme + "://" + parsed.Host + strings.TrimSuffix(parsed.Path, "/")
	return normalized, nil
}

func getHTML(rawUrl string) (string, error) {
	res, err := http.Get(rawUrl)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode > 400 {
		return "", fmt.Errorf("server responded with code %d", res.StatusCode)
	}

	contentType := res.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("server responded with Content-Type %s", contentType)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getURLsFromHTML(htmlBody string, rawBaseURL string) ([]string, error) {
	rootNode, err := html.Parse(strings.NewReader((htmlBody)))
	if err != nil {
		return []string{}, err
	}

	allHrefs := []string{}

	var traverseNodes func(node *html.Node)
	traverseNodes = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key != "href" {
					continue
				}

				href := attr.Val
				hrefRunes := []rune(href)
				if len(hrefRunes) == 0 || hrefRunes[0] == '/' {
					href = rawBaseURL + href
				}

				allHrefs = append(allHrefs, href)
			}
		}

		for curr := node.FirstChild; curr != nil; curr = curr.NextSibling {
			traverseNodes(curr)
		}
	}
	traverseNodes(rootNode)

	hrefTracker := map[string]struct{}{}
	uniqueHrefs := []string{}

	for _, href := range allHrefs {
		_, found := hrefTracker[href]
		if !found {
			uniqueHrefs = append(uniqueHrefs, href)
			hrefTracker[href] = struct{}{}
		}
	}

	return uniqueHrefs, nil
}
