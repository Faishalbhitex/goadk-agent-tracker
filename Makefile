# =========================
# Project Configuration
# =========================
APP_NAME        := agenttracker
CLI_NAME        := agenttracker-cli
BIN_DIR         := bin

CMD_ADK         := ./cmd/adk/main.go
CMD_CLI         := ./cmd/cli/main.go

GO              := go
GOFLAGS         := -v

# =========================
# Default target
# =========================
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make run-adk        Run ADK agent (dev, with inspector/web UI)"
	@echo "  make run-cli        Run interactive CLI"
	@echo "  make run-bot        Run Telegram bot"
	@echo "  make build          Build all binaries"
	@echo "  make build-adk      Build ADK binary"
	@echo "  make build-cli      Build CLI binary"
	@echo "  make build-bot      Build Telegram bot binary"
	@echo "  make logs               Tail all bot logs"
	@echo "  make logs-tools         Tail tool execution logs only"
	@echo "  make logs-errors        Tail error logs only"
	@echo "  make logs-interactions  Tail user/agent interaction logs"
	@echo "  make logs-today         Tail today's logs"
	@echo "  make clean          Remove built binaries"
	@echo "  make tidy           Go mod tidy"
	@echo "  make test           Run go test"
	@echo "  make tree           Show project tree (no .git)"

# =========================
# Run (Development)
# =========================
.PHONY: run-adk
run-adk:
	$(GO) run $(CMD_ADK)

.PHONY: run-cli
run-cli:
	$(GO) run $(CMD_CLI)

# =========================
# Build
# =========================
.PHONY: build
build: build-adk build-cli

.PHONY: build-adk
build-adk:
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(APP_NAME) $(CMD_ADK)

.PHONY: build-cli
build-cli:
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(CLI_NAME) $(CMD_CLI)

# =========================
# Telegram Bot
# =========================
BOT_NAME        := agenttracker-bot
CMD_BOT         := ./cmd/bot/main.go

.PHONY: run-bot
run-bot:
	@mkdir -p logs
	$(GO) run $(CMD_BOT)

.PHONY: build-bot
build-bot:
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(BOT_NAME) $(CMD_BOT)

.PHONY: logs
logs:
	@if ls logs/bot_tools_*.log >/dev/null 2>&1; then \
		tail -f logs/bot_tools_*.log; \
	else \
		echo "No log files found"; \
	fi

.PHONY: logs-today
logs-today:
	@tail -f logs/bot_tools_$(shell date +%Y%m%d).log

# =========================
# Utilities
# =========================
.PHONY: clean
clean:
	rm -rf $(BIN_DIR)

.PHONY: tidy
tidy:
	$(GO) mod tidy

.PHONY: test
test:
	$(GO) test ./...

.PHONY: tree
tree:
	tree -a -I .git
