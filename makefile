# Папки с бинарниками
BIN_DIR := bin

# Пути к main.go
PRODUCER := ./cmd/producer
CONSUMER := ./cmd/consumer

# Имена бинарников
PRODUCER_BIN := $(BIN_DIR)/producer
CONSUMER_BIN := $(BIN_DIR)/consumer

.PHONY: all run-producer run-consumer build-producer build-consumer clean

all: build-producer build-consumer

# =========================
# Сборка
# =========================
build-producer:
	@echo "🚀 Building producer..."
	@go build -o $(PRODUCER_BIN) $(PRODUCER)

build-consumer:
	@echo "🚀 Building consumer..."
	@go build -o $(CONSUMER_BIN) $(CONSUMER)

# =========================
# Запуск без сохранения бинарника
# =========================
run-producer:
	@echo "▶️ Running producer..."
	@go run $(PRODUCER)

run-consumer:
	@echo "▶️ Running consumer..."
	@go run $(CONSUMER)

# =========================
# Удалить бинарники
# =========================
clean:
	@echo "🧹 Cleaning binaries..."
	@rm -rf $(BIN_DIR)/*
