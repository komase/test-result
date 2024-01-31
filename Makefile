CURRENT_REVISION= $(shell git rev-parse --short HEAD)
LDFLAGS= "-s -w -X main.revision=$(CURRENT_REVISION)"

.PHONY: all
all: build

.PHONY: build
build:
		go build -ldflags=$(LDFLAGS) -o test-result

.PHONY: test
test:
		go test -v ./...

.PHONY: clean
clean:
	rm -rf test-result