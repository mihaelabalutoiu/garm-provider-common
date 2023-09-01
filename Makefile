SHELL := bash

GO_PACKAGES?=$(shell (go list ./... | grep -v 'vendor'))

.PHONY: go-test

go-test:
	go test $(GO_PACKAGES) -v $(TEST_ARGS) -timeout=15m -parallel=4
