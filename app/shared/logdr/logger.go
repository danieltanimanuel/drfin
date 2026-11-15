package logdr

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"runtime/debug"
)

func callerLine(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

func Info(msg string, attrs ...any) {
	all := append(attrs, "caller", callerLine(2))
	slog.Info(msg, all...)
}

func Warn(msg string, attrs ...any) {
	all := append(attrs, "caller", callerLine(2))
	slog.Warn(msg, all...)
}

func Debug(msg string, attrs ...any) {
	all := append(attrs, "caller", callerLine(2))
	slog.Debug(msg, all...)
}

func Error(msg string, err error, attrs ...any) {
	if err == nil {
		return
	}
	all := append(attrs,
		"error", err.Error(),
		"stack", string(debug.Stack()),
		"caller", callerLine(2),
	)
	slog.Error(msg, all...)
}