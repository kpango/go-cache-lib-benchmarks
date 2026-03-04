GO_VERSION := 1.26.0
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

bench: deps
	sleep 3
	go test -count=1 -timeout=30m -run=NONE -bench . -benchmem

init:
	GO111MODULE=on go mod init github.com/kpango/go-cache-lib-benchmarks
	GO111MODULE=on go mod tidy
	go get -u ./...

format:
	find ./ -type d -name .git -prune -o -type f -regex '.*[^\.pb]\.go' -print | xargs $(GOPATH)/bin/golines -w -m $(GOLINES_MAX_WIDTH)
	find ./ -type d -name .git -prune -o -type f -regex '.*[^\.pb]\.go' -print | xargs $(GOPATH)/bin/gofumpt -w
	find ./ -type d -name .git -prune -o -type f -regex '.*[^\.pb]\.go' -print | xargs $(GOPATH)/bin/strictgoimports -w
	find ./ -type d -name .git -prune -o -type f -regex '.*\.go' -print | xargs $(GOPATH)/bin/goimports -w
