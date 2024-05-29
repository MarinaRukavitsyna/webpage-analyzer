package utils

import (
	"errors"
	"net/http"

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

func AnalyzeURL(url string) (AnalysisResult, error) {
	resp, err := http.Get(url)
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

	result := AnalysisResult{
		HTMLVersion:       getHTMLVersion(doc),
		PageTitle:         getPageTitle(doc),
		Headings:          getHeadings(doc),
		InternalLinks:     getInternalLinks(doc, url),
		ExternalLinks:     getExternalLinks(doc),
		InaccessibleLinks: getInaccessibleLinks(doc),
		ContainsLoginForm: containsLoginForm(doc),
	}

	return result, nil
}

func getHeadings(doc *html.Node) map[string]int {
	panic("unimplemented")
}

func getInaccessibleLinks(doc *html.Node) int {
	panic("unimplemented")
}

func getExternalLinks(doc *html.Node) int {
	panic("unimplemented")
}

func containsLoginForm(doc *html.Node) bool {
	panic("unimplemented")
}

func getInternalLinks(doc *html.Node, url string) int {
	panic("unimplemented")
}

func getPageTitle(doc *html.Node) string {
	panic("unimplemented")
}

func getHTMLVersion(doc *html.Node) string {
	panic("unimplemented")
}
