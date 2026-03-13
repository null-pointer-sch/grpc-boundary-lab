# ──────────────────────────────────────────────────────────────
#  grpc-boundary-lab Master orchestrator  ·  Makefile
# ──────────────────────────────────────────────────────────────

# ── Colours ──────────────────────────────────────────────────
CYAN      := \033[36m
GREEN     := \033[32m
YELLOW    := \033[33m
BOLD      := \033[1m
RESET     := \033[0m

# ── Ports ────────────────────────────────────────────────────
FRONTEND_PORT      ?= 8082
GATEWAY_REST_PORT  ?= 8080
GATEWAY_GRPC_PORT  ?= 50052
BACKEND_REST_PORT  ?= 8081
BACKEND_GRPC_PORT  ?= 50051
BACKEND_TLS_PORT   ?= 50151
BACKEND_REST_TLS   ?= 8181

.DEFAULT_GOAL := help
.PHONY: help install build clean test check run stop kill-ports backend gateway frontend loadgen sweep docs docs-deploy certs

help: ## Show this help
	@printf "\n$(BOLD)$(CYAN)grpc-boundary-lab Orchestrator$(RESET)\n\n"
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | \
		awk -F ':.*## ' '{printf "  $(GREEN)%-18s$(RESET) %s\n", $$1, $$2}'
	@echo ""

## ── Global targets ──────────────────────────────────────────

all: stop install build check run ## Stop existing, install dependencies, build everything, run checks, and launch all containers

run: ## Start all containers in the background via Docker Compose
	@printf "\n$(BOLD)$(CYAN)Starting all containers$(RESET)\n"
	FRONTEND_PORT=$(FRONTEND_PORT) \
	GATEWAY_REST_PORT=$(GATEWAY_REST_PORT) \
	GATEWAY_GRPC_PORT=$(GATEWAY_GRPC_PORT) \
	BACKEND_REST_PORT=$(BACKEND_REST_PORT) \
	BACKEND_GRPC_PORT=$(BACKEND_GRPC_PORT) \
	BACKEND_TLS_PORT=$(BACKEND_TLS_PORT) \
	BACKEND_REST_TLS=$(BACKEND_REST_TLS) \
	docker compose up --build -d
	@printf "\n$(BOLD)$(GREEN)✔ Services are up!$(RESET)\n"
	@printf "  $(BOLD)Frontend UI:$(RESET)         http://localhost:$(FRONTEND_PORT)\n"
	@printf "  $(BOLD)Gateway API (REST):$(RESET)  http://localhost:$(GATEWAY_REST_PORT)/api/ping\n"
	@printf "  $(BOLD)Backend API (REST):$(RESET)  http://localhost:$(BACKEND_REST_PORT)/api/ping\n\n"

stop: ## Stop all containers and kill local processes holding ports
	@printf "\n$(BOLD)$(YELLOW)Stopping containers and clearing ports...$(RESET)\n"
	docker compose down --remove-orphans || true
	$(MAKE) kill-ports

kill-ports: ## Kill processes on configured app ports
	@printf "$(YELLOW)▸ killing processes on app ports…$(RESET)\n"
	@fuser -k $(FRONTEND_PORT)/tcp $(GATEWAY_REST_PORT)/tcp $(GATEWAY_GRPC_PORT)/tcp \
		$(BACKEND_REST_PORT)/tcp $(BACKEND_GRPC_PORT)/tcp $(BACKEND_TLS_PORT)/tcp \
		$(BACKEND_REST_TLS)/tcp 2>/dev/null || true

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
	rm -f coverage.out

test: ## Run tests (backend)
	$(MAKE) -C backend test

check: certs ## Run tests and linting
	$(MAKE) -C backend check
	# $(MAKE) -C frontend lint # Uncomment when linting is configured

certs: ## Generate and setup certificates in /tmp/certs
	@printf "\n$(BOLD)$(CYAN)Setting up certificates$(RESET)\n"
	@bash certs/gen.sh
	@mkdir -p /tmp/certs
	@cp certs/*.crt certs/*.key /tmp/certs/
	@printf "$(GREEN)✔ Certificates ready in /tmp/certs/$(RESET)\n"

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
