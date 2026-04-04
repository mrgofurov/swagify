// Example: Fiber Auto-Discovery
//
// This example shows how to use Swagify's Discover() feature to automatically
// generate API documentation for an existing Fiber app without migrating routes.
package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/swagify"
	"github.com/swagify/core"
)

// --- Models ---

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required" description:"Full name of the user" example:"John Doe"`
	Email string `json:"email" validate:"required,email" description:"Email address" example:"john@example.com"`
	Age   int    `json:"age,omitempty" description:"Age of the user" example:"30"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty" description:"Full name of the user" example:"Jane Doe"`
	Email *string `json:"email,omitempty" description:"Email address" example:"jane@example.com"`
}

type UserResponse struct {
	ID    int    `json:"id" description:"Unique user identifier" example:"1"`
	Name  string `json:"name" description:"Full name" example:"John Doe"`
	Email string `json:"email" description:"Email address" example:"john@example.com"`
	Age   int    `json:"age" description:"Age of the user" example:"30"`
}

type UserListResponse struct {
	Users []UserResponse `json:"users" description:"List of users"`
	Total int            `json:"total" description:"Total count" example:"42"`
}

type ErrorResponse struct {
	Error string `json:"error" description:"Error message" example:"Resource not found"`
	Code  int    `json:"code" description:"HTTP status code" example:"404"`
}

type ListUsersQuery struct {
	Page  int    `json:"page" description:"Page number" example:"1"`
	Limit int    `json:"limit" description:"Items per page" example:"20"`
	Sort  string `json:"sort,omitempty" description:"Sort field" example:"name"`
}

// --- Existing Handlers (unchanged from your original code) ---

func listUsers(c *fiber.Ctx) error {
	return c.JSON(UserListResponse{
		Users: []UserResponse{
			{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30},
			{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 25},
		},
		Total: 2,
	})
}

func getUser(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(UserResponse{
		ID:    1,
		Name:  fmt.Sprintf("User %s", id),
		Email: "user@example.com",
		Age:   30,
	})
}

func createUser(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(ErrorResponse{Error: "Invalid request", Code: 400})
	}
	return c.Status(201).JSON(UserResponse{
		ID:    1,
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	})
}

func updateUser(c *fiber.Ctx) error {
	return c.JSON(UserResponse{
		ID:    1,
		Name:  "Updated User",
		Email: "updated@example.com",
		Age:   25,
	})
}

func deleteUser(c *fiber.Ctx) error {
	return c.SendStatus(204)
}

func main() {
	// ====================================================
	// Step 1: Your existing Fiber app (nothing changes here)
	// ====================================================
	app := fiber.New(fiber.Config{
		AppName: "Swagify Fiber Discover Example",
	})
	app.Use(cors.New())

	// These are your EXISTING routes — no migration needed!
	app.Get("/users", listUsers)
	app.Get("/users/:id", getUser)
	app.Post("/users", createUser)
	app.Put("/users/:id", updateUser)
	app.Delete("/users/:id", deleteUser)

	// ====================================================
	// Step 2: Attach Swagify and discover routes
	// ====================================================
	api := swagify.NewFiber(app, swagify.FiberConfig{
		Info: &core.Info{
			Title:       "User Management API",
			Description: "API documented automatically using Swagify Discover.",
			Version:     "1.0.0",
		},
		Servers: []core.Server{
			{URL: "http://localhost:8080", Description: "Local development"},
		},
	})
	api.BasicAuth("admin", "admin123")
	api.AddTag("Users", "User management operations")

	// Auto-discover all existing routes
	api.Discover()

	// ====================================================
	// Step 3 (optional): Enrich specific routes with types
	// ====================================================
	// api.Enrich("GET /users",
	// 	swagify.Summary("List all users"),
	// 	swagify.Description("Returns a paginated list of all users in the system."),
	// 	swagify.Tags("Users"),
	// 	swagify.WithResponse(UserListResponse{}),
	// 	swagify.QueryParams(ListUsersQuery{}),
	// 	swagify.ErrorResponse(500, ErrorResponse{}, "Internal server error"),
	// )

	// api.Enrich("GET /users/:id",
	// 	swagify.Summary("Get user by ID"),
	// 	swagify.Description("Returns a single user by their unique identifier."),
	// 	swagify.Tags("Users"),
	// 	swagify.WithResponse(UserResponse{}),
	// 	swagify.ErrorResponse(404, ErrorResponse{}, "User not found"),
	// )

	// api.Enrich("POST /users",
	// 	swagify.Summary("Create a new user"),
	// 	swagify.Tags("Users"),
	// 	swagify.WithRequest(CreateUserRequest{}),
	// 	swagify.WithResponse(UserResponse{}),
	// 	swagify.SuccessStatus(201),
	// 	swagify.ErrorResponse(400, ErrorResponse{}, "Invalid request body"),
	// )

	// api.Enrich("PUT /users/:id",
	// 	swagify.Summary("Update a user"),
	// 	swagify.Tags("Users"),
	// 	swagify.WithRequest(UpdateUserRequest{}),
	// 	swagify.WithResponse(UserResponse{}),
	// 	swagify.ErrorResponse(404, ErrorResponse{}, "User not found"),
	// )

	// api.Enrich("DELETE /users/:id",
	// 	swagify.Summary("Delete a user"),
	// 	swagify.Tags("Users"),
	// 	swagify.ErrorResponse(404, ErrorResponse{}, "User not found"),
	// )

	// Register OpenAPI spec and docs UI
	api.RegisterOpenAPI("/openapi.json")
	api.RegisterDocs("/docs")

	log.Println("🚀 Server starting on http://localhost:8080")
	log.Println("📖 API Docs: http://localhost:8080/docs")
	log.Println("📋 OpenAPI Spec: http://localhost:8080/openapi.json")
	log.Fatal(app.Listen(":8080"))
}
