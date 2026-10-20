package flog

import (
	"bytes"
	"reflect"
	"testing"
)

func TestFormatMessageSeparator(t *testing.T) {
	tests := []struct {
		name, separator, msg, want string
		args, remaining            []any
	}{
		{"custom", "[]", "user=[] {}", "user=alice {}", []any{"alice", 42}, []any{42}},
		{"escapes", "[]", "[[]] [[[]]]", "[] [alice]", []any{"alice"}, nil},
		{"missing", "<>", "<> <>", "42 <>", []any{42}, nil},
		{"unicode", "「」", "「「「」「」」」", "「alice世界」", []any{"alice", "世界"}, nil},
		{"identical characters", "%%", "value=%%", "value=42", []any{42}, nil},
		{"fields", "[]", "[]", "alice", []any{Any("count", 3), "alice"}, []any{Any("count", 3)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, remaining := formatMessage(tt.msg, tt.args, tt.separator)
			if got != tt.want || len(remaining) != len(tt.remaining) || (len(remaining) > 0 && !reflect.DeepEqual(remaining, tt.remaining)) {
				t.Fatalf("got (%q, %v), want (%q, %v)", got, remaining, tt.want, tt.remaining)
			}
		})
	}
}

func TestFormatMessageInvalidSeparator(t *testing.T) {
	for _, separator := range []string{"", "{", "{}x", "\xff}"} {
		t.Run(separator, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("expected invalid separator to panic")
				}
			}()
			formatMessage("plain message", nil, separator)
		})
	}
}

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
		{"existing values", "user", []any{"alice", true}, "user alice true"},
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
