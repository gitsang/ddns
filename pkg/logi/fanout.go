package logi

import (
	"context"
	"errors"
	"log/slog"
	"slices"
)

var _ slog.Handler = (*FanOutHandler)(nil)

type FanOutHandler struct {
	handlers []slog.Handler
}

func NewFanOutHandler(handlers ...slog.Handler) slog.Handler {
	return &FanOutHandler{
		handlers: handlers,
	}
}

func (h *FanOutHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for i := range h.handlers {
		if h.handlers[i].Enabled(ctx, l) {
			return true
		}
	}

	return false
}

func (h *FanOutHandler) Handle(ctx context.Context, r slog.Record) error {
	var result error
	for i := range h.handlers {
		if h.handlers[i].Enabled(ctx, r.Level) {
			err := try(func() error {
				return h.handlers[i].Handle(ctx, r.Clone())
			})
			result = errors.Join(result, err)
		}
	}

	return result
}

func (h *FanOutHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, h := range h.handlers {
		newHandlers[i] = h.WithAttrs(slices.Clone(attrs))
	}
	return NewFanOutHandler(newHandlers...)
}

func (h *FanOutHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, h := range h.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return NewFanOutHandler(newHandlers...)
}
