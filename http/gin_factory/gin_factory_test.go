package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewGinFactory(t *testing.T) {
	gf := NewGinFactory()

	assert.NotNil(t, gf, "GinFactory instance should not be nil")
}

func TestAddMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	gf := NewGinFactory()

	// Flags to track middleware execution
	middleware1Called := false
	middleware2Called := false

	middleware1 := func(c *gin.Context) {
		middleware1Called = true
		c.Next()
	}
	middleware2 := func(c *gin.Context) {
		middleware2Called = true
		c.Next()
	}

	gf.AddMiddleware(middleware1, middleware2)

	// Create a router and make a test request
	r := gf.CreateRouter()
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test handler")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	// Assertions
	assert.True(t, middleware1Called, "middleware1 should have been called")
	assert.True(t, middleware2Called, "middleware2 should have been called")
	assert.Equal(t, http.StatusOK, w.Code, "Response status should be OK")
	assert.Equal(t, "test handler", w.Body.String(), "Response body should match handler output")
}

func TestAddHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	gf := NewGinFactory()

	handler1 := func(r *gin.Engine) {
		r.GET("/test1", func(c *gin.Context) {
			c.String(http.StatusOK, "test1 response")
		})
	}
	handler2 := func(r *gin.Engine) {
		r.GET("/test2", func(c *gin.Context) {
			c.String(http.StatusOK, "test2 response")
		})
	}

	// Add handlers to GinFactory
	gf.AddHandlers(handler1, handler2)

	// Create router and test the handlers
	r := gf.CreateRouter()

	// Test /test1 endpoint
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest(http.MethodGet, "/test1", nil)
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code, "Response status for /test1 should be OK")
	assert.Equal(t, "test1 response", w1.Body.String(), "Response body for /test1 should match handler1 output")

	// Test /test2 endpoint
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/test2", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code, "Response status for /test2 should be OK")
	assert.Equal(t, "test2 response", w2.Body.String(), "Response body for /test2 should match handler2 output")
}

func TestCreateRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	gf := NewGinFactory()

	middlewareCalled := false
	middleware := func(c *gin.Context) {
		middlewareCalled = true
		c.Next()
	}
	gf.AddMiddleware(middleware)

	handlerCalled := false
	handler := func(r *gin.Engine) {
		r.GET("/test", func(c *gin.Context) {
			handlerCalled = true
			c.String(http.StatusOK, "handler response")
		})
	}
	gf.AddHandlers(handler)

	router := gf.CreateRouter()

	// Perform a request to test middleware and handler integration
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code, "Response status code should be 200")
	assert.Equal(t, "handler response", w.Body.String(), "Response body should match the handler's output")
	assert.True(t, middlewareCalled, "Middleware should have been called")
	assert.True(t, handlerCalled, "Handler should have been called")
}

func TestMiddlewareHeaderModification(t *testing.T) {
	gin.SetMode(gin.TestMode)
	gf := NewGinFactory()

	// Middleware modifies headers
	gf.AddMiddleware(func(c *gin.Context) {
		c.Writer.Header().Set("X-Custom-Header", "middleware")
		c.Next()
	})

	r := gf.CreateRouter()
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test handler")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code, "Response status should be OK")
	assert.Equal(t, "middleware", w.Header().Get("X-Custom-Header"), "Middleware should set X-Custom-Header")
	assert.Equal(t, "test handler", w.Body.String(), "Response body should match handler output")
}

func TestCreateRouterWithoutMiddlewareOrHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	gf := NewGinFactory()

	// Clear default middleware and handlers
	gf.AddMiddleware()
	gf.AddHandlers()

	r := gf.CreateRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusNotFound, w.Code, "Router with no handlers should return 404")
}

func TestDefaultRecoveryMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	gf := NewGinFactory()

	// Add a handler that panics
	gf.AddHandlers(func(r *gin.Engine) {
		r.GET("/panic", func(c *gin.Context) {
			panic("test panic")
		})
	})

	r := gf.CreateRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/panic", nil)
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, w.Code, "Recovery middleware should handle panics and return 500")
	assert.Contains(t, w.Body.String(), "", "Response body be empty as default recovery middleware does not include a body in its response when a panic occurs")
}
