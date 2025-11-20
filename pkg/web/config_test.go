package web

import (
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		wantError bool
	}{
		{
			name: "valid config",
			config: &Config{
				Addr:       "127.0.0.1:8080",
				AdminUser:  "admin",
				AdminToken: "secret",
			},
			wantError: false,
		},
		{
			name: "empty addr",
			config: &Config{
				Addr:       "",
				AdminUser:  "admin",
				AdminToken: "secret",
			},
			wantError: true,
		},
		{
			name: "empty admin user",
			config: &Config{
				Addr:       "127.0.0.1:8080",
				AdminUser:  "",
				AdminToken: "secret",
			},
			wantError: true,
		},
		{
			name: "empty admin token",
			config: &Config{
				Addr:       "127.0.0.1:8080",
				AdminUser:  "admin",
				AdminToken: "",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
