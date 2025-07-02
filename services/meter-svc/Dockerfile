FROM rust:1.75-slim as builder

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    pkg-config \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy manifests
COPY Cargo.toml Cargo.lock ./

# Create dummy main to cache dependencies
RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN cargo build --release && rm -rf src

# Copy source code
COPY src/ ./src/

# Build application
RUN touch src/main.rs && cargo build --release

# Runtime image
FROM debian:bookworm-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    libssl3 \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd -m -u 1000 meter

# Copy binary
COPY --from=builder /app/target/release/meter-svc /usr/local/bin/meter-svc

# Set ownership
RUN chown meter:meter /usr/local/bin/meter-svc

USER meter

EXPOSE 8080

CMD ["meter-svc"]
