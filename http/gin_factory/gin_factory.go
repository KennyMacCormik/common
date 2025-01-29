// Package router provides a GinFactory for managing middleware and handlers in a modular way.
// It simplifies the creation of a Gin router with preconfigured middleware and route handlers.
package router

import "github.com/gin-gonic/gin"

// GinFactory is a factory for managing middleware and handlers in a Gin application.
// It provides methods for adding middleware, adding handlers, and creating a router instance.
type GinFactory struct {
	middleware []gin.HandlerFunc
	handlers   []func(router *gin.Engine)
}

// NewGinFactory initializes a new instance of GinFactory.
// It includes the default gin.Recovery middleware to handle panics gracefully.
func NewGinFactory() *GinFactory {
	return &GinFactory{middleware: []gin.HandlerFunc{gin.Recovery()}, handlers: make([]func(router *gin.Engine), 0)}
}

// AddMiddleware adds middleware to the GinFactory.
// Middleware is applied in the order it is added.
func (g *GinFactory) AddMiddleware(middleware ...gin.HandlerFunc) {
	g.middleware = append(g.middleware, middleware...)
}

// ResetMiddleware adds middleware to the GinFactory.
// Middleware is applied in the order it is added.
func (g *GinFactory) ResetMiddleware(middleware ...gin.HandlerFunc) {
	g.middleware = append(g.middleware, middleware...)
}

// AddHandlers adds route handlers to the GinFactory.
// Handlers are used to define specific routes and their behaviors.
func (g *GinFactory) AddHandlers(handlers ...func(router *gin.Engine)) {
	g.handlers = append(g.handlers, handlers...)
}

// CreateRouter creates a new gin.Engine instance with the configured middleware and handlers.
// The Gin router is initialized in release mode for optimal performance.
func (g *GinFactory) CreateRouter() *gin.Engine {
	router := gin.New()

	for _, m := range g.middleware {
		router.Use(m)
	}

	for _, h := range g.handlers {
		h(router)
	}

	return router
}
