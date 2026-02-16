package logger

import (
	"io"
	"log/slog"
	"os"
	"github.com/fatih/color"
)

var defaultHandler = NewCustomHandler(
	os.Stdout,
	&CustomOptions{
		Level: slog.LevelInfo,
		HidePC: true,
	},
)

var defaultLogger = slog.New(defaultHandler)
func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	color.Set(color.FgMagenta)
	defaultLogger.Warn(msg, args...)
	color.Unset()
}

func Error(msg string, args ...any) {
	color.Set(color.FgRed)
	defaultLogger.Error(msg, args...)
	color.Unset()
}

func SetLevel(level slog.Level) {
	defaultHandler.opts.Level = level
	defaultLogger = slog.New(defaultHandler)
}

func SetWriter(out io.Writer) {
	defaultHandler.out = out
	defaultLogger = slog.New(defaultHandler)
}

