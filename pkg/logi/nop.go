package logi

import (
	"context"
	"log/slog"
)

type nopHandler struct {
}

func (h nopHandler) Enabled(context.Context, slog.Level) bool {
	return false
}

func (h nopHandler) Handle(context.Context, slog.Record) error {
	return nil
}

func (h nopHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

func (h nopHandler) WithGroup(string) slog.Handler {
	return h
}

var NopHandler = nopHandler{}
