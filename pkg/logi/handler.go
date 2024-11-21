package logi

import (
	"fmt"
	"io"
	"log/slog"
	"path"
	"runtime"
	"strings"
	"time"
)

func parseLevel(level string) *slog.LevelVar {
	var slogLevel = new(slog.LevelVar)
	switch strings.ToLower(level) {
	case "debug", "dbg":
		slogLevel.Set(slog.LevelDebug)
	case "info", "inf":
		slogLevel.Set(slog.LevelInfo)
	case "warning", "warn":
		slogLevel.Set(slog.LevelWarn)
	case "error", "err":
		slogLevel.Set(slog.LevelError)
	default:
		slogLevel.Set(slog.LevelInfo)
	}
	return slogLevel
}

type ReplaceAttrFunc func(groups []string, a slog.Attr) slog.Attr

type VerbosityHandlers func(verbosity int)

type HandlerOptions struct {
	Format            string // console, json
	Level             string
	Attrs             map[string]any
	Writers           []io.Writer
	Verbosity         int
	VerbosityHandlers []VerbosityHandlers
	CallerSkip        int
	ReplaceAttrs      []ReplaceAttrFunc
}

func NewHandler(opt HandlerOptions) slog.Handler {
	// level
	slogLevel := parseLevel(opt.Level)

	// replace
	replaceChain := NewReplaceAttrChain()
	if len(opt.ReplaceAttrs) > 0 {
		for _, replaceAttr := range opt.ReplaceAttrs {
			replaceChain.Append(replaceAttr)
		}
	}

	// time
	replaceChain.Append(func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.String(a.Key, a.Value.Time().Format(time.RFC3339))
		}
		return a
	})

	// verbosity
	var addSource bool
	switch opt.Verbosity {
	case 0:
		addSource = false
	case 1:
		addSource = true
		replaceChain.Append(func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				pc, file, line, _ := runtime.Caller(opt.CallerSkip)
				fileName := path.Base(file)
				funcName := path.Base(runtime.FuncForPC(pc).Name())
				funcNames := strings.Split(funcName, ".")
				funcName = funcNames[len(funcNames)-1]
				return slog.String(slog.SourceKey, fmt.Sprintf("%s:%s:%d", fileName, funcName, line))
			}
			return a
		})
	case 2:
		addSource = true
		replaceChain.Append(func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				pc, file, line, _ := runtime.Caller(opt.CallerSkip)
				fileName := path.Base(file)
				funcName := path.Base(runtime.FuncForPC(pc).Name())
				return slog.String(slog.SourceKey, fmt.Sprintf("%s:%s:%d", fileName, funcName, line))
			}
			return a
		})
	case 3:
		addSource = true
		replaceChain.Append(func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				pc, file, line, _ := runtime.Caller(opt.CallerSkip)
				fileName := path.Base(file)
				funcName := path.Base(runtime.FuncForPC(pc).Name())
				return slog.Any(slog.SourceKey, slog.Source{
					Function: funcName,
					File:     fileName,
					Line:     line,
				})
			}
			return a
		})
	case 4:
		addSource = true
	}

	for _, handle := range opt.VerbosityHandlers {
		handle(opt.Verbosity)
	}

	// options
	handlerOpts := &slog.HandlerOptions{
		AddSource:   addSource,
		Level:       slogLevel,
		ReplaceAttr: replaceChain.ReplaceAttr,
	}

	// writer
	groupWriter := io.MultiWriter(opt.Writers...)

	// format
	var handler slog.Handler
	switch strings.ToLower(opt.Format) {
	case "console":
		handler = slog.NewTextHandler(groupWriter, handlerOpts)
	case "json":
		handler = slog.NewJSONHandler(groupWriter, handlerOpts)
	default:
		handler = slog.NewTextHandler(groupWriter, handlerOpts)
	}

	// attrs
	slogAttrs := make([]slog.Attr, 0, len(opt.Attrs))
	for key, value := range opt.Attrs {
		slogAttrs = append(slogAttrs, slog.Any(key, value))
	}
	handler = handler.WithAttrs(slogAttrs)

	return handler
}
