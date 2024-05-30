package utils

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type AnalysisResult struct {
	HTMLVersion       string
	PageTitle         string
	Headings          map[string]int
	InternalLinks     int
	ExternalLinks     int
	InaccessibleLinks int
	ContainsLoginForm bool
}

func AnalyzeURL(urlStr string) (AnalysisResult, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return AnalysisResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AnalysisResult{}, errors.New("Error fetching the URL: " + http.StatusText(resp.StatusCode))
	}

	// Parse the HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return AnalysisResult{}, err
	}

	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return AnalysisResult{}, err
	}

	result := AnalysisResult{
		HTMLVersion:       getHTMLVersion(doc),
		PageTitle:         getPageTitle(doc),
		Headings:          getHeadings(doc),
		InternalLinks:     getNumInternalLinks(doc, baseURL),
		ExternalLinks:     getNumExternalLinks(doc, baseURL),
		InaccessibleLinks: getNumInaccessibleLinks(doc),
		ContainsLoginForm: isContainLoginForm(doc),
	}

	return result, nil
}

// getHeadings returns a map of headings and their frequencies found in the provided HTML document
func getHeadings(doc *html.Node) map[string]int {
	headings := make(map[string]int)
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "h1", "h2", "h3", "h4", "h5", "h6":
				headings[n.Data]++
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)
	return headings
}

// isAccessible checks if a link is accessible
func isAccessible(link string) bool {
	resp, err := http.Head(link)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

// getNumInaccessibleLinks traverses the HTML document and returns the count of inaccessible links
func getNumInaccessibleLinks(doc *html.Node) int {
	var links []string
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	inaccessibleCount := 0
	for _, link := range links {
		if !isAccessible(link) {
			inaccessibleCount++
		}
	}
	return inaccessibleCount
}

// isExternalLink checks if a link is external
func isExternalLink(link string, baseURL *url.URL) bool {
	u, err := url.Parse(link)
	if err != nil {
		return false
	}
	return u.Host != "" && u.Host != baseURL.Host
}

// getNumExternalLinks returns the number of external links in the HTML document
func getNumExternalLinks(doc *html.Node, baseURL *url.URL) int {
	var links []string
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	externalCount := 0
	for _, link := range links {
		if isExternalLink(link, baseURL) {
			externalCount++
		}
	}
	return externalCount
}

// getNumInternalLinks returns the number of internal links in the HTML document
func getNumInternalLinks(doc *html.Node, baseURL *url.URL) int {
	var links []string
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	internalCount := 0
	for _, link := range links {
		if !isExternalLink(link, baseURL) {
			internalCount++
		}
	}
	return internalCount
}

// isContainLoginForm checks if the document contains a login form
func isContainLoginForm(doc *html.Node) bool {
	var traverse func(*html.Node) bool
	traverse = func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "form" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "input" {
					for _, attr := range c.Attr {
						if attr.Key == "type" && attr.Val == "password" {
							return true
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if traverse(c) {
				return true
			}
		}
		return false
	}
	return traverse(doc)
}

// getPageTitle extracts the global title of an HTML document from the provided html.Node
func getPageTitle(doc *html.Node) string {
	var title string
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil {
				title = n.FirstChild.Data
			}
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)
	return title
}

// getHTMLVersion determines the version of the HTML document by inspecting the doctype
// doctypes are taken from https://www.w3.org/QA/2002/04/valid-dtd-list.html
func getHTMLVersion(doc *html.Node) string {
	if doc.Type == html.DoctypeNode {
		doctype := strings.ToLower(doc.Data)
		if doctype == "html" && (doc.Attr == nil || len(doc.Attr) == 0) {
			return "HTML5"
		}

		if len(doc.Attr) > 0 {
			doctype = strings.ToLower(doc.Attr[0].Val)
		}

		switch {
		case strings.Contains(doctype, "4.01 transitional"):
			return "HTML 4.01 Transitional"
		case strings.Contains(doctype, "4.01 frameset"):
			return "HTML 4.01 Frameset"
		case strings.Contains(doctype, "4.01"):
			return "HTML 4.01 Strict"
		case strings.Contains(doctype, "xhtml 1.0 strict"):
			return "XHTML 1.0 Strict"
		case strings.Contains(doctype, "xhtml 1.0 transitional"):
			return "XHTML 1.0 Transitional"
		case strings.Contains(doctype, "xhtml 1.0 frameset"):
			return "XHTML 1.0 Frameset"
		case strings.Contains(doctype, "xhtml basic 1.1"):
			return "XHTML Basic 1.1"
		case strings.Contains(doctype, "xhtml 1.1"):
			return "XHTML 1.1"
		default:
			return "Unknown"
		}
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		if version := getHTMLVersion(c); version != "" {
			return version
		}
	}
	return "Unknown"
}
