package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

var (
	instance *Logger
	once     sync.Once
)

func GetLogger() *Logger {
	once.Do(func() {
		log := logrus.New()

		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			ForceColors:     true,
			TimestampFormat: "2006-01-02 15:04:05",
		})

		log.SetLevel(logrus.DebugLevel)
		log.SetOutput(os.Stdout)

		instance = &Logger{log}
	})
	return instance
}

func (l *Logger) logWithCaller(level logrus.Level, msg string, customLevel string) {
	colors := map[string]string{
		"error":   "\033[31m", // Красный для ошибок
		"warn":    "\033[33m", // Желтый для предупреждений
		"info":    "\033[34m", // Синий для информационных сообщений
		"success": "\033[32m", // Зеленый для успешных сообщений
	}

	var colorKey string
	switch level {
	case logrus.ErrorLevel:
		colorKey = "error"
	case logrus.WarnLevel:
		colorKey = "warn"
	case logrus.InfoLevel:
		if customLevel == "success" {
			colorKey = "success"
		} else {
			colorKey = "info"
		}
	default:
		colorKey = "info"
	}

	msg = fmt.Sprintf("%s%s\033[0m", colors[colorKey], msg)
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		l.Logger.Log(level, msg)
		return
	}
	fileLine := fmt.Sprintf("%s:%d", path.Base(file), line)
	coloredFileLine := fmt.Sprintf("\033[30m%s\033[0m", fileLine)
	l.Logger.Log(level, fmt.Sprintf("%s %s", coloredFileLine, msg))
}

func (l *Logger) Error(msg string) { l.logWithCaller(logrus.ErrorLevel, msg, "error") }
func (l *Logger) Info(msg string)  { l.logWithCaller(logrus.InfoLevel, msg, "info") }
func (l *Logger) Warn(msg string)  { l.logWithCaller(logrus.WarnLevel, msg, "warn") }

func (l *Logger) Errorf(format string, args ...any) {
	l.logWithCaller(logrus.ErrorLevel, fmt.Sprintf(format, args...), "error")
}
func (l *Logger) Infof(format string, args ...any) {
	l.logWithCaller(logrus.InfoLevel, fmt.Sprintf(format, args...), "info")
}
func (l *Logger) Warnf(format string, args ...any) {
	l.logWithCaller(logrus.WarnLevel, fmt.Sprintf(format, args...), "warn")
}

func (l *Logger) Successf(format string, args ...any) {
	l.logWithCaller(logrus.InfoLevel, fmt.Sprintf(format, args...), "success")
}

func (l *Logger) Success(msg string) { l.logWithCaller(logrus.InfoLevel, msg, "success") }
