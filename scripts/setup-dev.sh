# Page Magic Development Environment Setup Script
#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logo
echo -e "${BLUE}"
echo "  ____                   __  __             _      "
echo " |  _ \ __ _  __ _  ___  |  \/  | __ _  __ _(_) ___ "
echo " | |_) / _\` |/ _\` |/ _ \ | |\/| |/ _\` |/ _\` | |/ __|"
echo " |  __/ (_| | (_| |  __/ | |  | | (_| | (_| | | (__ "
echo " |_|   \__,_|\__, |\___| |_|  |_|\__,_|\__, |_|\___|"
echo "             |___/                    |___/        "
echo -e "${NC}"
echo -e "${GREEN}Setting up Page Magic development environment...${NC}"
echo

# Check if running as root
if [[ $EUID -eq 0 ]]; then
   echo -e "${RED}This script should not be run as root${NC}"
   exit 1
fi

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to install package with different package managers
install_package() {
    local package=$1

    if command_exists apt-get; then
        sudo apt-get update && sudo apt-get install -y "$package"
    elif command_exists yum; then
        sudo yum install -y "$package"
    elif command_exists brew; then
        brew install "$package"
    elif command_exists pacman; then
        sudo pacman -S --noconfirm "$package"
    else
        echo -e "${RED}Unable to install $package. Please install manually.${NC}"
        return 1
    fi
}

# Check and install dependencies
echo -e "${YELLOW}Checking dependencies...${NC}"

# Docker
if ! command_exists docker; then
    echo -e "${YELLOW}Installing Docker...${NC}"
    curl -fsSL https://get.docker.com -o get-docker.sh
    sh get-docker.sh
    sudo usermod -aG docker $USER
    rm get-docker.sh
    echo -e "${GREEN}Docker installed successfully${NC}"
else
    echo -e "${GREEN}Docker is already installed${NC}"
fi

# Docker Compose
if ! command_exists docker-compose; then
    echo -e "${YELLOW}Installing Docker Compose...${NC}"
    sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    echo -e "${GREEN}Docker Compose installed successfully${NC}"
else
    echo -e "${GREEN}Docker Compose is already installed${NC}"
fi

# Node.js (using NodeSource)
if ! command_exists node; then
    echo -e "${YELLOW}Installing Node.js...${NC}"
    curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
    sudo apt-get install -y nodejs
    echo -e "${GREEN}Node.js installed successfully${NC}"
else
    echo -e "${GREEN}Node.js is already installed ($(node --version))${NC}"
fi

# Go
if ! command_exists go; then
    echo -e "${YELLOW}Installing Go...${NC}"
    GO_VERSION="1.21.5"
    wget "https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz"
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf "go${GO_VERSION}.linux-amd64.tar.gz"
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    export PATH=$PATH:/usr/local/go/bin
    rm "go${GO_VERSION}.linux-amd64.tar.gz"
    echo -e "${GREEN}Go installed successfully${NC}"
else
    echo -e "${GREEN}Go is already installed ($(go version))${NC}"
fi

# Rust
if ! command_exists cargo; then
    echo -e "${YELLOW}Installing Rust...${NC}"
    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
    source ~/.cargo/env
    echo -e "${GREEN}Rust installed successfully${NC}"
else
    echo -e "${GREEN}Rust is already installed ($(rustc --version))${NC}"
fi

# Python 3 and pip
if ! command_exists python3; then
    echo -e "${YELLOW}Installing Python 3...${NC}"
    install_package python3
    install_package python3-pip
    echo -e "${GREEN}Python 3 installed successfully${NC}"
else
    echo -e "${GREEN}Python 3 is already installed ($(python3 --version))${NC}"
fi

# Git
if ! command_exists git; then
    echo -e "${YELLOW}Installing Git...${NC}"
    install_package git
    echo -e "${GREEN}Git installed successfully${NC}"
else
    echo -e "${GREEN}Git is already installed ($(git --version))${NC}"
fi

# Make
if ! command_exists make; then
    echo -e "${YELLOW}Installing Make...${NC}"
    install_package build-essential 2>/dev/null || install_package make
    echo -e "${GREEN}Make installed successfully${NC}"
else
    echo -e "${GREEN}Make is already installed${NC}"
fi

# VS Code (if not running in remote environment)
if [[ -z "$CODESPACES" && -z "$GITPOD_WORKSPACE_ID" ]]; then
    if ! command_exists code; then
        echo -e "${YELLOW}Installing VS Code...${NC}"
        wget -qO- https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > packages.microsoft.gpg
        sudo install -o root -g root -m 644 packages.microsoft.gpg /etc/apt/trusted.gpg.d/
        sudo sh -c 'echo "deb [arch=amd64,arm64,armhf signed-by=/etc/apt/trusted.gpg.d/packages.microsoft.gpg] https://packages.microsoft.com/repos/code stable main" > /etc/apt/sources.list.d/vscode.list'
        sudo apt-get update
        sudo apt-get install -y code
        echo -e "${GREEN}VS Code installed successfully${NC}"
    else
        echo -e "${GREEN}VS Code is already installed${NC}"
    fi
fi

echo
echo -e "${YELLOW}Setting up project dependencies...${NC}"

# Create necessary directories
mkdir -p .logs
mkdir -p .tmp
mkdir -p infrastructure/ssl
mkdir -p infrastructure/nats

# Set up environment files
if [ ! -f .env ]; then
    echo -e "${YELLOW}Creating .env file...${NC}"
    cp .env.example .env 2>/dev/null || cat > .env << 'EOF'
# Page Magic Environment Variables

# Database
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=pagemagic
POSTGRES_USER=pagemagic
POSTGRES_PASSWORD=pagemagic123

# Redis
REDIS_URL=redis://localhost:6379

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-in-production

# AI Services
OPENAI_API_KEY=your-openai-api-key
ANTHROPIC_API_KEY=your-anthropic-api-key

# Stripe
STRIPE_PUBLIC_KEY=pk_test_your-stripe-public-key
STRIPE_SECRET_KEY=sk_test_your-stripe-secret-key
STRIPE_WEBHOOK_SECRET=whsec_your-webhook-secret

# Email
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# Domain
NEXT_PUBLIC_DOMAIN=localhost:3000
NEXT_PUBLIC_API_URL=http://localhost:8080

# Development
NODE_ENV=development
RUST_LOG=debug
EOF
    echo -e "${GREEN}.env file created${NC}"
else
    echo -e "${GREEN}.env file already exists${NC}"
fi

# Install Node.js dependencies for frontend
if [ -f "apps/front-web/package.json" ]; then
    echo -e "${YELLOW}Installing frontend dependencies...${NC}"
    cd apps/front-web
    npm install
    cd ../..
    echo -e "${GREEN}Frontend dependencies installed${NC}"
fi

# Install Node.js dependencies for mobile app
if [ -f "apps/mobile-app/package.json" ]; then
    echo -e "${YELLOW}Installing mobile app dependencies...${NC}"
    cd apps/mobile-app
    npm install
    cd ../..
    echo -e "${GREEN}Mobile app dependencies installed${NC}"
fi

# Install Node.js dependencies for services
for service in prompt-svc i18n-svc; do
    if [ -f "services/$service/package.json" ]; then
        echo -e "${YELLOW}Installing $service dependencies...${NC}"
        cd "services/$service"
        npm install
        cd ../..
        echo -e "${GREEN}$service dependencies installed${NC}"
    fi
done

# Install Go dependencies
for service in auth-svc build-svc host-svc billing-svc; do
    if [ -f "services/$service/go.mod" ]; then
        echo -e "${YELLOW}Installing $service dependencies...${NC}"
        cd "services/$service"
        go mod download
        cd ../..
        echo -e "${GREEN}$service dependencies installed${NC}"
    fi
done

# Install Rust dependencies
for service in builder-svc meter-svc; do
    if [ -f "services/$service/Cargo.toml" ]; then
        echo -e "${YELLOW}Installing $service dependencies...${NC}"
        cd "services/$service"
        cargo build --release
        cd ../..
        echo -e "${GREEN}$service dependencies installed${NC}"
    fi
done

# Install Python dependencies
if [ -f "services/domain-svc/requirements.txt" ]; then
    echo -e "${YELLOW}Installing domain service dependencies...${NC}"
    cd services/domain-svc
    python3 -m venv venv
    source venv/bin/activate
    pip install -r requirements.txt
    deactivate
    cd ../..
    echo -e "${GREEN}Domain service dependencies installed${NC}"
fi

# Set up Git hooks
if [ -d ".git" ]; then
    echo -e "${YELLOW}Setting up Git hooks...${NC}"
    cp scripts/git-hooks/* .git/hooks/ 2>/dev/null || echo "No git hooks found"
    chmod +x .git/hooks/* 2>/dev/null || true
    echo -e "${GREEN}Git hooks set up${NC}"
fi

# Create development SSL certificates
echo -e "${YELLOW}Creating development SSL certificates...${NC}"
mkdir -p infrastructure/ssl
if [ ! -f "infrastructure/ssl/localhost.crt" ]; then
    openssl req -x509 -newkey rsa:4096 -keyout infrastructure/ssl/localhost.key -out infrastructure/ssl/localhost.crt -days 365 -nodes -subj "/C=US/ST=CA/L=San Francisco/O=Page Magic/CN=localhost"
    echo -e "${GREEN}SSL certificates created${NC}"
else
    echo -e "${GREEN}SSL certificates already exist${NC}"
fi

# Build development Docker images
echo -e "${YELLOW}Building Docker images...${NC}"
docker-compose -f docker-compose.dev.yml build --parallel
echo -e "${GREEN}Docker images built successfully${NC}"

# Create and start containers
echo -e "${YELLOW}Starting development environment...${NC}"
docker-compose -f docker-compose.dev.yml up -d postgres redis timescaledb minio nats
echo -e "${GREEN}Core services started${NC}"

# Wait for databases to be ready
echo -e "${YELLOW}Waiting for databases to be ready...${NC}"
sleep 30

# Run database migrations
echo -e "${YELLOW}Running database migrations...${NC}"
docker-compose -f docker-compose.dev.yml exec -T postgres psql -U pagemagic -d pagemagic -f /docker-entrypoint-initdb.d/01-init.sql 2>/dev/null || true
echo -e "${GREEN}Database migrations completed${NC}"

# Start all services
echo -e "${YELLOW}Starting all services...${NC}"
docker-compose -f docker-compose.dev.yml up -d
echo -e "${GREEN}All services started${NC}"

echo
echo -e "${GREEN}🎉 Page Magic development environment setup completed!${NC}"
echo
echo -e "${YELLOW}Available services:${NC}"
echo -e "  • Frontend: ${BLUE}http://localhost:3000${NC}"
echo -e "  • API Gateway: ${BLUE}http://localhost:8080${NC}"
echo -e "  • Auth Service: ${BLUE}http://localhost:3001${NC}"
echo -e "  • Prompt Service: ${BLUE}http://localhost:3002${NC}"
echo -e "  • Builder Service: ${BLUE}http://localhost:3003${NC}"
echo -e "  • Build Service: ${BLUE}http://localhost:3004${NC}"
echo -e "  • Host Service: ${BLUE}http://localhost:3005${NC}"
echo -e "  • Domain Service: ${BLUE}http://localhost:3006${NC}"
echo -e "  • i18n Service: ${BLUE}http://localhost:3007${NC}"
echo -e "  • Meter Service: ${BLUE}http://localhost:3008${NC}"
echo -e "  • Billing Service: ${BLUE}http://localhost:3009${NC}"
echo
echo -e "${YELLOW}Development tools:${NC}"
echo -e "  • Grafana: ${BLUE}http://localhost:3001${NC} (admin/admin123)"
echo -e "  • Prometheus: ${BLUE}http://localhost:9090${NC}"
echo -e "  • Redis Commander: ${BLUE}http://localhost:8081${NC}"
echo -e "  • pgAdmin: ${BLUE}http://localhost:8082${NC} (admin@pagemagic.dev/admin123)"
echo -e "  • MinIO Console: ${BLUE}http://localhost:9001${NC} (pagemagic/pagemagic123)"
echo -e "  • Mailhog: ${BLUE}http://localhost:8025${NC}"
echo -e "  • Jaeger: ${BLUE}http://localhost:16686${NC}"
echo
echo -e "${YELLOW}Useful commands:${NC}"
echo -e "  • Start all services: ${BLUE}make dev${NC}"
echo -e "  • Stop all services: ${BLUE}make down${NC}"
echo -e "  • View logs: ${BLUE}make logs${NC}"
echo -e "  • Build services: ${BLUE}make build${NC}"
echo -e "  • Run tests: ${BLUE}make test${NC}"
echo -e "  • Database shell: ${BLUE}make db-shell${NC}"
echo
echo -e "${GREEN}Happy coding! 🚀${NC}"

# Check if Docker group membership requires re-login
if ! docker ps >/dev/null 2>&1; then
    echo
    echo -e "${YELLOW}⚠️  You may need to log out and log back in for Docker group membership to take effect.${NC}"
    echo -e "${YELLOW}   Or run: ${BLUE}newgrp docker${NC}"
fi
