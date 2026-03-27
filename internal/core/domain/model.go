package domain

// AnalysisResult is the main result returned after analyzing a webpage
type AnalysisResult struct {
    HTMLVersion string            `json:"html_version"`
    Title       string            `json:"title"`
    Headings    map[string]int    `json:"headings"`   // h1, h2, etc.
    Links       LinkStats         `json:"links"`
    LoginForm   bool              `json:"login_form"`
}

// LinkStats contains counts of internal, external, and broken links
type LinkStats struct {
    Internal int `json:"internal"`
    External int `json:"external"`
    Broken   int `json:"broken"`
}