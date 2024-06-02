package main

import (
	"net/http"
	"net/url"
	"webpage-analyzer/cmd/api/analyzer"
)

type AnalysisRequest struct {
	URL string `json:"url"`
}

func (app *Config) Analyzer(w http.ResponseWriter, r *http.Request) {
	var requestPayload AnalysisRequest

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	parsedURL, err := url.ParseRequestURI(requestPayload.URL)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	result, err := analyzer.AnalyzeURL(parsedURL.String())
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:          false,
		StatusCode:     http.StatusOK,
		Message:        "OK",
		AnalysisResult: result,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}
