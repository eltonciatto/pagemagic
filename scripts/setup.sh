#!/bin/bash

# Page Magic - Script de Inicialização Rápida
# Este script configura o ambiente de desenvolvimento completo

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🪄 Page Magic - Configuração Inicial${NC}"
echo ""

# Verificar se Docker está rodando
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker não está rodando. Por favor, inicie o Docker primeiro.${NC}"
    exit 1
fi

# Verificar se Docker Compose está disponível
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}❌ Docker Compose não encontrado. Por favor, instale o Docker Compose.${NC}"
    exit 1
fi

# Criar arquivo .env se não existir
if [ ! -f .env ]; then
    echo -e "${YELLOW}📝 Criando arquivo .env...${NC}"
    cp .env.example .env
    echo -e "${GREEN}✅ Arquivo .env criado. Edite-o com suas configurações se necessário.${NC}"
fi

# Criar diretórios necessários
echo -e "${YELLOW}📁 Criando diretórios necessários...${NC}"
mkdir -p logs
mkdir -p data/postgres
mkdir -p data/redis
mkdir -p data/minio
mkdir -p tmp/builds

# Inicializar infraestrutura
echo -e "${YELLOW}🐳 Iniciando infraestrutura com Docker Compose...${NC}"
make infra-up

# Aguardar serviços ficarem prontos
echo -e "${YELLOW}⏳ Aguardando serviços ficarem prontos...${NC}"
sleep 10

# Verificar se PostgreSQL está pronto
echo -e "${YELLOW}🔍 Verificando PostgreSQL...${NC}"
for i in {1..30}; do
    if docker-compose -f infrastructure/docker/docker-compose.yml exec -T postgres pg_isready -U pagemagic > /dev/null 2>&1; then
        echo -e "${GREEN}✅ PostgreSQL está pronto!${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${RED}❌ PostgreSQL não ficou pronto em 30 tentativas${NC}"
        exit 1
    fi
    sleep 1
done

# Executar migrações do banco
echo -e "${YELLOW}🗄️  Executando migrações do banco...${NC}"
docker-compose -f infrastructure/docker/docker-compose.yml exec -T postgres psql -U pagemagic -d pagemagic -f /docker-entrypoint-initdb.d/migrations.sql || true

# Instalar dependências dos serviços
echo -e "${YELLOW}📦 Instalando dependências...${NC}"

# Go services
if command -v go &> /dev/null; then
    echo -e "${BLUE}  📦 Instalando dependências do auth-svc (Go)...${NC}"
    cd services/auth-svc && go mod download && cd ../..
else
    echo -e "${YELLOW}  ⚠️  Go não encontrado. Pule a instalação das dependências Go.${NC}"
fi

# Node.js services
if command -v npm &> /dev/null; then
    echo -e "${BLUE}  📦 Instalando dependências do prompt-svc (Node.js)...${NC}"
    cd services/prompt-svc && npm install && cd ../..
    
    echo -e "${BLUE}  📦 Instalando dependências do front-web (Next.js)...${NC}"
    cd apps/front-web && npm install && cd ../..
else
    echo -e "${YELLOW}  ⚠️  Node.js/npm não encontrado. Pule a instalação das dependências Node.js.${NC}"
fi

# Rust services
if command -v cargo &> /dev/null; then
    echo -e "${BLUE}  📦 Verificando dependências do builder-svc (Rust)...${NC}"
    cd services/builder-svc && cargo check && cd ../..
else
    echo -e "${YELLOW}  ⚠️  Rust/Cargo não encontrado. Pule a verificação das dependências Rust.${NC}"
fi

echo ""
echo -e "${GREEN}🎉 Configuração inicial concluída!${NC}"
echo ""
echo -e "${BLUE}📋 Próximos passos:${NC}"
echo -e "  1. ${YELLOW}make web-dev${NC}     - Inicia o front-end web (Next.js)"
echo -e "  2. ${YELLOW}make auth-dev${NC}    - Inicia o serviço de autenticação"
echo -e "  3. ${YELLOW}make prompt-dev${NC}  - Inicia o serviço de prompts IA"
echo -e "  4. ${YELLOW}make builder-dev${NC} - Inicia o serviço de construção"
echo ""
echo -e "${BLUE}📊 Monitoramento:${NC}"
echo -e "  • Grafana:    ${YELLOW}http://localhost:3000${NC}"
echo -e "  • Prometheus: ${YELLOW}http://localhost:9090${NC}"
echo -e "  • PgAdmin:    ${YELLOW}http://localhost:5050${NC}"
echo ""
echo -e "${BLUE}📖 Comandos úteis:${NC}"
echo -e "  • ${YELLOW}make help${NC}        - Ver todos os comandos disponíveis"
echo -e "  • ${YELLOW}make infra-down${NC}  - Parar infraestrutura"
echo -e "  • ${YELLOW}make dev-backend${NC} - Iniciar todos os serviços backend"
echo ""
echo -e "${GREEN}✨ Page Magic está pronto para o desenvolvimento!${NC}"
