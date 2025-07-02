# Page Magic - Makefile
# Automação de tarefas para desenvolvimento, build e deploy

.PHONY: help dev build test deploy clean

# Variáveis
DOCKER_REGISTRY ?= pagemagic
VERSION ?= latest
ENV ?= development

# Cores para output
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
RESET := \033[0m

help: ## Mostra esta ajuda
	@echo "$(BLUE)Page Magic - Comandos Disponíveis$(RESET)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(RESET) %s\n", $$1, $$2}'

# ==========================================
# SETUP E CONFIGURAÇÃO
# ==========================================

setup: ## Configuração inicial do projeto
	@echo "$(BLUE)Configurando Page Magic...$(RESET)"
	cp .env.example .env
	@echo "$(YELLOW)Edite o arquivo .env com suas configurações$(RESET)"
	make deps-install
	make infra-up

deps-install: ## Instala dependências de todos os projetos
	@echo "$(BLUE)Instalando dependências...$(RESET)"
	cd apps/front-web && npm install
	cd apps/mobile-app && npm install
	cd services/prompt-svc && npm install
	cd services/i18n-svc && npm install

# ==========================================
# INFRAESTRUTURA
# ==========================================

infra-up: ## Sobe a infraestrutura (Docker Compose)
	@echo "$(BLUE)Iniciando infraestrutura...$(RESET)"
	docker-compose -f infrastructure/docker/docker-compose.yml up -d
	@echo "$(GREEN)Infraestrutura iniciada!$(RESET)"

infra-down: ## Para a infraestrutura
	@echo "$(BLUE)Parando infraestrutura...$(RESET)"
	docker-compose -f infrastructure/docker/docker-compose.yml down

infra-logs: ## Mostra logs da infraestrutura
	docker-compose -f infrastructure/docker/docker-compose.yml logs -f

# ==========================================
# DESENVOLVIMENTO
# ==========================================

dev: ## Inicia todos os serviços em modo desenvolvimento
	@echo "$(BLUE)Iniciando ambiente de desenvolvimento...$(RESET)"
	make -j4 dev-web dev-services dev-proxy

dev-web: ## Inicia front-web em modo desenvolvimento
	@echo "$(BLUE)Iniciando front-web...$(RESET)"
	cd apps/front-web && npm run dev

dev-mobile: ## Inicia mobile-app com Expo
	@echo "$(BLUE)Iniciando mobile-app...$(RESET)"
	cd apps/mobile-app && npx expo start

dev-services: ## Inicia todos os serviços backend
	make -j8 dev-auth dev-prompt dev-builder dev-build dev-host dev-domain dev-meter dev-billing dev-i18n

dev-auth: ## Inicia auth-svc em modo desenvolvimento
	cd services/auth-svc && go run main.go

dev-prompt: ## Inicia prompt-svc em modo desenvolvimento
	cd services/prompt-svc && npm run dev

dev-builder: ## Inicia builder-svc em modo desenvolvimento
	cd services/builder-svc && cargo run

dev-build: ## Inicia build-svc em modo desenvolvimento
	cd services/build-svc && go run main.go

dev-host: ## Inicia host-svc em modo desenvolvimento
	cd services/host-svc && go run main.go

dev-domain: ## Inicia domain-svc em modo desenvolvimento
	cd services/domain-svc && python -m uvicorn main:app --reload

dev-meter: ## Inicia meter-svc em modo desenvolvimento
	cd services/meter-svc && cargo run

dev-billing: ## Inicia billing-svc em modo desenvolvimento
	cd services/billing-svc && go run main.go

dev-i18n: ## Inicia i18n-svc em modo desenvolvimento
	cd services/i18n-svc && npm run dev

dev-proxy: ## Inicia usage-proxy
	cd services/usage-proxy && nginx -c $(PWD)/services/usage-proxy/nginx.conf

# ==========================================
# BUILD E IMAGENS DOCKER
# ==========================================

build: ## Build de todos os serviços
ifdef SERVICE
	@echo "$(BLUE)Building $(SERVICE)...$(RESET)"
	docker build -t $(DOCKER_REGISTRY)/$(SERVICE):$(VERSION) services/$(SERVICE)
else
	@echo "$(BLUE)Building todos os serviços...$(RESET)"
	make build-services build-apps
endif

build-all: build-services build-apps ## Build completo (serviços + apps)

build-services: ## Build de todos os microserviços
	@echo "$(BLUE)Building microserviços...$(RESET)"
	$(MAKE) -j8 build-auth build-prompt build-builder build-build build-host build-domain build-meter build-billing build-i18n build-proxy

build-apps: ## Build de front-web e mobile-app
	@echo "$(BLUE)Building applications...$(RESET)"
	cd apps/front-web && npm run build
	cd apps/mobile-app && npx expo export

build-auth: ## Build auth-svc
	docker build -t $(DOCKER_REGISTRY)/auth-svc:$(VERSION) services/auth-svc

build-prompt: ## Build prompt-svc
	docker build -t $(DOCKER_REGISTRY)/prompt-svc:$(VERSION) services/prompt-svc

build-builder: ## Build builder-svc
	docker build -t $(DOCKER_REGISTRY)/builder-svc:$(VERSION) services/builder-svc

build-build: ## Build build-svc
	docker build -t $(DOCKER_REGISTRY)/build-svc:$(VERSION) services/build-svc

build-host: ## Build host-svc
	docker build -t $(DOCKER_REGISTRY)/host-svc:$(VERSION) services/host-svc

build-domain: ## Build domain-svc
	docker build -t $(DOCKER_REGISTRY)/domain-svc:$(VERSION) services/domain-svc

build-meter: ## Build meter-svc
	docker build -t $(DOCKER_REGISTRY)/meter-svc:$(VERSION) services/meter-svc

build-billing: ## Build billing-svc
	docker build -t $(DOCKER_REGISTRY)/billing-svc:$(VERSION) services/billing-svc

build-i18n: ## Build i18n-svc
	docker build -t $(DOCKER_REGISTRY)/i18n-svc:$(VERSION) services/i18n-svc

build-proxy: ## Build usage-proxy
	docker build -t $(DOCKER_REGISTRY)/usage-proxy:$(VERSION) services/usage-proxy

# ==========================================
# TESTES
# ==========================================

test: ## Executa todos os testes
	@echo "$(BLUE)Executando todos os testes...$(RESET)"
	make test-unit test-integration test-e2e

test-unit: ## Testes unitários
	@echo "$(BLUE)Executando testes unitários...$(RESET)"
	cd apps/front-web && npm run test
	cd apps/mobile-app && npm run test
	cd services/prompt-svc && npm test
	cd services/i18n-svc && npm test
	cd services/auth-svc && go test ./...
	cd services/build-svc && go test ./...
	cd services/host-svc && go test ./...
	cd services/billing-svc && go test ./...
	cd services/builder-svc && cargo test
	cd services/meter-svc && cargo test

test-integration: ## Testes de integração
	@echo "$(BLUE)Executando testes de integração...$(RESET)"
	docker-compose -f infrastructure/docker/docker-compose.test.yml up --abort-on-container-exit

test-e2e: ## Testes end-to-end
	@echo "$(BLUE)Executando testes E2E...$(RESET)"
	cd apps/front-web && npx playwright test

test-mobile: ## Testes mobile (Detox)
	@echo "$(BLUE)Executando testes mobile...$(RESET)"
	cd apps/mobile-app && npx detox test

test-load: ## Testes de carga (K6)
	@echo "$(BLUE)Executando testes de carga...$(RESET)"
	k6 run tests/load/script.js

# ==========================================
# DEPLOY
# ==========================================

deploy-staging: ## Deploy para staging
	@echo "$(BLUE)Deploy para staging...$(RESET)"
	kubectl apply -f infrastructure/k8s/staging/
	kubectl rollout status deployment -n pagemagic-staging

deploy-prod: ## Deploy para produção
	@echo "$(BLUE)Deploy para produção...$(RESET)"
	kubectl apply -f infrastructure/k8s/production/
	kubectl rollout status deployment -n pagemagic-prod

push: ## Push images para registry
ifdef SERVICE
	docker push $(DOCKER_REGISTRY)/$(SERVICE):$(VERSION)
else
	@echo "$(BLUE)Pushing todas as images...$(RESET)"
	docker push $(DOCKER_REGISTRY)/auth-svc:$(VERSION)
	docker push $(DOCKER_REGISTRY)/prompt-svc:$(VERSION)
	docker push $(DOCKER_REGISTRY)/builder-svc:$(VERSION)
	docker push $(DOCKER_REGISTRY)/build-svc:$(VERSION)
	docker push $(DOCKER_REGISTRY)/host-svc:$(VERSION)
	docker push $(DOCKER_REGISTRY)/domain-svc:$(VERSION)
	docker push $(DOCKER_REGISTRY)/meter-svc:$(VERSION)
	docker push $(DOCKER_REGISTRY)/billing-svc:$(VERSION)
	docker push $(DOCKER_REGISTRY)/i18n-svc:$(VERSION)
	docker push $(DOCKER_REGISTRY)/usage-proxy:$(VERSION)
endif

# ==========================================
# UTILITÁRIOS
# ==========================================

logs: ## Mostra logs de um serviço específico
ifdef SERVICE
	docker logs -f pagemagic_$(SERVICE)
else
	@echo "$(RED)Use: make logs SERVICE=<service-name>$(RESET)"
endif

shell: ## Acessa shell de um serviço
ifdef SERVICE
	docker exec -it pagemagic_$(SERVICE) /bin/sh
else
	@echo "$(RED)Use: make shell SERVICE=<service-name>$(RESET)"
endif

clean: ## Limpa containers, images e volumes
	@echo "$(BLUE)Limpando ambiente...$(RESET)"
	docker system prune -f
	docker volume prune -f

clean-all: ## Limpeza completa (cuidado!)
	@echo "$(RED)Limpeza completa...$(RESET)"
	docker system prune -af
	docker volume prune -f

db-migrate: ## Executa migrações do banco
	@echo "$(BLUE)Executando migrações...$(RESET)"
	docker exec pagemagic_postgres psql -U pagemagic -d pagemagic -f /docker-entrypoint-initdb.d/migrations.sql

db-seed: ## Popula banco com dados de teste
	@echo "$(BLUE)Populando banco...$(RESET)"
	docker exec pagemagic_postgres psql -U pagemagic -d pagemagic -f /docker-entrypoint-initdb.d/seed.sql

db-reset: ## Reset completo do banco
	@echo "$(RED)Resetando banco...$(RESET)"
	make infra-down
	docker volume rm pagemagic_postgres_data
	make infra-up
	sleep 10
	make db-migrate
	make db-seed

monitor: ## Abre dashboards de monitoramento
	@echo "$(BLUE)Abrindo dashboards...$(RESET)"
	@echo "Grafana: http://localhost:3000"
	@echo "Prometheus: http://localhost:9090"
	@echo "Loki: http://localhost:3100"

# ==========================================
# MOBILE ESPECÍFICO
# ==========================================

mobile-build-ios: ## Build iOS
	cd apps/mobile-app && npx eas build --platform ios

mobile-build-android: ## Build Android
	cd apps/mobile-app && npx eas build --platform android

mobile-submit-ios: ## Submit para App Store
	cd apps/mobile-app && npx eas submit --platform ios

mobile-submit-android: ## Submit para Google Play
	cd apps/mobile-app && npx eas submit --platform android

mobile-update: ## OTA Update
	cd apps/mobile-app && npx eas update

# ==========================================
# DESENVOLVIMENTO AVANÇADO
# ==========================================

lint: ## Executa linting em todos os projetos
	@echo "$(BLUE)Executando linting...$(RESET)"
	cd apps/front-web && npm run lint
	cd apps/mobile-app && npm run lint
	cd services/prompt-svc && npm run lint
	cd services/i18n-svc && npm run lint

format: ## Formata código
	@echo "$(BLUE)Formatando código...$(RESET)"
	cd apps/front-web && npm run format
	cd apps/mobile-app && npm run format
	cd services/prompt-svc && npm run format
	cd services/i18n-svc && npm run format
	cd services/builder-svc && cargo fmt
	cd services/meter-svc && cargo fmt

security-scan: ## Scan de segurança
	@echo "$(BLUE)Executando scan de segurança...$(RESET)"
	trivy fs .

docs: ## Gera documentação
	@echo "$(BLUE)Gerando documentação...$(RESET)"
	cd docs && npm run build

default: help
