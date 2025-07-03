# Contribuindo para o Page Magic 🪄

Obrigado por seu interesse em contribuir para o Page Magic! Este documento contém diretrizes para contribuir com o projeto.

## 📋 Índice

- [Como Contribuir](#como-contribuir)
- [Configuração do Ambiente](#configuração-do-ambiente)
- [Padrões de Código](#padrões-de-código)
- [Processo de Pull Request](#processo-de-pull-request)
- [Reportando Bugs](#reportando-bugs)
- [Sugerindo Features](#sugerindo-features)
- [Arquitetura do Projeto](#arquitetura-do-projeto)

## 🚀 Como Contribuir

### 1. Fork o Repositório

```bash
git clone https://github.com/[seu-username]/pagemagic.git
cd pagemagic
git remote add upstream https://github.com/eltonciatto/pagemagic.git
```

### 2. Configure o Ambiente

```bash
# Execute o script de setup
./scripts/setup.sh

# Configure suas variáveis de ambiente
cp .env.example .env
# Edite .env conforme necessário
```

### 3. Crie uma Branch

```bash
git checkout -b feature/sua-feature
# ou
git checkout -b fix/correção-bug
```

### 4. Faça suas Alterações

Siga os [padrões de código](#padrões-de-código) estabelecidos.

### 5. Teste suas Alterações

```bash
# Testes unitários
npm test

# Testes de integração
make test-integration

# Testes E2E
make test-e2e

# Lint
npm run lint

# Type check
npm run type-check
```

### 6. Commit e Push

```bash
git commit -m "feat: adiciona nova funcionalidade X"
git push origin feature/sua-feature
```

### 7. Abra um Pull Request

Abra um PR para a branch `develop` com uma descrição clara das alterações.

## 🛠️ Configuração do Ambiente

### Pré-requisitos

- **Node.js** 22+
- **Go** 1.23+
- **Rust** 1.75+
- **Python** 3.12+
- **Docker** & **Docker Compose**

### Instalação

```bash
# 1. Clone e configure
git clone https://github.com/eltonciatto/pagemagic.git
cd pagemagic

# 2. Execute o setup automatizado
./scripts/setup.sh

# 3. Inicie os serviços de desenvolvimento
make dev-backend  # Backend services
make dev-web      # Frontend web
make dev-mobile   # Mobile app (opcional)
```

### Variáveis de Ambiente

Configure as seguintes variáveis em `.env`:

```bash
# Database
DATABASE_URL=postgresql://pagemagic:pagemagic123@localhost:5432/pagemagic

# IA APIs
OPENAI_API_KEY=sk-your-key-here
ANTHROPIC_API_KEY=your-key-here

# Stripe
STRIPE_SECRET_KEY=sk_test_your-key
STRIPE_WEBHOOK_SECRET=whsec_your-secret

# Outros serviços
REDIS_URL=redis://localhost:6379
NATS_URL=nats://localhost:4222
```

## 📝 Padrões de Código

### Convenções de Commit

Usamos [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Tipos:**
- `feat`: Nova funcionalidade
- `fix`: Correção de bug
- `docs`: Alterações na documentação
- `style`: Formatação (não afeta o código)
- `refactor`: Refatoração
- `test`: Adição/alteração de testes
- `chore`: Tarefas de manutenção

**Exemplos:**
```
feat(auth): adiciona autenticação via magic link
fix(builder): corrige geração de CSS inválido
docs: atualiza README com instruções de setup
test(prompt): adiciona testes para geração IA
```

### Linguagens e Padrões

#### TypeScript/JavaScript

```bash
# Formatting
npm run format        # Prettier
npm run lint         # ESLint
npm run type-check   # TypeScript
```

**Regras:**
- Use TypeScript strict mode
- Prefira `const` sobre `let`
- Use arrow functions
- Documente funções públicas com JSDoc
- Evite `any`, use tipos específicos

#### Go

```bash
# Em cada serviço Go
go fmt ./...         # Formatting
go vet ./...         # Vet
golangci-lint run    # Linting
```

**Regras:**
- Siga o Go style guide
- Use context.Context para cancelamento
- Trate todos os erros
- Documente funções exportadas
- Use interfaces pequenas

#### Rust

```bash
# Em cada serviço Rust
cargo fmt           # Formatting
cargo clippy        # Linting
cargo test          # Tests
```

**Regras:**
- Use `rustfmt` padrão
- Evite `unwrap()` em produção
- Use `Result<T, E>` para erros
- Documente APIs públicas
- Prefira borrowing sobre ownership

#### Python

```bash
# Em cada serviço Python
black .             # Formatting
isort .             # Import sorting
flake8 .            # Linting
mypy .              # Type checking
```

**Regras:**
- Use type hints
- Siga PEP 8
- Use docstrings para funções
- Prefira async/await
- Use dataclasses/pydantic

### Estrutura de Arquivos

```
services/[service-name]/
├── cmd/                 # Entry points (Go)
├── src/                 # Source code
├── internal/            # Internal packages (Go)
├── pkg/                 # Public packages (Go)
├── configs/             # Configuration files
├── tests/               # Tests
├── Dockerfile           # Container definition
├── Makefile            # Service-specific commands
└── README.md           # Service documentation
```

## 🔄 Processo de Pull Request

### Checklist

Antes de abrir um PR, verifique:

- [ ] Código segue os padrões estabelecidos
- [ ] Testes passam (`make test`)
- [ ] Lint passa (`npm run lint`)
- [ ] Type check passa (`npm run type-check`)
- [ ] Documentação atualizada (se necessário)
- [ ] Changelog atualizado (para mudanças significativas)

### Revisão

Os PRs são revisados por:

1. **Maintainers** - Aprovação técnica
2. **Arquitetos** - Aprovação de design (para mudanças arquiteturais)
3. **QA** - Testes automatizados + manuais

### Critérios de Aceitação

- Funcionalidade implementada conforme especificação
- Código limpo e bem documentado
- Testes adequados (cobertura > 80%)
- Performance aceitável
- Segurança validada

## 🐛 Reportando Bugs

### Template de Bug Report

```markdown
## Descrição
Descrição clara do bug.

## Passos para Reproduzir
1. Faça X
2. Clique em Y
3. Veja o erro

## Comportamento Esperado
O que deveria acontecer.

## Comportamento Atual
O que realmente acontece.

## Ambiente
- OS: [ex: Ubuntu 22.04]
- Browser: [ex: Chrome 120]
- Node.js: [ex: 22.0.0]
- Versão: [ex: 1.0.0]

## Screenshots
Adicione screenshots se aplicável.

## Logs
```
Adicione logs relevantes
```

## Informações Adicionais
Contexto adicional sobre o problema.
```

### Bugs Críticos

Para bugs críticos (segurança, data loss, etc.):

1. **NÃO** abra issue pública
2. Envie email para: security@pagemagic.io
3. Inclua prova de conceito (se aplicável)
4. Aguarde resposta em 24h

## 💡 Sugerindo Features

### Template de Feature Request

```markdown
## Problema
Qual problema esta feature resolve?

## Solução Proposta
Descrição da solução ideal.

## Alternativas Consideradas
Outras abordagens que você considerou.

## Mockups/Wireframes
Adicione imagens se aplicável.

## Critérios de Aceitação
- [ ] Critério 1
- [ ] Critério 2

## Impacto Estimado
- Usuários afetados: [ex: todos, premium, etc.]
- Complexidade: [baixa/média/alta]
- Prioridade: [baixa/média/alta]
```

### Processo de Aprovação

1. **Community Discussion** - Issue aberta para discussão
2. **Product Review** - Análise pelo produto
3. **Technical Review** - Análise técnica/arquitetural
4. **Roadmap Planning** - Priorização no roadmap

## 🏗️ Arquitetura do Projeto

### Microserviços

| Serviço | Linguagem | Responsabilidade |
|---------|-----------|------------------|
| auth-svc | Go | Autenticação e autorização |
| prompt-svc | Node.js | Orquestração de IA |
| builder-svc | Rust | Construção de sites |
| build-svc | Go | Build e deploy |
| host-svc | Go | Hosting e containers |
| domain-svc | Python | Gestão de domínios |
| meter-svc | Rust | Medição de uso |
| billing-svc | Go | Cobrança e pagamentos |
| i18n-svc | Node.js | Internacionalização |

### Frontend

- **Web**: Next.js 15 + React 19 + TypeScript
- **Mobile**: React Native 0.74 + Expo SDK 51

### Infraestrutura

- **Database**: PostgreSQL 16 + TimescaleDB
- **Cache**: Redis
- **Message Broker**: NATS JetStream
- **Containers**: Docker + Docker Swarm
- **Monitoring**: Prometheus + Grafana
- **Reverse Proxy**: Nginx

### Fluxo de Dados

```mermaid
graph TD
    A[Frontend] --> B[Nginx]
    B --> C[Auth Service]
    B --> D[Prompt Service]
    D --> E[AI Models]
    B --> F[Builder Service]
    F --> G[Build Service]
    G --> H[Host Service]
    H --> I[Domain Service]

    C --> J[(PostgreSQL)]
    D --> J
    F --> J
    G --> J
    H --> J
    I --> J

    K[Usage Proxy] --> L[Meter Service]
    L --> M[(TimescaleDB)]
    L --> N[Billing Service]
    N --> O[Stripe]
```

## 📞 Suporte

### Canais de Comunicação

- **GitHub Issues**: Bugs e feature requests
- **GitHub Discussions**: Discussões gerais
- **Discord**: Chat em tempo real
- **Email**: contato@pagemagic.io

### Horários de Resposta

- **Issues críticos**: 24h
- **Pull Requests**: 48h
- **Issues gerais**: 72h
- **Feature requests**: 1 semana

## 📄 Licença

Ao contribuir, você concorda que suas contribuições serão licenciadas sob a [MIT License](LICENSE).

## 🙏 Reconhecimentos

Agradecemos a todos os contribuidores! Suas contribuições tornam o Page Magic melhor para todos.

---

**Obrigado por contribuir com o Page Magic! 🪄**
