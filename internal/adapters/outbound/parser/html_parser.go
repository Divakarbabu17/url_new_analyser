package parser

import (
	"strings"

	"url_new_analyser/internal/core/ports"

	"golang.org/x/net/html"
)

type HTMLParser struct{}

// NewHTMLParser creates a new parser
func NewHTMLParser() *HTMLParser {
	return &HTMLParser{}
}

// Parse converts HTML string to *html.Node
func (p *HTMLParser) Parse(htmlContent string) (*html.Node, error) {
	reader := strings.NewReader(htmlContent)
	doc, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// Ensure it implements Parser port
var _ ports.Parser = (*HTMLParser)(nil)
