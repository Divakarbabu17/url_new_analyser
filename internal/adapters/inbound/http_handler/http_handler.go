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
func (h *Handler) handleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse input JSON { "url": "https://example.com" }
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	// Call use case
	result, err := h.AnalyzeUseCase.Execute(req.URL)
	if err != nil {
		log.Println("Analyze error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// StartServer starts HTTP server with graceful shutdown
func StartServer(handler *Handler, addr string, stopChan <-chan struct{}) error {
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Run server in goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe failed: %v", err)
		}
	}()
	log.Println("Server started at", addr)

	// Wait for stop signal
	<-stopChan
	log.Println("Shutting down server...")

	// Graceful shutdown with 5-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}
