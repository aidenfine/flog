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
