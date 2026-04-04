// Example: Gin Auto-Discovery
//
// This example shows how to use Swagify's Discover() feature to automatically
// generate API documentation for an existing Gin app without migrating routes.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swagify"
	"github.com/swagify/core"
)

// --- Models ---

type CreateProductRequest struct {
	Name     string  `json:"name" validate:"required" description:"Product name" example:"Wireless Mouse"`
	Price    float64 `json:"price" validate:"required,min=0" description:"Product price" example:"29.99"`
	Category string  `json:"category" description:"Product category" example:"Electronics"`
}

type ProductResponse struct {
	ID       int     `json:"id" description:"Product ID" example:"1"`
	Name     string  `json:"name" description:"Product name" example:"Wireless Mouse"`
	Price    float64 `json:"price" description:"Product price" example:"29.99"`
	Category string  `json:"category" description:"Product category" example:"Electronics"`
}

type ProductListResponse struct {
	Products []ProductResponse `json:"products" description:"List of products"`
	Total    int               `json:"total" description:"Total count" example:"10"`
}

type ErrorResponse struct {
	Error string `json:"error" description:"Error message" example:"Not found"`
}

// --- Existing Handlers (unchanged from your original code) ---

func listProducts(c *gin.Context) {
	c.JSON(http.StatusOK, ProductListResponse{
		Products: []ProductResponse{
			{ID: 1, Name: "Wireless Mouse", Price: 29.99, Category: "Electronics"},
			{ID: 2, Name: "Keyboard", Price: 49.99, Category: "Electronics"},
		},
		Total: 2,
	})
}

func getProduct(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, ProductResponse{
		ID:       1,
		Name:     fmt.Sprintf("Product %s", id),
		Price:    29.99,
		Category: "Electronics",
	})
}

func createProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ProductResponse{
		ID:       1,
		Name:     req.Name,
		Price:    req.Price,
		Category: req.Category,
	})
}

func deleteProduct(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func main() {
	// ====================================================
	// Step 1: Your existing Gin app (nothing changes here)
	// ====================================================
	r := gin.Default()

	// These are your EXISTING routes — no migration needed!
	r.GET("/products", listProducts)
	r.GET("/products/:id", getProduct)
	r.POST("/products", createProduct)
	r.DELETE("/products/:id", deleteProduct)

	// ====================================================
	// Step 2: Attach Swagify and discover routes
	// ====================================================
	api := swagify.NewGin(r, swagify.GinConfig{
		Info: &core.Info{
			Title:       "Product Catalog API",
			Description: "API documented automatically using Swagify Discover.",
			Version:     "1.0.0",
		},
		Servers: []core.Server{
			{URL: "http://localhost:8080", Description: "Local development"},
		},
	})

	// Auto-discover all existing routes
	api.Discover()

	// ====================================================
	// Step 3 (optional): Enrich specific routes with types
	// ====================================================
	api.Enrich("GET /products",
		swagify.Summary("List all products"),
		swagify.Tags("Products"),
		swagify.WithResponse(ProductListResponse{}),
	)

	api.Enrich("POST /products",
		swagify.Summary("Create a product"),
		swagify.Tags("Products"),
		swagify.WithRequest(CreateProductRequest{}),
		swagify.WithResponse(ProductResponse{}),
		swagify.SuccessStatus(201),
	)

	api.Enrich("GET /products/:id",
		swagify.Summary("Get product by ID"),
		swagify.Tags("Products"),
		swagify.WithResponse(ProductResponse{}),
		swagify.ErrorResponse(404, ErrorResponse{}, "Product not found"),
	)

	api.Enrich("DELETE /products/:id",
		swagify.Summary("Delete a product"),
		swagify.Tags("Products"),
		swagify.ErrorResponse(404, ErrorResponse{}, "Product not found"),
	)

	// Register OpenAPI spec and docs UI
	api.RegisterOpenAPI("/openapi.json")
	api.RegisterDocs("/docs")

	log.Println("🚀 Server starting on http://localhost:8080")
	log.Println("📖 API Docs: http://localhost:8080/docs")
	log.Println("📋 OpenAPI Spec: http://localhost:8080/openapi.json")
	r.Run(":8080")
}
