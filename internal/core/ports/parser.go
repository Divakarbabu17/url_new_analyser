package ports

import "golang.org/x/net/html"

type Parser interface {
    Parse(htmlContent string) (*html.Node, error)
}