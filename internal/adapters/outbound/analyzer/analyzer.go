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
func AnalyzeDocument(doc *html.Node, baseURL string) *AnalysisData {
    data := &AnalysisData{
        Headings: make(map[string]int),
        Links:    []string{},
    }

    var f func(*html.Node)
    f = func(n *html.Node) {
        if n.Type == html.ElementNode {
            switch n.Data {
            case "html":
                // Simple check for HTML version
                data.HTMLVersion = detectHTMLVersion(n)
            case "title":
                if n.FirstChild != nil {
                    data.Title = n.FirstChild.Data
                }
            case "h1", "h2", "h3", "h4", "h5", "h6":
                data.Headings[n.Data]++
            case "form":
                if containsPasswordInput(n) {
                    data.LoginForm = true
                }
            case "a":
                for _, attr := range n.Attr {
                    if attr.Key == "href" && attr.Val != "" {
                        link := normalizeURL(attr.Val, baseURL)
                        data.Links = append(data.Links, link)
                    }
                }
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
    }

    f(doc)
    return data
}

// detectHTMLVersion tries to detect HTML version based on doctype
func detectHTMLVersion(n *html.Node) string {
    // For simplicity, assume HTML5 if no version detected
    for c := n.Parent; c != nil; c = c.PrevSibling {
        if c.Type == html.DoctypeNode {
            dt := strings.ToLower(c.Data)
            if strings.Contains(dt, "html 4.01") {
                return "HTML 4.01"
            }
            if strings.Contains(dt, "xhtml") {
                return "XHTML"
            }
        }
    }
    return "HTML5"
}

// containsPasswordInput checks if a form has an <input type="password">
func containsPasswordInput(form *html.Node) bool {
    var found bool
    var f func(*html.Node)
    f = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "input" {
            for _, attr := range n.Attr {
                if attr.Key == "type" && strings.ToLower(attr.Val) == "password" {
                    found = true
                    return
                }
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
    }
    f(form)
    return found
}

// normalizeURL resolves relative URLs to absolute using base URL
func normalizeURL(href, base string) string {
    baseParsed, err := url.Parse(base)
    if err != nil {
        return href
    }
    hrefParsed, err := url.Parse(href)
    if err != nil {
        return href
    }
    return baseParsed.ResolveReference(hrefParsed).String()
}