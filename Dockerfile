# - Stage 1 --------------------------------------------------------------------

    FROM golang:1.24.1-alpine AS builder

    # Enable CGO
    ENV CGO_ENABLED=1 \
        GOOS=linux \
        GOARCH=amd64

    WORKDIR /app

    # Install required dependencies for go-sqlite3
    RUN apk add --no-cache gcc musl-dev sqlite-dev

    # Copy dependencies files
    COPY go.mod go.sum ./

    # Download dependencies
    RUN go mod download

    # Copy app sources (packages)
    COPY . .

    # Build the app binary
    RUN go build -o app .

# - Stage 2 --------------------------------------------------------------------

    FROM alpine:latest AS runtime

    # Install SQLite runtime lib
    RUN apk add --no-cache sqlite-libs

    WORKDIR /root/

    # Copy app binary and SQLite DB
    COPY --from=builder /app/app .
    COPY --from=builder /app/data/players_sqlite3.db ./data/

    EXPOSE 9000

    CMD ["./app"]
