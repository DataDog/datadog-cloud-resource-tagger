MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))

ROOT_DIR := $(dir $(MAKEFILE_PATH))

# Define binaries directory
BIN_DIR := $(ROOT_DIR)/bin

BINARY := datadog-cloud-resource-tagger
PKG := github.com/Datadog/cloud-resource-tagger

build:
	go install && CGO_ENABLED=0 go build -o $(BIN_DIR)/$(BINARY) $(PKG)
	@echo "Build completed. Binaries are saved in $(BIN_DIR)"

update:
	go get -u
	go mod tidy

install:
	go install

test:
	@echo "Running tests..."
	go test ./... -v
	@echo "Tests completed successfully."

thirdparty-licenses:
	@echo "Retrieving third-party licenses..."
	go get github.com/google/go-licenses
	go install github.com/google/go-licenses
	@go-licenses report github.com/Datadog/cloud-resource-tagger > LICENSE-3rdparty.csv
	@echo "Third-party licenses retrieved and saved to LICENSE-3rdparty.csv"