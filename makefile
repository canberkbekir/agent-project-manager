SHELL := /bin/sh

BIN_DIR := bin
CONFIG  := configs/config.yaml

AGENTD_PKG   := ./cmd/agentd
AGENTCTL_PKG := ./cmd/agentctl

# Detect OS to handle executable suffixes and clean commands
ifeq ($(OS),Windows_NT)
	EXE := .exe
	MKDIR_P := if not exist $(BIN_DIR) mkdir $(BIN_DIR)
	RM_RF := rmdir /s /q
	RM_F := del /q
else
	EXE :=
	MKDIR_P := mkdir -p $(BIN_DIR)
	RM_RF := rm -rf
	RM_F := rm -f
endif

AGENTD_BIN   := $(BIN_DIR)/agentd$(EXE)
AGENTCTL_BIN := $(BIN_DIR)/agentctl$(EXE)

GO ?= go
GOFLAGS ?=
LDFLAGS ?=

.PHONY: all tidy fmt vet test build build-agentd build-agentctl \
        run-agentd run-agentctl debug-agentd debug-agentctl clean help \
        swagger swagger-install

all: build

tidy:
	$(GO) mod tidy

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

test:
	$(GO) test ./...

build: build-agentd build-agentctl

build-agentd:
	@$(MKDIR_P)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(AGENTD_BIN) $(AGENTD_PKG)

build-agentctl:
	@$(MKDIR_P)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(AGENTCTL_BIN) $(AGENTCTL_PKG)

run-agentd:
	$(GO) run $(GOFLAGS) $(AGENTD_PKG) --config $(CONFIG)

run-agentctl:
	$(GO) run $(GOFLAGS) $(AGENTCTL_PKG) --help

# Debug (requires Delve: go install github.com/go-delve/delve/cmd/dlv@latest)
# Builds with optimizations disabled for better debugging.
debug-agentd:
	dlv debug $(AGENTD_PKG) --build-flags="-gcflags=all=-N -l" -- --config $(CONFIG)

debug-agentctl:
	dlv debug $(AGENTCTL_PKG) --build-flags="-gcflags=all=-N -l" -- --help

clean:
	-@$(RM_RF) $(BIN_DIR) 2>nul || true
	-@$(RM_RF) docs 2>nul || true

swagger-install:
	@echo "Installing swag..."
	@$(GO) install github.com/swaggo/swag/cmd/swag@latest

swagger: swagger-install
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/agentd/main.go -o docs --parseDependency --parseInternal
	@echo "Swagger docs generated in docs/"

help:
	@echo Targets:
	@echo \  tidy\ \ \ \ \ \ \ \ - go mod tidy
	@echo \  fmt\ \ \ \ \ \ \ \ \ - go fmt ./...
	@echo \  vet\ \ \ \ \ \ \ \ \ - go vet ./...
	@echo \  test\ \ \ \ \ \ \ \ - go test ./...
	@echo \  build\ \ \ \ \ \ \ - build binaries into $(BIN_DIR)
	@echo \  run-agentd\ \ \ \ - run agentd with config
	@echo \  run-agentctl\ \ \ - run agentctl help
	@echo \  debug-agentd\ \ \ - debug agentd with dlv
	@echo \  debug-agentctl\ \ - debug agentctl with dlv
	@echo \  clean\ \ \ \ \ \ \ - remove build artifacts
	@echo \  swagger\ \ \ \ \ \ - generate Swagger API documentation