# 🎯 Page Magic - Guia de Implementação FOCADA

## 🚨 FOCO TOTAL: FLUXO END-TO-END

**Meta:** Usuário digita descrição → Site online em <30 segundos

### 📋 FLUXO CRÍTICO

```
[INPUT] "Site para cafeteria moderna" 
   ↓ 
[PROMPT-SVC] → JSON estruturado
   ↓
[BUILDER-SVC] → HTML/CSS responsivo  
   ↓
[BUILD-SVC] → Container Docker
   ↓
[HOST-SVC] → Site online
   ↓
[OUTPUT] https://abc123.pagemagic.dev
```

---

## 🎯 IMPLEMENTAÇÃO POR DIA

### **DIA 1: PROMPT-SVC - IA FUNCIONANDO**

**Objetivo:** `/v1/generate` recebe texto → retorna JSON válido

**Implementar:**

1. **OpenAI client** básico
2. **Prompt engineering** para JSON estruturado
3. **Validation** da resposta
4. **Error handling** básico

**Files a criar/editar:**
```
services/prompt-svc/src/
  services/openai.ts     # Cliente OpenAI
  services/generator.ts  # Lógica de geração
  validation/schema.ts   # Validação JSON
  controllers/generate.ts # Endpoint /v1/generate
```

**JSON Schema de saída:**
```json
{
  "title": "string",
  "description": "string", 
  "sections": [
    {
      "type": "hero|features|testimonials|cta",
      "title": "string",
      "content": "string",
      "cta": { "text": "string", "link": "string" }
    }
  ],
  "theme": {
    "primaryColor": "#hex",
    "fontFamily": "string"
  }
}
```

**Critério de sucesso:** POST /v1/generate com "site para cafeteria" retorna JSON válido

---

### **DIA 2: BUILDER-SVC - JSON→HTML**

**Objetivo:** JSON estruturado → HTML/CSS responsivo

**Implementar:**

1. **Template engine** (Tera/Handlebars)
2. **Templates base** para cada tipo de seção
3. **CSS framework** embarcado (Tailwind)
4. **API endpoint** para conversão

**Files a criar:**
```
services/builder-svc/templates/
  base.html           # Layout principal
  sections/
    hero.html        # Hero section
    features.html    # Features grid
    testimonials.html # Testimonials
    cta.html         # Call to action

src/
  templates.rs       # Template rendering
  handlers/build.rs  # Endpoint /v1/build
```

**Critério de sucesso:** JSON do dia 1 → HTML responsivo funcional

---

### **DIA 3: BUILD-SVC - HTML→CONTAINER**

**Objetivo:** HTML/CSS → Container Docker deployável

**Implementar:**

1. **Static build** (copiar HTML/CSS/assets)
2. **Dockerfile generation** dinâmico
3. **Docker build** via API
4. **Registry push** (local registry)

**Files a criar:**
```
services/build-svc/internal/
  builders/static.go    # Static site builder
  docker/client.go     # Docker API client
  storage/registry.go  # Registry operations
  handlers/build.go    # Endpoint /v1/build
```

**Dockerfile template:**
```dockerfile
FROM nginx:alpine
COPY dist/ /usr/share/nginx/html/
EXPOSE 80
```

**Critério de sucesso:** HTML do dia 2 → Container rodando localmente

---

### **DIA 4: HOST-SVC - CONTAINER→ONLINE**

**Objetivo:** Container → Site acessível publicamente

**Implementar:**

1. **Docker run** com porta dinâmica
2. **Nginx proxy** reverso automático  
3. **Subdomain generation** (abc123.pagemagic.dev)
4. **Basic health checks**

**Files a criar:**
```
services/host-svc/internal/
  orchestrator/docker.go   # Docker container management
  proxy/nginx.go          # Nginx config generation
  dns/subdomain.go        # Subdomain management
  handlers/deploy.go      # Endpoint /v1/deploy
```

**Nginx config template:**
```nginx
server {
  server_name {subdomain}.pagemagic.dev;
  location / {
    proxy_pass http://localhost:{port};
  }
}
```

**Critério de sucesso:** Container do dia 3 → Site acessível em URL

---

### **DIA 5: FRONT-WEB - INTERFACE BÁSICA**

**Objetivo:** Interface web para o fluxo completo

**Implementar:**

1. **Generator page** (textarea + botão)
2. **Preview system** (iframe)
3. **Deploy button** funcional
4. **Basic navigation**

**Pages a criar:**
```
apps/front-web/src/app/
  generate/page.tsx      # Página principal
  preview/[id]/page.tsx  # Preview do site
  dashboard/page.tsx     # Lista de sites

components/
  Generator.tsx          # Formulário de geração
  SitePreview.tsx       # Preview component
  DeployButton.tsx      # Botão de deploy
```

**Flow:**
```
/ → /generate → /preview/abc123 → Click Deploy → Site Online
```

**Critério de sucesso:** Fluxo completo funciona na interface

---

### **DIA 6-7: INTEGRAÇÃO E POLISH**

**Implementar:**

1. **Service communication** entre todos
2. **Error handling** propagation
3. **Loading states** na interface
4. **Basic logging**

---

## 🎯 TESTES DE VALIDAÇÃO

### **Teste End-to-End Diário:**

```bash
# Dia 1: Prompt service
curl -X POST localhost:3001/v1/generate \
  -d '{"description": "site para cafeteria moderna"}'

# Dia 2: Builder service  
curl -X POST localhost:3002/v1/build \
  -d '{"json": {...}}'

# Dia 3: Build service
curl -X POST localhost:3003/v1/build \
  -d '{"html": "...", "css": "..."}'

# Dia 4: Host service
curl -X POST localhost:3004/v1/deploy \
  -d '{"containerId": "abc123"}'

# Dia 5: Frontend
# Abrir http://localhost:3000/generate e testar flow
```

### **Teste de Integração (Fim Semana 1):**

1. Abrir `http://localhost:3000`
2. Digitar "site para cafeteria moderna"
3. Clicar "Gerar Site"
4. Ver preview carregando
5. Clicar "Publicar Site"
6. Receber URL: `https://abc123.pagemagic.dev`
7. Acessar URL e ver site funcionando

---

## 🚫 O QUE **NÃO** IMPLEMENTAR

### **Ignorar completamente por agora:**

- ❌ OAuth (login simples serve)
- ❌ Database complex (memória/arquivos serve)
- ❌ vLLM cluster (OpenAI serve)
- ❌ Docker Swarm (docker run serve)
- ❌ SSL automático (HTTP serve)
- ❌ Editor visual (textarea serve)
- ❌ Mobile app
- ❌ Billing/Stripe
- ❌ Advanced themes
- ❌ Observability complexa

### **Regra de Ouro:**

**"Se não é essencial para o fluxo principal funcionar, NÃO implementar agora"**

---

## 🏁 DEFINIÇÃO DE SUCESSO

### **Fim da Semana 1:**

✅ **Demo funcionando:** Entrada de texto → Site online acessível
✅ **Performance aceitável:** <60 segundos end-to-end
✅ **Reliability básica:** >70% success rate
✅ **Interface utilizável:** Fluxo claro e intuitivo

### **Critério PASS/FAIL:**

**PASS:** Conseguir fazer demo para stakeholder mostrando fluxo completo
**FAIL:** Qualquer parte do fluxo não funciona

---

## 📞 PRÓXIMA AÇÃO IMEDIATA

**🎯 COMEÇAR AGORA:**

```bash
cd /workspaces/pagemagic/services/prompt-svc
```

1. **Instalar OpenAI SDK**
2. **Criar endpoint /v1/generate** 
3. **Testar com prompt simples**
4. **Validar JSON de saída**

**Time-box:** 4 horas para ter primeiro endpoint funcionando

---

## 🎯 MANTRA PARA A SEMANA

**"FUNCIONAL > BONITO"**
**"SIMPLES > SOFISTICADO"**
**"FLUXO > FEATURES"**

O objetivo é ter algo que **FUNCIONA** de ponta a ponta, mesmo que básico. Features avançadas só depois do fluxo principal estar estável.
