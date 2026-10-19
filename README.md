# flog

A customizable development logger for Go applications, with console output, optional timestamps, and named fields.

Install with Go 1.25.5 or later:

```sh
go get github.com/aidenfine/flog
```

Run tests and benchmarks with `make test` and `make bench`.
Use `make profile` to collect CPU and memory profiles in `profiles/`, then
`make pprof-cpu` or `make pprof-memory` to open the profile viewer.
Override the benchmark duration with `make profile BENCHTIME=10s`.
Benchmark stderr is saved to `profiles/benchmark.stderr.log` during profiling.

Simple example
```go
import "github.com/aidenfine/flog"

port := 8080
logger, err := flog.New(false, "app", flog.Config{ShowTime: true})
if err != nil {
    panic(err)
}
defer logger.Close()
logger.Info("server started: ", port)
logger.Debug("handling request")
logger.Error("request failed")
```

Use `{}` placeholders in `Info`, `Debug`, and `Error` to insert values in order:

```go
logger.Info("user={}, is user active: {}", "alice", true)
// [INFO] [app] - user=alice, is user active: true
```

Values use Go's default `%v` formatting, so strings appear without quotes.
Write `{{` and `}}` for literal braces (for example, `{{}}` logs `{}`).
Unfilled placeholders remain `{}`; extra arguments are appended as before.
Explicit `flog.Any(...)` fields stay separate and do not consume placeholders.
Formatting is skipped for disabled log levels. Calls without placeholders or
escaped braces keep their existing behavior.
<!--
Logs go to stderr. An empty `LogLevel` enables all levels; `DEBUG`, `INFO`, or `ERROR` suppresses that level and lower levels. JSON encoding is not yet supported. Set `OutputPath: []string{"app.log"}` to also append logs to files. `New` returns an error if a file cannot be opened; call `Close` to release files when finished. Concurrent log calls are serialized. Logging after `Close` is ignored. Output destinations are written in order (stderr, then files); a write failure stops that entry from reaching later destinations.-->

The package is very flexible allowing for loosely structured logs as well.
```go
import "github.com/aidenfine/flog"

port := 8080
logger, err := flog.New(false, "app", flog.Config{ShowTime: true})
if err != nil {
    panic(err)
}
defer logger.Close()
logger.Info("server on port: ", port, "with the config of", flog.Config{
	OutputPath: nil,
	ShowTime:   true,
})
```
