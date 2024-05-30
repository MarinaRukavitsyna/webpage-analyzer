package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"webpage-analyzer/utils"
)

// Mock the AnalyzeURL function in utils package
func mockAnalyzeURL(urlStr string) (utils.AnalysisResult, error) {
	switch urlStr {
	case "https://example.com/html5":
		return utils.AnalysisResult{
			HTMLVersion:          "HTML5",
			PageTitle:            "Test Page",
			Headings:             map[string]int{"h1": 1, "h2": 1},
			NumInternalLinks:     1,
			NumExternalLinks:     1,
			NumInaccessibleLinks: 0,
			IsContainLoginForm:   true,
		}, nil
	case "https://example.com/no-headings":
		return utils.AnalysisResult{
			HTMLVersion:          "HTML5",
			PageTitle:            "No Headings",
			Headings:             map[string]int{},
			NumInternalLinks:     1,
			NumExternalLinks:     0,
			NumInaccessibleLinks: 0,
			IsContainLoginForm:   false,
		}, nil
	case "https://example.com/inaccessible-link":
		return utils.AnalysisResult{
			HTMLVersion:          "HTML5",
			PageTitle:            "Inaccessible Link",
			Headings:             map[string]int{},
			NumInternalLinks:     0,
			NumExternalLinks:     1,
			NumInaccessibleLinks: 1,
			IsContainLoginForm:   false,
		}, nil
	case "https://example.com/html4":
		return utils.AnalysisResult{
			HTMLVersion:          "HTML 4.01 Strict",
			PageTitle:            "HTML 4.01 Strict",
			Headings:             map[string]int{"h1": 1, "h2": 1},
			NumInternalLinks:     0,
			NumExternalLinks:     0,
			NumInaccessibleLinks: 0,
			IsContainLoginForm:   false,
		}, nil
	default:
		return utils.AnalysisResult{}, errors.New("unknown URL")
	}
}

func TestAnalyzePage(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedBody   AnalysisResult
	}{
		{
			name:           "HTML5 Document",
			requestBody:    `{"url":"https://example.com/html5"}`,
			expectedStatus: http.StatusOK,
			expectedBody: AnalysisResult{
				HTMLVersion:          "HTML5",
				PageTitle:            "Test Page",
				Headings:             map[string]int{"h1": 1, "h2": 1},
				NumInternalLinks:     1,
				NumExternalLinks:     1,
				NumInaccessibleLinks: 0,
				IsContainLoginForm:   true,
			},
		},
		{
			name:           "No Headings",
			requestBody:    `{"url":"https://example.com/no-headings"}`,
			expectedStatus: http.StatusOK,
			expectedBody: AnalysisResult{
				HTMLVersion:          "HTML5",
				PageTitle:            "No Headings",
				Headings:             map[string]int{},
				NumInternalLinks:     1,
				NumExternalLinks:     0,
				NumInaccessibleLinks: 0,
				IsContainLoginForm:   false,
			},
		},
		{
			name:           "Inaccessible Link",
			requestBody:    `{"url":"https://example.com/inaccessible-link"}`,
			expectedStatus: http.StatusOK,
			expectedBody: AnalysisResult{
				HTMLVersion:          "HTML5",
				PageTitle:            "Inaccessible Link",
				Headings:             map[string]int{},
				NumInternalLinks:     0,
				NumExternalLinks:     1,
				NumInaccessibleLinks: 1,
				IsContainLoginForm:   false,
			},
		},
		{
			name:           "HTML 4.01 Strict",
			requestBody:    `{"url":"https://example.com/html4"}`,
			expectedStatus: http.StatusOK,
			expectedBody: AnalysisResult{
				HTMLVersion:          "HTML 4.01 Strict",
				PageTitle:            "HTML 4.01 Strict",
				Headings:             map[string]int{"h1": 1, "h2": 1},
				NumInternalLinks:     0,
				NumExternalLinks:     0,
				NumInaccessibleLinks: 0,
				IsContainLoginForm:   false,
			},
		},
		{
			name:           "Invalid URL Format",
			requestBody:    `{"url":"invalid-url"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody: AnalysisResult{
				ErrorMessage: "Invalid URL format",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/analyze", bytes.NewBufferString(tt.requestBody))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			handler := AnalyzePage(mockAnalyzeURL)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			var result AnalysisResult
			err = json.NewDecoder(rr.Body).Decode(&result)
			if err != nil {
				t.Fatal(err)
			}

			if result.HTMLVersion != tt.expectedBody.HTMLVersion {
				t.Errorf("HTMLVersion = %v, want %v", result.HTMLVersion, tt.expectedBody.HTMLVersion)
			}
			if result.PageTitle != tt.expectedBody.PageTitle {
				t.Errorf("PageTitle = %v, want %v", result.PageTitle, tt.expectedBody.PageTitle)
			}
			if !equalMaps(result.Headings, tt.expectedBody.Headings) {
				t.Errorf("Headings = %v, want %v", result.Headings, tt.expectedBody.Headings)
			}
			if result.NumInternalLinks != tt.expectedBody.NumInternalLinks {
				t.Errorf("NumInternalLinks = %v, want %v", result.NumInternalLinks, tt.expectedBody.NumInternalLinks)
			}
			if result.NumExternalLinks != tt.expectedBody.NumExternalLinks {
				t.Errorf("NumExternalLinks = %v, want %v", result.NumExternalLinks, tt.expectedBody.NumExternalLinks)
			}
			if result.NumInaccessibleLinks != tt.expectedBody.NumInaccessibleLinks {
				t.Errorf("NumInaccessibleLinks = %v, want %v", result.NumInaccessibleLinks, tt.expectedBody.NumInaccessibleLinks)
			}
			if result.IsContainLoginForm != tt.expectedBody.IsContainLoginForm {
				t.Errorf("IsContainLoginForm = %v, want %v", result.IsContainLoginForm, tt.expectedBody.IsContainLoginForm)
			}
			if result.ErrorMessage != tt.expectedBody.ErrorMessage {
				t.Errorf("ErrorMessage = %v, want %v", result.ErrorMessage, tt.expectedBody.ErrorMessage)
			}
		})
	}
}

// Helper function to compare two maps
func equalMaps(a, b map[string]int) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
