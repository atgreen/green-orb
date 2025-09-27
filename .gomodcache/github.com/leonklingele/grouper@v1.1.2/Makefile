BINARY_OUT := ./grouper
COVERAGE_OUT := ./go.coverage

.PHONY: all
all: build test lint

.PHONY: build
build:
	go build -v -o ${BINARY_OUT}

.PHONY: clean
clean:
	go clean

.PHONY: test
test:
	go test -v -race ./...

.PHONY: test-cover
test-cover:
	go test -v -race -covermode=atomic -coverprofile=${COVERAGE_OUT} ./...

.PHONY: test-cover-web
test-cover-web: test-cover
	go tool cover -html=${COVERAGE_OUT}

.PHONY: lint
lint: golint

GOLANGCI_OUT_FORMAT ?= colored-line-number
.PHONY: golint
golint:
	@golangci-lint run -v \
	--out-format $(GOLANGCI_OUT_FORMAT) \
	--enable-all
