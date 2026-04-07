.PHONY: build run test lint deps clean

# Binary name
APP_NAME=incident.exe

build: ## Build the CLI executable
	go build -o $(APP_NAME) main.go

run: build ## Run the built CLI executable
	./$(APP_NAME)

test: ## Run the suite of tests
	go test ./... -v

lint: ## Run linters and checks
	go vet ./...
	cspell "**/*"

deps: ## Download and tidy dependencies
	go mod download
	go mod tidy

clean: ## Clean up built binary
	rm -f $(APP_NAME)

help: ## Show this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
