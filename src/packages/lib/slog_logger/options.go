package logger

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"reflect"
	"sync"
)

var defaultTimeFormat = "[15:04:05.000]"

type Logger struct {
	*slog.Logger
}

func NewLogger(opt *Option) *Logger {

	defaultOpt := &Option{
		Output:           os.Stdout,
		Colorize:         false,
		OutputEmptyAttrs: false,
		TimeFormat:       defaultTimeFormat,
	}

	resultOpt := mergeStruct(defaultOpt, opt).(*Option)
	fmt.Printf("Options: %+v \n", resultOpt)

	logger := slog.New(New(&slog.HandlerOptions{AddSource: true}, resultOpt))

	return &Logger{logger}
}

func New(handlerOptions *slog.HandlerOptions, opt *Option) *Handler {
	if handlerOptions == nil {
		handlerOptions = &slog.HandlerOptions{AddSource: true}
	}

	buf := &bytes.Buffer{}

	return &Handler{
		innerHandler: slog.NewJSONHandler(buf, &slog.HandlerOptions{
			Level:       handlerOptions.Level,
			AddSource:   handlerOptions.AddSource,
			ReplaceAttr: suppressDefaults(handlerOptions.ReplaceAttr),
		}),
		replaceAttr: handlerOptions.ReplaceAttr,
		buffer:      buf,
		mu:          &sync.Mutex{},
		option:      opt,
	}
}

type Option struct {
	Output           io.Writer
	Level            string
	Format           string
	Colorize         bool
	OutputEmptyAttrs bool
	TimeFormat       string
	CustomColors     CustomColors
	Source           Source
}

type Source struct {
	Add        bool
	ShowLine   bool
	ShowFunc   bool
	PathMode   string
	TrimPrefix string
}

type CustomColors struct {
	Debug int
	Warn  int
	Info  int
	Error int
}

func suppressDefaults(next func([]string, slog.Attr) slog.Attr) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey || a.Key == slog.LevelKey || a.Key == slog.MessageKey {
			return slog.Attr{}
		}
		if next == nil {
			return a
		}
		return next(groups, a)
	}
}

func mergeStruct(defaultOpt, customOpt any) any {
	if customOpt == nil {
		return defaultOpt
	}

	defVal := reflect.ValueOf(defaultOpt).Elem()
	custVal := reflect.ValueOf(customOpt).Elem()

	result := reflect.New(defVal.Type()).Elem()
	result.Set(defVal) // Копируем дефолтные значения

	for i := 0; i < defVal.NumField(); i++ {
		field := custVal.Field(i)
		if !field.IsZero() {
			result.Field(i).Set(field)
		}
	}

	return result.Addr().Interface()
}
