# ──────────────────────────────────────────────────────────────
#  grpc-boundary-lab frontend  ·  Makefile
# ──────────────────────────────────────────────────────────────

# ── Colours ──────────────────────────────────────────────────
CYAN      := \033[36m
GREEN     := \033[32m
BOLD      := \033[1m
RESET     := \033[0m

.DEFAULT_GOAL := help
.PHONY: help install dev build preview clean

help: ## Show this help
	@printf "\n$(BOLD)$(CYAN)Frontend Commands$(RESET)\n\n"
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | \
		awk -F ':.*## ' '{printf "  $(GREEN)%-18s$(RESET) %s\n", $$1, $$2}'
	@echo ""

install: ## Install dependencies
	@printf "$(CYAN)▸ installing dependencies…$(RESET)\n"
	npm install

dev: ## Start Vite development server
	npm run dev

build: ## Build for production
	@printf "$(CYAN)▸ building frontend…$(RESET)\n"
	npm run build
	@printf "$(GREEN)✔ frontend built in dist/$(RESET)\n"

preview: ## Preview the production build
	npm run preview

clean: ## Remove node_modules and dist overrides
	rm -rf node_modules dist
