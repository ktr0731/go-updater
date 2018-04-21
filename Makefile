SHELL := /bin/bash

.PHONY: dep
dep:
ifeq ($(shell which dep 2>/dev/null),)
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif

.PHONY: deps
deps: dep
	dep ensure

.PHONY: build
build: deps
	go build

.PHONY: test
test:
	go test -race -v ./...
