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
        // CREATE routes (existing)
        categories.POST("", productController.CreateCategory)
        categories.POST("/subcategory", productController.CreateSubCategory)
        categories.POST("/sub-subcategory", productController.CreateSubSubCategory)
        
        // GET routes (new)
        categories.GET("", productController.ListCategories)                    // Get all categories
        categories.GET("/:id", productController.GetCategoryByID)              // Get category by ID
        categories.GET("/:id/products", productController.GetProductsByCategory) // existing
        categories.GET("/:category_id/subcategories", productController.GetSubCategoriesByCategory) // Get subcategories by category
        categories.GET("/hierarchy", productController.GetCategoryHierarchy) 
		
		categories.DELETE("/:id", productController.DeleteCategory)// existing
    }

    // SubCategory routes
    subcategories := v1.Group("/subcategories")
    {
        subcategories.GET("/:id", productController.GetSubCategoryByID)                    // Get specific subcategory
        subcategories.GET("/:subcategory_id/sub-subcategories", productController.GetSubSubCategoriesBySubCategory) // Get sub-subcategories

		subcategories.GET("/:id/products", productController.GetProductsBySubCategory)

		subcategories.DELETE("/:id", productController.DeleteSubCategory)

    }

    // SubSubCategory routes
    subSubcategories := v1.Group("/sub-subcategories")
    {
        subSubcategories.GET("/:id", productController.GetSubSubCategoryByID)            // Get specific sub-subcategory
	   subSubcategories.GET("/:id/products", productController.GetProductsBySubSubCategory)
	   subSubcategories.DELETE("/:id", productController.DeleteSubSubCategory)

    }

    // Product routes (existing)
    products := v1.Group("/products")
    {
        products.POST("", productController.CreateProduct)
        products.GET("/:id", productController.GetProductByID)
        products.GET("", productController.ListProducts)
        products.PUT("/:id", productController.UpdateProduct)
        products.DELETE("/:id", productController.DeleteProduct)
        products.GET("/search", productController.SearchProducts)
    }
}