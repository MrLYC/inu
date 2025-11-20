package middleware

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestBasicAuth_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(BasicAuth("admin", "secret123"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	credentials := base64.StdEncoding.EncodeToString([]byte("admin:secret123"))
	req.Header.Set("Authorization", "Basic "+credentials)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestBasicAuth_MissingAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(BasicAuth("admin", "secret123"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestBasicAuth_InvalidFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(BasicAuth("admin", "secret123"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer token123")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestBasicAuth_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(BasicAuth("admin", "secret123"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	credentials := base64.StdEncoding.EncodeToString([]byte("admin:wrongpassword"))
	req.Header.Set("Authorization", "Basic "+credentials)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestBasicAuth_InvalidBase64(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(BasicAuth("admin", "secret123"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Basic invalid-base64!")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestBasicAuth_InvalidCredentialsFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(BasicAuth("admin", "secret123"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	// No colon in credentials
	credentials := base64.StdEncoding.EncodeToString([]byte("adminonly"))
	req.Header.Set("Authorization", "Basic "+credentials)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}
