# flog

A customizable development logger for Go applications, with console output, optional timestamps, and named fields.

Install with Go 1.25.5 or later:

```sh
go get github.com/aidenfine/flog
```

Simple example
```go
logger.Info("new user","user" user.name, "is user active", true).Caller()
// [INFO] [app] - new user {"user": "alice", "is user active": true} - user_service.go:421
```

```go
import "github.com/aidenfine/flog"

port := 8080
logger, err := flog.New(false, "app", flog.Config{ShowTime: true, Separator: "{}"})
if err != nil {
    panic(err)
}
defer logger.Close()

logger.Info("server started: ", 
	"port", port
)
logger.Debug("handling request").Caller() // Print caller with log message
logger.Error("an api request failed!", 
	"http_error", 429,
	"reason", "rate_limit",
	""
)
```

## Configuration

Build a `flog.Config` and pass it as the third argument to `flog.New`:

```go
cfg := flog.Config{
    LogLevel:   "INFO",
    OutputPath: []string{"app.log"},
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
| `ShowTime` | Include a timestamp when `true`. Defaults to `false`. | `false`|
| `TimeEncoding` | Change encoding type on timestamps, using Go's `time` library. | `RFC3339` |
| `Encoding` | Use `"console"` or leave empty for console output. `"json"` is not supported yet and falls back to console with a warning. | `"console"` |

## Performance

Latest benchmark result:

| Benchmark | Iterations | Time per op | Bytes per op | Allocations per op |
|---|---:|---:|---:|---:|
| BenchmarkLogMessageWith10Fields-12 | 2,875,400 | 2,084 ns/op | 24 B/op | 1 allocs/op |
| BenchmarkLogMessageWith10StrictFields-12 | 2,980,706 | 2,094 ns/op | 24 B/op | 1 allocs/op |

The `-12` suffix indicates `GOMAXPROCS=12`.

Run tests and benchmarks with `make test` and `make bench`.
Use `make profile` to collect CPU and memory profiles in `profiles/`.
Run `make pprof-cpu` or `make pprof-memory` to collect fresh profiles and open
the corresponding profile viewer.

Override the benchmark duration with `make profile BENCHTIME=10s`.
Profiles default to the loose key-value benchmark. Add `STRICT=1` to profile
strict fields, for example `make pprof-memory STRICT=1` or `make pprof-cpu STRICT=1`.
`make bench` runs both benchmarks; `make bench STRICT=1` runs only strict fields.
Use `BENCH` to override the benchmark selection explicitly.


Benchmark stderr is saved to `profiles/benchmark.stderr.log` during profiling.

---
