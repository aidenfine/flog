package flog_test

import (
	"testing"

	"github.com/aidenfine/flog"
)

func BenchmarkLogMessageWith10Fields(b *testing.B) {
	logger, err := flog.New(true, "my-log", flog.Config{
		OutputPath: nil,
		ShowTime:   true,
	})
	if err != nil {
		b.Fatal(err)
	}
	b.Cleanup(func() {
		if err := logger.Close(); err != nil {
			b.Error(err)
		}
	})

	fields := []any{
		"request_id", "req-123",
		"method", "GET",
		"path", "/api/users",
		"status", 200,
		"duration_ms", 12.5,
		"user_id", 42,
		"remote_ip", "1.1.1.1.1",
		"bytes_sent", 1024,
		"cached", false,
		"service", "api",
	}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		logger.Info("request completed", fields...)
	}
}

func BenchmarkLogMessageWith10StrictFields(b *testing.B) {
	logger, err := flog.New(true, "my-log", flog.Config{
		OutputPath: nil,
		ShowTime:   true,
	})
	if err != nil {
		b.Fatal(err)
	}
	b.Cleanup(func() {
		if err := logger.Close(); err != nil {
			b.Error(err)
		}
	})

	fields := []any{
		flog.String("request_id", "req-123"),
		flog.String("method", "GET"),
		flog.String("path", "/api/users"),
		flog.Int("status", 200),
		flog.Int("duration_ms", 12),
		flog.Int("user_id", 42),
		flog.String("remote_ip", "1.1.1.1.1"),
		flog.Int("bytes_sent", 1024),
		flog.Bool("cached", false),
		flog.String("service", "api"),
	}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		logger.Info("request completed", fields...)
	}
}
