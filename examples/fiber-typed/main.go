// Fiber typed handlers example demonstrating swagify generics-based type inference.
package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/mrgofurov/swagify"
	"github.com/mrgofurov/swagify/core"
	"github.com/mrgofurov/swagify/router"
)

// --- Models ---

type CreateProductRequest struct {
	Name        string   `json:"name" validate:"required" description:"Product name" example:"Gaming Laptop"`
	Description string   `json:"description" description:"Product description" example:"High-performance gaming laptop"`
	Price       float64  `json:"price" validate:"required,min=0" description:"Price in USD" example:"1299.99"`
	Category    string   `json:"category" validate:"required,oneof=electronics clothing food" description:"Product category"`
	Tags        []string `json:"tags,omitempty" description:"Product tags" example:"gaming"`
	InStock     bool     `json:"in_stock" description:"Whether the product is available" example:"true"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty" description:"Product name"`
	Description *string  `json:"description,omitempty" description:"Product description"`
	Price       *float64 `json:"price,omitempty" validate:"min=0" description:"Price in USD"`
	InStock     *bool    `json:"in_stock,omitempty" description:"Whether the product is available"`
}

type ProductResponse struct {
	ID          int       `json:"id" description:"Unique product identifier" example:"1"`
	Name        string    `json:"name" description:"Product name" example:"Gaming Laptop"`
	Description string    `json:"description" description:"Product description"`
	Price       float64   `json:"price" description:"Price in USD" example:"1299.99"`
	Category    string    `json:"category" description:"Product category" example:"electronics"`
	Tags        []string  `json:"tags" description:"Product tags"`
	InStock     bool      `json:"in_stock" description:"Availability status" example:"true"`
	CreatedAt   time.Time `json:"created_at" description:"Creation timestamp"`
	UpdatedAt   time.Time `json:"updated_at" description:"Last update timestamp"`
}

type ProductListResponse struct {
	Products []ProductResponse `json:"products" description:"List of products"`
	Total    int               `json:"total" description:"Total count" example:"100"`
	Page     int               `json:"page" description:"Current page" example:"1"`
	Pages    int               `json:"pages" description:"Total pages" example:"10"`
}

type ErrorResponse struct {
	Error string `json:"error" description:"Error message"`
	Code  int    `json:"code" description:"Error code"`
}

// --- Typed Handlers (request/response types inferred via generics) ---

func createProduct(c *fiber.Ctx, req CreateProductRequest) (ProductResponse, error) {
	return ProductResponse{
		ID:          1,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		Tags:        req.Tags,
		InStock:     req.InStock,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func updateProduct(c *fiber.Ctx, req UpdateProductRequest) (ProductResponse, error) {
	return ProductResponse{
		ID:          1,
		Name:        "Updated Product",
		Description: "Updated description",
		Price:       999.99,
		Category:    "electronics",
		InStock:     true,
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now(),
	}, nil
}

type GetProductQuery struct {
	// Empty — just uses path param
}

func getProduct(c *fiber.Ctx, _ GetProductQuery) (ProductResponse, error) {
	return ProductResponse{
		ID:          1,
		Name:        "Gaming Laptop",
		Description: "A high-performance gaming laptop",
		Price:       1299.99,
		Category:    "electronics",
		Tags:        []string{"gaming", "laptop"},
		InStock:     true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

type ListProductsQuery struct {
	Page     int      `json:"page" description:"Page number" example:"1"`
	Limit    int      `json:"limit" description:"Items per page" example:"20"`
	Category string   `json:"category,omitempty" description:"Filter by category" example:"electronics"`
	MinPrice *float64 `json:"min_price,omitempty" description:"Minimum price filter"`
	MaxPrice *float64 `json:"max_price,omitempty" description:"Maximum price filter"`
}

func listProducts(c *fiber.Ctx, query ListProductsQuery) (ProductListResponse, error) {
	return ProductListResponse{
		Products: []ProductResponse{
			{
				ID: 1, Name: "Gaming Laptop", Price: 1299.99,
				Category: "electronics", InStock: true,
				CreatedAt: time.Now(), UpdatedAt: time.Now(),
			},
		},
		Total: 1,
		Page:  query.Page,
		Pages: 1,
	}, nil
}

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Swagify Typed Handlers Example",
	})

	app.Use(cors.New())

	api := swagify.NewFiber(app, swagify.FiberConfig{
		Info: &core.Info{
			Title:       "Product Catalog API",
			Description: "A type-safe product catalog API using swagify's generic typed handlers.",
			Version:     "2.0.0",
		},
		Servers: []core.Server{
			{URL: "http://localhost:8081", Description: "Local development"},
		},
	})

	api.AddTag("Products", "Product management endpoints")

	// Typed routes — request and response schemas are automatically inferred!
	router.TypedGET(api, "/products", listProducts,
		swagify.Summary("List products"),
		swagify.Description("Returns a filterable, paginated list of products."),
		swagify.Tags("Products"),
		swagify.ErrorResponse(500, ErrorResponse{}, "Internal server error"),
	)

	router.TypedGET(api, "/products/:id", getProduct,
		swagify.Summary("Get product by ID"),
		swagify.Tags("Products"),
		swagify.ErrorResponse(404, ErrorResponse{}, "Product not found"),
	)

	router.TypedPOST(api, "/products", createProduct,
		swagify.Summary("Create a product"),
		swagify.Description("Creates a new product in the catalog."),
		swagify.Tags("Products"),
		swagify.ErrorResponse(400, ErrorResponse{}, "Invalid request"),
		swagify.ErrorResponse(422, ErrorResponse{}, "Validation failed"),
	)

	router.TypedPUT(api, "/products/:id", updateProduct,
		swagify.Summary("Update a product"),
		swagify.Tags("Products"),
		swagify.ErrorResponse(400, ErrorResponse{}, "Invalid request"),
		swagify.ErrorResponse(404, ErrorResponse{}, "Product not found"),
	)

	api.RegisterOpenAPI("/openapi.json")
	api.RegisterDocs("/docs")

	log.Println("🚀 Server starting on http://localhost:8081")
	log.Println("📖 API Docs: http://localhost:8081/docs")
	log.Fatal(app.Listen(":8081"))
}
