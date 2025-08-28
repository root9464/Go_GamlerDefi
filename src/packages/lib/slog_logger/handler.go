package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"sync"
)

const (
	timeFormat   = "[15:04:05.000]"
	messageDelim = " "
)

type Handler struct {
	innerHandler slog.Handler
	replaceAttr  func([]string, slog.Attr) slog.Attr
	buffer       *bytes.Buffer
	mu           *sync.Mutex
	option       *Option
}

type HandlerOption func(*Handler)

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.innerHandler.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{
		innerHandler: h.innerHandler.WithAttrs(attrs),
		replaceAttr:  h.replaceAttr,
		buffer:       h.buffer,
		mu:           h.mu,
		option:       h.option,
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		innerHandler: h.innerHandler.WithGroup(name),
		replaceAttr:  h.replaceAttr,
		buffer:       h.buffer,
		mu:           h.mu,
		option:       h.option,
	}
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	var builder strings.Builder

	if timeStr := h.formatTime(r); timeStr != "" {
		builder.WriteString(timeStr)
		builder.WriteString(messageDelim)
	}

	if levelStr := h.formatLevel(r); levelStr != "" {
		builder.WriteString(levelStr)
		builder.WriteString(messageDelim)
	}

	if msgStr := h.formatMessage(r); msgStr != "" {
		builder.WriteString(msgStr)
		builder.WriteString(messageDelim)
	}

	if attrsStr, err := h.formatAttrs(ctx, r); err != nil {
		return err
	} else if attrsStr != "" {
		builder.WriteString(attrsStr)
	}

	_, err := fmt.Fprintln(h.option.Output, builder.String())
	return err
}

func (h *Handler) formatAttrs(ctx context.Context, r slog.Record) (string, error) {
	attrs, err := h.computeAttrs(ctx, r)
	if err != nil {
		return "", err
	}

	if !h.option.OutputEmptyAttrs && len(attrs) == 0 {
		return "", nil
	}

	attrsBytes, err := json.MarshalIndent(attrs, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal attrs: %w", err)
	}

	if h.option.Colorize {
		return colorize(cyan, string(attrsBytes)), nil
	}
	return string(attrsBytes), nil
}

func (h *Handler) computeAttrs(ctx context.Context, r slog.Record) (map[string]any, error) {
	h.mu.Lock()
	defer func() {
		h.buffer.Reset()
		h.mu.Unlock()
	}()

	if err := h.innerHandler.Handle(ctx, r); err != nil {
		return nil, fmt.Errorf("inner handler: %w", err)
	}

	var attrs map[string]any
	if err := json.Unmarshal(h.buffer.Bytes(), &attrs); err != nil {
		return nil, fmt.Errorf("unmarshal attrs: %w", err)
	}

	if h.option.Source.Add {
		if srcRaw, ok := attrs[slog.SourceKey]; ok {
			if srcMap, ok := srcRaw.(map[string]any); ok {
				file, _ := srcMap["file"].(string)
				line, _ := srcMap["line"].(float64)
				function, _ := srcMap["function"].(string)

				switch h.option.Source.PathMode {
				case "filename":
					file = filepath.Base(file)
				case "relative":
					if prefix := h.option.Source.TrimPrefix; prefix != "" {
						if idx := strings.Index(file, prefix); idx != -1 {
							file = file[idx:]
						}
					}
				}

				newSource := map[string]any{}
				newSource["file"] = file
				if h.option.Source.ShowLine {
					newSource["line"] = line
				}

				if h.option.Source.ShowFunc {
					newSource["function"] = function
				}

				attrs[slog.SourceKey] = newSource
			}
		}
	}

	return attrs, nil
}

func (h *Handler) formatLevel(r slog.Record) string {
	color := getLevelColor(r.Level)
	if h.option.CustomColors.Debug != 0 && r.Level == slog.LevelDebug {
		color = h.option.CustomColors.Debug
	}
	if h.option.CustomColors.Info != 0 && r.Level == slog.LevelInfo {
		color = h.option.CustomColors.Info
	}
	if h.option.CustomColors.Warn != 0 && r.Level == slog.LevelWarn {
		color = h.option.CustomColors.Warn
	}
	if h.option.CustomColors.Error != 0 && r.Level == slog.LevelError {
		color = h.option.CustomColors.Error
	}

	return h.formatAttr(slog.Attr{
		Key:   slog.LevelKey,
		Value: slog.AnyValue(r.Level),
	}, color)
}

func (h *Handler) formatTime(r slog.Record) string {
	return h.formatAttr(slog.Attr{
		Key:   slog.TimeKey,
		Value: slog.StringValue(r.Time.Format(h.option.TimeFormat)),
	}, lightGray)
}

func (h *Handler) formatMessage(r slog.Record) string {
	return h.formatAttr(slog.Attr{
		Key:   slog.MessageKey,
		Value: slog.StringValue(r.Message),
	}, white)
}

func (h *Handler) formatAttr(attr slog.Attr, color int) string {
	if h.replaceAttr != nil {
		attr = h.replaceAttr(nil, attr)
	}

	if attr.Equal(slog.Attr{}) {
		return ""
	}

	str := attr.Value.String()
	if attr.Key == slog.LevelKey {
		str += ":"
	}

	if h.option.Colorize {
		return colorize(color, str)
	}
	return str
}
