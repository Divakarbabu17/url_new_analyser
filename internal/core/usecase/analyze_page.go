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

