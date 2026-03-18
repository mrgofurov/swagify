// Fiber basic CRUD example demonstrating swagify with standard Fiber handlers.
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
	Age   int    `json:"age,omitempty" validate:"min=1,max=150" description:"Age of the user" example:"30"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty" description:"Full name of the user" example:"Jane Doe"`
	Email *string `json:"email,omitempty" validate:"email" description:"Email address" example:"jane@example.com"`
	Age   *int    `json:"age,omitempty" validate:"min=1,max=150" description:"Age of the user" example:"25"`
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
	Error   string `json:"error" description:"Error message" example:"Resource not found"`
	Code    int    `json:"code" description:"HTTP status code" example:"404"`
	Details string `json:"details,omitempty" description:"Additional error details"`
}

type MessageResponse struct {
	Message string `json:"message" description:"Response message" example:"Operation successful"`
}

// --- Query Params ---

type ListUsersQuery struct {
	Page  int    `json:"page" description:"Page number" example:"1"`
	Limit int    `json:"limit" description:"Items per page" example:"20"`
	Sort  string `json:"sort,omitempty" description:"Sort field" example:"name"`
}

// --- Handlers ---

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
	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(ErrorResponse{Error: "Invalid request", Code: 400})
	}
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
	app := fiber.New(fiber.Config{
		AppName: "Swagify Fiber Basic Example",
	})

	app.Use(cors.New())

	// Create swagify adapter with configuration
	api := swagify.NewFiber(app, swagify.FiberConfig{
		Info: &core.Info{
			Title:       "User Management API",
			Description: "A comprehensive user management API built with Fiber and documented with Swagify.",
			Version:     "1.0.0",
			Contact: &core.Contact{
				Name:  "API Support",
				Email: "support@example.com",
			},
			License: &core.License{
				Name: "MIT",
				URL:  "https://opensource.org/licenses/MIT",
			},
		},
		Servers: []core.Server{
			{URL: "http://localhost:8080", Description: "Local development"},
		},
	})

	// Add tags
	api.AddTag("Users", "User management operations")

	// Register routes with documentation
	api.GET("/users", listUsers,
		swagify.Summary("List all users"),
		swagify.Description("Returns a paginated list of all users in the system."),
		swagify.Tags("Users"),
		swagify.WithResponse(UserListResponse{}),
		swagify.QueryParams(ListUsersQuery{}),
		swagify.ErrorResponse(500, ErrorResponse{}, "Internal server error"),
	)

	api.GET("/users/:id", getUser,
		swagify.Summary("Get user by ID"),
		swagify.Description("Returns a single user by their unique identifier."),
		swagify.Tags("Users"),
		swagify.WithResponse(UserResponse{}),
		swagify.ErrorResponse(404, ErrorResponse{}, "User not found"),
	)

	api.POST("/users", createUser,
		swagify.Summary("Create a new user"),
		swagify.Description("Creates a new user with the provided information."),
		swagify.Tags("Users"),
		swagify.WithRequest(CreateUserRequest{}),
		swagify.WithResponse(UserResponse{}),
		swagify.SuccessStatus(201),
		swagify.ErrorResponse(400, ErrorResponse{}, "Invalid request body"),
		swagify.ErrorResponse(422, ErrorResponse{}, "Validation failed"),
	)

	api.PUT("/users/:id", updateUser,
		swagify.Summary("Update a user"),
		swagify.Description("Updates an existing user's information."),
		swagify.Tags("Users"),
		swagify.WithRequest(UpdateUserRequest{}),
		swagify.WithResponse(UserResponse{}),
		swagify.ErrorResponse(400, ErrorResponse{}, "Invalid request body"),
		swagify.ErrorResponse(404, ErrorResponse{}, "User not found"),
	)

	api.DELETE("/users/:id", deleteUser,
		swagify.Summary("Delete a user"),
		swagify.Description("Permanently deletes a user from the system."),
		swagify.Tags("Users"),
		swagify.WithResponse(MessageResponse{}),
		swagify.ErrorResponse(404, ErrorResponse{}, "User not found"),
	)

	// Register OpenAPI spec and docs UI
	api.RegisterOpenAPI("/openapi.json")
	api.RegisterDocs("/docs")

	log.Println("🚀 Server starting on http://localhost:8080")
	log.Println("📖 API Docs: http://localhost:8080/docs")
	log.Println("📋 OpenAPI Spec: http://localhost:8080/openapi.json")
	log.Fatal(app.Listen(":8080"))
}
