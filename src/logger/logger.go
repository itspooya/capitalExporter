package logger

import (
	"go.uber.org/zap"
	"log"
)

// InitLogger initializes the logger and returns a SugaredLogger for easy logging
func InitLogger(debug bool) *zap.SugaredLogger {
	var logger *zap.Logger
	var err error
	if debug {
		// Use the development config with debug level logging
		logger, err = zap.NewDevelopment()
	} else {
		// Use the production config with info level logging
		logger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	sugar := logger.Sugar()
	return sugar
}
