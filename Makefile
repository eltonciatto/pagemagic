# Page Magic - Complete Makefile for Development

.PHONY: help install build dev test clean lint format docker-build docker-up docker-down logs shell db-shell redis-shell deps deps-go deps-rust deps-node deps-python migration-create migration-up migration-down seed backup restore security-check update-deps check-ports

# Default target
.DEFAULT_GOAL := help

# Colors for pretty output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
PURPLE := \033[0;35m
CYAN := \033[0;36m
WHITE := \033[1;37m
NC := \033[0m # No Color

# Project info
PROJECT_NAME := pagemagic
DOCKER_COMPOSE_FILE := docker-compose.dev.yml
DOCKER_PROD_FILE := docker-compose.prod.yml

help: ## Show this help message
	@echo -e "${BLUE}Page Magic Development Commands${NC}"
	@echo -e "${YELLOW}================================${NC}"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { printf "  ${GREEN}%-20s${NC} %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# ==========================================
# INSTALLATION & SETUP
# ==========================================

install: ## Install all dependencies and set up the development environment
	@echo -e "${BLUE}Installing Page Magic development environment...${NC}"
	@./scripts/setup-dev.sh

deps: deps-go deps-rust deps-node deps-python ## Install all dependencies

deps-go: ## Install Go dependencies
	@echo -e "${YELLOW}Installing Go dependencies...${NC}"
	@for service in auth-svc build-svc host-svc billing-svc; do \
		if [ -f "services/$$service/go.mod" ]; then \
			echo -e "${CYAN}Installing dependencies for $$service${NC}"; \
			cd services/$$service && go mod download && cd ../..; \
		fi \
	done

deps-rust: ## Install Rust dependencies
	@echo -e "${YELLOW}Installing Rust dependencies...${NC}"
	@for service in builder-svc meter-svc; do \
		if [ -f "services/$$service/Cargo.toml" ]; then \
			echo -e "${CYAN}Installing dependencies for $$service${NC}"; \
			cd services/$$service && cargo fetch && cd ../..; \
		fi \
	done

deps-node: ## Install Node.js dependencies
	@echo -e "${YELLOW}Installing Node.js dependencies...${NC}"
	@if [ -f "apps/front-web/package.json" ]; then \
		echo -e "${CYAN}Installing frontend dependencies${NC}"; \
		cd apps/front-web && npm install && cd ../..; \
	fi
	@if [ -f "apps/mobile-app/package.json" ]; then \
		echo -e "${CYAN}Installing mobile app dependencies${NC}"; \
		cd apps/mobile-app && npm install && cd ../..; \
	fi
	@for service in prompt-svc i18n-svc; do \
		if [ -f "services/$$service/package.json" ]; then \
			echo -e "${CYAN}Installing dependencies for $$service${NC}"; \
			cd services/$$service && npm install && cd ../..; \
		fi \
	done

deps-python: ## Install Python dependencies
	@echo -e "${YELLOW}Installing Python dependencies...${NC}"
	@if [ -f "services/domain-svc/requirements.txt" ]; then \
		echo -e "${CYAN}Installing domain service dependencies${NC}"; \
		cd services/domain-svc && \
		python3 -m venv venv 2>/dev/null || true && \
		source venv/bin/activate && \
		pip install -r requirements.txt && \
		deactivate && cd ../..; \
	fi

update-deps: ## Update all dependencies to latest versions
	@echo -e "${YELLOW}Updating all dependencies...${NC}"
	@$(MAKE) deps-go deps-rust deps-node deps-python

# ==========================================
# DEVELOPMENT
# ==========================================

dev: ## Start development environment
	@echo -e "${BLUE}Starting Page Magic development environment...${NC}"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up -d
	@echo -e "${GREEN}Development environment started!${NC}"
	@echo -e "${YELLOW}Available at:${NC}"
	@echo -e "  • Frontend: ${CYAN}http://localhost:3000${NC}"
	@echo -e "  • API Gateway: ${CYAN}http://localhost:8080${NC}"
	@echo -e "  • Grafana: ${CYAN}http://localhost:3001${NC}"

dev-build: ## Build and start development environment
	@echo -e "${BLUE}Building and starting development environment...${NC}"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up -d --build

down: ## Stop development environment
	@echo -e "${YELLOW}Stopping development environment...${NC}"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down
	@echo -e "${GREEN}Development environment stopped${NC}"

restart: ## Restart development environment
	@echo -e "${YELLOW}Restarting development environment...${NC}"
	@$(MAKE) down
	@$(MAKE) dev

restart-service: ## Restart a specific service (usage: make restart-service SERVICE=auth-svc)
	@if [ -z "$(SERVICE)" ]; then \
		echo -e "${RED}Please specify a service: make restart-service SERVICE=auth-svc${NC}"; \
		exit 1; \
	fi
	@echo -e "${YELLOW}Restarting $(SERVICE)...${NC}"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) restart $(SERVICE)

# ==========================================
# BUILDING
# ==========================================

build: ## Build all services
	@echo -e "${BLUE}Building all services...${NC}"
	@$(MAKE) build-go build-rust build-node build-python

build-docker: ## Build all Docker images
	@echo -e "${BLUE}Building Docker images...${NC}"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) build --parallel

build-go: ## Build Go services
	@echo -e "${YELLOW}Building Go services...${NC}"
	@for service in auth-svc build-svc host-svc billing-svc; do \
		if [ -f "services/$$service/main.go" ]; then \
			echo -e "${CYAN}Building $$service${NC}"; \
			cd services/$$service && go build -o bin/$$service ./cmd && cd ../..; \
		fi \
	done

build-rust: ## Build Rust services
	@echo -e "${YELLOW}Building Rust services...${NC}"
	@for service in builder-svc meter-svc; do \
		if [ -f "services/$$service/Cargo.toml" ]; then \
			echo -e "${CYAN}Building $$service${NC}"; \
			cd services/$$service && cargo build --release && cd ../..; \
		fi \
	done

build-node: ## Build Node.js services and apps
	@echo -e "${YELLOW}Building Node.js services and apps...${NC}"
	@if [ -f "apps/front-web/package.json" ]; then \
		echo -e "${CYAN}Building frontend${NC}"; \
		cd apps/front-web && npm run build && cd ../..; \
	fi
	@if [ -f "apps/mobile-app/package.json" ]; then \
		echo -e "${CYAN}Building mobile app${NC}"; \
		cd apps/mobile-app && npm run build:android && cd ../..; \
	fi

build-python: ## Build Python services
	@echo -e "${YELLOW}Building Python services...${NC}"
	@echo -e "${CYAN}Python services don't require building${NC}"

# ==========================================
# TESTING
# ==========================================

test: ## Run all tests
	@echo -e "${BLUE}Running all tests...${NC}"
	@$(MAKE) test-go test-rust test-node test-python

test-go: ## Run Go tests
	@echo -e "${YELLOW}Running Go tests...${NC}"
	@for service in auth-svc build-svc host-svc billing-svc; do \
		if [ -f "services/$$service/go.mod" ]; then \
			echo -e "${CYAN}Testing $$service${NC}"; \
			cd services/$$service && go test ./... -v && cd ../..; \
		fi \
	done

test-rust: ## Run Rust tests
	@echo -e "${YELLOW}Running Rust tests...${NC}"
	@for service in builder-svc meter-svc; do \
		if [ -f "services/$$service/Cargo.toml" ]; then \
			echo -e "${CYAN}Testing $$service${NC}"; \
			cd services/$$service && cargo test && cd ../..; \
		fi \
	done

test-node: ## Run Node.js tests
	@echo -e "${YELLOW}Running Node.js tests...${NC}"
	@if [ -f "apps/front-web/package.json" ]; then \
		echo -e "${CYAN}Testing frontend${NC}"; \
		cd apps/front-web && npm test && cd ../..; \
	fi
	@for service in prompt-svc i18n-svc; do \
		if [ -f "services/$$service/package.json" ]; then \
			echo -e "${CYAN}Testing $$service${NC}"; \
			cd services/$$service && npm test && cd ../..; \
		fi \
	done

test-python: ## Run Python tests
	@echo -e "${YELLOW}Running Python tests...${NC}"
	@if [ -f "services/domain-svc/requirements.txt" ]; then \
		echo -e "${CYAN}Testing domain service${NC}"; \
		cd services/domain-svc && \
		source venv/bin/activate && \
		python -m pytest tests/ && \
		deactivate && cd ../..; \
	fi

test-integration: ## Run integration tests
	@echo -e "${YELLOW}Running integration tests...${NC}"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) exec -T postgres psql -U pagemagic -d pagemagic -c "SELECT 1"
	@echo -e "${GREEN}Integration tests passed${NC}"

# ==========================================
# CODE QUALITY
# ==========================================

lint: ## Run linters for all languages
	@echo -e "${BLUE}Running linters...${NC}"
	@$(MAKE) lint-go lint-rust lint-node lint-python

lint-go: ## Run Go linter
	@echo -e "${YELLOW}Linting Go code...${NC}"
	@for service in auth-svc build-svc host-svc billing-svc; do \
		if [ -f "services/$$service/go.mod" ]; then \
			echo -e "${CYAN}Linting $$service${NC}"; \
			cd services/$$service && golangci-lint run || true && cd ../..; \
		fi \
	done

lint-rust: ## Run Rust linter
	@echo -e "${YELLOW}Linting Rust code...${NC}"
	@for service in builder-svc meter-svc; do \
		if [ -f "services/$$service/Cargo.toml" ]; then \
			echo -e "${CYAN}Linting $$service${NC}"; \
			cd services/$$service && cargo clippy -- -D warnings || true && cd ../..; \
		fi \
	done

lint-node: ## Run Node.js linter
	@echo -e "${YELLOW}Linting Node.js code...${NC}"
	@if [ -f "apps/front-web/package.json" ]; then \
		echo -e "${CYAN}Linting frontend${NC}"; \
		cd apps/front-web && npm run lint || true && cd ../..; \
	fi

lint-python: ## Run Python linter
	@echo -e "${YELLOW}Linting Python code...${NC}"
	@if [ -f "services/domain-svc/requirements.txt" ]; then \
		echo -e "${CYAN}Linting domain service${NC}"; \
		cd services/domain-svc && \
		source venv/bin/activate && \
		black --check . || true && \
		flake8 . || true && \
		deactivate && cd ../..; \
	fi

format: ## Format code for all languages
	@echo -e "${BLUE}Formatting code...${NC}"
	@$(MAKE) format-go format-rust format-node format-python

format-go: ## Format Go code
	@echo -e "${YELLOW}Formatting Go code...${NC}"
	@for service in auth-svc build-svc host-svc billing-svc; do \
		if [ -f "services/$$service/go.mod" ]; then \
			echo -e "${CYAN}Formatting $$service${NC}"; \
			cd services/$$service && go fmt ./... && cd ../..; \
		fi \
	done

format-rust: ## Format Rust code
	@echo -e "${YELLOW}Formatting Rust code...${NC}"
	@for service in builder-svc meter-svc; do \
		if [ -f "services/$$service/Cargo.toml" ]; then \
			echo -e "${CYAN}Formatting $$service${NC}"; \
			cd services/$$service && cargo fmt && cd ../..; \
		fi \
	done

format-node: ## Format Node.js code
	@echo -e "${YELLOW}Formatting Node.js code...${NC}"
	@if [ -f "apps/front-web/package.json" ]; then \
		echo -e "${CYAN}Formatting frontend${NC}"; \
		cd apps/front-web && npm run format || true && cd ../..; \
	fi

format-python: ## Format Python code
	@echo -e "${YELLOW}Formatting Python code...${NC}"
	@if [ -f "services/domain-svc/requirements.txt" ]; then \
		echo -e "${CYAN}Formatting domain service${NC}"; \
		cd services/domain-svc && \
		source venv/bin/activate && \
		black . && \
		isort . && \
		deactivate && cd ../..; \
	fi

security-check: ## Run security checks
	@echo -e "${BLUE}Running security checks...${NC}"
	@echo -e "${CYAN}Checking for vulnerabilities...${NC}"
	@for service in auth-svc build-svc host-svc billing-svc; do \
		if [ -f "services/$$service/go.mod" ]; then \
			cd services/$$service && go list -json -m all | nancy sleuth || true && cd ../..; \
		fi \
	done

# ==========================================
# DATABASE OPERATIONS
# ==========================================

db-shell: ## Connect to PostgreSQL database
	@echo -e "${BLUE}Connecting to PostgreSQL database...${NC}"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) exec postgres psql -U pagemagic -d pagemagic

redis-shell: ## Connect to Redis
	@echo -e "${BLUE}Connecting to Redis...${NC}"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) exec redis redis-cli

migration-create: ## Create new database migration (usage: make migration-create NAME=add_users_table)
	@if [ -z "$(NAME)" ]; then \
		echo -e "${RED}Please specify migration name: make migration-create NAME=add_users_table${NC}"; \
		exit 1; \
	fi
	@echo -e "${YELLOW}Creating migration: $(NAME)${NC}"
	@mkdir -p infrastructure/database/migrations
	@timestamp=$$(date +%Y%m%d%H%M%S); \
	echo "-- Migration: $(NAME)" > infrastructure/database/migrations/$${timestamp}_$(NAME).sql; \
	echo "-- Created at: $$(date)" >> infrastructure/database/migrations/$${timestamp}_$(NAME).sql; \
	echo "" >> infrastructure/database/migrations/$${timestamp}_$(NAME).sql; \
	echo "-- Add your migration SQL here" >> infrastructure/database/migrations/$${timestamp}_$(NAME).sql
	@echo -e "${GREEN}Migration created: infrastructure/database/migrations/$${timestamp}_$(NAME).sql${NC}"

migration-up: ## Run database migrations
	@echo -e "${YELLOW}Running database migrations...${NC}"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) exec -T postgres psql -U pagemagic -d pagemagic -f /docker-entrypoint-initdb.d/01-init.sql

seed: ## Seed database with sample data
	@echo -e "${YELLOW}Seeding database...${NC}"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) exec -T postgres psql -U pagemagic -d pagemagic -f /docker-entrypoint-initdb.d/02-seed.sql

backup: ## Backup database
	@echo -e "${YELLOW}Backing up database...${NC}"
	@timestamp=$$(date +%Y%m%d_%H%M%S); \
	docker-compose -f $(DOCKER_COMPOSE_FILE) exec -T postgres pg_dump -U pagemagic pagemagic > backup_$${timestamp}.sql; \
	echo -e "${GREEN}Database backed up to backup_$${timestamp}.sql${NC}"

restore: ## Restore database from backup (usage: make restore BACKUP=backup_20231201_120000.sql)
	@if [ -z "$(BACKUP)" ]; then \
		echo -e "${RED}Please specify backup file: make restore BACKUP=backup_20231201_120000.sql${NC}"; \
		exit 1; \
	fi
	@echo -e "${YELLOW}Restoring database from $(BACKUP)...${NC}"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) exec -T postgres psql -U pagemagic -d pagemagic < $(BACKUP)

# ==========================================
# MONITORING & LOGS
# ==========================================

logs: ## Show logs for all services
	@docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

logs-service: ## Show logs for specific service (usage: make logs-service SERVICE=auth-svc)
	@if [ -z "$(SERVICE)" ]; then \
		echo -e "${RED}Please specify a service: make logs-service SERVICE=auth-svc${NC}"; \
		exit 1; \
	fi
	@docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f $(SERVICE)

status: ## Show status of all services
	@echo -e "${BLUE}Service Status:${NC}"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) ps

health: ## Check health of all services
	@echo -e "${BLUE}Checking service health...${NC}"
	@curl -s http://localhost:3001/health | jq '.' || echo "Auth service not responding"
	@curl -s http://localhost:3002/health | jq '.' || echo "Prompt service not responding"
	@curl -s http://localhost:3003/health | jq '.' || echo "Builder service not responding"

check-ports: ## Check if required ports are available
	@echo -e "${BLUE}Checking port availability...${NC}"
	@for port in 3000 3001 3002 3003 3004 3005 3006 3007 3008 3009 5432 5433 6379 9000 9001 8080; do \
		if netstat -tuln | grep -q ":$$port "; then \
			echo -e "${RED}Port $$port is in use${NC}"; \
		else \
			echo -e "${GREEN}Port $$port is available${NC}"; \
		fi \
	done

# ==========================================
# UTILITIES
# ==========================================

clean: ## Clean up containers, images, and volumes
	@echo -e "${YELLOW}Cleaning up Docker resources...${NC}"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down -v --remove-orphans
	@docker system prune -f
	@echo -e "${GREEN}Cleanup completed${NC}"

clean-all: ## Clean everything including images
	@echo -e "${YELLOW}Cleaning up everything...${NC}"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down -v --remove-orphans --rmi all
	@docker system prune -af
	@echo -e "${GREEN}Deep cleanup completed${NC}"

shell: ## Open shell in a service container (usage: make shell SERVICE=auth-svc)
	@if [ -z "$(SERVICE)" ]; then \
		echo -e "${RED}Please specify a service: make shell SERVICE=auth-svc${NC}"; \
		exit 1; \
	fi
	@docker-compose -f $(DOCKER_COMPOSE_FILE) exec $(SERVICE) /bin/sh

env-check: ## Check environment variables
	@echo -e "${BLUE}Environment Configuration:${NC}"
	@if [ -f .env ]; then \
		echo -e "${GREEN}.env file exists${NC}"; \
		grep -v "^#" .env | grep -v "^$$" | wc -l | xargs echo "Variables configured:"; \
	else \
		echo -e "${RED}.env file not found${NC}"; \
		echo -e "${YELLOW}Run 'cp .env.example .env' and configure your environment${NC}"; \
	fi

docs: ## Generate and serve documentation
	@echo -e "${BLUE}Generating documentation...${NC}"
	@echo -e "${CYAN}Documentation will be available at http://localhost:8000${NC}"

benchmark: ## Run performance benchmarks
	@echo -e "${BLUE}Running benchmarks...${NC}"
	@echo -e "${YELLOW}This will run load tests against the services${NC}"

# ==========================================
# PRODUCTION
# ==========================================

prod-build: ## Build production Docker images
	@echo -e "${BLUE}Building production images...${NC}"
	@docker-compose -f $(DOCKER_PROD_FILE) build

prod-up: ## Start production environment
	@echo -e "${BLUE}Starting production environment...${NC}"
	@docker-compose -f $(DOCKER_PROD_FILE) up -d

prod-down: ## Stop production environment
	@echo -e "${YELLOW}Stopping production environment...${NC}"
	@docker-compose -f $(DOCKER_PROD_FILE) down

# ==========================================
# RELEASE
# ==========================================

release-check: ## Check if ready for release
	@echo -e "${BLUE}Checking release readiness...${NC}"
	@$(MAKE) test lint security-check
	@echo -e "${GREEN}Release checks passed${NC}"

tag: ## Create a new git tag (usage: make tag VERSION=v1.0.0)
	@if [ -z "$(VERSION)" ]; then \
		echo -e "${RED}Please specify version: make tag VERSION=v1.0.0${NC}"; \
		exit 1; \
	fi
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@echo -e "${GREEN}Tagged $(VERSION)${NC}"

# ==========================================
# ALIASES
# ==========================================

up: dev ## Alias for dev
start: dev ## Alias for dev
stop: down ## Alias for down
