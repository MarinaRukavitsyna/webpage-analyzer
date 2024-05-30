package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"webpage-analyzer/utils"
)

type AnalysisRequest struct {
	URL string `json:"url"`
}

type AnalysisResult struct {
	HTMLVersion          string         `json:"htmlVersion"`
	PageTitle            string         `json:"pageTitle"`
	Headings             map[string]int `json:"headings"`
	NumInternalLinks     int            `json:"numInternalLinks"`
	NumExternalLinks     int            `json:"numExternalLinks"`
	NumInaccessibleLinks int            `json:"numInaccessibleLinks"`
	IsContainLoginForm   bool           `json:"isContainLoginForm"`
	ErrorMessage         string         `json:"error,omitempty"`
}

func AnalyzePage(analyzeURL utils.AnalyzeURLFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AnalysisRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
			return
		}

		parsedURL, err := url.ParseRequestURI(req.URL)
		if err != nil {
			http.Error(w, `{"error":"Invalid URL format"}`, http.StatusBadRequest)
			return
		}

		result, err := analyzeURL(parsedURL.String())
		if err != nil {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
