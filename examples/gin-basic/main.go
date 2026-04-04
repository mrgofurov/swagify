// Gin basic CRUD example demonstrating swagify with Gin framework.
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrgofurov/swagify"
	"github.com/mrgofurov/swagify/core"
)

// --- Models ---

type CreateBookRequest struct {
	Title  string  `json:"title" validate:"required" description:"Book title" example:"The Go Programming Language"`
	Author string  `json:"author" validate:"required" description:"Book author" example:"Alan Donovan"`
	ISBN   string  `json:"isbn" validate:"required" description:"ISBN number" example:"978-0134190440"`
	Year   int     `json:"year" validate:"min=1000,max=2100" description:"Publication year" example:"2015"`
	Price  float64 `json:"price" validate:"min=0" description:"Price in USD" example:"34.99"`
}

type UpdateBookRequest struct {
	Title  *string  `json:"title,omitempty" description:"Book title"`
	Author *string  `json:"author,omitempty" description:"Book author"`
	Price  *float64 `json:"price,omitempty" validate:"min=0" description:"Price in USD"`
}

type BookResponse struct {
	ID     int     `json:"id" description:"Unique book identifier" example:"1"`
	Title  string  `json:"title" description:"Book title" example:"The Go Programming Language"`
	Author string  `json:"author" description:"Book author" example:"Alan Donovan"`
	ISBN   string  `json:"isbn" description:"ISBN number" example:"978-0134190440"`
	Year   int     `json:"year" description:"Publication year" example:"2015"`
	Price  float64 `json:"price" description:"Price in USD" example:"34.99"`
}

type BookListResponse struct {
	Books []BookResponse `json:"books" description:"List of books"`
	Total int            `json:"total" description:"Total count" example:"100"`
}

type ErrorResponse struct {
	Error string `json:"error" description:"Error message" example:"Resource not found"`
	Code  int    `json:"code" description:"HTTP status code" example:"404"`
}

type MessageResponse struct {
	Message string `json:"message" description:"Response message" example:"Book deleted successfully"`
}

// --- Handlers ---

func listBooks(c *gin.Context) {
	c.JSON(http.StatusOK, BookListResponse{
		Books: []BookResponse{
			{ID: 1, Title: "The Go Programming Language", Author: "Alan Donovan", ISBN: "978-0134190440", Year: 2015, Price: 34.99},
			{ID: 2, Title: "Go in Action", Author: "William Kennedy", ISBN: "978-1617291784", Year: 2015, Price: 29.99},
		},
		Total: 2,
	})
}

func getBook(c *gin.Context) {
	c.JSON(http.StatusOK, BookResponse{
		ID: 1, Title: "The Go Programming Language", Author: "Alan Donovan",
		ISBN: "978-0134190440", Year: 2015, Price: 34.99,
	})
}

func createBook(c *gin.Context) {
	var req CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error(), Code: 400})
		return
	}
	c.JSON(http.StatusCreated, BookResponse{
		ID: 1, Title: req.Title, Author: req.Author,
		ISBN: req.ISBN, Year: req.Year, Price: req.Price,
	})
}

func updateBook(c *gin.Context) {
	var req UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error(), Code: 400})
		return
	}
	c.JSON(http.StatusOK, BookResponse{
		ID: 1, Title: "Updated Book", Author: "Updated Author",
		ISBN: "978-0000000000", Year: 2024, Price: 39.99,
	})
}

func deleteBook(c *gin.Context) {
	c.JSON(http.StatusOK, MessageResponse{Message: "Book deleted successfully"})
}

func main() {
	r := gin.Default()

	api := swagify.NewGin(r, swagify.GinConfig{
		Info: &core.Info{
			Title:       "Bookstore API",
			Description: "A bookstore CRUD API built with Gin and documented with Swagify.",
			Version:     "1.0.0",
			Contact: &core.Contact{
				Name:  "Bookstore Team",
				Email: "books@example.com",
			},
		},
		Servers: []core.Server{
			{URL: "http://localhost:8082", Description: "Local development"},
		},
	})

	api.AddTag("Books", "Book management operations")

	// Register routes
	api.GET("/books", listBooks,
		swagify.Summary("List all books"),
		swagify.Description("Returns a list of all books in the bookstore."),
		swagify.Tags("Books"),
		swagify.WithResponse(BookListResponse{}),
	)

	api.GET("/books/:id", getBook,
		swagify.Summary("Get book by ID"),
		swagify.Tags("Books"),
		swagify.WithResponse(BookResponse{}),
		swagify.ErrorResponse(404, ErrorResponse{}, "Book not found"),
	)

	api.POST("/books", createBook,
		swagify.Summary("Create a book"),
		swagify.Description("Adds a new book to the bookstore catalog."),
		swagify.Tags("Books"),
		swagify.WithRequest(CreateBookRequest{}),
		swagify.WithResponse(BookResponse{}),
		swagify.SuccessStatus(201),
		swagify.ErrorResponse(400, ErrorResponse{}, "Invalid request"),
	)

	api.PUT("/books/:id", updateBook,
		swagify.Summary("Update a book"),
		swagify.Tags("Books"),
		swagify.WithRequest(UpdateBookRequest{}),
		swagify.WithResponse(BookResponse{}),
		swagify.ErrorResponse(400, ErrorResponse{}, "Invalid request"),
		swagify.ErrorResponse(404, ErrorResponse{}, "Book not found"),
	)

	api.DELETE("/books/:id", deleteBook,
		swagify.Summary("Delete a book"),
		swagify.Tags("Books"),
		swagify.WithResponse(MessageResponse{}),
		swagify.ErrorResponse(404, ErrorResponse{}, "Book not found"),
	)

	// Register OpenAPI and Docs
	api.RegisterOpenAPI("/openapi.json")
	api.RegisterDocs("/docs")

	log.Println("🚀 Server starting on http://localhost:8082")
	log.Println("📖 API Docs: http://localhost:8082/docs")
	log.Fatal(r.Run(":8082"))
}
