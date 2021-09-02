BASE_PATH := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
MKFILE_PATH := $(BASE_PATH)/Makefile

.DEFAULT_GOAL := help

all: clean echo-server echo-client ## Clean up and build binaries for server and client

echo-server: ## Build binaries for echo server
	go build -o bin/echo-server ./cmd/echo-server

echo-client: ## Build binaries for echo client
	go build -o bin/echo-client ./cmd/echo-client

clean: ## Remove binaries
	@rm -rf bin
	@find $(BASE_PATH) -name ".DS_Store" -depth -exec rm {} \;

help: ## Show help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: all help echo-server echo-client clean
