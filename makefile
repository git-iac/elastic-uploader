BINARY_NAME = ./build/elastic-uploader
CMD_PATH = ./cmd

.PHONY: all build run tidy lint clean help
# Default target: build the application
all: build

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(CMD_PATH)
	@echo "Build complete: $(BINARY_NAME)"

run: build
	@echo "Running $(BINARY_NAME) with args: $(ARGS)"
	@./$(BINARY_NAME) $(ARGS)

tidy:
	@echo "Running go mod tidy..."
	@go mod tidy

clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)

