package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func normalizeURL(input string) (string, error) {
	parsed, err := url.Parse(input)
	if err != nil {
		return "", err
	}

	return parsed.Host + strings.TrimSuffix(parsed.Path, "/"), nil
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
	if contentType != "text/html" {
		return "", fmt.Errorf("server responded with Content-Type %s", contentType)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseUrl := args[0]
	fmt.Println("starting crawl of: " + baseUrl)

	html, err := getHTML(baseUrl)
	if err != nil {
		panic(fmt.Sprintf("Unable to retrieve HTML for URL %s", baseUrl))
	}

	fmt.Println(html)
}
