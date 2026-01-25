# GitHub Copilot Instructions

> **⚡ Token Efficiency Note**: This is a minimal pointer file (~500 tokens, auto-loaded by Copilot).
> For complete operational details, reference: `#file:AGENTS.md` (~2,500 tokens, loaded on-demand)
> For specialized knowledge, use: `#file:SKILLS/<skill-name>/SKILL.md` (loaded on-demand when needed)

## 🎯 Quick Context

**Project**: Gin-based REST API demonstrating idiomatic Go patterns
**Stack**: Go 1.25 • Gin • GORM • SQLite • Docker • testify
**Pattern**: Controller → Service → ORM (layered architecture)
**Philosophy**: Learning-focused PoC emphasizing simplicity and Go best practices

## 📐 Core Conventions

- **Naming**: camelCase (unexported), PascalCase (exported)
- **Error Handling**: Explicit error returns, no panics in handlers
- **Pointers**: Use for structs in handlers, minimize copying
- **Testing**: testify/assert for readable assertions
- **Formatting**: gofmt (standard), goimports for imports

## 🏗️ Architecture at a Glance

```
Controller → Service → GORM → Database
     ↓           ↓
  Cache      Validation
```

- **Controllers**: HTTP handlers with Gin context
- **Services**: Business logic and GORM operations
- **Routes**: Gin router with cache middleware
- **Models**: Structs with JSON/GORM tags
- **Cache**: gin-contrib/cache (10min/1hr TTL)

## ✅ Copilot Should

- Generate idiomatic Go code with proper error handling
- Use GORM APIs correctly (`db.First()`, `db.Create()`, etc.)
- Follow Gin patterns (`c.JSON()`, `c.BindJSON()`)
- Write table-driven tests with testify
- Apply struct tags for JSON/GORM mappings
- Use proper HTTP status codes (http.StatusXXX)
- Handle errors explicitly at every layer

## 🚫 Copilot Should Avoid

- Using `panic()` in production code
- Ignoring errors
- Global variables for state
- Not using pointers for large structs
- Missing error checks on database operations
- Using `fmt.Print` instead of proper logging

## ⚡ Quick Commands

```bash
# Run with hot reload (if installed)
air

# Run normally
go run main.go

# Test
go test ./... -v

# Docker
docker compose up

# Swagger: http://localhost:9000/swagger/index.html
```

## 📚 Need More Detail?

**For operational procedures**: Load `#file:AGENTS.md`
**For Docker expertise**: *(Planned)* `#file:SKILLS/docker-containerization/SKILL.md`
**For testing patterns**: *(Planned)* `#file:SKILLS/testing-patterns/SKILL.md`

---

💡 **Why this structure?** Copilot auto-loads this file on every chat (~500 tokens). Loading `AGENTS.md` or `SKILLS/` explicitly gives you deep context only when needed, saving 80% of your token budget!
