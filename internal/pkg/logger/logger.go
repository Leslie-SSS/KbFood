package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// Logger is the global logger instance
	Logger zerolog.Logger
)

// Config holds logger configuration
type Config struct {
Level      string `envconfig:"LEVEL" default:"info"`
	Format     string `envconfig:"FORMAT" default:"json"`   // json or console
	Output     string `envconfig:"OUTPUT" default:"stdout"` // stdout or file path
	TimeFormat string `envconfig:"TIME_FORMAT" default:""`
}

// Init initializes the global logger
func Init(cfg *Config) error {
	// Set log level
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Set time format
	if cfg.TimeFormat != "" {
		zerolog.TimeFieldFormat = cfg.TimeFormat
	} else {
		zerolog.TimeFieldFormat = time.RFC3339
	}

	// Set output
	var output io.Writer = os.Stdout
	if cfg.Output == "console" || cfg.Format == "console" {
		// Console output with colors
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
			NoColor:    false,
		}
	}

	Logger = zerolog.New(output).With().Timestamp().Logger()
	log.Logger = Logger

	return nil
}

// Debug logs a debug message
func Debug(msg string) {
	log.Debug().Msg(msg)
}

// Debugf logs a debug message with formatting
func Debugf(format string, args ...interface{}) {
	log.Debug().Msgf(format, args...)
}

// Info logs an info message
func Info(msg string) {
	log.Info().Msg(msg)
}

// Infof logs an info message with formatting
func Infof(format string, args ...interface{}) {
	log.Info().Msgf(format, args...)
}

// Warn logs a warning message
func Warn(msg string) {
	log.Warn().Msg(msg)
}

// Warnf logs a warning message with formatting
func Warnf(format string, args ...interface{}) {
	log.Warn().Msgf(format, args...)
}

// Error logs an error message
func Error(msg string) {
	log.Error().Msg(msg)
}

// Errorf logs an error message with formatting
func Errorf(format string, args ...interface{}) {
	log.Error().Msgf(format, args...)
}

// Err logs an error with context
func Err(err error, msg string) {
	log.Error().Err(err).Msg(msg)
}

// Errf logs an error with formatting
func Errf(err error, format string, args ...interface{}) {
	log.Error().Err(err).Msgf(format, args...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string) {
	log.Fatal().Msg(msg)
}

// Fatalf logs a fatal message with formatting and exits
func Fatalf(format string, args ...interface{}) {
	log.Fatal().Msgf(format, args...)
}

// With creates a new logger with additional context
func With() zerolog.Context {
	return log.With()
}
