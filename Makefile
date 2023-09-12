.PHONY: all
all: lint test

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	go test -cover `go list ./... | grep -v -e /internal`