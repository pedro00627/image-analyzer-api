package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestNewRateLimiter(t *testing.T) {
	limiter := NewRateLimiter(rate.Every(time.Second), 5)
	assert.NotNil(t, limiter)
	assert.Implements(t, (*RateLimiterInterface)(nil), limiter)
}

func TestRateLimiter_AllowsRequestsWithinLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create rate limiter: 2 requests per second
	limiter := NewRateLimiter(rate.Every(500*time.Millisecond), 2)

	router := gin.New()
	router.Use(limiter.Limit())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// First request should succeed
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Second request should succeed (within burst)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRateLimiter_BlocksRequestsExceedingLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create rate limiter: 2 requests per second, burst of 2
	limiter := NewRateLimiter(rate.Every(500*time.Millisecond), 2)

	router := gin.New()
	router.Use(limiter.Limit())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// Make 2 requests (should succeed - within burst)
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// Third request should be blocked
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.Contains(t, w.Body.String(), "Rate limit exceeded")
}

func TestRateLimiter_DifferentIPsHaveIndependentLimits(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create rate limiter: 1 request per second, burst of 1
	limiter := NewRateLimiter(rate.Every(time.Second), 1)

	router := gin.New()
	router.Use(limiter.Limit())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// First IP makes request (should succeed)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "192.168.1.1:1234"
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Second IP makes request (should also succeed - independent limit)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "192.168.1.2:1234"
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// First IP makes second request immediately (should be blocked)
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/test", nil)
	req3.RemoteAddr = "192.168.1.1:1234"
	router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusTooManyRequests, w3.Code)
}

func TestRateLimiter_TokensRegenerateOverTime(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create rate limiter: 1 request per 100ms, burst of 1
	limiter := NewRateLimiter(rate.Every(100*time.Millisecond), 1)

	router := gin.New()
	router.Use(limiter.Limit())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// First request (should succeed)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "192.168.1.1:1234"
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Immediate second request (should be blocked)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "192.168.1.1:1234"
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusTooManyRequests, w2.Code)

	// Wait for token to regenerate
	time.Sleep(150 * time.Millisecond)

	// Third request after waiting (should succeed)
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/test", nil)
	req3.RemoteAddr = "192.168.1.1:1234"
	router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)
}

func TestRateLimiter_ErrorResponseFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	limiter := NewRateLimiter(rate.Every(time.Second), 1)

	router := gin.New()
	router.Use(limiter.Limit())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// Make 2 requests to exceed limit
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		router.ServeHTTP(w, req)
	}

	// Verify error response format
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.Contains(t, w.Body.String(), `"success":false`)
	assert.Contains(t, w.Body.String(), `"error":"Rate limit exceeded. Please try again later."`)
}
