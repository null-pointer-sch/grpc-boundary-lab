# ──────────────────────────────────────────────────────────────
#  grpc-boundary-lab  ·  Makefile
# ──────────────────────────────────────────────────────────────

BIN       := bin
GO        := go
MODULE    := ./cmd

# ── Colours ──────────────────────────────────────────────────
CYAN      := \033[36m
GREEN     := \033[32m
YELLOW    := \033[33m
BOLD      := \033[1m
RESET     := \033[0m

# ── Server defaults ─────────────────────────────────────────
BACKEND_PORT  ?= 50051
BACKEND_REST_PORT ?= 8081
GATEWAY_PORT  ?= 50052
GATEWAY_REST_PORT  ?= 8080
TARGET_HOST   ?= 127.0.0.1

# ── Loadgen defaults ────────────────────────────────────────
REQUESTS    ?= 50000
WARMUP      ?= 2000
CONCURRENCY ?= 32
DEADLINE_MS ?= 20000
RUNS        ?= 1

# ── Sweep defaults ──────────────────────────────────────────
SWEEP_CONC  ?= 1 4 16 32 64 128

# ─────────────────────────────────────────────────────────────
#  Targets
# ─────────────────────────────────────────────────────────────

.DEFAULT_GOAL := help

.PHONY: help build clean test vet check
.PHONY: backend gateway
.PHONY: loadgen loadgen-backend loadgen-gw loadgen-2k sweep
.PHONY: docs docs-build docs-deploy

## ── Dev ─────────────────────────────────────────────────────

help: ## Show this help
	@printf "\n$(BOLD)$(CYAN)grpc-boundary-lab$(RESET)\n\n"
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | \
		awk -F ':.*## ' '{printf "  $(GREEN)%-18s$(RESET) %s\n", $$1, $$2}'
	@printf "\n$(BOLD)Loadgen vars$(RESET) (override with make … VAR=val):\n"
	@printf "  TARGET_HOST=$(TARGET_HOST)  REQUESTS=$(REQUESTS)  WARMUP=$(WARMUP)\n"
	@printf "  CONCURRENCY=$(CONCURRENCY)  DEADLINE_MS=$(DEADLINE_MS)  RUNS=$(RUNS)\n\n"

build: ## Build all binaries → bin/
	@printf "$(CYAN)▸ building…$(RESET)\n"
	@mkdir -p $(BIN)
	$(GO) build -buildvcs=false -o $(BIN)/backend  $(MODULE)/backend
	$(GO) build -buildvcs=false -o $(BIN)/gateway  $(MODULE)/gateway
	$(GO) build -buildvcs=false -o $(BIN)/loadgen  $(MODULE)/loadgen
	@printf "$(GREEN)✔ binaries in $(BIN)/$(RESET)\n"

clean: ## Remove built binaries
	rm -rf $(BIN)

test: ## Run all tests
	$(GO) test ./...

vet: ## Run go vet
	$(GO) vet ./...

check: test vet ## Run tests + vet

## ── Servers ─────────────────────────────────────────────────

backend: ## Start backend server
	BACKEND_PORT=$(BACKEND_PORT) $(GO) run $(MODULE)/backend

gateway: ## Start gateway server
	GATEWAY_PORT=$(GATEWAY_PORT) BACKEND_HOST=$(TARGET_HOST) BACKEND_PORT=$(BACKEND_PORT) \
		$(GO) run $(MODULE)/gateway

## ── Load generation ─────────────────────────────────────────

define RUN_LOADGEN
	WARMUP=$(WARMUP) \
	CONCURRENCY=$(CONCURRENCY) \
	DEADLINE_MS=$(DEADLINE_MS) \
	RUNS=$(RUNS) \
	REQUESTS=$(REQUESTS) \
	TARGET_HOST=$(TARGET_HOST) \
	TARGET_PORT=$(1) \
	MODE=$(2) \
	$(GO) run $(MODULE)/loadgen
endef

loadgen: loadgen-backend ## Alias for loadgen-backend

loadgen-backend: ## Loadgen → backend (gRPC direct)
	$(call RUN_LOADGEN,$(BACKEND_PORT),grpc)

loadgen-backend-rest: ## Loadgen → backend (REST direct)
	$(call RUN_LOADGEN,$(BACKEND_REST_PORT),rest)

loadgen-gw: ## Loadgen → gateway (gRPC proxy)
	$(call RUN_LOADGEN,$(GATEWAY_PORT),grpc)

loadgen-gw-rest: ## Loadgen → gateway (REST proxy)
	$(call RUN_LOADGEN,$(GATEWAY_REST_PORT),rest)

loadgen-2k: ## Quick 2 k-request run → backend
	$(MAKE) --no-print-directory loadgen-backend REQUESTS=2000

## ── Benchmarks (Gateway Router) ─────────────────────────────

bench-rest: ## Benchmark REST routing to gateway (TLS off)
	PROTOCOL=rest TLS=0 $(call RUN_LOADGEN,$(GATEWAY_REST_PORT),rest)

bench-rest-tls: ## Benchmark REST routing to gateway (TLS on)
	PROTOCOL=rest TLS=1 $(call RUN_LOADGEN,$(GATEWAY_REST_PORT),rest)

bench-grpc: ## Benchmark gRPC routing to gateway (TLS off)
	PROTOCOL=grpc TLS=0 $(call RUN_LOADGEN,$(GATEWAY_PORT),grpc)

bench-grpc-tls: ## Benchmark gRPC routing to gateway (TLS on)
	PROTOCOL=grpc TLS=1 $(call RUN_LOADGEN,$(GATEWAY_PORT),grpc)

bench-matrix: ## Run all 4 benchmark modes and summarize results
	bash bench-matrix.sh

sweep: ## Sweep concurrency levels
	@printf "\n$(BOLD)$(YELLOW)sweep$(RESET)  REQUESTS=$(REQUESTS) WARMUP=$(WARMUP) DEADLINE_MS=$(DEADLINE_MS) RUNS=$(RUNS)\n"
	@for c in $(SWEEP_CONC); do \
	  printf "\n$(BOLD)──── backend (gRPC)  c=$$c ────$(RESET)\n"; \
	  $(MAKE) --no-print-directory loadgen-backend CONCURRENCY=$$c; \
	  printf "\n$(BOLD)──── backend (REST)  c=$$c ────$(RESET)\n"; \
	  $(MAKE) --no-print-directory loadgen-backend-rest CONCURRENCY=$$c; \
	  printf "\n$(BOLD)──── gateway (gRPC)  c=$$c ────$(RESET)\n"; \
	  $(MAKE) --no-print-directory loadgen-gw      CONCURRENCY=$$c; \
	  printf "\n$(BOLD)──── gateway (REST)  c=$$c ────$(RESET)\n"; \
	  $(MAKE) --no-print-directory loadgen-gw-rest CONCURRENCY=$$c; \
	done

## ── Docs ────────────────────────────────────────────────────

docs: ## Serve docs locally
	cd docs && poetry run mkdocs serve

docs-build: ## Build docs site
	cd docs && poetry run mkdocs build

docs-deploy: ## Deploy docs to GitHub Pages
	cd docs && poetry run mkdocs gh-deploy