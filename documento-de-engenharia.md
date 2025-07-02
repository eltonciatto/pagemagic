# **Page Magic**

### Especificação Técnica – Plataforma **No-Code** de Landing Pages, Blogs e Sites guiada por IA

*(documento de engenharia — desktop + mobile Android / iOS)*

---

## 0. Escopo em uma frase

**O usuário descreve o que quer → a IA gera o site completo → o usuário edita visualmente → clica em “Publicar” → a plataforma cria o contêiner, conecta o domínio, mede o uso (Stripe) e cobra — tudo via web ou app móvel.**

---

## 1. Arquitetura macro

```mermaid
graph TD
subgraph Front-ends
A1[front-web (Next.js)]:::fe
A2[mobile-app (React Native)]:::fe
end
A1 -->|REST/GRPC| P(prompt-svc)
A2 -->|REST/GRPC| P

A1 --> B(builder-svc)
A2 --> B

B --> C(build-svc)
C --> D(host-svc)
D --> E(domain-svc)

A1 --> F(usage-proxy)
A2 --> F
F --> G(meter-svc)
G --> H(billing-svc)
H --> I[Stripe]

A1 --> J(auth-svc)
A2 --> J
auth-svc -->|JWT| A1
auth-svc -->|JWT| A2

D --> K[(Observability)]

subgraph IA
P --> L[vLLM cluster]
P --> M[OpenAI / Claude adapters]
end

subgraph Data
auth-svc --> Postgres[(PostgreSQL 16)]
builder-svc --> Postgres
meter-svc --> TSDB[(TimescaleDB)]
billing-svc --> Postgres
end
```

* **Mensageria:** NATS JetStream
* **Observabilidade:** Prometheus + Grafana + Loki
* **Deploy:** Coolify (Docker Swarm) via GitOps (Flux)
* **Edge:** Cloudflare proxied, Workers para cache estático

---

## 2. Serviços detalhados

| Sigla             | Linguagem                       | Responsabilidade                                           |
| ----------------- | ------------------------------- | ---------------------------------------------------------- |
| **front-web**     | Next.js 15.3 • React 19         | UI desktop e PWA, editor *no-code*, dashboards             |
| **mobile-app**    | React Native 0.74 (Expo SDK 51) | Experiência nativa Android/iOS, preview WebView, push, IAP |
| **auth-svc**      | Go 1.23                         | Magic-link, OAuth (NextAuth adapter)                       |
| **prompt-svc**    | Node 22                         | Orquestra prompts (LangChainJS)                            |
| **builder-svc**   | Rust (Actix-web)                | JSON → AST → HTML/React                                    |
| **build-svc**     | Go                              | Gera imagem Docker via Buildx + Turbopack                  |
| **host-svc**      | Go                              | Cria/pausa/escala contêineres (Swarm API)                  |
| **domain-svc**    | Python (FastAPI)                | Compra domínios, gerencia DNS/ACME                         |
| **usage-proxy**   | Nginx + Lua                     | Injeta eventos de uso                                      |
| **meter-svc**     | Rust                            | Agrega eventos, envia Stripe Meters                        |
| **billing-svc**   | Go                              | Webhooks Stripe → credit-control & dunning                 |
| **i18n-svc**      | Node                            | Administra bundles JSON de tradução                        |
| **observability** | Helm charts                     | Prometheus, Grafana, Alertmanager                          |

---

## 3. Fluxos primários

### 3.1 Geração de página (web + mobile)

1. **POST `/generate`** (front-web ou mobile-app)
   Payload → *prompt-svc*
2. *prompt-svc* chama **vLLM** (Llama-3 70B) ou API externa.
3. Retorno `{sections[],css,js,i18n}` → *builder-svc* → AST.
4. **Streaming multipart** devolve preview:
   *Web*: iframe side-by-side.
   *Mobile*: WebView no modo “preview” + overlay de edição.

> Todas strings → chaves i18n (`"hero.title"`) + mapa `i18n[key]=texto`.

### 3.2 Publicação

Idêntico para web e mobile:

1. **POST `/sites/:id/build`**
2. *build-svc* compila, gera **image\_id**.
3. *host-svc* cria serviço `stack_site_<id>`.
4. *domain-svc* aponta DNS, ACME DNS-01 (Caddy).
5. Evento `site.published` → apps mostram status **Publicado**.

### 3.3 Medição & cobrança

* **usage-proxy** adiciona header `x-meter-event`.
* *meter-svc* envia lotes de 60 s → Stripe Meters (`page_generate`, `ai_token`, `container_hours`, `storage_gb`).
* *billing-svc* trata `invoice.finalized`, bloqueia se `paid=false`.

---

## 4. APIs essenciais

### 4.1 prompt-svc

```http
POST /v1/generate
Authorization: Bearer <jwt>
{
  "project_id": "proj_x",
  "locale": "en",
  "prompt": "Create a landing page for a vegan bakery"
}
→ 202
{ "generation_id": "gen_abc", "status": "queued" }
```

### 4.2 builder-svc

```http
POST /v1/ast
{
  "sections": [ { "id":"hero", "html":"<header>...</header>" } ],
  "css":"...", "js":"..."
}
→ 200 { "ast_id":"ast_xyz", "warnings":[] }
```

### 4.3 mobile-app SDK extras

| Endpoint                                    | Uso                                              |
| ------------------------------------------- | ------------------------------------------------ |
| `GET /v1/sites/:id/preview?platform=mobile` | Retorna bundle otimizado (lazy JS)               |
| `POST /v1/push/register`                    | Salva `expo_push_token` p/ notificações de quota |

---

## 5. Modelo de dados (resumido)

*Conforme diagrama UML anterior* + **tabela mobile\_push\_tokens**

```
mobile_push_tokens: { id PK, user_id FK, token, platform, created_at }
```

---

## 6. Camada mobile — detalhes

| Área                  | Implementação                                                          |
| --------------------- | ---------------------------------------------------------------------- |
| **Stack**             | React Native 0.74 (Expo Router 3), TypeScript, Zustand                 |
| **Preview**           | `WebView` (<iframe>) com bridge postMessage para salvar edições        |
| **Uploads**           | `expo-image-picker` + direct-PUT MinIO (presigned URL)                 |
| **Push**              | Expo Notifications; eventos de pagamento, limite 90 % e build pronto   |
| **Pagamentos in-app** | Links Stripe Checkout no *webview modal* (política App Store OK)       |
| **Deep Linking**      | `pagemagic://site/123` abre editor                                     |
| **Offline**           | Cache AST em SQLite; re-sync quando online                             |
| **A/B Quick share**   | “Compartilhar preview” gera link expiring 24 h (Cloudflare Signed URL) |

---

## 7. Stripe Meters (2025+)

| Meter           | key               | Coleta                     | Fórmula |
| --------------- | ----------------- | -------------------------- | ------- |
| Gerações        | `page_generate`   | prompt-svc                 | `count` |
| Tokens          | `ai_token`        | prompt-svc lê header usage | `sum`   |
| Horas contêiner | `container_hours` | host-svc tick 60 min       | `sum`   |
| Storage         | `storage_gb`      | Cron diário MinIO          | `last`  |

*Criados em `/v1/billing/meter`; `usage_records` legado **deprecado**.*

---

## 8. Detalhes IA

| Item              | Valor                                                                 |
| ----------------- | --------------------------------------------------------------------- |
| **Modelo padrão** | Llama-3 70B-Instr (FP16) via vLLM 0.8                                 |
| **GPU**           | 1 × A100-80 GB → \~780 tok/s                                          |
| **Fallback**      | INT4-AWQ em T4 spot                                                   |
| **Rotas**         | `/generate`, `/continue`, `/rewrite` (saída validada por JSON Schema) |

---

## 9. Pipeline CI/CD

1. **PR** → GitHub Actions: lint + Vitest + Playwright (web) + Detox (mobile).
2. `make image SERVICE=x` → push registry.
3. Flux atualiza Helm/Swarm.
4. Smoke test staging.
5. Tag **main** → produção; apps mobile via **EAS Update** OTA + builds store semestrais.

---

## 10. Segurança & Compliance

* CSP rigorosa; zero inline JS no builder.
* OWASP Top-10 — WAF Cloudflare + scanners ZAP.
* PCI-DSS: cartões somente em Stripe (token).
* LGPD/GDPR: PII criptografado (pgcrypto).
* Mobile: Keystore/Keychain para JWT; TLS 1.3 obrigatório.
* Contêineres: usuário não-root, seccomp, memória min 256 MiB.

---

## 11. Cronograma granular (inclui mobile)

| Semana | Entrega                                      | Owner    |
| ------ | -------------------------------------------- | -------- |
| 1-2    | auth-svc + infra base                        | DevOps   |
| 3-4    | prompt-svc mínimo + vLLM                     | IA       |
| 5-6    | builder-svc + front-web preview              | FE/BE    |
| 7      | build-svc + host-svc                         | Platform |
| 8      | domain-svc                                   | Platform |
| 9      | meter-svc + usage-proxy                      | BE       |
| 10     | billing-svc + Stripe Webhooks                | BE       |
| 11     | **mobile-app MVP** (login, geração, preview) | Mobile   |
| 12     | Dashboards (Grafana + front)                 | FE/SRE   |
| 13     | Push notifications + quota                   | Mobile   |
| 14     | QA geral, hardening, docs                    | QA       |
| 15     | Lançamento stores + monitoramento            | PM/SRE   |

---

### Conclusão

Esta versão “**Page Magic**” inclui:

* **Back-end micro-serviços** claramente segregados.
* **Apps** web e **mobile nativos** integrados ao mesmo backend.
* **Contratos de API**, modelo de dados, fluxo de build & deploy.
* **Medição Stripe Meters** 100 % alinhada à política 2025.
* Cronograma prático de 15 semanas para MVP *multi-plataforma*.

Com o documento, o time de desenvolvimento tem tudo que precisa para começar — minimizando dúvidas, evitando retrabalho e garantindo que Page Magic seja escalável, competitivo e financeiramente sustentável.
