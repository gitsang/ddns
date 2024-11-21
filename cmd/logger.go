package main

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/docker/go-units"
	"github.com/gitsang/ddns/pkg/logi"
	timex "github.com/gitsang/ddns/pkg/util/time"
	"github.com/natefinch/lumberjack"
)

func NewLogHandler(defaultConf LogConfig, fancoutConfs ...LogConfig) []slog.Handler {
	configs := []LogConfig{defaultConf}
	configs = append(configs, fancoutConfs...)

	handlers := make([]slog.Handler, 0)
	for _, config := range configs {
		writers := make([]io.Writer, 0)
		if config.Output.Stdout.Enable {
			writers = append(writers, os.Stdout)
		}
		if config.Output.Stderr.Enable {
			writers = append(writers, os.Stderr)
		}
		if config.Output.File.Enable {
			maxSize, err := units.FromHumanSize(config.Output.File.MaxSize)
			if err != nil {
				panic(err)
			}
			maxAgeDur, err := timex.ParseDuration(config.Output.File.MaxAge)
			if err != nil {
				panic(err)
			}
			writers = append(writers, &lumberjack.Logger{
				Filename:   config.Output.File.Path,
				MaxSize:    int(maxSize / units.MB),
				MaxAge:     int(maxAgeDur / (24 * time.Hour)),
				MaxBackups: config.Output.File.MaxBackups,
				LocalTime:  false,
				Compress:   config.Output.File.Compress,
			})
		}

		handler := logi.NewHandler(
			logi.HandlerOptions{
				Format:     config.Format,
				Level:      config.Level,
				Writers:    writers,
				Verbosity:  config.Verbosity,
				CallerSkip: 7,
			},
		)

		handlers = append(handlers, handler)
	}

	return handlers
}
