# - Stage 1 --------------------------------------------------------------------

    FROM golang:1.24-alpine AS builder

    # Enable CGO for sqlite3 support
    ENV CGO_ENABLED=1
    RUN apk add --no-cache gcc musl-dev sqlite-dev

    WORKDIR /app

    # Copy modules and checksums files
    COPY go.mod go.sum ./

    # Download modules
    RUN go mod download

    # Copy application sources (packages)
    COPY main.go ./
    COPY controller/ ./controller/
    COPY data/ ./data/
    COPY docs/ ./docs/
    COPY model/ ./model/
    COPY route/ ./route/
    COPY service/ ./service/
    COPY swagger/ ./swagger/
    COPY assets/ ./assets/
    COPY README.md ./

    # Build the application binary
    RUN go build -o app .

    # - Stage 2 --------------------------------------------------------------------

    FROM golang:1.24-alpine AS runtime

    WORKDIR /app

    # Add non-root user for security hardening
    RUN adduser --disabled-password --gecos "" gin && \
        mkdir -p /app/data /app/assets && \
        chown -R gin:gin /app

    # Copy application binary and database file
    COPY --from=builder /app/app .
    COPY --from=builder /app/data/players_sqlite3.db ./data/

    # Copy README and assets for GHCR package metadata
    COPY --from=builder /app/README.md ./
    COPY --from=builder /app/assets/ ./assets/

    USER gin

    EXPOSE 9000

    CMD ["./app"]
