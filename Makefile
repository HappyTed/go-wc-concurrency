BINARY_NAME := go-wc

.PHONY: all build test

# --- По умолчанию ---
all: build

# --- Сборка ---
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o ./bin/$(BINARY_NAME)

# --- Тесты ---
test:
	@echo " > Running a static code analyzer..."
	go vet
	@echo " > Running tests..."
	go test -v -cover ./...
