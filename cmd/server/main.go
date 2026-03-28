package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"url_new_analyser/internal/adapters/inbound/http_handler"
	"url_new_analyser/internal/adapters/outbound/fetcher"
	"url_new_analyser/internal/adapters/outbound/linkchecker"
	"url_new_analyser/internal/adapters/outbound/parser"
	"url_new_analyser/internal/core/usecase"
)

func main() {
	// 1️⃣ Initialize outbound adapters
	httpFetcher := fetcher.NewHTTPFetcher()
	htmlParser := parser.NewHTMLParser()
	linkChecker := linkchecker.NewWorkerPoolLinkChecker(10) // 10 workers

	// 2️⃣ Initialize core use case
	analyzeUseCase := usecase.NewAnalyzePageUseCase(httpFetcher, htmlParser, linkChecker)

	// 3️⃣ Initialize HTTP handler (inbound adapter)
	analyzeHandler := http_handler.NewAnalyzeHandler(analyzeUseCase)

	// 4️⃣ Create HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(analyzeHandler.HandleAnalyze),
	}

	// 5️⃣ Start server in a goroutine
	go func() {
		log.Printf("Server started on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// 6️⃣ Wait for interrupt signal to shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop // Block until signal received
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 7️⃣ Shutdown server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %v", err)
	}

	// 8️⃣ Stop link checker workers gracefully
	linkChecker.Stop()

	log.Println("Server exited gracefully")
}