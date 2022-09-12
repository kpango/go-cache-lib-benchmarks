GO_VERSION:=$(shell go version)

clean:
	go clean ./...
	go clean -modcache
	rm -rf ./*.log
	rm -rf ./*.svg
	rm -rf ./go.*
	rm -rf pprof
	rm -rf bench
	rm -rf vendor

bench: clean init
	sleep 3
	go test -count=1 -run=NONE -bench . -benchmem

init:
	GO111MODULE=on go mod init
	GO111MODULE=on go mod tidy
