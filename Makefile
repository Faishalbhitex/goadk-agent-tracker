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
	@echo "  make build          Build all binaries"
	@echo "  make build-adk      Build ADK binary"
	@echo "  make build-cli      Build CLI binary"
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
