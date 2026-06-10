# tryit — top-level Makefile
# Phase 1 targets only. Phase 2/3 add release-server, release-ext, etc.

SHELL := /bin/bash
SERVER_DIR := server
EXT_DIR    := extension

.PHONY: help dev server-dev server-build server-test server-lint \
        ext-install ext-dev ext-build contract-check reset-pairing clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?##' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?##"}{printf "  %-18s %s\n",$$1,$$2}'

dev: ## Run server + extension dev (two foreground processes — run in two terminals or use `make -j2 dev`)
	@echo "Run 'make server-dev' in one terminal and 'make ext-dev' in another."

server-dev: ## Run the Go server (prints pairing token on first start)
	cd $(SERVER_DIR) && go run ./cmd/tryit

server-build: ## Build the Go server binary into server/tryit
	cd $(SERVER_DIR) && go build -o tryit ./cmd/tryit

server-test: ## Run Go unit tests
	cd $(SERVER_DIR) && go test ./...

server-lint: ## Run golangci-lint (install via brew: brew install golangci-lint)
	cd $(SERVER_DIR) && golangci-lint run

ext-install: ## Install extension dependencies
	cd $(EXT_DIR) && npm install

ext-dev: ## Build extension in watch mode -> extension/dist
	cd $(EXT_DIR) && npm run dev

ext-build: ## One-shot extension build
	cd $(EXT_DIR) && npm run build

contract-check: ## Round-trip a fixture through the Go struct + JSON schema (D4)
	cd $(SERVER_DIR) && go test ./api/... -run TestContractRoundTrip -v

reset-pairing: ## Forget the pairing token and bound origin
	rm -f $(HOME)/.tryit/pair.json
	@echo "Pairing reset. Restart the server to generate a fresh token."

clean:
	rm -rf $(SERVER_DIR)/tryit $(EXT_DIR)/dist $(EXT_DIR)/node_modules
