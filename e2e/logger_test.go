package e2e_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/aidenfine/flog"
)

func TestLoggerOutput(t *testing.T) {
	for _, level := range []string{"DEBUG", "INFO", "WARN", "ERROR", "PANIC"} {
		for _, fields := range []string{"none", "key-value", "typed"} {
			t.Run(level+"/"+fields, func(t *testing.T) {
				want := "[" + level + "] [e2e] - hello"
				if fields != "none" {
					want += ` {"user": "alice", "attempt": 3, "active": true}`
				}
				runLogger(t, level, fields, "", want+"\n")
			})
		}
	}
}

func TestLoggerLevelFiltering(t *testing.T) {
	levels := []string{"DEBUG", "INFO", "WARN", "ERROR", "PANIC"}
	for minimumIndex, minimum := range levels {
		for levelIndex, level := range levels {
			t.Run(minimum+"/"+level, func(t *testing.T) {
				want := ""
				if levelIndex >= minimumIndex {
					want = "[" + level + "] [e2e] - hello\n"
				}
				runLogger(t, level, "none", minimum, want)
			})
		}
	}
}

func runLogger(t *testing.T, level, fields, minimum, want string) {
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
	cmd := exec.Command(executable, "-test.run=^TestLoggerProcess$", "--", level, fields, minimum, paths[0], paths[1])
	cmd.Env = append(os.Environ(), "FLOG_E2E_CHILD=1")
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("logger process: %v\nstdout: %s\nstderr: %s", err, &stdout, &stderr)
	}
	if stdout.Len() != 0 {
		t.Errorf("unexpected stdout: %q", stdout.String())
	}
	if got := stderr.String(); got != want {
		t.Errorf("stderr:\n got %q\nwant %q", got, want)
	}
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if got := string(data); got != "existing\n"+want {
			t.Errorf("%s:\n got %q\nwant %q", path, got, "existing\n"+want)
		}
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
	if len(args) != 6 {
		t.Fatal("expected level, fields, minimum level, and two output paths")
	}
	level, style, minimum := args[1], args[2], args[3]
	logger, err := flog.New(false, "e2e", flog.Config{LogLevel: minimum, OutputPath: args[4:]})
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
