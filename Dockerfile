# ------------------------------------------------------------------------------
# Stage 1: Builder
# This stage builds the application and its dependencies.
# ------------------------------------------------------------------------------
FROM golang:1.24-alpine3.21 AS builder

# Enable CGO for SQLite support
ENV CGO_ENABLED=1

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy modules and checksums files
COPY go.mod go.sum      ./

# Download modules
RUN go mod download

# Copy application sources (packages)
COPY main.go            ./
COPY controller/        ./controller/
COPY data/              ./data/
COPY docs/              ./docs/
COPY model/             ./model/
COPY route/             ./route/
COPY service/           ./service/
COPY swagger/           ./swagger/

# Build the application binary
RUN go build -o app .

# ------------------------------------------------------------------------------
# Stage 2: Runtime
# This stage creates the final, minimal image to run the application.
# ------------------------------------------------------------------------------
FROM alpine:3.21 AS runtime

# Install curl for health check
RUN apk add --no-cache curl

WORKDIR /app

# Metadata labels for the image. These are useful for registries and inspection.
LABEL org.opencontainers.image.title="🧪 RESTful API with Go and Gin"
LABEL org.opencontainers.image.description="Proof of Concept for a RESTful API made with Go and Gin"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.source="https://github.com/nanotaboada/go-samples-gin-restful"
LABEL org.sonarsource.docker.dockerfile="/Dockerfile"

# https://rules.sonarsource.com/docker/RSPEC-6504/

# Copy application binary and database file
COPY --from=builder     /app/app                    .

# Copy metadata docs for container registries (e.g.: GitHub Container Registry)
COPY --chmod=444        README.md                   ./
COPY --chmod=555        assets/                     ./assets/

# Copy entrypoint and healthcheck scripts
COPY --chmod=555        scripts/entrypoint.sh       ./entrypoint.sh
COPY --chmod=555        scripts/healthcheck.sh      ./healthcheck.sh

# Copy pre-seeded SQLite database as init bundle
COPY --chmod=555        storage/                    ./docker-compose/

# Add non-root user and prepare volume mount point
RUN adduser --disabled-password --gecos "" gin && \
    mkdir -p /storage && \
    chown -R gin:gin /storage

USER gin

EXPOSE 9000

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD ["./healthcheck.sh"]

ENTRYPOINT ["./entrypoint.sh"]
CMD ["./app"]
