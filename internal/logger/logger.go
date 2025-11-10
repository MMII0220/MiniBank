package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger

// init автоматически инициализирует логгер при импорте пакета
func init() {
	output := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}

	logger = zerolog.New(output).With().
		Timestamp().
		Str("service", "minibank").
		Logger()

	logger.Info().Msg("Logger initialized")
}

// GetLogger возвращает глобальный логгер (публичная функция)
func GetLogger() zerolog.Logger {
	return logger
}
