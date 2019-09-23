SHELL=/bin/sh

export GO111MODULE=on

.PHONY: ci
ci: deps lint build

.PHONY: deps
deps:
	go mod download
	go mod vendor

.PHONY: build
build:
	docker build -t ivchtest .

.PHONY: lint
lint:
	GO111MODULE=off go get github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run

.PHONY: run
run: build
	docker run ivchtest

.PHONY: clean
clean:
	docker rmi ivchtest -f