// Package ui provides embedded API documentation UI assets
// and handler registration for multiple web frameworks.
package ui

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	_ "embed"
)

//go:embed scalar/index.html
var scalarHTML string

// templateData holds the data injected into the HTML template.
type templateData struct {
	Title   string
	SpecURL string
}

// renderHTML renders the docs HTML with the given spec URL.
func renderHTML(specURL string, title ...string) (string, error) {
	t := "API Documentation"
	if len(title) > 0 && title[0] != "" {
		t = title[0]
	}

	tmpl, err := template.New("docs").Parse(scalarHTML)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, templateData{
		Title:   t,
		SpecURL: specURL,
	})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// RegisterFiber registers the docs UI route on a Fiber app.
func RegisterFiber(app *fiber.App, docsPath, specURL string) {
	html, err := renderHTML(specURL)
	if err != nil {
		panic("swagify: failed to render docs UI template: " + err.Error())
	}

	app.Get(docsPath, func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(html)
	})
}

// RegisterGin registers the docs UI route on a Gin engine.
func RegisterGin(engine *gin.Engine, docsPath, specURL string) {
	html, err := renderHTML(specURL)
	if err != nil {
		panic("swagify: failed to render docs UI template: " + err.Error())
	}

	engine.GET(docsPath, func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})
}

// RegisterHTTP registers the docs UI route on a standard http.ServeMux.
func RegisterHTTP(mux *http.ServeMux, docsPath, specURL string) {
	html, err := renderHTML(specURL)
	if err != nil {
		panic("swagify: failed to render docs UI template: " + err.Error())
	}

	mux.HandleFunc("GET "+docsPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	})
}
