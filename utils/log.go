package utils

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Global logger instance which is used as a base to create all other scoped loggers
var logger *zerolog.Logger

func GetLogger(scope string) zerolog.Logger {
	if logger == nil && os.Getenv("NO_FILE_LOGGING") == "" {
		// Setup lumberjack for log rotation
		logWriter := &lumberjack.Logger{
			Filename:   filepath.Join(Getenv("LOG_DIR", "."), Getenv("LOG_FILENAME", "app.log")),
			MaxSize:    int(Unwrap(strconv.ParseInt(Getenv("MAX_LOG_SIZE", "10"), 10, 32))),
			MaxBackups: int(Unwrap(strconv.ParseInt(Getenv("MAX_LOG_BACKUPS", "10"), 10, 32))),
			MaxAge:     int(Unwrap(strconv.ParseInt(Getenv("MAX_LOG_AGE", "180"), 10, 32))),
			Compress:   true,
		}

		consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
		multi := zerolog.MultiLevelWriter(logWriter, consoleWriter)
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

		level := Getenv("LOG_LEVEL", "info")

		switch strings.ToLower(level) {
		case "trace":
			zerolog.SetGlobalLevel(zerolog.TraceLevel)
		case "debug":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case "warn":
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		case "error":
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		case "panic":
			zerolog.SetGlobalLevel(zerolog.PanicLevel)
		default:
			// By default we log everything in info mode
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}

		newLogger := log.With().Caller().Logger().Output(multi)
		logger = &newLogger
	} else if logger == nil {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		newLogger := log.With().Caller().Logger()
		logger = &newLogger
	}

	return logger.With().
		Str("scope", scope).
		Logger()
}
