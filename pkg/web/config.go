package web

import (
	"fmt"
)

// Config holds the configuration for the web server
type Config struct {
	// Addr is the address to listen on (e.g., "127.0.0.1:8080")
	Addr string
	// AdminUser is the admin username for HTTP Basic Auth
	AdminUser string
	// AdminToken is the admin password/token for HTTP Basic Auth
	AdminToken string
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Addr == "" {
		return fmt.Errorf("addr cannot be empty")
	}
	// If AdminToken is empty, no authentication is required
	// If AdminToken is set, AdminUser must also be set
	if c.AdminToken != "" && c.AdminUser == "" {
		return fmt.Errorf("admin-user cannot be empty when admin-token is set")
	}
	return nil
}

// IsAuthEnabled returns true if authentication is enabled
func (c *Config) IsAuthEnabled() bool {
	return c.AdminToken != ""
}
