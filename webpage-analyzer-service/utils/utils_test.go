package utils

import (
	"net/url"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

// Helper function to parse HTML from a string
func parseHTML(input string) (*html.Node, error) {
	return html.Parse(strings.NewReader(input))
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

// Table-driven tests for getHeadings
func TestGetHeadings(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected map[string]int
	}{
		{
			name:     "NoHeadings",
			html:     `<html><head><title>Test</title></head><body></body></html>`,
			expected: map[string]int{},
		},
		{
			name: "SingleHeading",
			html: `<html><head><title>Test</title></head><body><h1>Heading 1</h1></body></html>`,
			expected: map[string]int{
				"h1": 1,
			},
		},
		{
			name: "MultipleHeadings",
			html: `<html><head><title>Test</title></head><body>
				<h1>Heading 1</h1>
				<h2>Heading 2</h2>
				<h3>Heading 3</h3>
				<h1>Another Heading 1</h1>
			</body></html>`,
			expected: map[string]int{
				"h1": 2,
				"h2": 1,
				"h3": 1,
			},
		},
		{
			name: "NestedHeadings",
			html: `<html><head><title>Test</title></head><body>
				<div><h1>Heading 1</h1></div>
				<section>
					<h2>Heading 2</h2>
					<div><h3>Heading 3</h3></div>
				</section>
			</body></html>`,
			expected: map[string]int{
				"h1": 1,
				"h2": 1,
				"h3": 1,
			},
		},
		{
			name: "AllHeadingLevels",
			html: `<html><head><title>Test</title></head><body>
				<h1>Heading 1</h1>
				<h2>Heading 2</h2>
				<h3>Heading 3</h3>
				<h4>Heading 4</h4>
				<h5>Heading 5</h5>
				<h6>Heading 6</h6>
			</body></html>`,
			expected: map[string]int{
				"h1": 1,
				"h2": 1,
				"h3": 1,
				"h4": 1,
				"h5": 1,
				"h6": 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parseHTML(tt.html)
			if err != nil {
				t.Fatalf("Failed to parse HTML: %v", err)
			}
			got := getHeadings(doc)
			if !equalMaps(got, tt.expected) {
				t.Errorf("getHeadings() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Table-driven tests for getNumInaccessibleLinks
func TestGetNumInaccessibleLinks(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected int
	}{
		{
			name:     "NoLinks",
			html:     `<html><head><title>Test</title></head><body></body></html>`,
			expected: 0,
		},
		{
			name: "AllAccessibleLinks",
			html: `<html><head><title>Test</title></head><body>
				<a href="https://www.google.com">Google</a>
				<a href="https://www.example.com">Example</a>
			</body></html>`,
			expected: 0,
		},
		{
			name: "SomeInaccessibleLinks",
			html: `<html><head><title>Test</title></head><body>
				<a href="https://www.google.com">Google</a>
				<a href="https://www.nonexistentwebsite.com">Nonexistent</a>
			</body></html>`,
			expected: 1,
		},
		{
			name: "AllInaccessibleLinks",
			html: `<html><head><title>Test</title></head><body>
				<a href="https://www.nonexistentwebsite1.com">Nonexistent1</a>
				<a href="https://www.nonexistentwebsite2.com">Nonexistent2</a>
			</body></html>`,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parseHTML(tt.html)
			if err != nil {
				t.Fatalf("Failed to parse HTML: %v", err)
			}
			got := getNumInaccessibleLinks(doc)
			if got != tt.expected {
				t.Errorf("getInaccessibleLinks() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Table-driven tests for getNumExternalLinks
func TestGetNumExternalLinks(t *testing.T) {
	baseURL, _ := url.Parse("https://www.example.com")

	tests := []struct {
		name     string
		html     string
		expected int
	}{
		{
			name:     "NoLinks",
			html:     `<html><head><title>Test</title></head><body></body></html>`,
			expected: 0,
		},
		{
			name: "NoExternalLinks",
			html: `<!DOCTYPE html><html><head><title>Test</title></head><body>
				<a href="/internal-link">Internal Link</a>
				<a href="/another-internal-link">Another Internal Link</a>
			</body></html>`,
			expected: 0,
		},
		{
			name: "SomeExternalLinks",
			html: `<!DOCTYPE html><html><head><title>Test</title></head><body>
				<a href="https://www.google.com">Google</a>
				<a href="/internal-link">Internal Link</a>
				<a href="https://www.example.org">Example</a>
				<a href="https://www.example.com/InternalLink">Internal Link Example</a>
			</body></html>`,
			expected: 2,
		},
		{
			name: "AllExternalLinks",
			html: `<html><head><title>Test</title></head><body>
				<h1>Example Domain</h1>
				<p>This domain is for use in illustrative examples in documents. You may use this
				domain in literature without prior coordination or asking for permission.</p>
				<p><a href="https://www.iana.org/domains/example">More information...</a></p>
				<a href="https://www.google.com">Google</a>
				<a href="https://www.example.org">Example</a>
			</body></html>`,
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parseHTML(tt.html)
			if err != nil {
				t.Fatalf("Failed to parse HTML: %v", err)
			}
			got := getNumExternalLinks(doc, baseURL)
			if got != tt.expected {
				t.Errorf("getExternalLinks() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Table-driven tests for getNumInternalLinks
func TestGetNumInternalLinks(t *testing.T) {
	baseURL, _ := url.Parse("https://www.example.com")

	tests := []struct {
		name     string
		html     string
		expected int
	}{
		{
			name:     "NoLinks",
			html:     `<!DOCTYPE html><html><head><title>Test</title></head><body></body></html>`,
			expected: 0,
		},
		{
			name: "NoInternalLinks",
			html: `<html><head><title>Test</title></head><body>
				<a href="https://www.google.com">Google</a>
				<a href="https://www.example.org">Example</a>
			</body></html>`,
			expected: 0,
		},
		{
			name: "SomeInternalLinks",
			html: `<!DOCTYPE html><html><head><title>Test</title></head><body>
				<a href="/internal-link">Internal Link</a>
				<a href="https://www.example.com/another-internal-link">Another Internal Link</a>
				<a href="https://www.google.com">Google</a>
			</body></html>`,
			expected: 2,
		},
		{
			name: "AllInternalLinks",
			html: `<!DOCTYPE html><html><head><title>Test</title></head><body>
				<a href="/internal-link">Internal Link</a>
				<a href="https://www.example.com/another-internal-link">Another Internal Link</a>
			</body></html>`,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parseHTML(tt.html)
			if err != nil {
				t.Fatalf("Failed to parse HTML: %v", err)
			}
			got := getNumInternalLinks(doc, baseURL)
			if got != tt.expected {
				t.Errorf("getInternalLinks() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Table-driven tests for isContainLoginForm
func TestIsContainLoginForm(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected bool
	}{
		{
			name:     "NoForms",
			html:     `<html><head><title>Test</title></head><body></body></html>`,
			expected: false,
		},
		{
			name: "FormWithoutPassword",
			html: `<html><head><title>Test</title></head><body>
				<form>
					<input type="text" name="username">
					<input type="submit" value="Submit">
				</form>
			</body></html>`,
			expected: false,
		},
		{
			name: "FormWithPassword",
			html: `<html><head><title>Test</title></head><body>
				<form>
					<input type="text" name="username">
					<input type="password" name="password">
					<input type="submit" value="Login">
				</form>
			</body></html>`,
			expected: true,
		},
		{
			name: "NestedFormWithPassword",
			html: `<html><head><title>Test</title></head><body>
				<div>
					<form>
						<input type="text" name="username">
						<input type="password" name="password">
						<input type="submit" value="Login">
					</form>
				</div>
			</body></html>`,
			expected: true,
		},
		{
			name: "MultipleFormsOneWithPassword",
			html: `<html><head><title>Test</title></head><body>
				<form>
					<input type="text" name="username">
					<input type="submit" value="Submit">
				</form>
				<form>
					<input type="text" name="username">
					<input type="password" name="password">
					<input type="submit" value="Login">
				</form>
			</body></html>`,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parseHTML(tt.html)
			if err != nil {
				t.Fatalf("Failed to parse HTML: %v", err)
			}
			got := isContainLoginForm(doc)
			if got != tt.expected {
				t.Errorf("containsLoginForm() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Table-driven tests for getPageTitle
func TestGetPageTitle(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "BasicTitle",
			html:     `<html><head><title>Test Title</title></head><body></body></html>`,
			expected: "Test Title",
		},
		{
			name:     "NoTitle",
			html:     `<html><head></head><body></body></html>`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parseHTML(tt.html)
			if err != nil {
				t.Fatalf("Failed to parse HTML: %v", err)
			}
			got := getPageTitle(doc)
			if got != tt.expected {
				t.Errorf("getPageTitle() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Table-driven tests for getHTMLVersion
func TestGetHTMLVersion(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "HTML5",
			html:     `<!DOCTYPE html><html><head><title>Test</title></head><body></body></html>`,
			expected: "HTML5",
		},
		{
			name: "XHTML 1.0 Transitional",
			html: `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
				<html xmlns="http://www.w3.org/1999/xhtml"><head><title>Test</title></head><body></body></html>`,
			expected: "XHTML 1.0 Transitional",
		},
		{
			name: "HTML 4.01 Transitional",
			html: `<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
				<html><head><title>Test</title></head><body></body></html>`,
			expected: "HTML 4.01 Transitional",
		},
		{
			name: "XHTML 1.0 Frameset",
			html: `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Frameset//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-frameset.dtd">
				<html><head><title>Test</title></head><body></body></html>`,
			expected: "XHTML 1.0 Frameset",
		},
		{
			name: "XHTML 1.1",
			html: `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">
				<html><head><title>Test</title></head><body></body></html>`,
			expected: "XHTML 1.1",
		},
		{
			name: "XHTML Basic 1.1",
			html: `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML Basic 1.1//EN" "http://www.w3.org/TR/xhtml-basic/xhtml-basic11.dtd">
				<html><head><title>Test</title></head><body></body></html>`,
			expected: "XHTML Basic 1.1",
		},
		{
			name: "Unknown DOCTYPE",
			html: `<!DOCTYPE something unknown>
				<html><head><title>Test</title></head><body></body></html>`,
			expected: "Unknown",
		},
		{
			name:     "No DOCTYPE",
			html:     `<html><head><title>Test</title></head><body></body></html>`,
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parseHTML(tt.html)
			if err != nil {
				t.Fatalf("Failed to parse HTML: %v", err)
			}
			got := getHTMLVersion(doc)
			if got != tt.expected {
				t.Errorf("getHTMLVersion() = %v, want %v", got, tt.expected)
			}
		})
	}
}
