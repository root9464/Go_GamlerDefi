package logger

import (
	"fmt"
	"log/slog"
	"strconv"
)

const (
	reset = "\033[0m"

	black        = 30
	red          = 31
	green        = 32
	yellow       = 33
	blue         = 34
	magenta      = 35
	cyan         = 36
	lightGray    = 37
	darkGray     = 90
	lightRed     = 91
	lightGreen   = 92
	lightYellow  = 93
	lightBlue    = 94
	lightMagenta = 95
	lightCyan    = 96
	white        = 97
)

func colorize(colorCode int, v string) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), v, reset)
}

func getLevelColor(level slog.Level) int {
	switch {
	case level <= slog.LevelDebug:
		return lightGray
	case level <= slog.LevelInfo:
		return blue
	case level < slog.LevelWarn:
		return lightBlue
	case level < slog.LevelError:
		return lightYellow
	case level <= slog.LevelError+1:
		return lightRed
	default:
		return lightMagenta
	}
}
