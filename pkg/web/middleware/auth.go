package middleware

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// BasicAuth returns a middleware that implements HTTP Basic Authentication
func BasicAuth(username, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("Authentication failed: missing Authorization header from %s", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Authentication required",
				"code":    401,
			})
			c.Abort()
			return
		}

		// Parse "Basic <base64>" format
		const prefix = "Basic "
		if !strings.HasPrefix(authHeader, prefix) {
			log.Printf("Authentication failed: invalid Authorization format from %s", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Invalid authentication format",
				"code":    401,
			})
			c.Abort()
			return
		}

		// Decode base64
		encoded := authHeader[len(prefix):]
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			log.Printf("Authentication failed: invalid base64 encoding from %s", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Invalid authentication encoding",
				"code":    401,
			})
			c.Abort()
			return
		}

		// Split username:password
		credentials := string(decoded)
		parts := strings.SplitN(credentials, ":", 2)
		if len(parts) != 2 {
			log.Printf("Authentication failed: invalid credentials format from %s", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Invalid credentials format",
				"code":    401,
			})
			c.Abort()
			return
		}

		// Verify credentials
		providedUser, providedPass := parts[0], parts[1]
		if providedUser != username || providedPass != password {
			log.Printf("Authentication failed: invalid username or password (user: %s) from %s", providedUser, c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Invalid username or password",
				"code":    401,
			})
			c.Abort()
			return
		}

		// Authentication successful
		c.Next()
	}
}
