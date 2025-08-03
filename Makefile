
# Makefile for TFX

.DEFAULT_GOAL := test

GO_TEST_FLAGS := -covermode=atomic -coverpkg=./... -coverprofile=coverage.out
GO_RACE_FLAGS := -race
GO_VERBOSE    := -v

.PHONY: test test-verbose test-race coverage clean demo build-demo

test:
	go test ./... -short $(GO_TEST_FLAGS)

#  Test only one package
test-one:
	@echo "Usage: make test-one PKG=path/to/package"
	@echo "Example: make test-one PKG=./logfx"
	@echo "Running tests for package: $(PKG)"
	go test $(PKG) $(GO_TEST_FLAGS)

test-verbose:
	go test ./... $(GO_TEST_FLAGS) $(GO_VERBOSE)

test-race:
	go test ./... $(GO_TEST_FLAGS) $(GO_RACE_FLAGS)

coverage:
	go tool cover -html=coverage.out

clean:
	rm -f coverage.out
	rm -f bin/demo

build-demo:
	mkdir -p bin
	go build -o bin/demo ./cmd/demo

demo: build-demo
	./bin/demo

fix:
	goimports -w .
	gofumpt -w .
	gci write -s standard -s default -s "prefix($(shell go list -m))" .
	go mod tidy
	go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -fix ./...
# 	golangci-lint run --fix || true

tidy:
	go mod tidy