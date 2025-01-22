package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitializeLogger configures the global logger settings
func InitializeLogger() {
	// Customize the logger format (human-readable for development, JSON for production)
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

// Info logs an info-level message
func Info(msg string, fields map[string]interface{}) {
	event := log.Info()
	for key, value := range fields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

// Error logs an error-level message
func Error(msg string, fields map[string]interface{}) {
	event := log.Error()
	for key, value := range fields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}
