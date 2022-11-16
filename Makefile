
generate:
	go generate github.com/ekhvalov/otus-banners-rotation/internal/environment/server/grpc

test:
	go test -race -count 100 ./internal/...

test-simple:
	go test ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: generate test test-simple lint
