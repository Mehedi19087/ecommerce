package catalog

import (
	"github.com/gin-gonic/gin"
)

func SetupCatalogRoutes(
	router *gin.Engine,
	productController *ProductController) {
	v1 := router.Group("/api/v1")
	//category routes
	categories := v1.Group("/categories")
	{
		categories.POST("", productController.CreateCategory)
		categories.GET("/:id/products", productController.GetProductsByCategory)
	}

	// Product routes
	products := v1.Group("/products")
	{
		// No auth checks - all endpoints accessible to everyone
		products.POST("", productController.CreateProduct)
		products.GET("/:id", productController.GetProductByID)
		products.GET("", productController.ListProducts)
		products.PUT("/:id", productController.UpdateProduct)
		products.DELETE("/:id", productController.DeleteProduct)
		products.GET("/search", productController.SearchProducts)
	}

}
