package e2e_test

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/aidenfine/flog"
)

func TestLoggerOutput(t *testing.T) {
	for _, encoding := range []string{"console", "json"} {
		for _, level := range []string{"DEBUG", "INFO", "WARN", "ERROR", "PANIC"} {
			for _, fields := range []string{"none", "key-value", "typed"} {
				t.Run(encoding+"/"+level+"/"+fields, func(t *testing.T) {
					want := "[" + level + "] [e2e] - hello"
					if fields != "none" {
						want += ` {"user": "alice", "attempt": 3, "active": true}`
					}
					if encoding == "json" {
						want = `{"msg":"hello","level":"` + level + `","logger":"e2e"`
						if fields != "none" {
							want += `,"fields":{"user":"alice","attempt":3,"active":true}`
						}
						want += "}"
					}
					runLogger(t, encoding, level, fields, "", want+"\n")
				})
			}
		}
	}
}

func TestLoggerLevelFiltering(t *testing.T) {
	levels := []string{"DEBUG", "INFO", "WARN", "ERROR", "PANIC"}
	for _, encoding := range []string{"console", "json"} {
		for minimumIndex, minimum := range levels {
			for levelIndex, level := range levels {
				t.Run(encoding+"/"+minimum+"/"+level, func(t *testing.T) {
					want := ""
					if levelIndex >= minimumIndex {
						want = "[" + level + "] [e2e] - hello\n"
						if encoding == "json" {
							want = `{"msg":"hello","level":"` + level + `","logger":"e2e"}` + "\n"
						}
					}
					runLogger(t, encoding, level, "none", minimum, want)
				})
			}
		}
	}
}

// Compare decoded JSON so whitespace and object key order are not part of the contract.
func assertLogOutput(t *testing.T, encoding, got, want string) {
	t.Helper()
	if encoding != "json" || want == "" {
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
		return
	}
	if !bytes.HasSuffix([]byte(got), []byte{'\n'}) || bytes.Count([]byte(got), []byte{'\n'}) != 1 {
		t.Fatalf("expected one newline-terminated JSON record, got %q", got)
	}
	var actual, expected map[string]any
	if err := json.Unmarshal([]byte(got), &actual); err != nil {
		t.Fatalf("invalid JSON output %q: %v", got, err)
	}
	if err := json.Unmarshal([]byte(want), &expected); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %s, want %s", got, want)
	}
}

func runLogger(t *testing.T, encoding, level, fields, minimum, want string) {
	t.Helper()
	dir := t.TempDir()
	paths := []string{filepath.Join(dir, "first.log"), filepath.Join(dir, "second.log")}
	for _, path := range paths {
		if err := os.WriteFile(path, []byte("existing\n"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(executable, "-test.run=^TestLoggerProcess$", "--", encoding, level, fields, minimum, paths[0], paths[1])
	cmd.Env = append(os.Environ(), "FLOG_E2E_CHILD=1")
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("logger process: %v\nstdout: %s\nstderr: %s", err, &stdout, &stderr)
	}
	if stdout.Len() != 0 {
		t.Errorf("unexpected stdout: %q", stdout.String())
	}
	assertLogOutput(t, encoding, stderr.String(), want)
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.HasPrefix(data, []byte("existing\n")) {
			t.Fatalf("%s: existing contents lost: %q", path, data)
		}
		assertLogOutput(t, encoding, string(data[len("existing\n"):]), want)
	}
}

func TestLoggerProcess(t *testing.T) {
	if os.Getenv("FLOG_E2E_CHILD") != "1" {
		t.Skip("subprocess helper")
	}
	args := os.Args
	for len(args) > 0 && args[0] != "--" {
		args = args[1:]
	}
	if len(args) != 7 {
		t.Fatal("expected encoding, level, fields, minimum level, and two output paths")
	}
	encoding, level, style, minimum := args[1], args[2], args[3], args[4]
	logger, err := flog.New(false, "e2e", flog.Config{Encoding: encoding, LogLevel: minimum, OutputPath: args[5:]})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { logger.Close() })
	var fields []any
	switch style {
	case "key-value":
		fields = []any{"user", "alice", "attempt", 3, "active", true}
	case "typed":
		fields = []any{flog.String("user", "alice"), flog.Int("attempt", 3), flog.Bool("active", true)}
	}
	methods := map[string]func(string, ...any){
		"DEBUG": logger.Debug,
		"INFO":  logger.Info,
		"WARN":  logger.Warn,
		"ERROR": logger.Error,
		"PANIC": logger.Panic,
	}
	panicked := false
	func() {
		defer func() { panicked = recover() != nil }()
		methods[level]("hello", fields...)
	}()
	if panicked != (level == "PANIC") {
		t.Fatalf("%s: panicked = %v", level, panicked)
	}
	if err := logger.Close(); err != nil {
		t.Fatal(err)
	}
	if err := logger.Close(); err != nil {
		t.Fatalf("second Close: %v", err)
	}
	methods[level]("after close")
	if t.Failed() {
		os.Exit(1)
	}
	os.Exit(0)
}
