# Swagify

**Swagify** is a modern, code-first OpenAPI 3.1 documentation package for Go web frameworks such as **Fiber** and **Gin**.

It helps Go developers build beautiful API documentation with a developer experience closer to **FastAPI**, **Django REST Swagger**, and other modern API tooling — without ugly Swagger comments or annotation-heavy setup.

With Swagify, you define your routes, request models, and response models in Go, and it generates:

- OpenAPI 3.1 JSON
- interactive API docs UI
- component schemas
- request/response documentation
- query, path, and header parameter documentation
- tags, summaries, descriptions, and security metadata

## Why Swagify?

Traditional Swagger tooling in Go often relies on comment-based generation like this:

```go
// @Summary Create user
// @Description create user
// @Tags users
// @Accept json
// @Produce json
// @Param body body CreateUserRequest true "User data"
// @Success 200 {object} User
// @Router /users [post]
```
That works, but it is verbose, fragile, and not very pleasant to maintain.

Swagify takes a different approach:

- code-first
- type-driven
- clean route options
- minimal boilerplate
- modern docs UI
- framework adapters for Fiber and Gin

## Features

- OpenAPI 3.1 generation
- Fiber support
- Gin support
- route-level summaries, descriptions, tags, and operation IDs
- request and response model documentation
- query, path, and header parameter documentation
- custom success and error responses
- security documentation helpers
- modern embedded docs UI
- no filesystem-relative static docs dependency
- no Swagger comments required

## Installation

```bash
go get github.com/mrgofurov/swagify
```

Replace the module path above with your real repository path.

## Quick Start

### Fiber

```go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mrgofurov/swagify"
	"github.com/mrgofurov/swagify/core"
)

type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func getUser(c *fiber.Ctx) error {
	id := c.Params("id")

	return c.JSON(UserResponse{
		ID:    1,
		Name:  "Ali",
		Email: "ali@example.com",
		_ = id,
	})
}

func main() {
	app := fiber.New()

	api := swagify.NewFiber(app, swagify.FiberConfig{
		Info: &core.Info{
			Title:       "User API",
			Description: "Example API documented with Swagify",
			Version:     "1.0.0",
		},
	})

	api.GET("/users/:id", getUser,
		swagify.Summary("Get user by ID"),
		swagify.Description("Returns a single user by its unique identifier."),
		swagify.Tags("Users"),
		swagify.WithResponse(UserResponse{}),
		swagify.ErrorResponse(404, ErrorResponse{}, "User not found"),
	)

	api.RegisterOpenAPI("/openapi.json")
	api.RegisterDocs("/docs")

	log.Fatal(app.Listen(":8080"))
}
```

Open:

Docs UI: http://localhost:8080/docs

OpenAPI JSON: http://localhost:8080/openapi.json

### Gin

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrgofurov/swagify"
)

type UserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func getUser(c *gin.Context) {
	c.JSON(http.StatusOK, UserResponse{
		ID:   1,
		Name: "Ali",
	})
}

func main() {
	r := gin.Default()

	api := swagify.NewGin(r)

	api.GET("/users/:id", getUser,
		swagify.Summary("Get user by ID"),
		swagify.Tags("Users"),
		swagify.WithResponse(UserResponse{}),
	)

	api.RegisterOpenAPI("/openapi.json")
	api.RegisterDocs("/docs")

	r.Run(":8080")
}
```
## Basic Usage

### Documenting a GET endpoint

```go
api.GET("/users/:id", getUser,
	swagify.Summary("Get user by ID"),
	swagify.Description("Returns a single user by their unique identifier."),
	swagify.Tags("Users"),
	swagify.WithResponse(UserResponse{}),
	swagify.ErrorResponse(404, ErrorResponse{}, "User not found"),
)
```

### Documenting a POST endpoint

```go
api.POST("/users", createUser,
	swagify.Summary("Create a new user"),
	swagify.Description("Creates a new user with the provided information."),
	swagify.Tags("Users"),
	swagify.WithRequest(CreateUserRequest{}),
	swagify.WithResponse(UserResponse{}),
	swagify.SuccessStatus(201),
	swagify.ErrorResponse(400, ErrorResponse{}, "Invalid request body"),
)
```

## Auto-Discovery (Zero Migration)

Already have a Fiber or Gin project with hundreds of routes? No need to rewrite them. Use `Discover()` to auto-generate documentation from your existing routes:

### Fiber

```go
// Your existing app — nothing changes here
app := fiber.New()
app.Get("/users", listUsers)
app.Get("/users/:id", getUser)
app.Post("/users", createUser)
app.Put("/users/:id", updateUser)
app.Delete("/users/:id", deleteUser)

// Attach Swagify and auto-discover all routes
api := swagify.NewFiber(app, swagify.FiberConfig{
    Info: &core.Info{
        Title:   "My API",
        Version: "1.0.0",
    },
})
api.Discover()

api.RegisterOpenAPI()
api.RegisterDocs()
log.Fatal(app.Listen(":8080"))
```

### Gin

```go
r := gin.Default()
r.GET("/products", listProducts)
r.POST("/products", createProduct)

api := swagify.NewGin(r)
api.Discover()

api.RegisterOpenAPI()
api.RegisterDocs()
r.Run(":8080")
```

### Enriching Discovered Routes

After discovery, you can optionally enrich specific routes with request/response types and metadata:

```go
api.Discover()

// Add rich documentation to specific routes
api.Enrich("GET /users",
    swagify.Summary("List all users"),
    swagify.Tags("Users"),
    swagify.WithResponse(UserListResponse{}),
    swagify.QueryParams(ListUsersQuery{}),
)

api.Enrich("POST /users",
    swagify.Summary("Create a new user"),
    swagify.Tags("Users"),
    swagify.WithRequest(CreateUserRequest{}),
    swagify.WithResponse(UserResponse{}),
)
```

### Discovery Options

Control which routes are discovered:

```go
api.Discover(swagify.DiscoverOptions{
    // Only document routes under /api
    IncludePaths: []string{"/api"},

    // Skip health/metrics endpoints
    ExcludePaths: []string{"/internal", "/metrics"},
})
```

### What Discover Auto-Generates

Even without enrichment, `Discover()` generates useful documentation:

| Route | Auto Summary | Auto Tag |
|-------|-------------|----------|
| `GET /users` | List users | Users |
| `GET /users/:id` | Get user by id | Users |
| `POST /users` | Create user | Users |
| `PUT /users/:id` | Update user by id | Users |
| `DELETE /users/:id` | Delete user by id | Users |
| `GET /api/v1/orders` | List orders | Orders |

## Protecting Docs with BasicAuth

Protect your `/docs` and `/openapi.json` endpoints with HTTP Basic Authentication:

```go
api := swagify.NewFiber(app)

// Simple — just username and password
api.BasicAuth("admin", "secret123")

api.RegisterOpenAPI()  // protected
api.RegisterDocs()     // protected
```

With a custom realm (shown in the browser prompt):

```go
api.BasicAuth("admin", "secret123", swagify.DocsAuthConfig{
    Realm: "My API Documentation",
})
```

> **Note:** `BasicAuth()` must be called **before** `RegisterOpenAPI()` and `RegisterDocs()`. It only protects the docs endpoints — your API routes remain unaffected.