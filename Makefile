build:
	go build -v -o ./bin/rotator ./cmd/rotator

build-img-base:
	docker build \
		--tag $(DOCKER_IMG_BASE) \
		--file build/base.Dockerfile .

build-img:
	docker build \
		--tag "otus-golang/rotator:develop" \
		--force-rm=true \
		--file build/Dockerfile .

generate:
	go generate github.com/ekhvalov/otus-banners-rotation/internal/app
	go generate github.com/ekhvalov/otus-banners-rotation/internal/environment/storage/redis
	go generate github.com/ekhvalov/otus-banners-rotation/internal/environment/server/grpc

test:
	go test -race -count 100 ./internal/...

test-integration: build-img
	cd deployments/integration-tests && docker-compose up --exit-code-from tester && docker-compose down --volumes

test-simple:
	go test ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build build-img generate test test-integration test-simple lint
