# Documentação Técnica - Page Magic

## Arquitetura Geral

Page Magic é uma plataforma no-code para criação de sites guiada por IA, construída com arquitetura de microserviços.

### Stack Tecnológico

- **Frontend Web**: Next.js 15.3 + React 19 + TypeScript + Tailwind CSS
- **Mobile**: React Native 0.74 + Expo SDK 51
- **Backend**: 
  - Go (auth-svc, build-svc, host-svc, billing-svc)
  - Node.js + TypeScript (prompt-svc, i18n-svc)
  - Rust (builder-svc, meter-svc)
  - Python (domain-svc)
  - Nginx + Lua (usage-proxy)
- **Database**: PostgreSQL 16 + TimescaleDB
- **Cache**: Redis
- **Message Broker**: NATS JetStream
- **Observability**: Prometheus + Grafana + Loki
- **Containerização**: Docker + Docker Compose
- **Orquestração**: Docker Swarm (desenvolvimento) / Kubernetes (produção)

## Estrutura de Diretórios

```
pagemagic/
├── apps/                          # Aplicações front-end
│   ├── front-web/                 # Next.js web app
│   └── mobile-app/                # React Native app
├── services/                      # Microserviços
│   ├── auth-svc/                  # Autenticação (Go)
│   ├── prompt-svc/                # Orquestração IA (Node.js)
│   ├── builder-svc/               # Construção sites (Rust)
│   ├── build-svc/                 # Build Docker (Go)
│   ├── host-svc/                  # Hosting (Go)
│   ├── domain-svc/                # Domínios (Python)
│   ├── usage-proxy/               # Proxy métricas (Nginx+Lua)
│   ├── meter-svc/                 # Medição (Rust)
│   ├── billing-svc/               # Cobrança (Go)
│   └── i18n-svc/                  # i18n (Node.js)
├── infrastructure/                # Infraestrutura
│   ├── database/                  # Schemas SQL
│   ├── docker/                    # Docker configs
│   ├── k8s/                       # Kubernetes manifests
│   └── observability/             # Monitoring configs
├── shared/                        # Código compartilhado
│   ├── types/                     # TypeScript types
│   ├── schemas/                   # JSON schemas
│   └── utils/                     # Utilities
└── scripts/                       # Scripts de automação
```

## Serviços Implementados

### 1. auth-svc (Go)

**Responsabilidades:**
- Autenticação via magic link
- Gerenciamento de usuários
- JWT tokens (access + refresh)
- OAuth (Google, GitHub, Apple)

**Endpoints principais:**
- `POST /api/v1/auth/magic-link` - Enviar magic link
- `POST /api/v1/auth/verify` - Verificar magic link
- `POST /api/v1/auth/refresh` - Renovar token
- `GET /api/v1/profile` - Obter perfil do usuário

**Tecnologias:**
- Gin (HTTP framework)
- PostgreSQL (database)
- JWT-go (tokens)
- bcrypt (hashing)

### 2. prompt-svc (Node.js + TypeScript)

**Responsabilidades:**
- Orquestração de prompts para IA
- Integração com múltiplos provedores (OpenAI, Anthropic, vLLM)
- Cache de respostas
- Rate limiting

**Endpoints principais:**
- `POST /api/v1/prompts/generate` - Gerar site via IA
- `POST /api/v1/prompts/improve` - Melhorar conteúdo
- `GET /api/v1/prompts/templates` - Templates de prompts

**Tecnologias:**
- Express.js
- LangChain
- Redis (cache)
- NATS (messaging)

### 3. builder-svc (Rust)

**Responsabilidades:**
- Conversão JSON → AST → HTML/React
- Geração de código otimizado
- Processamento de templates
- Validação de estruturas

**Endpoints principais:**
- `POST /api/v1/build/website` - Build completo do site
- `POST /api/v1/build/component` - Build de componente
- `GET /api/v1/build/status/:id` - Status do build

**Tecnologias:**
- Axum (HTTP framework)
- Serde (serialização)
- Handlebars/Tera (templates)
- SWC (transpilação)

## Configuração de Desenvolvimento

### Pré-requisitos

- Docker & Docker Compose
- Go 1.22+
- Node.js 18+
- Rust 1.75+
- Python 3.11+

### Início Rápido

```bash
# 1. Clonar repositório
git clone https://github.com/pagemagic/pagemagic.git
cd pagemagic

# 2. Executar script de setup
./scripts/setup.sh

# 3. Iniciar serviços
make web-dev        # Frontend
make auth-dev       # Auth service
make prompt-dev     # Prompt service
make builder-dev    # Builder service
```

### Comandos de Desenvolvimento

```bash
# Infraestrutura
make infra-up          # Iniciar infraestrutura
make infra-down        # Parar infraestrutura

# Desenvolvimento
make web-dev           # Frontend Next.js
make auth-dev          # Auth service
make prompt-dev        # Prompt service  
make builder-dev       # Builder service
make dev-backend       # Todos os serviços backend

# Build
make build-services    # Build todos os serviços
make auth-build        # Build auth-svc
make prompt-build      # Build prompt-svc
make builder-build     # Build builder-svc

# Testes
make test-services     # Testes todos os serviços
make auth-test         # Testes auth-svc
make prompt-test       # Testes prompt-svc
make builder-test      # Testes builder-svc
```

## Padrões de Desenvolvimento

### Estrutura de Serviço (Go)

```
service-name/
├── main.go                 # Entry point
├── go.mod                  # Dependências
├── Dockerfile             # Container config
├── internal/              # Código interno
│   ├── app/               # Aplicação principal
│   ├── config/            # Configurações
│   ├── handlers/          # HTTP handlers
│   ├── models/            # Modelos de dados
│   ├── repository/        # Camada de dados
│   └── services/          # Lógica de negócio
└── cmd/                   # Comandos CLI (opcional)
```

### Estrutura de Serviço (Node.js)

```
service-name/
├── package.json           # Dependências
├── tsconfig.json          # TypeScript config
├── Dockerfile             # Container config
├── src/
│   ├── index.ts           # Entry point
│   ├── config/            # Configurações
│   ├── routes/            # Rotas HTTP
│   ├── controllers/       # Controllers
│   ├── services/          # Lógica de negócio
│   ├── models/            # Modelos de dados
│   ├── middleware/        # Middlewares
│   └── utils/             # Utilitários
└── tests/                 # Testes
```

### Estrutura de Serviço (Rust)

```
service-name/
├── Cargo.toml             # Dependências
├── Dockerfile             # Container config
├── src/
│   ├── main.rs            # Entry point
│   ├── config.rs          # Configurações
│   ├── handlers.rs        # HTTP handlers
│   ├── models.rs          # Modelos de dados
│   ├── services.rs        # Lógica de negócio
│   └── utils.rs           # Utilitários
└── tests/                 # Testes
```

## Comunicação Entre Serviços

### HTTP/REST

- Comunicação síncrona para operações críticas
- Autenticação via JWT
- Rate limiting por serviço
- Circuit breaker para resiliência

### NATS JetStream

- Comunicação assíncrona para eventos
- Streams persistentes
- At-least-once delivery
- Dead letter queues

### Exemplo de Event

```json
{
  "type": "site.generated",
  "version": "1.0",
  "timestamp": "2024-07-02T10:00:00Z",
  "data": {
    "site_id": "uuid",
    "user_id": "uuid",
    "generation_time": 1500,
    "tokens_used": 2500
  }
}
```

## Banco de Dados

### Schema Principal (PostgreSQL)

```sql
-- Usuários
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    email_verified BOOLEAN DEFAULT FALSE,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    avatar_url VARCHAR(255),
    status user_status DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Projetos/Sites
CREATE TABLE projects (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    domain VARCHAR(255),
    status project_status DEFAULT 'draft',
    site_data JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Magic Links (auth)
CREATE TABLE magic_links (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

### TimescaleDB (Métricas)

```sql
-- Eventos de uso
CREATE TABLE usage_events (
    time TIMESTAMPTZ NOT NULL,
    user_id UUID,
    event_type VARCHAR(50),
    resource_id UUID,
    quantity DOUBLE PRECISION,
    metadata JSONB
);

SELECT create_hypertable('usage_events', 'time');
```

## Observabilidade

### Métricas (Prometheus)

```go
// Exemplo de métricas em Go
var (
    httpRequests = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
        },
        []string{"method", "endpoint"},
    )
)
```

### Logs (Estruturados)

```json
{
  "timestamp": "2024-07-02T10:00:00Z",
  "level": "info",
  "service": "auth-svc",
  "trace_id": "abc123",
  "user_id": "user123",
  "message": "User authenticated successfully",
  "duration": 150
}
```

### Health Checks

Todos os serviços expõem:
- `GET /health` - Health check básico
- `GET /metrics` - Métricas Prometheus
- `GET /ready` - Readiness probe

## Segurança

### Autenticação/Autorização

- JWT tokens com RS256
- Refresh tokens com rotação
- Magic links com TTL curto (15min)
- Rate limiting por usuário/IP

### Comunicação

- TLS 1.3 obrigatório em produção
- mTLS entre serviços internos
- API Gateway com autenticação centralizada

### Dados Sensíveis

- Criptografia PII com AES-256-GCM
- Hashing senhas com bcrypt (custo 12)
- Rotação automática de chaves
- Audit logs para operações críticas

## Deploy e CI/CD

### Desenvolvimento

- Docker Compose para infraestrutura local
- Hot reload em todos os serviços
- Volumes para persistência de dados

### Staging/Produção

- Kubernetes com Helm charts
- GitOps com FluxCD
- Blue-green deployments
- Rollback automático em falhas

### Pipeline CI/CD

```yaml
# .github/workflows/ci.yml (exemplo)
name: CI/CD Pipeline

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run tests
        run: make test-services
      
  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Build services
        run: make build-services
      - name: Push images
        run: make push-images
        
  deploy:
    needs: build
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to staging
        run: make deploy-staging
```

## Monitoramento e Alertas

### Dashboards Grafana

- **Overview**: Métricas gerais do sistema
- **Services**: Métricas por serviço
- **Business**: Métricas de negócio (gerações, usuários)
- **Infrastructure**: Recursos computacionais

### Alertas

- Latência > 500ms (95th percentile)
- Error rate > 1%
- CPU/Memory > 80%
- Disk space < 10%
- Falhas de deploy

## Contribuição

### Fluxo de Desenvolvimento

1. Fork do repositório
2. Criar branch feature (`git checkout -b feature/amazing-feature`)
3. Implementar mudanças com testes
4. Commit seguindo conventional commits
5. Push e criar Pull Request
6. Code review e aprovação
7. Merge para main

### Conventional Commits

```
feat: add new endpoint for site generation
fix: resolve memory leak in builder service
docs: update API documentation
test: add integration tests for auth flow
chore: update dependencies
```

### Testes

- **Unit**: Cobertura mínima 80%
- **Integration**: Principais fluxos
- **E2E**: Cenários críticos do usuário
- **Load**: Performance benchmarks

## Troubleshooting

### Problemas Comuns

**Serviços não iniciam:**
```bash
# Verificar logs
docker-compose logs service-name

# Verificar portas
netstat -tulpn | grep :8080

# Resetar infraestrutura
make infra-down && make infra-up
```

**Banco de dados inacessível:**
```bash
# Verificar conexão
docker-compose exec postgres pg_isready

# Conectar manualmente
docker-compose exec postgres psql -U pagemagic

# Verificar logs
docker-compose logs postgres
```

**Build failures:**
```bash
# Limpar caches
make clean

# Rebuild completo
make build-services

# Verificar dependências
make deps-install
```

---

Para mais informações, consulte:
- [README.md](./README.md) - Visão geral do projeto
- [API Documentation](./docs/api/) - Documentação das APIs
- [Deployment Guide](./docs/deployment/) - Guia de deploy
