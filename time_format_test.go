package flog

import (
	"bytes"
	"testing"
	"time"

	"github.com/aidenfine/flog/src"
)

type fixedTimeClock struct {
	src.Clock
	now time.Time
}

func (c fixedTimeClock) Now() time.Time { return c.now }

func TestLoggerTimeFormat(t *testing.T) {
	for _, tt := range []struct {
		name   string
		show   bool
		layout string
		prefix string
	}{
		{"default", true, "", "2026-11-02T15:04:05Z "},
		{"nanoseconds", true, time.RFC3339Nano, "2026-11-02T15:04:05.123456789Z "},
		{"custom", true, "15:04:05.000", "15:04:05.123 "},
		{"hidden", false, time.RFC3339Nano, ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(false, "app", Config{ShowTime: tt.show})
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if err := logger.Close(); err != nil {
					t.Error(err)
				}
			})
			logger.clock = fixedTimeClock{now: time.Date(2026, time.November, 2, 15, 4, 5, 123456789, time.UTC)}
			var out bytes.Buffer
			logger.output = &out
			logger.Info("hello")
			if want := tt.prefix + "[INFO] [app] - hello\n"; out.String() != want {
				t.Errorf("output = %q, want %q", out.String(), want)
			}
		})
	}
}
