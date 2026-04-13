// Basic CRUD example — FastAPI-style with automatic schema inference.
// Run: go run ./examples/basic
// Docs: http://localhost:8080/docs
package main

import (
	"fmt"
	"log"

	"github.com/mrgofurov/swagify"
	"github.com/mrgofurov/swagify/core"
)

// --- Models ---

type CreateUserRequest struct {
	Name  string `json:"name" description:"Full name" example:"Alice"`
	Email string `json:"email" description:"Email address" example:"alice@example.com"`
	Age   int    `json:"age,omitempty" description:"Age" example:"30"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty" description:"Full name" example:"Bob"`
	Email *string `json:"email,omitempty" description:"Email address" example:"bob@example.com"`
}

type User struct {
	ID    int    `json:"id" description:"User ID" example:"1"`
	Name  string `json:"name" description:"Full name" example:"Alice"`
	Email string `json:"email" description:"Email address" example:"alice@example.com"`
	Age   int    `json:"age" description:"Age" example:"30"`
}

type UserList struct {
	Users []User `json:"users"`
	Total int    `json:"total" example:"1"`
}

type ListUsersQuery struct {
	Page  int    `json:"page" description:"Page number" example:"1"`
	Limit int    `json:"limit" description:"Items per page" example:"20"`
	Sort  string `json:"sort,omitempty" description:"Sort field" example:"name"`
}

// --- Handlers ---

func listUsers(ctx *swagify.Ctx, q ListUsersQuery) (UserList, error) {
	return UserList{
		Users: []User{{ID: 1, Name: "Alice", Email: "alice@example.com", Age: 30}},
		Total: 1,
	}, nil
}

func getUser(ctx *swagify.Ctx) (User, error) {
	id := ctx.Param("id")
	return User{ID: 1, Name: fmt.Sprintf("User %s", id), Email: "user@example.com"}, nil
}

func createUser(ctx *swagify.Ctx, req CreateUserRequest) (User, error) {
	return User{ID: 1, Name: req.Name, Email: req.Email, Age: req.Age}, nil
}

func updateUser(ctx *swagify.Ctx, req UpdateUserRequest) (User, error) {
	return User{ID: 1, Name: "Updated", Email: "updated@example.com"}, nil
}

func deleteUser(ctx *swagify.Ctx) error {
	return nil
}

func main() {
	api := swagify.New(swagify.Config{
		Title:       "User API",
		Description: "Simple user management API",
		Version:     "1.0.0",
		Servers:     []core.Server{{URL: "http://localhost:8080"}},
	})

	// All schemas, summaries, and tags are inferred from the handler signatures.
	// Add opts only when you need to override something.
	api.GET("/users", listUsers)
	api.GET("/users/{id}", getUser)
	api.POST("/users", createUser)
	api.PUT("/users/{id}", updateUser)
	api.DELETE("/users/{id}", deleteUser)

	log.Println("Listening on http://localhost:8080")
	log.Println("Docs:    http://localhost:8080/docs")
	log.Fatal(api.Run(":8080"))
}
