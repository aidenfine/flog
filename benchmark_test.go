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
		flog.Any("request_id", "req-123"),
		flog.Any("method", "GET"),
		flog.Any("path", "/api/users"),
		flog.Any("status", 200),
		flog.Any("duration_ms", 12.5),
		flog.Any("user_id", 42),
		flog.Any("remote_ip", "1.1.1.1.1"),
		flog.Any("bytes_sent", 1024),
		flog.Any("cached", false),
		flog.Any("service", "api"),
	}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		logger.Info("request completed", fields...)
	}
}
