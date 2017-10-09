VERSION=0.1.0
COMMIT=$(shell git rev-parse --verify HEAD)

PACKAGES=$(shell go list ./... | grep -v /vendor/ | grep -v /cmd/)
BUILD_FLAGS=-ldflags "-X main.VERSION=$(VERSION) -X main.COMMIT=$(COMMIT)"

.PHONY: all
all: build

.PHONY: build
build: vendor
	go build $(BUILD_FLAGS) .

.PHONY: test
test: vendor
	go test -v $(PACKAGES)
	go vet $(PACKAGES)

.PHONY: clean
clean:
	rm -rf secretctl
	rm -rf dist

dist:
	mkdir -p dist
	
	GOARCH=amd64 GOOS=darwin go build $(BUILD_FLAGS) .
	tar -czf dist/secretctl_darwin_amd64.tar.gz secretctl
	rm -rf secretctl
	
	GOARCH=amd64 GOOS=linux go build $(BUILD_FLAGS) .
	tar -czf dist/secretctl_linux_amd64.tar.gz secretctl
	rm -rf secretctl
	
	GOARCH=arm64 GOOS=linux go build $(BUILD_FLAGS) .
	tar -czf dist/secretctl_linux_arm64.tar.gz secretctl
	rm -rf secretctl
	
	GOARCH=arm GOOS=linux go build $(BUILD_FLAGS) .
	tar -czf dist/secretctl_linux_arm.tar.gz secretctl
	rm -rf secretctl

vendor:
	glide install
