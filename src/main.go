package main

import (
	"capitalExporter/account"
	"capitalExporter/config"
	"capitalExporter/logger"
	"capitalExporter/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

func main() {
	// Initialize logger
	cfg, err := config.LoadConfiguration()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	sugar := logger.InitLogger(cfg.Debug)
	defer func(sugar *zap.SugaredLogger) {
		err := sugar.Sync()
		if err != nil {
			log.Fatalf("Error syncing logger: %v", err)
		}
	}(sugar) // Generate session

	sessionTokens, err := account.GenerateSession(cfg)
	if err != nil {
		sugar.Fatal("Error generating session:", err)
	}

	// Update Metrics initially
	err = metrics.UpdateMetrics(sessionTokens)
	if err != nil {
		sugar.Fatal("Error updating metrics:", err)
	}

	sugar.Infow("Metrics updated", "time", time.Now().Format("2006-01-02 15:04:05.999999999 -0700 MST"))

	// Start the HTTP server for Prometheus
	go startServer(cfg.Port)

	// Set up a ticker to update metrics
	ticker := time.NewTicker(time.Duration(cfg.Interval) * time.Second)
	for range ticker.C {
		err = metrics.UpdateMetrics(sessionTokens)
		if err != nil {
			sugar.Fatal("Error updating metrics:", err)
		}
		sugar.Infow("Metrics updated", "time", time.Now().Format("2006-01-02 15:04:05.999999999 -0700 MST"))
	}
}

func startServer(port string) {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
