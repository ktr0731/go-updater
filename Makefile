SHELL := /bin/bash

.PHONY: build
build:
	go build

.PHONY: test
test:
	go test -race -v ./...
