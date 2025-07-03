# Page Magic 🪄

### Plataforma No-Code de Landing Pages, Blogs e Sites guiada por IA

Page Magic é uma plataforma completa que permite aos usuários criar sites profissionais através de descrições em linguagem natural. A IA gera o site completo, o usuário edita visualmente e publica com um clique.

## 🎯 Status do Projeto

**�️ Progresso Geral:** 45% implementado | 55% pendente
**⚡ Arquitetura:** ✅ 100% completa
**🔧 Funcionalidades Core:** ⚠️ 35% funcionais

### 📋 Documentação de Status
- 📄 [**Análise Detalhada de Pendências**](MISSING_IMPLEMENTATIONS.md) - Análise completa do que falta
- 🎯 [**Próximos Passos Prioritários**](NEXT_STEPS.md) - Roadmap executivo focado no MVP
- 📊 [**Relatório de Progresso Completo**](PROGRESS_REPORT.md) - Status detalhado de cada serviço
- 📈 [**Status de Implementação**](STATUS.md) - Overview de todos os componentes

## 🚀 Início Rápido

```bash
# Clonar o repositório
git clone https://github.com/eltonciatto/pagemagic.git
cd pagemagic

# Configurar ambiente
cp .env.example .env
# Editar .env com suas configurações

# Instalar dependências
npm install

# Inicializar infraestrutura
make infra-up

# Build de todos os serviços
make build-all

# Iniciar em modo desenvolvimento
make dev
```

## 🎯 Status de Implementação

### ✅ **COMPLETO** - Arquitetura & Estrutura
- [x] **Arquitetura de microserviços** - 11 serviços estruturados
- [x] **Infraestrutura base** - Docker, PostgreSQL, Redis, NATS
- [x] **Makefile completo** - Automação de build/dev/test/deploy
- [x] **Schemas e tipos** - TypeScript types, validação, utilitários
- [x] **Documentação técnica** - Specs completas e guias
- [x] **CI/CD Pipeline** - GitHub Actions workflow completo
- [x] **Configurações** - Nginx, Prometheus, Docker Compose

### ✅ **IMPLEMENTADO** - Serviços Base (Estrutura)
- [x] **auth-svc** (Go) - Autenticação JWT, magic links, repositórios
- [x] **prompt-svc** (Node.js) - Base para integração IA
- [x] **builder-svc** (Rust) - Base para build system
- [x] **billing-svc** (Go) - Estrutura Stripe integration
- [x] **domain-svc** (Python) - Base DNS/SSL management
- [x] **front-web** (Next.js) - Interface base e componentes

### ⚠️ **EM PROGRESSO** - Funcionalidades Core
- [ ] **Integração IA real** - vLLM cluster + streaming
- [ ] **Pipeline de build** - JSON → AST → HTML/React
- [ ] **Orquestração containers** - Docker Swarm + deploy
- [ ] **Proxy de uso** - Nginx + Lua + medição
- [ ] **Editor visual** - WYSIWYG + preview real-time
- [ ] **App mobile** - React Native navigation + features

### ❌ **FALTANDO** - Infraestrutura Produção
- [ ] **Observabilidade completa** - Grafana dashboards + alerting
- [ ] **Testes** - Unit + integration + E2E (0% coverage)
- [ ] **Security hardening** - CSP + WAF + compliance
- [ ] **Performance optimization** - Caching + CDN setup

## 🏗️ Arquitetura

### Front-ends
- **front-web**: Next.js 15.3 + React 19 (Desktop/PWA)
- **mobile-app**: React Native 0.74 + Expo SDK 51 (iOS/Android)

### Microserviços
- **auth-svc** (Go): Autenticação magic-link + OAuth
- **prompt-svc** (Node.js): Orquestração de prompts IA (LangChain)
- **builder-svc** (Rust): Conversão JSON → AST → HTML/React
- **build-svc** (Go): Geração de imagens Docker + Turbopack
- **host-svc** (Go): Gerenciamento de contêineres (Docker Swarm)
- **domain-svc** (Python): Compra de domínios + DNS/ACME
- **usage-proxy** (Nginx + Lua): Coleta de métricas de uso
- **meter-svc** (Rust): Agregação de eventos + Stripe Meters
- **billing-svc** (Go): Webhooks Stripe + controle de créditos
- **i18n-svc** (Node.js): Gerenciamento de traduções

### Infraestrutura
- **Database**: PostgreSQL 16 + TimescaleDB
- **Message Broker**: NATS JetStream
- **Observability**: Prometheus + Grafana + Loki
- **Deploy**: Docker Swarm + GitOps (Flux)
- **Edge**: Cloudflare + Workers
- **AI**: vLLM cluster (Llama-3 70B) + OpenAI/Claude adapters

## 🚀 Início Rápido

```bash
# Clonar o repositório
git clone https://github.com/eltonciatto/pagemagic.git
cd pagemagic

# Configurar ambiente
cp .env.example .env
# Editar .env com suas configurações

# Inicializar infraestrutura
make infra-up

# Build de todos os serviços
make build-all

# Iniciar todos os serviços
make start-all
```

## 📁 Estrutura do Projeto

```
pagemagic/
├── apps/
│   ├── front-web/          # Next.js web app
│   └── mobile-app/         # React Native app
├── services/
│   ├── auth-svc/           # Serviço de autenticação
│   ├── prompt-svc/         # Orquestração IA
│   ├── builder-svc/        # Construção de sites
│   ├── build-svc/          # Build Docker
│   ├── host-svc/           # Hosting
│   ├── domain-svc/         # Domínios
│   ├── usage-proxy/        # Proxy de métricas
│   ├── meter-svc/          # Medição
│   ├── billing-svc/        # Cobrança
│   └── i18n-svc/          # Internacionalização
├── infrastructure/
│   ├── database/           # Schemas SQL
│   ├── docker/            # Docker configs
│   ├── k8s/               # Kubernetes manifests
│   └── observability/     # Monitoring configs
└── shared/
    ├── types/             # TypeScript types
    ├── schemas/           # JSON schemas
    └── utils/             # Utilities comuns
```

## 🛠️ Comandos de Desenvolvimento

```bash
# Desenvolvimento
make dev-web              # Inicia front-web em modo dev
make dev-mobile           # Inicia mobile-app com Expo
make dev-service SVC=auth # Inicia serviço específico

# Build
make build SERVICE=auth   # Build de serviço específico
make build-all           # Build de todos os serviços

# Testes
make test                # Executa todos os testes
make test-e2e           # Testes end-to-end
make test-mobile        # Testes mobile (Detox)

# Deploy
make deploy-staging     # Deploy para staging
make deploy-prod       # Deploy para produção
```

## 🔧 Configuração

### Variáveis de Ambiente

Copie `.env.example` para `.env` e configure:

- Database URLs (PostgreSQL, TimescaleDB)
- API keys (OpenAI, Stripe, Cloudflare)
- JWT secrets
- Message broker URLs (NATS)

### IA e Modelos

Configure os provedores de IA em `services/prompt-svc/config/`:
- vLLM local (Llama-3 70B)
- OpenAI API
- Anthropic Claude

## 📊 Monitoramento

- **Grafana**: http://localhost:3000
- **Prometheus**: http://localhost:9090
- **Loki**: http://localhost:3100

## 🔒 Segurança

- CSP rigorosa (zero inline JS)
- OWASP Top-10 compliance
- PCI-DSS (pagamentos via Stripe)
- LGPD/GDPR (criptografia PII)
- TLS 1.3 obrigatório

## 📱 Mobile

O app móvel inclui:
- Login/cadastro
- Geração via IA
- Editor visual (WebView + bridge)
- Preview em tempo real
- Push notifications
- Pagamentos in-app (Stripe)
- Deep linking
- Modo offline

## 💰 Modelo de Negócio

Medição via Stripe Meters:
- **Gerações**: Contagem de sites gerados
- **Tokens IA**: Tokens consumidos
- **Container Hours**: Horas de hosting
- **Storage**: GB de armazenamento

## 🧪 Testes

- **Unit**: Jest + Vitest
- **Integration**: Supertest
- **E2E Web**: Playwright
- **E2E Mobile**: Detox
- **Load**: K6

## 📝 Licença

MIT License - veja [LICENSE](LICENSE) para detalhes.

## 🤝 Contribuindo

1. Fork o projeto
2. Crie sua feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📞 Suporte

- 📧 Email: suporte@pagemagic.io
- 💬 Discord: [Page Magic Community](https://discord.gg/pagemagic)
- 📚 Docs: [docs.pagemagic.io](https://docs.pagemagic.io)
