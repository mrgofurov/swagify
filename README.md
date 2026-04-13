# Swagify

**Swagify** is a modern, code-first API framework for Go that provides a developer experience closer to **FastAPI**, **Django REST Swagger**, and other modern API tooling — without the need for ugly Swagger comments, annotation-heavy setup, or massive external dependencies.

With Swagify, you define your routes, request models, and response models in plain Go, and it automatically generates and serves:

- Complete OpenAPI 3.1 JSON
- Interactive API docs UI (Swagger)
- Component schemas directly from Go structs
- Request/response documentation
- Query, path, and header parameter documentation

## Why Swagify?

Traditional Swagger tooling in Go often relies on comment-based generation:

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
That works, but it is verbose, fragile, and unpleasant to maintain.

Previously, Swagify required Gin or Fiber to run. **Not anymore.** Swagify is now its own standalone API framework built purely on Go's standard library (`net/http`). 

- **Zero dependencies** (just Go `1.22+`)
- **Code-first and type-driven**
- **Automatic schema inference**
- **Extremely fast** and simple to use

## Features

- **Zero external dependencies**: Built directly on Go's `net/http`
- **FastAPI-like Developer Experience**: Handlers naturally declare their request and response types, and Swagify automatically maps them to HTTP logic and OpenAPI specifications.
- **OpenAPI 3.1 generation** out-of-the-box
- **Automatic Route Discovery**: Add a route, and its schema, path, method, and tags are documented instantly.
- **Modern embedded Docs UI** (`/docs` endpoint)

## Installation

```bash
go get github.com/mrgofurov/swagify
```

## Quick Start

Here is a full API implementation demonstrating how concise and powerful Swagify can be:

```go
package main

import (
	"fmt"
	"log"

	"github.com/mrgofurov/swagify"
	"github.com/mrgofurov/swagify/core"
)

// --- Models ---
// Fields, descriptions, and examples seamlessly propagate to the OpenAPI specs.

type CreateUserRequest struct {
	Name  string `json:"name" description:"Full name" example:"Alice"`
	Email string `json:"email" description:"Email address" example:"alice@example.com"`
}

type User struct {
	ID    int    `json:"id" description:"User ID" example:"1"`
	Name  string `json:"name" description:"Full name" example:"Alice"`
	Email string `json:"email" description:"Email address" example:"alice@example.com"`
}

type ListUsersQuery struct {
	Page  int `json:"page" description:"Page number" example:"1"`
	Limit int `json:"limit" description:"Items per page" example:"20"`
}

type UserList struct {
	Users []User `json:"users"`
	Total int    `json:"total" example:"1"`
}

// --- Handlers ---
// Handlers declare your request query params/bodies and response bodies as native Go types.

func listUsers(ctx *swagify.Ctx, q ListUsersQuery) (UserList, error) {
	return UserList{
		Users: []User{{ID: 1, Name: "Alice", Email: "alice@example.com"}},
		Total: 1,
	}, nil
}

func getUser(ctx *swagify.Ctx) (User, error) {
    // Path params are accessed intuitively
	id := ctx.Param("id")
	return User{ID: 1, Name: fmt.Sprintf("User %s", id), Email: "user@example.com"}, nil
}

func createUser(ctx *swagify.Ctx, req CreateUserRequest) (User, error) {
	return User{ID: 1, Name: req.Name, Email: req.Email}, nil
}

func main() {
	api := swagify.New(swagify.Config{
		Title:       "User API",
		Description: "Simple user management API",
		Version:     "1.0.0",
		Servers:     []core.Server{{URL: "http://localhost:8080"}},
	})

	// All schemas, summaries, and tags are inferred from handlers
	api.GET("/users", listUsers)
	api.GET("/users/{id}", getUser)
	api.POST("/users", createUser)

	log.Println("Listening on http://localhost:8080")
	log.Println("Docs:    http://localhost:8080/docs")
	log.Fatal(api.Run(":8080")) // Start standard library server!
}
```

Open:
- **Docs UI**: http://localhost:8080/docs
- **OpenAPI JSON**: http://localhost:8080/openapi.json

## Handler Signatures

Swagify scans your handler's function signature to automatically determine how to parse requests, serialize responses, and build your API specification.

| Signature | What's Inferred |
|-----------|-----------------|
| `func(http.ResponseWriter, *http.Request)` | Plain native handler, no automatic schema. |
| `func(*swagify.Ctx) error` | No body handling (use `ctx.Param("id")` for path variables). |
| `func(*swagify.Ctx) (Res, error)` | Automatically generates a response schema. |
| `func(*swagify.Ctx, Req) (Res, error)` | Automatically generates request **and** response schemas. |

### Request Parsing Behavior

How `Req` is parsed depends automatically on the HTTP method:

- **`GET` / `DELETE`**: `Req` is parsed and hydrated from **URL Query Parameters**.
- **`POST` / `PUT` / `PATCH`**: `Req` is parsed from the **JSON Body**.
- **Path Parameters**: Extracted directly using `ctx.Param("name")`.

## Manual Overrides

Even though most metadata is inferred naturally from the handlers and structs, you can always enhance any route explicitly using decorators:

```go
api.GET("/users/{id}", getUser,
    swagify.Summary("Get user by ID"),
    swagify.Description("Returns a single user by their unique identifier."),
    swagify.Tags("Users"),
    swagify.ErrorResponse(404, ErrorResponse{}, "User not found"),
)
```