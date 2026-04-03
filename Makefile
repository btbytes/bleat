BINARY_NAME := bleat
VERSION_FILE := VERSION

VERSION := $(shell cat $(VERSION_FILE))
MAJOR_MINOR := $(word 1,$(subst ., ,$(VERSION))).$(word 2,$(subst ., ,$(VERSION)))
PATCH := $(shell git rev-list --count HEAD 2>/dev/null || echo 0)
FULL_VERSION := $(MAJOR_MINOR).$(PATCH)

LDFLAGS := -ldflags="-s -w"

.PHONY: build clean release

build:
	@go build $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "Built $(BINARY_NAME)"

release:
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o release/$(BINARY_NAME)-darwin-arm64 .
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o release/$(BINARY_NAME)-darwin-amd64 .
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o release/$(BINARY_NAME)-linux-amd64 .
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o release/$(BINARY_NAME)-linux-arm64 .
	@echo "Built all release binaries"

clean:
	@rm -f $(BINARY_NAME)
	@rm -rf release/
	@echo "Cleaned $(BINARY_NAME)"
