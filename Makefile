GO ?= go
BENCH ?= ^BenchmarkLogMessageWith10Fields$$
BENCHTIME ?= 5s
PROFILE_DIR ?= profiles

.PHONY: test bench profile pprof-cpu pprof-memory

test:
	$(GO) test ./...

bench:
	mkdir -p "$(PROFILE_DIR)"
	$(GO) test -c -o "$(PROFILE_DIR)/flog.test" .
	"$(abspath $(PROFILE_DIR))/flog.test" -test.run='^$$' -test.bench='$(BENCH)' -test.benchtime=$(BENCHTIME) -test.benchmem 2>"$(PROFILE_DIR)/benchmark.stderr.log"

profile:
	mkdir -p "$(PROFILE_DIR)"
	$(GO) test -c -o "$(PROFILE_DIR)/flog.test" .
	"$(abspath $(PROFILE_DIR))/flog.test" -test.run='^$$' -test.bench='$(BENCH)' -test.benchtime=$(BENCHTIME) -test.benchmem -test.cpuprofile="$(PROFILE_DIR)/cpu.prof" -test.memprofile="$(PROFILE_DIR)/mem.prof" 2>"$(PROFILE_DIR)/benchmark.stderr.log"

pprof-cpu: profile
	$(GO) tool pprof -http=localhost:8080 "$(PROFILE_DIR)/flog.test" "$(PROFILE_DIR)/cpu.prof"

pprof-memory: profile
	$(GO) tool pprof -http=localhost:8081 -alloc_space "$(PROFILE_DIR)/flog.test" "$(PROFILE_DIR)/mem.prof"
