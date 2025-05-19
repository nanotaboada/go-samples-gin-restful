# 🧪 RESTful API with Go and Gin

## Status

[![Go CI](https://github.com/nanotaboada/go-samples-gin-restful/actions/workflows/go.yml/badge.svg)](https://github.com/nanotaboada/go-samples-gin-restful/actions/workflows/go.yml)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=nanotaboada_go-samples-gin-restful&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=nanotaboada_go-samples-gin-restful)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/e2b234e8182d4a0f8efb9da619e1dc26)](https://app.codacy.com/gh/nanotaboada/go-samples-gin-restful/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![codecov](https://codecov.io/gh/nanotaboada/go-samples-gin-restful/graph/badge.svg?token=i37VDcDWwx)](https://codecov.io/gh/nanotaboada/go-samples-gin-restful)
[![Go Report Card](https://goreportcard.com/badge/github.com/nanotaboada/go-samples-gin-restful)](https://goreportcard.com/report/github.com/nanotaboada/go-samples-gin-restful)
[![CodeFactor](https://www.codefactor.io/repository/github/nanotaboada/go-samples-gin-restful/badge)](https://www.codefactor.io/repository/github/nanotaboada/go-samples-gin-restful)
[![codebeat badge](https://codebeat.co/badges/d09d49d6-1fb4-4c7d-8009-aaef12c6b02b)](https://codebeat.co/projects/github-com-nanotaboada-go-samples-gin-restful-master)

## About

Proof of Concept for a REST API made with [Go](https://go.dev/) and [Gin](https://gin-gonic.com/).

## Structure

![Simplified, conceptual project structure and main application flow](assets/images/structure.svg)

_Figure: Simplified, conceptual project structure and main application flow. Not all dependencies are shown._

## Start

```console
go run .
```

## Documentation

```console
http://localhost:9000/swagger/index.html
```

![API Documentation](assets/images/swagger.png)

## Container

### Docker Compose

This setup uses [Docker Compose](https://docs.docker.com/compose/) to build and run the app and manage a persistent SQLite database stored in a Docker volume.

#### Build the image

```bash
docker compose build
```

#### Start the app

```bash
docker compose up
```

> On first run, the container copies a pre-seeded SQLite database into a persistent volume
> On subsequent runs, that volume is reused and the data is preserved

#### Stop the app

```bash
docker compose down
```

#### Optional: database reset

```bash
docker compose down -v
```

> This removes the volume and will reinitialize the database from the built-in seed file the next time you `up`.

## Credits

The solution has been coded using [Visual Studio Code](https://code.visualstudio.com/) with the [Go](https://marketplace.visualstudio.com/items?itemName=golang.go) extension.

## Terms

All trademarks, registered trademarks, service marks, product names, company names, or logos mentioned on this repository are the property of their respective owners. All usage of such terms herein is for identification purposes only and constitutes neither an endorsement nor a recommendation of those items. Furthermore, the use of such terms is intended to be for educational and informational purposes only.
