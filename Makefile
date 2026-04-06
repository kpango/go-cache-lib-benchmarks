GO_VERSION :="1.27.3"
GOPATH := $(eval GOPATH := $(shell go env GOPATH))$(GOPATH)
GOLINES_MAX_WIDTH     ?= 200

ROOTDIR = $(eval ROOTDIR := $(or $(shell git rev-parse --show-toplevel), $(PWD)))$(ROOTDIR)

all: clean install lint test bench

clean:
	go clean ./...
	go clean -modcache
	rm -rf \
	    $(ROOTDIR)/*.log \
	    $(ROOTDIR)/*.svg \
	    $(ROOTDIR)/go.mod \
	    $(ROOTDIR)/go.sum \
	    $(ROOTDIR)/pprof \
	    $(ROOTDIR)/bench \
	    $(ROOTDIR)/vendor

.PHONY: deps
## install Go package dependencies
deps: \
	clean \
	init
	head -n -1 $(ROOTDIR)/go.mod.default | awk 'NR>=6 && $$0 !~ /(upgrade|latest|master|main)/' | sort
	rm -rf $(ROOTDIR)/vendor \
	    $(ROOTDIR)/go.sum \
	    $(ROOTDIR)/go.mod 2>/dev/null
	cp $(ROOTDIR)/go.mod.default $(ROOTDIR)/go.mod
	sed -E "s/^go [0-9]+\.[0-9]+(\.[0-9]+)?/go $(GO_VERSION)/; s/#.*//" $(ROOTDIR)/go.mod > $(ROOTDIR)/go.mod.tmp
	mv $(ROOTDIR)/go.mod.tmp $(ROOTDIR)/go.mod
	GOTOOLCHAIN=go$(GO_VERSION) go mod tidy
	GOTOOLCHAIN=go$(GO_VERSION) go get -u all 2>/dev/null || true
	GOTOOLCHAIN=go$(GO_VERSION) go get -u ./... 2>/dev/null || true

BENCH_LOG ?= $(ROOTDIR)/bench.log

# bench-run runs the benchmarks, tees the output to $(BENCH_LOG) and prints a
# list of failed benchmarks. A failing library is reported as FAILED instead of
# aborting the run, so this target does not exit non-zero for them.
.PHONY: bench-run
bench-run:
	sleep 3
	go test -count=1 -timeout=30m -run=NONE -bench . -benchmem $(BENCH_ARGS) 2>&1 | tee $(BENCH_LOG)
	@echo ""
	@if grep -qE '^[[:space:]]*--- FAIL' $(BENCH_LOG); then \
	    echo "FAILED benchmarks:"; \
	    grep -E '^[[:space:]]*--- FAIL' $(BENCH_LOG) | sed -E 's/^[[:space:]]*--- FAIL: /  FAILED  /; s/ \([0-9.]+s\)$$//' | sort -u; \
	else \
	    echo "All benchmarks passed."; \
	fi

.PHONY: bench
bench: deps
	$(MAKE) bench-run

.PHONY: bench-custom
bench-custom: deps
	$(MAKE) bench-run BENCH_ARGS='-benchparallelism="$(BENCH_PARALLELISM)"'

.PHONY: visualize
## generate charts (images/*.svg|html, README charts) from $(BENCH_LOG)
visualize:
	@test -f $(BENCH_LOG) || { echo "$(BENCH_LOG) not found: run 'make bench' first"; exit 1; }
	cd $(ROOTDIR)/visualize && go run . -f $(BENCH_LOG)

init:
	GO111MODULE=on go mod init github.com/kpango/go-cache-lib-benchmarks
	GO111MODULE=on go mod tidy
	go get -u ./...

update-graph:
	cd $(ROOTDIR)/visualize && go run . -f ../README.md

format:
	find ./ -type d -name .git -prune -o -type f -regex '.*[^\.pb]\.go' -print | xargs $(GOPATH)/bin/golines -w -m $(GOLINES_MAX_WIDTH)
	find ./ -type d -name .git -prune -o -type f -regex '.*[^\.pb]\.go' -print | xargs $(GOPATH)/bin/gofumpt -w
	find ./ -type d -name .git -prune -o -type f -regex '.*[^\.pb]\.go' -print | xargs $(GOPATH)/bin/strictgoimports -w
	find ./ -type d -name .git -prune -o -type f -regex '.*\.go' -print | xargs $(GOPATH)/bin/goimports -w
