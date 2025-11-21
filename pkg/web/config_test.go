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
			name: "valid config with auth",
			config: &Config{
				Addr:       "127.0.0.1:8080",
				AdminUser:  "admin",
				AdminToken: "secret",
			},
			wantError: false,
		},
		{
			name: "valid config without auth",
			config: &Config{
				Addr:       "127.0.0.1:8080",
				AdminUser:  "",
				AdminToken: "",
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
			name: "token set but no user",
			config: &Config{
				Addr:       "127.0.0.1:8080",
				AdminUser:  "",
				AdminToken: "secret",
			},
			wantError: true,
		},
		{
			name: "user set but no token (auth disabled)",
			config: &Config{
				Addr:       "127.0.0.1:8080",
				AdminUser:  "admin",
				AdminToken: "",
			},
			wantError: false,
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

func TestConfig_IsAuthEnabled(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		want   bool
	}{
		{
			name: "auth enabled with token",
			config: &Config{
				AdminUser:  "admin",
				AdminToken: "secret",
			},
			want: true,
		},
		{
			name: "auth disabled - empty token",
			config: &Config{
				AdminUser:  "admin",
				AdminToken: "",
			},
			want: false,
		},
		{
			name: "auth disabled - both empty",
			config: &Config{
				AdminUser:  "",
				AdminToken: "",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.IsAuthEnabled()
			if got != tt.want {
				t.Errorf("IsAuthEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}
