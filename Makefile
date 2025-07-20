
# Makefile for TFX

.DEFAULT_GOAL := test

GO_TEST_FLAGS := -covermode=atomic -coverpkg=./... -coverprofile=coverage.out
GO_RACE_FLAGS := -race
GO_VERBOSE    := -v

.PHONY: test test-verbose test-race coverage clean demo build-demo

test:
	go test ./... $(GO_TEST_FLAGS)

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
