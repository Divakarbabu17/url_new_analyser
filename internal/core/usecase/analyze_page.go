package usecase

import (
	"errors"
	"net/url"

	"url_new_analyser/internal/core/analyzer"
	"url_new_analyser/internal/core/domain"
	"url_new_analyser/internal/core/ports"
)

type AnalyzePageUseCase struct {
	Fetcher     ports.Fetcher
	Parser      ports.Parser
	LinkChecker ports.LinkChecker
}

func NewAnalyzePageUseCase(
	fetcher ports.Fetcher,
	parser ports.Parser,
	linkChecker ports.LinkChecker,
) *AnalyzePageUseCase {
	return &AnalyzePageUseCase{
		Fetcher:     fetcher,
		Parser:      parser,
		LinkChecker: linkChecker,
	}
}

func (uc *AnalyzePageUseCase) Execute(rawURL string) (*domain.AnalysisResult, error) {
	// 1. Validate URL
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, errors.New("invalid URL")
	}

	// 2. Fetch HTML
	htmlBody, statusCode, err := uc.Fetcher.Fetch(parsedURL.String())
	if err != nil {
		return nil, err
	}

	if statusCode >= 400 {
		return nil, errors.New("failed to fetch URL: bad status code")
	}

	// 3. Parse HTML
	doc, err := uc.Parser.Parse(htmlBody)
	if err != nil {
		return nil, err
	}

	// 4. Analyze HTML (PURE LOGIC — no interface)
	analysisData := analyzer.AnalyzeDocument(doc, parsedURL.String())

	// 5. Check links (async)
	linkResults := uc.LinkChecker.CheckLinks(analysisData.Links)

	// 6. Count internal / external / broken
	internal := 0
	external := 0
	broken := 0

	for _, res := range linkResults {
		if !res.OK {
			broken++
		}

		linkURL, err := url.Parse(res.URL)
		if err != nil {
			continue
		}

		if linkURL.Host == parsedURL.Host {
			internal++
		} else {
			external++
		}
	}

	// 7. Build result
	result := &domain.AnalysisResult{
		HTMLVersion: analysisData.HTMLVersion,
		Title:       analysisData.Title,
		Headings:    analysisData.Headings,
		Links: domain.LinkStats{
			Internal: internal,
			External: external,
			Broken:   broken,
		},
		LoginForm: analysisData.LoginForm,
	}

	return result, nil
}
