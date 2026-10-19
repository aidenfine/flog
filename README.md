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
