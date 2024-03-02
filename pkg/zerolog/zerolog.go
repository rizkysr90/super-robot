package logger

import (
	"api-iad-ams/internal/config"
	"os"

	"github.com/rs/zerolog"
)

func New(cfg config.Config) zerolog.Logger {
	// level, err := zerolog.ParseLevel(strings.ToLower(strings.TrimSpace(cfg.LogLevel)))
	// if err != nil {
	// 	level = zerolog.InfoLevel
	// }
	// if level == zerolog.Disabled {
	// 	return zerolog.Nop()
	// }
	level := zerolog.DebugLevel
	// output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}
	output := os.Stdout

	zerolog.DefaultContextLogger = nil
	zerolog.SetGlobalLevel(level)
	return zerolog.New(output).With().Timestamp().Logger()
}
