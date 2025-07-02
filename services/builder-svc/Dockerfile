# Build stage
FROM rust:1.75-bullseye AS builder

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    pkg-config \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy Cargo files
COPY Cargo.toml Cargo.lock ./

# Create empty src/main.rs for dependency caching
RUN mkdir src && echo "fn main() {}" > src/main.rs

# Build dependencies (cached layer)
RUN cargo build --release
RUN rm src/main.rs

# Copy source code
COPY src/ ./src/

# Build the application
RUN touch src/main.rs && cargo build --release

# Runtime stage
FROM debian:bullseye-slim

WORKDIR /app

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    libssl1.1 \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd --create-home --shell /bin/bash app

# Copy binary from builder stage
COPY --from=builder /app/target/release/builder-svc /app/builder-svc

# Copy templates and assets (if they exist)
COPY --chown=app:app templates/ /app/templates/
COPY --chown=app:app assets/ /app/assets/

# Switch to non-root user
USER app

# Expose port
EXPOSE 8082

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8082/health || exit 1

# Start the application
CMD ["./builder-svc"]
