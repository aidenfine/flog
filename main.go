package flog

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/aidenfine/flog/src"
)

type Logger struct {
	production bool
	// onPanic    bool
	// onFatal    bool
	name   string
	config Config
	clock  src.Clock
	output io.Writer
	files  []*os.File
	mu     sync.Mutex
	fields []src.Field
	closed bool
}

type Config struct {
	LogLevel   string // min level at which logs will show.
	OutputPath []string
	Separator  string
	// errorOutputPath []string
	// stacktrace bool
	ShowTime bool
	Encoding string // json or console. json not currently supported
}

// New creates a logger and opens its output files in append mode.
func New(production bool, name string, cfg Config) (*Logger, error) {

	// set default
	if cfg.Separator == "" {
		cfg.Separator = "{}"
	}

	// this does nothing as of now since Encoding is currently not used yet.
	if cfg.Encoding == "" {
		cfg.Encoding = "console"
	}

	l := &Logger{
		production: production,
		name:       name,
		config:     cfg,
		clock:      src.SystemClock,
	}
	outputs := []io.Writer{os.Stderr}
	for _, path := range cfg.OutputPath {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("open log file %q: %w", path, err), l.Close())
		}
		l.files = append(l.files, file)
		outputs = append(outputs, file)
	}
	l.output = io.MultiWriter(outputs...)
	return l, nil
}

func (l *Logger) Info(msg string, fields ...any) {
	l.log("INFO", msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...any) {
	l.log("DEBUG", msg, fields...)
}

func (l *Logger) Error(msg string, fields ...any) {
	l.log("ERROR", msg, fields...)
}

func (l *Logger) log(level, msg string, fields ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.closed {
		return
	}
	if e := l.check(msg, level, fields...); e != nil {
		if _, err := src.WriteLog(l.output, e); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write log: %v\n", err)
		}
		// clear fields
		clear(l.fields)
		l.fields = l.fields[:0]
	}
}

// Close releases the output files opened by New. It is safe to call repeatedly.
// Log calls after Close are ignored. Standard error is never closed.
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.closed {
		return nil
	}
	l.closed = true
	var errs []error
	for _, file := range l.files {
		if err := file.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (l *Logger) check(msg string, level string, fields ...any) *src.Entity {
	if src.LogLevelToInt[level] < src.LogLevelToInt[l.config.LogLevel] {
		return nil
	}
	separator := l.config.Separator
	if separator == "" {
		separator = "{}"
	}
	msg, fields = formatMessage(msg, fields, separator)

	if l.config.Encoding == "json" {
		fmt.Printf("[WARN] [flog] - 'json' encoding not yet supported defaulting to 'console'.")
	}

	var t *time.Time = nil
	if l.config.ShowTime {
		now := l.clock.Now()
		t = &now
	}

	l.fields = src.AppendFields(l.fields[:0], fields...)
	e := &src.Entity{
		Msg:        msg,
		LoggerName: l.name,
		Time:       t,
		Fields:     l.fields,
		LogLevel:   level,
	}

	// not sure how much this will log.
	// May need to do an offset capture??
	// if l.config.stacktrace {
	// 	stack := debug.Stack()
	// 	fmt.Println(string(stack))
	// }

	return e

}
