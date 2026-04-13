// Nested schemas example — demonstrates complex nested types with swagify.
// Run: go run ./examples/nested-schemas
// Docs: http://localhost:8081/docs
package main

import (
	"log"
	"time"

	"github.com/mrgofurov/swagify"
	"github.com/mrgofurov/swagify/core"
)

type Address struct {
	Street  string `json:"street" description:"Street address" example:"123 Main St"`
	City    string `json:"city" description:"City" example:"San Francisco"`
	State   string `json:"state" description:"State" example:"CA"`
	Country string `json:"country" description:"Country code" example:"US"`
}

type ContactInfo struct {
	Email string  `json:"email" description:"Email address" example:"john@company.com"`
	Phone *string `json:"phone,omitempty" description:"Phone number" example:"+1-555-0100"`
}

type Employee struct {
	ID         int       `json:"id" description:"Employee ID" example:"1"`
	FirstName  string    `json:"first_name" description:"First name" example:"John"`
	LastName   string    `json:"last_name" description:"Last name" example:"Doe"`
	Title      string    `json:"title" description:"Job title" example:"Senior Engineer"`
	Contact    ContactInfo `json:"contact"`
	Address    *Address    `json:"address,omitempty"`
	Skills     []string    `json:"skills" description:"Skills list"`
	StartDate  time.Time   `json:"start_date"`
	IsActive   bool        `json:"is_active" example:"true"`
}

type EmployeeList struct {
	Employees []Employee `json:"employees"`
	Total     int        `json:"total" example:"1"`
}

type CreateEmployeeRequest struct {
	FirstName string      `json:"first_name" description:"First name"`
	LastName  string      `json:"last_name" description:"Last name"`
	Title     string      `json:"title" description:"Job title"`
	Contact   ContactInfo `json:"contact"`
	Address   *Address    `json:"address,omitempty"`
	Skills    []string    `json:"skills,omitempty"`
}

func listEmployees(ctx *swagify.Ctx) (EmployeeList, error) {
	return EmployeeList{
		Employees: []Employee{{
			ID: 1, FirstName: "John", LastName: "Doe", Title: "Senior Engineer",
			Contact:   ContactInfo{Email: "john@company.com"},
			Skills:    []string{"Go", "Kubernetes"},
			StartDate: time.Now(),
			IsActive:  true,
		}},
		Total: 1,
	}, nil
}

func getEmployee(ctx *swagify.Ctx) (Employee, error) {
	return Employee{
		ID: 1, FirstName: "John", LastName: "Doe", Title: "Senior Engineer",
		Contact:   ContactInfo{Email: "john@company.com"},
		Skills:    []string{"Go", "Kubernetes"},
		StartDate: time.Now(),
		IsActive:  true,
	}, nil
}

func createEmployee(ctx *swagify.Ctx, req CreateEmployeeRequest) (Employee, error) {
	return Employee{
		ID: 1, FirstName: req.FirstName, LastName: req.LastName,
		Title: req.Title, Contact: req.Contact, Skills: req.Skills,
		StartDate: time.Now(), IsActive: true,
	}, nil
}

func main() {
	api := swagify.New(swagify.Config{
		Title:       "Employee API",
		Description: "Demonstrates deeply nested schemas, slices, pointers, and time.Time.",
		Version:     "1.0.0",
		Servers:     []core.Server{{URL: "http://localhost:8081"}},
	})

	api.GET("/employees", listEmployees)
	api.GET("/employees/{id}", getEmployee)
	api.POST("/employees", createEmployee)

	log.Println("Listening on http://localhost:8081")
	log.Println("Docs:    http://localhost:8081/docs")
	log.Fatal(api.Run(":8081"))
}
