package main

import (
	"github.com/simon4545/goshop/database"
	"github.com/simon4545/goshop/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize database
	database.InitDB()

	// Create a new Echo instance
	e := echo.New()

	// Routes
	e.GET("/product-categories", handlers.GetProductCategories)

	e.POST("/products", handlers.AddProduct)
	e.DELETE("/products/:id", handlers.DeleteProduct)
	e.GET("/products/:id", handlers.GetProduct)
	e.GET("/products/sku/:sku", handlers.GetProductBySKU)

	e.POST("/skus", handlers.AddSKUs)
	e.DELETE("/skus/:sku", handlers.DeleteSKU)

	e.GET("/news/categories", handlers.GetNewsCategories)
	e.POST("/news/categories", handlers.CreateNewsCategory)
	e.PUT("/news/categories/:id", handlers.UpdateNewsCategory)
	e.DELETE("/news/categories/:id", handlers.DeleteNewsCategory)

	e.POST("/news", handlers.CreateNews)
	e.PUT("/news/:id", handlers.UpdateNews)
	e.DELETE("/news/:id", handlers.DeleteNews)
	e.GET("/news/:id", handlers.GetNews)

	e.POST("/images", handlers.UploadImage)
	e.DELETE("/images/:id", handlers.DeleteImage)
	e.GET("/images/:id", handlers.GetImage)

	// Serve static files
	e.Static("/uploads", "./uploads")

	// Start server
	e.Logger.Fatal(e.Start(":8081"))
}
