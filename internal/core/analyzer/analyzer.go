package analyzer

import (
    "net/url"
    "strings"

    "golang.org/x/net/html"
)

// AnalysisData holds extracted information from a webpage
type AnalysisData struct {
    HTMLVersion string
    Title       string
    Headings    map[string]int
    Links       []string
    LoginForm   bool
}

// AnalyzeDocument traverses the HTML tree and extracts required data
func AnalyzeDocument(doc *html.Node, baseURL string) *AnalysisData {
    data := &AnalysisData{
        Headings: make(map[string]int),
        Links:    []string{},
    }

    var traverse func(*html.Node)
    traverse = func(n *html.Node) {
        if n.Type == html.ElementNode {
            switch n.Data {

            // Detect HTML version
            case "html":
                data.HTMLVersion = detectHTMLVersion(n)

            // Extract title
            case "title":
                if n.FirstChild != nil {
                    data.Title = strings.TrimSpace(n.FirstChild.Data)
                }

            // Count headings
            case "h1", "h2", "h3", "h4", "h5", "h6":
                data.Headings[n.Data]++

            // Detect login form
            case "form":
                if containsPasswordInput(n) {
                    data.LoginForm = true
                }

            // Extract links
            case "a":
                for _, attr := range n.Attr {
                    if attr.Key == "href" && attr.Val != "" {
                        link := normalizeURL(attr.Val, baseURL)
                        data.Links = append(data.Links, link)
                    }
                }
            }
        }

        // Traverse children
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            traverse(c)
        }
    }

    traverse(doc)
    return data
}

//////////////////////////////////////////////////////////
// 🔽 Helper Functions (Pure Logic)
//////////////////////////////////////////////////////////

// detectHTMLVersion inspects doctype node to determine HTML version
func detectHTMLVersion(n *html.Node) string {
    // Walk upwards and check for doctype
    for p := n.Parent; p != nil; p = p.Parent {
        for c := p.FirstChild; c != nil; c = c.NextSibling {
            if c.Type == html.DoctypeNode {
                dt := strings.ToLower(c.Data)

                if strings.Contains(dt, "html 4.01") {
                    return "HTML 4.01"
                }
                if strings.Contains(dt, "xhtml") {
                    return "XHTML"
                }
                return "HTML5"
            }
        }
    }

    // Default assumption
    return "HTML5"
}

// containsPasswordInput checks if a form contains <input type="password">
func containsPasswordInput(form *html.Node) bool {
    var found bool

    var traverse func(*html.Node)
    traverse = func(n *html.Node) {
        if found {
            return
        }

        if n.Type == html.ElementNode && n.Data == "input" {
            for _, attr := range n.Attr {
                if attr.Key == "type" && strings.ToLower(attr.Val) == "password" {
                    found = true
                    return
                }
            }
        }

        for c := n.FirstChild; c != nil; c = c.NextSibling {
            traverse(c)
        }
    }

    traverse(form)
    return found
}

// normalizeURL converts relative URLs to absolute URLs
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