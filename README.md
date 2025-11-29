# 🧪 RESTful API with Go and Gin

[![Go CI](https://github.com/nanotaboada/go-samples-gin-restful/actions/workflows/go.yml/badge.svg)](https://github.com/nanotaboada/go-samples-gin-restful/actions/workflows/go.yml)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=nanotaboada_go-samples-gin-restful&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=nanotaboada_go-samples-gin-restful)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/e2b234e8182d4a0f8efb9da619e1dc26)](https://app.codacy.com/gh/nanotaboada/go-samples-gin-restful/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![codecov](https://codecov.io/gh/nanotaboada/go-samples-gin-restful/graph/badge.svg?token=i37VDcDWwx)](https://codecov.io/gh/nanotaboada/go-samples-gin-restful)
[![Go Report Card](https://goreportcard.com/badge/github.com/nanotaboada/go-samples-gin-restful)](https://goreportcard.com/report/github.com/nanotaboada/go-samples-gin-restful)
[![CodeFactor](https://www.codefactor.io/repository/github/nanotaboada/go-samples-gin-restful/badge)](https://www.codefactor.io/repository/github/nanotaboada/go-samples-gin-restful)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Proof of Concept for a RESTful API built with [Go](https://github.com/golang/go) and [Gin](https://github.com/gin-gonic/gin). Manage football player data with SQLite, GORM, Swagger documentation, and in-memory caching.

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Architecture](#architecture)
- [API Reference](#api-reference)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Testing](#testing)
- [Docker](#docker)
- [Environment Variables](#environment-variables)
- [Command Summary](#command-summary)
- [Contributing](#contributing)
- [Legal](#legal)

## Features

- 🔌 **RESTful CRUD operations** for football player data
- 📚 **Interactive API documentation** with Swagger UI and Postman collection
- ⚡ **In-memory caching** (1-hour TTL)
- 🩺 **Health check endpoint** for monitoring
- 🐳 **Full containerization support**
- 💿 **Relational database with ORM**
- ✅ **Comprehensive integration tests**
- 🔄 **CI/CD pipeline** with automated testing and deployment

## Tech Stack

| Category | Technology |
|----------|------------|
| **Language** | [Go 1.24.1](https://github.com/golang/go) |
| **Web Framework** | [Gin](https://github.com/gin-gonic/gin) |
| **ORM** | [GORM](https://github.com/go-gorm/gorm) |
| **Database** | [SQLite](https://github.com/sqlite/sqlite) |
| **API Documentation** | [Swagger/OpenAPI](https://github.com/swaggo/swag) |
| **Caching** | [gin-contrib/cache](https://github.com/gin-contrib/cache) |
| **Testing** | [testify](https://github.com/stretchr/testify) |
| **Containerization** | [Docker](https://github.com/docker) & [Docker Compose](https://github.com/docker/compose) |

## Project Structure

```text
/
├── main.go                 # Entry point: DB connection, route setup, server start
├── controller/             # HTTP handlers (request/response logic)
│   └── player_controller.go
├── service/                # Business logic (ORM interactions)
│   └── player_service.go
├── route/                  # Route configuration and middleware
│   ├── player_route.go     # Route setup with caching middleware
│   └── path.go             # Path constants
├── model/                  # Data structures
│   └── player_model.go
├── data/                   # Database connection
│   └── player_data.go
├── swagger/                # Swagger configuration
│   └── swagger.go
├── docs/                   # Auto-generated Swagger docs (DO NOT EDIT)
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── tests/                  # Integration tests
│   ├── main_test.go
│   ├── player_fake.go
│   └── players.json
├── storage/                # SQLite database file (pre-seeded)
├── scripts/                # Container entrypoint & healthcheck
└── .github/workflows/      # CI/CD pipelines
```

## Architecture

Layered architecture with clear separation of concerns. Dependencies flow from data layer through services and controllers to routes. External dependencies (Gin and GORM) integrate at their respective layers. Integration tests (dotted lines) validate the complete application.

```mermaid
%%{init: {'theme':'base'}}%%
graph LR
    %% Core Application Flow
    main[main.go]
    route[route]
    controller[controller]
    service[service]
    data[data]
    model[model]

    %% External Dependencies
    Gin[Gin]
    GORM[GORM]

    %% Tests and Documentation
    tests[tests]
    docs[docs]

    %% Main Application Flow
    data --> service
    service --> controller
    controller --> route
    route --> main

    %% Model connections
    model --> data
    model --> controller

    %% External Dependencies connections
    GORM --> data
    Gin --> controller
    Gin --> route

    %% Tests and Documentation connections
    main -.-> tests
    docs --> route

    %% Styling
    classDef coreModule fill:#b3d9ff,stroke:#333333,stroke-width:2px
    classDef dependency fill:#ffcccc,stroke:#333333,stroke-width:2px
    classDef testModule fill:#ccffcc,stroke:#333333,stroke-width:2px
    classDef docsModule fill:#ffffcc,stroke:#333333,stroke-width:2px

    class main,route,controller,service,data,model coreModule
    class Gin,GORM dependency
    class tests testModule
    class docs docsModule
```

## API Reference

Interactive API documentation is available via Swagger UI at `http://localhost:9000/swagger/index.html` when the server is running.

> 💡 **Note:** The Swagger documentation is automatically generated from code annotations using [swaggo/swag](https://github.com/swaggo/swag). To regenerate after making changes, run `swag init`.

**Quick Reference:**

- `GET /players` - List all players
- `GET /players/:id` - Get player by ID
- `GET /players/squadnumber/:squadnumber` - Get player by squad number
- `POST /players` - Create new player
- `PUT /players/:id` - Update player
- `DELETE /players/:id` - Remove player
- `GET /health` - Health check

For complete endpoint documentation with request/response schemas, explore the [interactive Swagger UI](http://localhost:9000/swagger/index.html). You can also access the OpenAPI JSON specification at `http://localhost:9000/swagger.json`.

### Postman Collection

A pre-configured Postman collection is available at [`postman-collections/go-samples-gin-restful.postman_collection.json`](postman-collections/go-samples-gin-restful.postman_collection.json).

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.24.0 or higher**
- **Docker & Docker Compose** (optional, for containerized deployment)

## Quick Start

### Clone the repository

```bash
git clone https://github.com/nanotaboada/go-samples-gin-restful.git
cd go-samples-gin-restful
```

### Install dependencies

```bash
go mod download
```

### Start the development server

```bash
go run .
```

The server will start on `http://localhost:9000`.

### Access the application

- **API:** `http://localhost:9000`
- **Swagger Documentation:** `http://localhost:9000/swagger/index.html`
- **Health Check:** `http://localhost:9000/health`

## Testing

Run the test suite with coverage:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v ./... -coverprofile=coverage.out

# Run tests with detailed coverage for specific packages
go test -v ./... \
  -coverpkg=github.com/nanotaboada/go-samples-gin-restful/service,github.com/nanotaboada/go-samples-gin-restful/controller,github.com/nanotaboada/go-samples-gin-restful/route \
  -covermode=atomic \
  -coverprofile=coverage.out

# View coverage report
go tool cover -html=coverage.out
```

Tests are located in the `tests/` directory and use testify for integration testing. Coverage reports are generated for controllers, services, and routes only.

**Coverage targets:** 80% minimum for service, controller, and route packages.

## Docker

This project includes full Docker support with multi-stage builds and Docker Compose for easy deployment.

### Build the Docker image

```bash
docker compose build
```

### Start the application

```bash
docker compose up
```

> 💡 **Note:** On first run, the container copies a pre-seeded SQLite database into a persistent volume. On subsequent runs, that volume is reused and the data is preserved.

### Stop the application

```bash
docker compose down
```

### Reset the database

To remove the volume and reinitialize the database from the built-in seed file:

```bash
docker compose down -v
```

The containerized application runs on port 9000 and includes health checks that monitor the `/health` endpoint every 30 seconds.

## Environment Variables

The application can be configured using the following environment variables (declared in [`compose.yaml`](https://github.com/nanotaboada/go-samples-gin-restful/blob/master/compose.yaml)):

```bash
# Database storage path (default: ./storage/players-sqlite3.db)
# In Docker: /storage/players-sqlite3.db
STORAGE_PATH=./storage/players-sqlite3.db

# Gin framework mode: debug, release, or test (default: debug)
# In Docker: release
GIN_MODE=release
```

## Command Summary

| Command | Description |
|---------|-------------|
| `go run .` | Start development server |
| `go build` | Build the application |
| `go test ./...` | Run all tests |
| `go test -v ./... -coverprofile=coverage.out` | Run tests with coverage report |
| `go fmt ./...` | Format code |
| `go mod tidy` | Clean up dependencies |
| `swag init` | Regenerate Swagger documentation |
| `docker compose build` | Build Docker image |
| `docker compose up` | Start Docker container |
| `docker compose down` | Stop Docker container |
| `docker compose down -v` | Stop and remove Docker volume |

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on the code of conduct and the process for submitting pull requests.

**Key guidelines:**

- Follow [Conventional Commits](https://www.conventionalcommits.org/) for commit messages
- Ensure all tests pass (`go test ./...`)
- Run `go fmt` before committing
- Keep changes small and focused

## Legal

This is a proof-of-concept project intended for educational and demonstration purposes. All trademarks, registered trademarks, service marks, product names, company names, or logos mentioned are the property of their respective owners and are used for identification purposes only.
