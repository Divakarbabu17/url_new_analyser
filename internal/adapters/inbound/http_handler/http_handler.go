package http_handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"url_new_analyser/internal/core/usecase"
)

// Handler holds dependencies for HTTP requests
type Handler struct {
	AnalyzeUseCase *usecase.AnalyzePageUseCase
}

// NewHandler creates a new HTTP handler
func NewHandler(analyzeUseCase *usecase.AnalyzePageUseCase) *Handler {
	return &Handler{
		AnalyzeUseCase: analyzeUseCase,
	}
}

// RegisterRoutes sets up HTTP endpoints
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/analyze", h.handleAnalyze)
}

// handleAnalyze handles /analyze POST requests

// StartServer starts HTTP server with graceful shutdown

