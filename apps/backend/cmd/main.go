package main

import (
	"context"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/omidgz/multi-tenant-store-admin-panel/apps/backend/internal/config"
	"github.com/omidgz/multi-tenant-store-admin-panel/apps/backend/internal/handler"
	"github.com/omidgz/multi-tenant-store-admin-panel/apps/backend/internal/middleware"
	"github.com/omidgz/multi-tenant-store-admin-panel/apps/backend/internal/repository"
	"github.com/omidgz/multi-tenant-store-admin-panel/apps/backend/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Create database connection pool
	pool, err := pgxpool.New(context.Background(), cfg.GetPostgresConnectionString())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer pool.Close()

	// Initialize repository, service, and handler
	repo := repository.NewPostgresProductRepository(pool)
	productService := service.NewProductService(repo)
	productHandler := handler.NewProductHandler(productService)

	// Set up Gin router
	r := gin.Default()
	r.Use(gin.Recovery())

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "X-Tenant-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.Use(middleware.TenantMiddleware(cfg))

	// Routes
	api := r.Group("/api")
	{
		products := api.Group("/products")
		{
			products.POST("", productHandler.CreateProduct)
			products.GET("", productHandler.GetProducts)
			products.GET("/:id", productHandler.GetProduct)
			products.PUT("/:id", productHandler.UpdateProduct)
			products.DELETE("/:id", productHandler.DeleteProduct)
		}
	}

	log.Printf("Server starting on port %s...", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
