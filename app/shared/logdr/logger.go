package logdr

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
)

const traceDepth = 5

// frames returns up to depth frames as "file:line func".
// skip=0 is the function that called frames; skip=1 is its caller.
func frames(skip, depth int) []string {
	pcs := make([]uintptr, depth)
	n := runtime.Callers(skip+2, pcs) // 0=Callers, 1=frames
	if n == 0 {
		return nil
	}
	it := runtime.CallersFrames(pcs[:n])
	out := make([]string, 0, n)
	for {
		f, more := it.Next()
		fn := f.Function
		if i := strings.LastIndex(fn, "/"); i >= 0 {
			fn = fn[i+1:]
		}
		out = append(out, fmt.Sprintf("%s:%d %s", filepath.Base(f.File), f.Line, fn))
		if !more {
			break
		}
	}
	return out
}

func callerLine(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

func Info(msg string, attrs ...any) {
	all := append(attrs, "caller", callerLine(2), "trace", frames(1, traceDepth))
	slog.Info(msg, all...)
}

func Warn(msg string, attrs ...any) {
	all := append(attrs, "caller", callerLine(2), "trace", frames(1, traceDepth))
	slog.Warn(msg, all...)
}

func Debug(msg string, attrs ...any) {
	all := append(attrs, "caller", callerLine(2), "trace", frames(1, traceDepth))
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