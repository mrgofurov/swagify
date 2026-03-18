// Nested schemas example demonstrating complex type handling.
package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/swagify"
	"github.com/swagify/core"
)

// --- Deeply nested model hierarchy ---

type Address struct {
	Street  string `json:"street" description:"Street address" example:"123 Main St"`
	City    string `json:"city" description:"City name" example:"San Francisco"`
	State   string `json:"state" description:"State/province" example:"CA"`
	ZipCode string `json:"zip_code" description:"Postal code" example:"94102"`
	Country string `json:"country" description:"Country code" example:"US"`
}

type ContactInfo struct {
	Email   string  `json:"email" validate:"required,email" description:"Primary email" example:"john@company.com"`
	Phone   string  `json:"phone,omitempty" description:"Phone number" example:"+1-555-0100"`
	Website *string `json:"website,omitempty" description:"Personal website" format:"uri"`
}

type SocialLinks struct {
	Twitter  *string `json:"twitter,omitempty" description:"Twitter handle" example:"@johndoe"`
	LinkedIn *string `json:"linkedin,omitempty" description:"LinkedIn profile URL" format:"uri"`
	GitHub   *string `json:"github,omitempty" description:"GitHub username" example:"johndoe"`
}

type CompanyInfo struct {
	Name     string   `json:"name" description:"Company name" example:"Acme Corp"`
	Industry string   `json:"industry" description:"Industry sector" example:"Technology"`
	Size     string   `json:"size" validate:"oneof=startup small medium large enterprise" description:"Company size"`
	Address  Address  `json:"address" description:"Company headquarters"`
}

type EmployeeProfile struct {
	ID          int            `json:"id" description:"Employee ID" example:"1"`
	FirstName   string         `json:"first_name" validate:"required" description:"First name" example:"John"`
	LastName    string         `json:"last_name" validate:"required" description:"Last name" example:"Doe"`
	Title       string         `json:"title" description:"Job title" example:"Senior Engineer"`
	Department  string         `json:"department" description:"Department name" example:"Engineering"`
	Contact     ContactInfo    `json:"contact" description:"Contact information"`
	Social      *SocialLinks   `json:"social,omitempty" description:"Social media links"`
	Company     CompanyInfo    `json:"company" description:"Company information"`
	HomeAddress *Address       `json:"home_address,omitempty" description:"Home address"`
	Skills      []string       `json:"skills" description:"List of skills"`
	Projects    []ProjectBrief `json:"projects" description:"Active projects"`
	Metadata    map[string]string `json:"metadata,omitempty" description:"Custom metadata key-value pairs"`
	StartDate   time.Time      `json:"start_date" description:"Employment start date"`
	IsActive    bool           `json:"is_active" description:"Whether the employee is currently active" example:"true"`
}

type ProjectBrief struct {
	ID          int       `json:"id" description:"Project ID" example:"101"`
	Name        string    `json:"name" description:"Project name" example:"Project Alpha"`
	Role        string    `json:"role" description:"Role in the project" example:"Lead Developer"`
	StartDate   time.Time `json:"start_date" description:"Project start date"`
	EndDate     *time.Time `json:"end_date,omitempty" description:"Project end date (null if ongoing)"`
}

type CreateEmployeeRequest struct {
	FirstName   string       `json:"first_name" validate:"required" description:"First name"`
	LastName    string       `json:"last_name" validate:"required" description:"Last name"`
	Title       string       `json:"title" validate:"required" description:"Job title"`
	Department  string       `json:"department" validate:"required" description:"Department"`
	Contact     ContactInfo  `json:"contact" description:"Contact info"`
	Social      *SocialLinks `json:"social,omitempty" description:"Social links"`
	Company     CompanyInfo  `json:"company" description:"Company info"`
	HomeAddress *Address     `json:"home_address,omitempty" description:"Home address"`
	Skills      []string     `json:"skills,omitempty" description:"Skills list"`
}

type EmployeeListResponse struct {
	Employees []EmployeeProfile `json:"employees" description:"List of employees"`
	Total     int               `json:"total" description:"Total employee count" example:"500"`
}

type ErrorResponse struct {
	Error string `json:"error" description:"Error message"`
	Code  int    `json:"code" description:"Error code"`
}

// --- Handlers ---

func listEmployees(c *fiber.Ctx) error {
	return c.JSON(EmployeeListResponse{
		Employees: []EmployeeProfile{
			{
				ID: 1, FirstName: "John", LastName: "Doe",
				Title: "Senior Engineer", Department: "Engineering",
				Contact: ContactInfo{Email: "john@company.com"},
				Company: CompanyInfo{
					Name: "Acme Corp", Industry: "Technology", Size: "large",
					Address: Address{City: "San Francisco", State: "CA", Country: "US"},
				},
				Skills:    []string{"Go", "Python", "Kubernetes"},
				Projects:  []ProjectBrief{{ID: 1, Name: "Project Alpha", Role: "Lead", StartDate: time.Now()}},
				StartDate: time.Now().Add(-365 * 24 * time.Hour),
				IsActive:  true,
			},
		},
		Total: 1,
	})
}

func getEmployee(c *fiber.Ctx) error {
	return c.JSON(EmployeeProfile{
		ID: 1, FirstName: "John", LastName: "Doe",
		Title: "Senior Engineer", Department: "Engineering",
		Contact: ContactInfo{Email: "john@company.com", Phone: "+1-555-0100"},
		Company: CompanyInfo{
			Name: "Acme Corp", Industry: "Technology", Size: "large",
			Address: Address{Street: "123 Main St", City: "San Francisco", State: "CA", ZipCode: "94102", Country: "US"},
		},
		Skills:    []string{"Go", "Python", "Kubernetes"},
		Projects:  []ProjectBrief{{ID: 1, Name: "Project Alpha", Role: "Lead", StartDate: time.Now()}},
		Metadata:  map[string]string{"badge_number": "EMP001", "floor": "5"},
		StartDate: time.Now().Add(-365 * 24 * time.Hour),
		IsActive:  true,
	})
}

func createEmployee(c *fiber.Ctx) error {
	var req CreateEmployeeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(ErrorResponse{Error: err.Error(), Code: 400})
	}
	return c.Status(201).JSON(EmployeeProfile{
		ID: 1, FirstName: req.FirstName, LastName: req.LastName,
		Title: req.Title, Department: req.Department,
		Contact: req.Contact, Company: req.Company,
		Skills: req.Skills, StartDate: time.Now(), IsActive: true,
	})
}

func main() {
	app := fiber.New()
	app.Use(cors.New())

	api := swagify.NewFiber(app, swagify.FiberConfig{
		Info: &core.Info{
			Title:       "Employee Directory API",
			Description: "Demonstrates deeply nested schemas, arrays, maps, pointers, and time.Time handling.",
			Version:     "1.0.0",
		},
		Servers: []core.Server{
			{URL: "http://localhost:8084", Description: "Local development"},
		},
	})

	api.AddTag("Employees", "Employee management with complex nested schemas")

	api.GET("/employees", listEmployees,
		swagify.Summary("List employees"),
		swagify.Description("Returns all employees with their nested profile data."),
		swagify.Tags("Employees"),
		swagify.WithResponse(EmployeeListResponse{}),
	)

	api.GET("/employees/:id", getEmployee,
		swagify.Summary("Get employee by ID"),
		swagify.Description("Returns a single employee with full nested profile."),
		swagify.Tags("Employees"),
		swagify.WithResponse(EmployeeProfile{}),
		swagify.ErrorResponse(404, ErrorResponse{}, "Employee not found"),
	)

	api.POST("/employees", createEmployee,
		swagify.Summary("Create an employee"),
		swagify.Description("Creates a new employee with nested company, address, and contact data."),
		swagify.Tags("Employees"),
		swagify.WithRequest(CreateEmployeeRequest{}),
		swagify.WithResponse(EmployeeProfile{}),
		swagify.SuccessStatus(201),
		swagify.ErrorResponse(400, ErrorResponse{}, "Invalid request"),
	)

	api.RegisterOpenAPI("/openapi.json")
	api.RegisterDocs("/docs")

	log.Println("🚀 Server starting on http://localhost:8084")
	log.Println("📖 API Docs: http://localhost:8084/docs")
	log.Fatal(app.Listen(":8084"))
}
