# ──────────────────────────────────────────────────────────────
#  grpc-boundary-lab Master orchestrator  ·  Makefile
# ──────────────────────────────────────────────────────────────

# ── Colours ──────────────────────────────────────────────────
CYAN      := \033[36m
GREEN     := \033[32m
YELLOW    := \033[33m
BOLD      := \033[1m
RESET     := \033[0m

.DEFAULT_GOAL := help
.PHONY: help install build clean test check run stop kill-ports backend gateway frontend loadgen sweep docs docs-deploy

help: ## Show this help
	@printf "\n$(BOLD)$(CYAN)grpc-boundary-lab Orchestrator$(RESET)\n\n"
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | \
		awk -F ':.*## ' '{printf "  $(GREEN)%-18s$(RESET) %s\n", $$1, $$2}'
	@echo ""

## ── Global targets ──────────────────────────────────────────

all: install build check run ## Install dependencies, build everything, run checks, and launch all containers

run: ## Start all containers in the background via Docker Compose
	@printf "\n$(BOLD)$(CYAN)Starting all containers$(RESET)\n"
	docker compose up --build -d

stop: ## Stop all containers and kill local processes holding ports
	@printf "\n$(BOLD)$(YELLOW)Stopping containers and clearing ports...$(RESET)\n"
	docker compose down --remove-orphans || true
	$(MAKE) kill-ports

kill-ports: ## Kill processes on 8080, 8081, 50051, 50052
	@printf "$(YELLOW)▸ killing processes on app ports…$(RESET)\n"
	@fuser -k 8080/tcp 8081/tcp 50051/tcp 50052/tcp 2>/dev/null || true

install: ## Install dependencies (frontend)
	$(MAKE) -C frontend install

build: ## Build both backend and frontend
	@printf "\n$(BOLD)$(CYAN)Building Backend$(RESET)\n"
	$(MAKE) -C backend build
	@printf "\n$(BOLD)$(CYAN)Building Frontend$(RESET)\n"
	$(MAKE) -C frontend build

clean: ## Clean generated build files across projects
	$(MAKE) -C backend clean
	$(MAKE) -C frontend clean

test: ## Run tests (backend)
	$(MAKE) -C backend test

check: ## Run tests and linting
	$(MAKE) -C backend check
	# $(MAKE) -C frontend lint # Uncomment when linting is configured

## ── Application Dev ─────────────────────────────────────────

backend: ## Run the backend server
	$(MAKE) -C backend backend

gateway: ## Run the gateway server
	$(MAKE) -C backend gateway

frontend: ## Run the frontend dev server
	$(MAKE) -C frontend dev

## ── Benchmarking ────────────────────────────────────────────

loadgen: ## Run the load generator straight to backend
	$(MAKE) -C backend loadgen

loadgen-gw: ## Run the load generator to the gateway proxy
	$(MAKE) -C backend loadgen-gw

sweep: ## Run a performance sweep
	$(MAKE) -C backend sweep

## ── Documentation ───────────────────────────────────────────

docs: ## Serve docs locally
	$(MAKE) -C backend docs

docs-deploy: ## Deploy the documentation
	$(MAKE) -C backend docs-deploy
