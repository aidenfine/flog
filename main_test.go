package flog_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/aidenfine/flog"
)

func TestFileOutput(t *testing.T) {
	dir := t.TempDir()
	paths := []string{filepath.Join(dir, "first.log"), filepath.Join(dir, "second.log")}
	for _, path := range paths {
		if err := os.WriteFile(path, []byte("existing\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	logger, err := flog.New(false, "app", flog.Config{OutputPath: paths})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { logger.Close() })
	logger.Info("one")
	logger.Debug("two")
	logger.Error("three")
	if err := logger.Close(); err != nil {
		t.Fatal(err)
	}
	if err := logger.Close(); err != nil {
		t.Fatal(err)
	}
	logger.Info("after close")
	want := "existing\n[INFO] [app] - one\n[DEBUG] [app] - two\n[ERROR] [app] - three\n"
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != want {
			t.Errorf("%s: got %q, want %q", path, data, want)
		}
	}
}

func TestFileOpenError(t *testing.T) {
	logger, err := flog.New(false, "app", flog.Config{OutputPath: []string{
		filepath.Join(t.TempDir(), "valid.log"),
		filepath.Join(t.TempDir(), "missing", "invalid.log"),
	}})
	if err == nil || logger != nil {
		t.Fatalf("got logger %v, error %v", logger, err)
	}
}

func TestConcurrentFileOutput(t *testing.T) {
	path := filepath.Join(t.TempDir(), "concurrent.log")
	logger, err := flog.New(false, "app", flog.Config{OutputPath: []string{path}, LogLevel: "INFO"})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { logger.Close() })
	logger.Debug("filtered")
	var wg sync.WaitGroup
	for i := range 20 {
		wg.Add(1)
		go func(i int) { defer wg.Done(); logger.Info(fmt.Sprintf("entry %d", i)) }(i)
	}
	wg.Wait()
	if err := logger.Close(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	seen := make(map[string]bool)
	for _, line := range lines {
		seen[line] = true
	}
	if len(lines) != 20 {
		t.Fatalf("got %d entries, want 20", len(lines))
	}
	for i := 0; i < 20; i++ {
		if !seen[fmt.Sprintf("[INFO] [app] - entry %d", i)] {
			t.Errorf("missing entry %d", i)
		}
	}
}

func Test(t *testing.T) {
	logger, err := flog.New(true, "my-log", flog.Config{
		OutputPath: []string{"mylog.log", "mylog2.log", "log3.txt"},
		ShowTime:   true,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := logger.Close(); err != nil {
			t.Error(err)
		}
	})

	port := 8080
	logger.Info("server on port: ", port, "with the config of", flog.Config{
		OutputPath: nil,
		ShowTime:   true,
	})
}
