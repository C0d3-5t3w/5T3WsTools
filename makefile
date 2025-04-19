.PHONY: deps

deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download
	@go mod vendor
	@go mod verify
	@echo "Dependencies installed."

