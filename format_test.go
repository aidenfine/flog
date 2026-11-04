package flog

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
	"time"
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

func TestKeyValueLogging(t *testing.T) {
	var out bytes.Buffer
	logger := &Logger{name: "app", output: &out}
	logger.Info("failed to fetch URL", "url", "https://example.com", "attempt", 3, "backoff", time.Second)
	want := `[INFO] [app] - failed to fetch URL {"url": "https://example.com", "attempt": 3, "backoff": "1s"}` + "\n"
	if out.String() != want {
		t.Fatalf("got %q, want %q", out.String(), want)
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

func TestCallerOutsideFields(t *testing.T) {
	var out bytes.Buffer
	logger := &Logger{name: "app", output: &out}
	logger.Caller().Info("first", "attempt", 3)
	logger.Info("second", "attempt", 4)
	lines := strings.Split(strings.TrimSuffix(out.String(), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("unexpected output: %q", out.String())
	}
	for i, line := range lines {
		if strings.Contains(line, "!BADKEY") || !strings.Contains(line, "} - format_test.go:") {
			t.Errorf("line %d: %q", i, line)
		}
	}
}
