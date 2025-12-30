package logger

import (
	"io"
	"os"
	"strings"

	"agent-project-manager/internal/config"
	"github.com/sirupsen/logrus"
)

var (
	defaultLogger *logrus.Logger
)

// Config holds logger configuration
type Config struct {
	Level        string // debug, info, warn, error, fatal
	Format       string // text, json
	Output       string // stdout, stderr, or file path
	ReportCaller bool   // include caller information
}

func Init() {
	cfg := Config{
		Level:        "info",
		Format:       "text",
		Output:       "stdout",
		ReportCaller: false,
	}

	if appCfg, err := config.Load(); err == nil {
		if appCfg.Logger.Level != "" {
			cfg.Level = appCfg.Logger.Level
		}
		if appCfg.Logger.Format != "" {
			cfg.Format = appCfg.Logger.Format
		}
		if appCfg.Logger.Output != "" {
			cfg.Output = appCfg.Logger.Output
		}
		cfg.ReportCaller = appCfg.Logger.ReportCaller
	}

	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Set format
	switch strings.ToLower(cfg.Format) {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}

	// Set output
	var output io.Writer
	switch strings.ToLower(cfg.Output) {
	case "stderr":
		output = os.Stderr
	case "stdout", "":
		output = os.Stdout
	default:
		// File path
		file, err := os.OpenFile(cfg.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			// Fallback to stdout if file can't be opened
			output = os.Stdout
		} else {
			output = file
		}
	}
	logger.SetOutput(output)

	// Set report caller
	logger.SetReportCaller(cfg.ReportCaller)

	defaultLogger = logger
}

// Get returns the default logger instance
func Get() *logrus.Logger {
	if defaultLogger == nil {
		Init()
	}
	return defaultLogger
}

// WithField creates an entry with a single field
func WithField(key string, value interface{}) *logrus.Entry {
	return Get().WithField(key, value)
}

// WithFields creates an entry with multiple fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Get().WithFields(fields)
}

// WithError creates an entry with an error field
func WithError(err error) *logrus.Entry {
	return Get().WithError(err)
}

// Debug logs a message at level Debug
func Debug(args ...interface{}) {
	Get().Debug(args...)
}

// Debugf logs a formatted message at level Debug
func Debugf(format string, args ...interface{}) {
	Get().Debugf(format, args...)
}

// Info logs a message at level Info
func Info(args ...interface{}) {
	Get().Info(args...)
}

// Infof logs a formatted message at level Info
func Infof(format string, args ...interface{}) {
	Get().Infof(format, args...)
}

// Warn logs a message at level Warn
func Warn(args ...interface{}) {
	Get().Warn(args...)
}

// Warnf logs a formatted message at level Warn
func Warnf(format string, args ...interface{}) {
	Get().Warnf(format, args...)
}

// Error logs a message at level Error
func Error(args ...interface{}) {
	Get().Error(args...)
}

// Errorf logs a formatted message at level Error
func Errorf(format string, args ...interface{}) {
	Get().Errorf(format, args...)
}

// Fatal logs a message at level Fatal and exits
func Fatal(args ...interface{}) {
	Get().Fatal(args...)
}

// Fatalf logs a formatted message at level Fatal and exits
func Fatalf(format string, args ...interface{}) {
	Get().Fatalf(format, args...)
}

// Panic logs a message at level Panic
func Panic(args ...interface{}) {
	Get().Panic(args...)
}

// Panicf logs a formatted message at level Panic
func Panicf(format string, args ...interface{}) {
	Get().Panicf(format, args...)
}
