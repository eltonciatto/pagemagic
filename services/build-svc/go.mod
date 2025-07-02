module pagemagic/build-svc

go 1.23

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/docker/docker v24.0.7+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/moby/buildkit v0.12.4
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	github.com/prometheus/client_golang v1.18.0
	github.com/nats-io/nats.go v1.31.0
	github.com/redis/go-redis/v9 v9.4.0
	go.uber.org/zap v1.26.0
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/containerd/containerd v1.7.10
	github.com/distribution/reference v0.5.0
	github.com/docker/cli v24.0.7+incompatible
	github.com/docker/buildx v0.11.2
	github.com/opencontainers/image-spec v1.1.0-rc5
	github.com/pkg/errors v0.9.1
)
