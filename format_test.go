package flog

import (
	"bytes"
	"testing"
)

func TestPlaceholderLogging(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		args []any
		want string
	}{
		{"values", "user={}, is he active: {}", []any{"alice", true}, "user=alice, is he active: true"},
		{"missing", "{} {}", []any{42}, "42 {}"},
		{"extra", "user={}", []any{"alice", 42}, "user=alice 42"},
		{"no arguments", "empty={}", nil, "empty={}"},
		{"escapes", "{{}} {{{}}} }} {", []any{"alice"}, "{} {alice} } {"},
		{"escaped without arguments", "{{}}", nil, "{}"},
		{"unicode", "你好 {}", []any{"世界"}, "你好 世界"},
		{"nil", "value={}", []any{nil}, "value=<nil>"},
		{"literal replacement", "{} {}", []any{"{}", "%s"}, "{} %s"},
		{"existing values", "user", []any{"alice", true}, "user \"alice\" true"},
		{"fields", "user={}", []any{Any("count", 3), "alice", Any("active", true)}, "user=alice 3 true"},
		{"only fields", "user={}", []any{Any("count", 3)}, "user={} 3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			logger := &Logger{name: "app", output: &out}
			for _, method := range []struct {
				level string
				log   func(string, ...any)
			}{{"INFO", logger.Info}, {"DEBUG", logger.Debug}, {"ERROR", logger.Error}} {
				out.Reset()
				method.log(tt.msg, tt.args...)
				want := "[" + method.level + "] [app] - " + tt.want + "\n"
				if got := out.String(); got != want {
					t.Errorf("%s: got %q, want %q", method.level, got, want)
				}
			}
		})
	}
}

type countedStringer struct{ calls *int }

func (v countedStringer) String() string {
	*v.calls++
	return "value"
}

func TestDisabledLogsSkipFormatting(t *testing.T) {
	var out bytes.Buffer
	calls := 0
	logger := &Logger{config: Config{LogLevel: "ERROR"}, output: &out}
	logger.Debug("{}", countedStringer{&calls})
	logger.closed = true
	logger.Error("{}", countedStringer{&calls})
	if calls != 0 || out.Len() != 0 {
		t.Fatalf("disabled logs formatted %d values and wrote %q", calls, out.String())
	}
}
