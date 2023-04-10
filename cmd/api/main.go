package main

import (
	"log"

	"github.com/NPG27/supermarket_dop/cmd/api/handlers"
	"github.com/NPG27/supermarket_dop/internal/middleware"
	"github.com/NPG27/supermarket_dop/internal/repository"
	"github.com/NPG27/supermarket_dop/internal/service"
	"github.com/NPG27/supermarket_dop/pkg/store"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	storage := store.NewStore("./data/products.json")
	productRepo, errProductRepo := repository.NewProductRepository(storage)
	if errProductRepo != nil {
		log.Fatalf("Error initializing repository: %v", errProductRepo)
	}
	productService := service.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	products := router.Group("/products")
	products.Use(middleware.VerifyToken())
	{
		products.GET("", productHandler.GetAllProducts)
		products.GET("/:id", productHandler.GetProductByID)
		products.GET("/filter", productHandler.GetProductByPriceGreaterThan)
		products.POST("", productHandler.CreateProduct)
		products.PATCH("/:id", productHandler.PatchProduct)
		products.PUT("/:id", productHandler.UpdateProduct)
		products.DELETE("/:id", productHandler.DeleteProduct)
	}

	if err := router.Run(); err != nil {
		panic(err)
	}
}
