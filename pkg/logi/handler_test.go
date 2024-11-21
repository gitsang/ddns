package logi

import (
	"encoding/json"
	"io"
	"log/slog"
	"testing"

	"gotest.tools/v3/assert"
)

type TestingWriter struct {
	t         *testing.T
	expectMap map[string]any
}

func (w *TestingWriter) Write(p []byte) (n int, err error) {
	var actualMap map[string]any
	err = json.Unmarshal(p, &actualMap)
	assert.NilError(w.t, err)
	actualJsonBytes, err := json.Marshal(actualMap)
	assert.NilError(w.t, err)
	expectJsonBytes, err := json.Marshal(w.expectMap)
	assert.NilError(w.t, err)
	assert.Equal(w.t, string(actualJsonBytes), string(expectJsonBytes))
	return len(p), nil
}

func NewTestingWriter(t *testing.T, expectMap map[string]any) *TestingWriter {
	return &TestingWriter{
		t:         t,
		expectMap: expectMap,
	}
}

func TestLogger(t *testing.T) {
	t.Run("remove attr", func(t *testing.T) {
		w := NewTestingWriter(t, map[string]any{"level": "high", "msg": "hello"})

		handler := NewHandler(HandlerOptions{
			Format:  "json",
			Level:   "info",
			Writers: []io.Writer{w},
			ReplaceAttrs: []ReplaceAttrFunc{
				func(groups []string, a slog.Attr) slog.Attr {
					if a.Key == slog.TimeKey {
						return slog.Attr{}
					}
					if a.Key == slog.LevelKey {
						return slog.Attr{Key: "level", Value: slog.StringValue("high")}
					}
					return a
				}},
		})
		logger := slog.New(handler)
		logger.Info("hello")
		logger.Debug("hello")
	})
	removeTimeKey := []ReplaceAttrFunc{
		func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}
	t.Run("with attrs", func(t *testing.T) {
		w := NewTestingWriter(t, map[string]any{"level": "INFO", "msg": "hello", "too": "young", "sometime": "naive"})
		handler := NewHandler(HandlerOptions{
			Format: "json",
			Level:  "info",
			Attrs: map[string]any{
				"too":      "young",
				"sometime": "naive",
			},
			Writers:      []io.Writer{w},
			ReplaceAttrs: removeTimeKey,
		})
		logger := slog.New(handler)
		logger.Info("hello")
		logger.Debug("hello")
	})
}
