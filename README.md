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
logger, err := flog.New(false, "app", flog.Config{ShowTime: true, Separator: "{}"})
if err != nil {
    panic(err)
}
defer logger.Close()
logger.Info("server started: ", port)
logger.Debug("handling request")
logger.Error("request failed")
```

## Configuration

Build a `flog.Config` and pass it as the third argument to `flog.New`:

```go
cfg := flog.Config{
    LogLevel:   "INFO",
    OutputPath: []string{"app.log"},
    Separator:  "{}",
    ShowTime:   true,
    Encoding:   "console",
}

logger, err := flog.New(false, "app", cfg)
if err != nil {
    panic(err)
}
defer logger.Close()
```

| Field | Purpose | Default |
| --- | --- | --- |
| `LogLevel` | Minimum level to log: `"DEBUG"`, `"INFO"`, or `"ERROR"`. For example, `"INFO"` includes info and error messages. An empty string enables all levels. | `"INFO"` |
| `OutputPath` | Files to create or append logs to. Logs also go to stderr. Leave empty for stderr only. | `None` |
| `Separator` | Placeholder used to insert values into messages, such as `"{}"` or `"[]"`. Defaults to `"{}"` when empty. Custom values must contain exactly two characters; invalid values panic when formatting a log. | `"{}"` |
| `ShowTime` | Include a timestamp when `true`. Defaults to `false`. | `false`|
| `Encoding` | Use `"console"` or leave empty for console output. `"json"` is not supported yet and falls back to console with a warning. | `"console"` |

`New` returns an error if an output file cannot be opened; call `Close` when finished to release the files.

So if `Separator` is set to `{}` use `{}` placeholders in `Info`, `Debug`, and `Error` to insert values in order:

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

The package is very flexible allowing for loosely structured logs as well.
```go
import "github.com/aidenfine/flog"

port := 8080
logger, err := flog.New(false, "app", flog.Config{ShowTime: true, Separator: "{}"})
if err != nil {
    panic(err)
}
defer logger.Close()
logger.Info("server on port: ", port, "with the config of", server.config)
```
