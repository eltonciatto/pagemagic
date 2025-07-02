FROM node:20-alpine

WORKDIR /app

# Install system dependencies
RUN apk add --no-cache \
    python3 \
    make \
    g++

# Copy package files
COPY package*.json ./
COPY tsconfig.json ./

# Install dependencies
RUN npm ci --only=production

# Copy source code
COPY src/ ./src/

# Build application
RUN npm run build

# Create non-root user
RUN addgroup -g 1000 i18n && \
    adduser -D -s /bin/sh -u 1000 -G i18n i18n

# Create directories and set permissions
RUN mkdir -p logs && \
    chown -R i18n:i18n /app

USER i18n

EXPOSE 8080

CMD ["node", "dist/index.js"]
