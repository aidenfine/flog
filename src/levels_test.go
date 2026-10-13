package src

import "testing"

// this test are not really "needed" the logic is pretty straightforward but oh well free test cov %

func TestCapsLevelToLower(t *testing.T) {
	tests := []struct {
		name  string
		level string
		want  string
	}{
		{name: "error", level: "ERROR", want: "error"},
		{name: "debug", level: "DEBUG", want: "debug"},
		{name: "info", level: "INFO", want: "info"},
		{name: "empty"},
		{name: "unknown", level: "WARN"},
		{name: "lowercase", level: "error"},
		{name: "mixed case", level: "Info"},
		{name: "whitespace", level: " DEBUG "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := capsLevelToLower(tt.level); got != tt.want {
				t.Errorf("capsLevelToLower(%q) = %q, want %q", tt.level, got, tt.want)
			}
		})
	}
}

func TestLowerLevelToCaps(t *testing.T) {
	tests := []struct {
		name  string
		level string
		want  string
	}{
		{name: "error", level: "error", want: "ERROR"},
		{name: "debug", level: "debug", want: "DEBUG"},
		{name: "info", level: "info", want: "INFO"},
		{name: "empty"},
		{name: "unknown", level: "warn"},
		{name: "uppercase", level: "ERROR"},
		{name: "mixed case", level: "Info"},
		{name: "whitespace", level: " debug "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lowerLevelToCaps(tt.level); got != tt.want {
				t.Errorf("lowerLevelToCaps(%q) = %q, want %q", tt.level, got, tt.want)
			}
		})
	}
}

func TestLogLevelToInt(t *testing.T) {
	tests := []struct {
		level string
		want  int
		ok    bool
	}{
		{level: "ERROR", want: 3, ok: true},
		{level: "INFO", want: 2, ok: true},
		{level: "DEBUG", want: 1, ok: true},
		{level: "WARN"},
		{level: "error"},
		{level: ""},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			got, ok := LogLevelToInt[tt.level]
			if got != tt.want || ok != tt.ok {
				t.Errorf("LogLevelToInt[%q] = (%d, %t), want (%d, %t)", tt.level, got, ok, tt.want, tt.ok)
			}
		})
	}
}
