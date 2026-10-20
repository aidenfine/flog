GO ?= go
BENCH ?= ^BenchmarkLogMessageWith10Fields$$
BENCHTIME ?= 5s
PROFILE_DIR ?= profiles
BASE ?=
COMPARE_ADDR ?= localhost:8082

.PHONY: test bench profile pprof-cpu pprof-memory pprof-compare

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

# Compare existing memory profiles without regenerating or overwriting them.
pprof-compare:
	@test -n "$(BASE)" || { echo 'Usage: make pprof-compare BASE=path/to/baseline/mem.prof [PROFILE_DIR=profiles]'; exit 1; }
	@test -f "$(BASE)" || { echo 'Baseline profile not found: $(BASE)'; exit 1; }
	@test -f "$(PROFILE_DIR)/mem.prof" || { echo 'Current memory profile not found; run make profile first.'; exit 1; }
	@test -f "$(PROFILE_DIR)/flog.test" || { echo 'Current profile binary not found; run make profile first.'; exit 1; }
	$(GO) tool pprof -http="$(COMPARE_ADDR)" -alloc_space -base="$(BASE)" "$(PROFILE_DIR)/flog.test" "$(PROFILE_DIR)/mem.prof"
