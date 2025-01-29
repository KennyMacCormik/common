# GinFactory

GinFactory is a modular utility for managing middleware and handlers in a Gin application. It provides a simple interface for configuring and creating a Gin router with preconfigured middleware and route handlers, improving code organization and reusability.

## Features

- **Middleware Management**: Easily add or reset middleware in a specified order.
- **Handler Management**: Add modular route handlers to define application behavior.
- **Default Recovery Middleware**: Includes Gin's `gin.Recovery` middleware for handling panics gracefully.
- **Lightweight and Flexible**: Designed to integrate seamlessly with Gin applications.

## Installation

To use GinFactory, you need to install the required dependencies:

### Dependencies Installation

```bash
# Install Gin framework
go get -u github.com/gin-gonic/gin

# Install Testify for testing
go get -u github.com/stretchr/testify
```

### Package Installation

Clone or install the package directly:

```bash
# Install the router package
go get -u github.com/KennyMacCormik/common/http/gin_factory
```

## Usage

Below are examples of how to use GinFactory to manage middleware and route handlers:

### Create a New GinFactory

```go
import "github.com/KennyMacCormik/common/http/gin_factory"

// Initialize GinFactory
factory := router.NewGinFactory()
```

### Add Middleware

```go
factory.AddMiddleware(func(c *gin.Context) {
    c.Writer.Header().Set("X-Custom-Header", "middleware")
    c.Next()
})
```

### Add Route Handlers

```go
factory.AddHandlers(func(r *gin.Engine) {
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })
})
```

### Create a Router

```go
router := factory.CreateRouter()
router.Run() // Start the Gin server
```

## API Documentation

### Functions

#### `func NewGinFactory() *GinFactory`
Initializes a new instance of `GinFactory` with default recovery middleware.

#### `func (g *GinFactory) AddMiddleware(middleware ...gin.HandlerFunc)`
Adds one or more middleware functions to the factory. Middleware is applied in the order it is added.

#### `func (g *GinFactory) ResetMiddleware(middleware ...gin.HandlerFunc)`
Replaces all middleware in the factory with the specified middleware functions.

#### `func (g *GinFactory) AddHandlers(handlers ...func(router *gin.Engine))`
Adds one or more route handlers to the factory.

#### `func (g *GinFactory) CreateRouter() *gin.Engine`
Creates and returns a new Gin router instance with the configured middleware and handlers applied.

## Type Descriptions

### `type GinFactory`
A factory for managing middleware and handlers in a Gin application.

- **Methods**:
    - `AddMiddleware`
    - `ResetMiddleware`
    - `AddHandlers`
    - `CreateRouter`

## License

This package is licensed under the [MIT License](https://opensource.org/licenses/MIT).

## Thanks

Special thanks to the contributors and maintainers of the following packages:

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Testify](https://github.com/stretchr/testify)

