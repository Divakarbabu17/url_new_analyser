package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"url_new_analyser/internal/adapters/inbound/http_handler"
	"url_new_analyser/internal/adapters/outbound/fetcher"
	"url_new_analyser/internal/adapters/outbound/linkchecker"
	"url_new_analyser/internal/adapters/outbound/parser"
	"url_new_analyser/internal/core/usecase"
)

func main() {
	// -------------------------------
	// 1️⃣ Initialize outbound adapters
	// -------------------------------
	httpFetcher := fetcher.NewHTTPFetcher(10 * time.Second)              // HTTP fetcher with 10s timeout
	htmlParser := parser.NewHTMLParser()                                 // HTML parser
	linkChecker := linkchecker.NewWorkerPoolLinkChecker(10, 10*time.Second) // 10 concurrent workers, 10s timeout per link

	// -------------------------------
	// 2️⃣ Initialize core use case
	// -------------------------------
	analyzeUseCase := usecase.NewAnalyzePageUseCase(httpFetcher, htmlParser, linkChecker)

	// -------------------------------
	// 3️⃣ Initialize HTTP handler (inbound adapter)
	// -------------------------------
	analyzeHandler := http_handler.NewHandler(analyzeUseCase)

	// -------------------------------
	// 4️⃣ Create HTTP server with routes
	// -------------------------------
	mux := http.NewServeMux()
	analyzeHandler.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// -------------------------------
	// 5️⃣ Start server in a separate goroutine
	// -------------------------------
	go func() {
		log.Printf("Server started at %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// -------------------------------
	// 6️⃣ Wait for OS interrupt signal (Ctrl+C) for graceful shutdown
	// -------------------------------
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop // Block until signal received
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Gracefully shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %v", err)
	}

	// Stop worker pool
	linkChecker.Stop() // ensures all workers exit cleanly

	log.Println("Server exited gracefully")
}