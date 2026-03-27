package analyzer

import (
    "net/url"
    "strings"

    "golang.org/x/net/html"
)

// AnalysisData holds analysis results from a web page
type AnalysisData struct {
    HTMLVersion string
    Title       string
    Headings    map[string]int
    Links       []string
    LoginForm   bool
}

// AnalyzeDocument performs the HTML analysis (pure logic)


// detectHTMLVersion tries to detect HTML version based on doctype


// containsPasswordInput checks if a form has an <input type="password">


// normalizeURL resolves relative URLs to absolute using base URL


