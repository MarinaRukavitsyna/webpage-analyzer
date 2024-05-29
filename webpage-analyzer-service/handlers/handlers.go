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

	result, err := utils.AnalyzeURL(parsedURL.String())
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
