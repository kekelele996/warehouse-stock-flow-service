package util

import (
	"context"
	"log/slog"
	"os"
)

// NewLogger 创建结构化 JSON 日志器。
func NewLogger(level string) *slog.Logger {
	var lv slog.Level
	switch level {
	case "debug":
		lv = slog.LevelDebug
	case "warn":
		lv = slog.LevelWarn
	case "error":
		lv = slog.LevelError
	default:
		lv = slog.LevelInfo
	}
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lv})
	return slog.New(handler)
}

// Info 使用全局默认日志器输出 Info 级别日志。
func Info(ctx context.Context, template string, args ...any) {
	slog.Default().InfoContext(ctx, template, args...)
}

// Warn 使用全局默认日志器输出 Warn 级别日志。
func Warn(ctx context.Context, template string, args ...any) {
	slog.Default().WarnContext(ctx, template, args...)
}

// Error 使用全局默认日志器输出 Error 级别日志。
func Error(ctx context.Context, template string, args ...any) {
	slog.Default().ErrorContext(ctx, template, args...)
}
