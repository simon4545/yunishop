package main

import (
	"github.com/simon4545/goshop/database"
	"github.com/simon4545/goshop/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize database
	database.InitDB()
	render := handlers.NewRenderer()
	render.AddDirectory("templates")
	e := echo.New()
	e.Renderer = render
	// Create a new Echo instance
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	e.POST("/create-order", handlers.CreatePayOrder)
	e.POST("/capture-order/:id", handlers.CaptureOrder)

	// Routes
	e.GET("/product-categories", handlers.GetProductCategories)

	e.POST("/products", handlers.AddProduct)
	e.DELETE("/products/:id", handlers.DeleteProduct)
	e.GET("/products/:id", handlers.GetProduct)
	e.PUT("/products/:id", handlers.UpdateProduct)

	e.GET("/news/categories", handlers.GetNewsCategories)
	e.POST("/news/categories", handlers.CreateNewsCategory)
	e.PUT("/news/categories/:id", handlers.UpdateNewsCategory)
	e.DELETE("/news/categories/:id", handlers.DeleteNewsCategory)

	e.POST("/news", handlers.CreateNews)
	e.PUT("/news/:id", handlers.UpdateNews)
	e.DELETE("/news/:id", handlers.DeleteNews)
	e.GET("/news/:id", handlers.GetNews)

	e.POST("/images", handlers.UploadImage)

	e.GET("/", handlers.GetProductsByCategory)
	// Serve static files
	e.Static("/uploads", "./uploads")
	e.Static("/", "static")
	// Start server
	e.Logger.Fatal(e.Start(":8082"))
}
