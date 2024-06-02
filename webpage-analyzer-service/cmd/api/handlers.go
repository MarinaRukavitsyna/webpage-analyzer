package main

import (
	"log"
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
		log.Println(err)
		app.errorJSON(w, err)
		http.Error(w, `{"error":"Invalid input"}`, http.StatusBadRequest)
		return
	}

	parsedURL, err := url.ParseRequestURI(requestPayload.URL)
	if err != nil {
		http.Error(w, `{"error":"Invalid URL format"}`, http.StatusBadRequest)
		return
	}

	result, err := analyzer.AnalyzeURL(parsedURL.String())
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: result.HTMLVersion,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}
