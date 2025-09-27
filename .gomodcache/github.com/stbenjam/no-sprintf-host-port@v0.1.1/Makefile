all: nosprintfhostport.so nosprintfhostport
.PHONY: lint test

clean:
	rm -f nosprintfhostport.so nosprintfhostport

test:
	go test ./...

lint:
	golangci-lint run ./...

nosprintfhostport:
	go build ./cmd/nosprintfhostport

nosprintfhostport.so:
	go build -buildmode=plugin ./plugin/nosprintfhostport.go
