.PHONY: help run build clean test download-htmx

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

download-htmx: ## Download HTMX library
	@echo "Downloading HTMX..."
	@curl -sL https://unpkg.com/htmx.org@1.9.10/dist/htmx.min.js -o ui/static/js/htmx.min.js || \
	 wget -q https://unpkg.com/htmx.org@1.9.10/dist/htmx.min.js -O ui/static/js/htmx.min.js || \
	 echo "Failed to download HTMX. Please download manually."
	@echo "HTMX downloaded successfully!"

run: ## Run the application in development mode
	@echo "Starting application..."
	@go run main.go

build: ## Build the application binary
	@echo "Building application..."
	@go build -o perma-app main.go
	@echo "Build complete: ./perma-app"

build-optimized: ## Build optimized binary (smaller size)
	@echo "Building optimized binary..."
	@go build -ldflags="-s -w" -o perma-app main.go
	@echo "Optimized build complete: ./perma-app"

clean: ## Remove built binaries and database
	@echo "Cleaning up..."
	@rm -f perma-app app.db app.db-shm app.db-wal
	@echo "Clean complete!"

test: ## Run tests (placeholder)
	@echo "Running tests..."
	@go test -v ./...

deps: ## Download Go dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies ready!"

setup: deps download-htmx ## Complete setup (deps + HTMX)
	@echo "Setup complete! Run 'make run' to start the application."

.DEFAULT_GOAL := help
