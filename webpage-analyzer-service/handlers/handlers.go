package handlers

import (
	"net/http"
)

type AnalysisRequest struct {
	URL string `json:"url"`
}

type AnalysisResult struct {
	HTMLVersion       string         `json:"htmlVersion"`
	PageTitle         string         `json:"pageTitle"`
	Headings          map[string]int `json:"headings"`
	InternalLinks     int            `json:"internalLinks"`
	ExternalLinks     int            `json:"externalLinks"`
	InaccessibleLinks int            `json:"inaccessibleLinks"`
	ContainsLoginForm bool           `json:"containsLoginForm"`
	ErrorMessage      string         `json:"error,omitempty"`
}

func AnalyzePage(w http.ResponseWriter, r *http.Request) {

}
