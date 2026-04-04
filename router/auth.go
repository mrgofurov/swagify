package router

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
)

// DocsAuthConfig holds credentials for protecting the docs and OpenAPI endpoints.
type DocsAuthConfig struct {
	// Username for basic auth access to docs.
	Username string

	// Password for basic auth access to docs.
	Password string

	// Realm is the authentication realm shown in the browser prompt.
	// Default: "Swagify Docs"
	Realm string
}

// fiberBasicAuth creates a Fiber middleware that enforces HTTP Basic Auth.
func fiberBasicAuth(cfg DocsAuthConfig) fiber.Handler {
	realm := cfg.Realm
	if realm == "" {
		realm = "Swagify Docs"
	}

	expectedUser := sha256.Sum256([]byte(cfg.Username))
	expectedPass := sha256.Sum256([]byte(cfg.Password))

	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			c.Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			return c.SendStatus(401)
		}

		if !strings.HasPrefix(auth, "Basic ") {
			c.Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			return c.SendStatus(401)
		}

		decoded, err := base64.StdEncoding.DecodeString(auth[6:])
		if err != nil {
			c.Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			return c.SendStatus(401)
		}

		parts := strings.SplitN(string(decoded), ":", 2)
		if len(parts) != 2 {
			c.Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			return c.SendStatus(401)
		}

		userHash := sha256.Sum256([]byte(parts[0]))
		passHash := sha256.Sum256([]byte(parts[1]))

		userMatch := subtle.ConstantTimeCompare(userHash[:], expectedUser[:]) == 1
		passMatch := subtle.ConstantTimeCompare(passHash[:], expectedPass[:]) == 1

		if !userMatch || !passMatch {
			c.Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			return c.SendStatus(401)
		}

		return c.Next()
	}
}

// ginBasicAuth creates a Gin middleware that enforces HTTP Basic Auth.
func ginBasicAuth(cfg DocsAuthConfig) gin.HandlerFunc {
	realm := cfg.Realm
	if realm == "" {
		realm = "Swagify Docs"
	}

	expectedUser := sha256.Sum256([]byte(cfg.Username))
	expectedPass := sha256.Sum256([]byte(cfg.Password))

	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.Header("WWW-Authenticate", `Basic realm="`+realm+`"`)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(auth, "Basic ") {
			c.Header("WWW-Authenticate", `Basic realm="`+realm+`"`)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(auth[6:])
		if err != nil {
			c.Header("WWW-Authenticate", `Basic realm="`+realm+`"`)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(string(decoded), ":", 2)
		if len(parts) != 2 {
			c.Header("WWW-Authenticate", `Basic realm="`+realm+`"`)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userHash := sha256.Sum256([]byte(parts[0]))
		passHash := sha256.Sum256([]byte(parts[1]))

		userMatch := subtle.ConstantTimeCompare(userHash[:], expectedUser[:]) == 1
		passMatch := subtle.ConstantTimeCompare(passHash[:], expectedPass[:]) == 1

		if !userMatch || !passMatch {
			c.Header("WWW-Authenticate", `Basic realm="`+realm+`"`)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}

// httpBasicAuth creates a standard http middleware that enforces HTTP Basic Auth.
func httpBasicAuth(cfg DocsAuthConfig) func(http.Handler) http.Handler {
	realm := cfg.Realm
	if realm == "" {
		realm = "Swagify Docs"
	}

	expectedUser := sha256.Sum256([]byte(cfg.Username))
	expectedPass := sha256.Sum256([]byte(cfg.Password))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Basic ") {
				w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(auth[6:])
			if err != nil {
				w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(string(decoded), ":", 2)
			if len(parts) != 2 {
				w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userHash := sha256.Sum256([]byte(parts[0]))
			passHash := sha256.Sum256([]byte(parts[1]))

			userMatch := subtle.ConstantTimeCompare(userHash[:], expectedUser[:]) == 1
			passMatch := subtle.ConstantTimeCompare(passHash[:], expectedPass[:]) == 1

			if !userMatch || !passMatch {
				w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
