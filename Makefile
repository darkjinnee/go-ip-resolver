.PHONY: build run clean test

# Переменные
BINARY_NAME=ip-resolver
BUILD_DIR=bin

BINARY_NAME_REALISE=ip-resolver-linux
BUILD_DIR_REALISE=bin/linux

CONFIG_PATH=configs/domains.json

# Релизная сборка
realise:
	@echo "Building realise $(BINARY_NAME_REALISE)..."
	@mkdir -p $(BUILD_DIR_REALISE)
	@go build -o $(BUILD_DIR_REALISE)/$(BINARY_NAME_REALISE) ./cmd/ip-resolver

# Сборка
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/ip-resolver

# Запуск
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME) -config $(CONFIG_PATH)

# Запуск с кастомными параметрами
run-dev: build
	@echo "Running $(BINARY_NAME) in development mode..."
	@./$(BUILD_DIR)/$(BINARY_NAME) -config $(CONFIG_PATH) -addr :8080 -timeout 10s

# Тесты
test:
	@echo "Running tests..."
	@go test -v ./...

# Очистка
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

# Установка зависимостей
deps:
	@echo "Installing dependencies..."
	@go mod tidy

# Форматирование кода
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Линтинг
lint:
	@echo "Running linter..."
	@go vet ./...

# Полная проверка
check: fmt lint test
	@echo "All checks passed!"
