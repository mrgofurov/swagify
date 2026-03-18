// Auth example demonstrating JWT and API key security documentation with swagify.
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/swagify"
	"github.com/swagify/core"
)

// --- Models ---

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" description:"User email" example:"admin@example.com"`
	Password string `json:"password" validate:"required,min=8" description:"User password" example:"secretpassword"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token" description:"JWT access token" example:"eyJhbGciOiJIUzI1NiIs..."`
	RefreshToken string `json:"refresh_token" description:"JWT refresh token" example:"eyJhbGciOiJIUzI1NiIs..."`
	TokenType    string `json:"token_type" description:"Token type" example:"Bearer"`
	ExpiresIn    int    `json:"expires_in" description:"Token expiry in seconds" example:"3600"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" description:"Refresh token"`
}

type UserProfile struct {
	ID    int    `json:"id" description:"User ID" example:"1"`
	Name  string `json:"name" description:"Full name" example:"Admin User"`
	Email string `json:"email" description:"Email" example:"admin@example.com"`
	Role  string `json:"role" description:"User role" example:"admin"`
}

type ErrorResponse struct {
	Error   string `json:"error" description:"Error message"`
	Code    int    `json:"code" description:"HTTP status code"`
}

// --- Handlers ---

func login(c *fiber.Ctx) error {
	return c.JSON(TokenResponse{
		AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	})
}

func refreshToken(c *fiber.Ctx) error {
	return c.JSON(TokenResponse{
		AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	})
}

func getProfile(c *fiber.Ctx) error {
	return c.JSON(UserProfile{
		ID:    1,
		Name:  "Admin User",
		Email: "admin@example.com",
		Role:  "admin",
	})
}

func updateProfile(c *fiber.Ctx) error {
	return c.JSON(UserProfile{
		ID:    1,
		Name:  "Updated Admin",
		Email: "admin@example.com",
		Role:  "admin",
	})
}

func main() {
	app := fiber.New()
	app.Use(cors.New())

	api := swagify.NewFiber(app, swagify.FiberConfig{
		Info: &core.Info{
			Title:       "Auth API",
			Description: "An authentication API demonstrating JWT bearer and API key security documentation.",
			Version:     "1.0.0",
		},
		Servers: []core.Server{
			{URL: "http://localhost:8083", Description: "Local development"},
		},
		SecuritySchemes: map[string]core.SecurityScheme{
			"bearerAuth": {
				Type:         "http",
				Scheme:       "bearer",
				BearerFormat: "JWT",
				Description:  "JWT Bearer token authentication. Use the /auth/login endpoint to obtain a token.",
			},
			"apiKeyAuth": {
				Type:        "apiKey",
				Name:        "X-API-Key",
				In:          "header",
				Description: "API key authentication via the X-API-Key header.",
			},
		},
	})

	api.AddTag("Authentication", "Login and token management")
	api.AddTag("Profile", "User profile management (requires authentication)")

	// Public routes (no auth required)
	api.POST("/auth/login", login,
		swagify.Summary("User login"),
		swagify.Description("Authenticates a user and returns JWT tokens."),
		swagify.Tags("Authentication"),
		swagify.WithRequest(LoginRequest{}),
		swagify.WithResponse(TokenResponse{}),
		swagify.ErrorResponse(401, ErrorResponse{}, "Invalid credentials"),
	)

	api.POST("/auth/refresh", refreshToken,
		swagify.Summary("Refresh token"),
		swagify.Description("Exchanges a refresh token for a new access token."),
		swagify.Tags("Authentication"),
		swagify.WithRequest(RefreshRequest{}),
		swagify.WithResponse(TokenResponse{}),
		swagify.ErrorResponse(401, ErrorResponse{}, "Invalid refresh token"),
	)

	// Protected routes (bearer auth required)
	api.GET("/profile", getProfile,
		swagify.Summary("Get current user profile"),
		swagify.Description("Returns the profile of the currently authenticated user."),
		swagify.Tags("Profile"),
		swagify.WithResponse(UserProfile{}),
		swagify.SecurityBearer(),
		swagify.ErrorResponse(401, ErrorResponse{}, "Unauthorized"),
	)

	api.PUT("/profile", updateProfile,
		swagify.Summary("Update profile"),
		swagify.Description("Updates the current user's profile information."),
		swagify.Tags("Profile"),
		swagify.WithRequest(UserProfile{}),
		swagify.WithResponse(UserProfile{}),
		swagify.SecurityBearer(),
		swagify.ErrorResponse(401, ErrorResponse{}, "Unauthorized"),
		swagify.ErrorResponse(400, ErrorResponse{}, "Invalid request"),
	)

	api.RegisterOpenAPI("/openapi.json")
	api.RegisterDocs("/docs")

	log.Println("🚀 Server starting on http://localhost:8083")
	log.Println("📖 API Docs: http://localhost:8083/docs")
	log.Fatal(app.Listen(":8083"))
}
