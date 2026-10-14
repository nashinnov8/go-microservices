package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nashinnov8/go-ecommerce/pkg/middleware"
	"github.com/nashinnov8/go-ecommerce/services/product/internal/transport/http/handler"
)

func NewRouter(log *slog.Logger, env string, productHandler *handler.ProductHandler) *gin.Engine {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Global middlewares
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.StructuredLogger(log))
	r.Use(middleware.CORS())

	// Health check endpoint for Kubernetes/Docker probes
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "up",
			"service": "product-service",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		products := v1.Group("/products")
		{
			products.POST("", productHandler.CreateProduct)
			products.GET("", productHandler.ListProducts)
			products.GET("/:id", productHandler.GetProductByID)
			products.PUT("/:id", productHandler.UpdateProduct)
			products.PATCH("/:id/stock", productHandler.UpdateStock)
			products.DELETE("/:id", productHandler.DeleteProduct)
		}
	}

	return r
}
