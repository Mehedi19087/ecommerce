package catalog

import (
    "github.com/gin-gonic/gin"
)

func SetupCatalogRoutes(
    router *gin.Engine,
    productController *ProductController) {

    
    v1 := router.Group("/api/v1")

    v1.POST("/upload", productController.UploadImage)
    
    //category routes
    categories := v1.Group("/categories")
    {
        // CREATE routes
        categories.POST("", productController.CreateCategory)
        categories.POST("/subcategory", productController.CreateSubCategory)
        categories.POST("/sub-subcategory", productController.CreateSubSubCategory)
        
        // GET routes - ✅ FIXED: Use consistent parameter names
        categories.GET("", productController.ListCategories)                    
        categories.GET("/hierarchy", productController.GetCategoryHierarchy)    // Move this before :id routes
        categories.GET("/:id", productController.GetCategoryByID)              
        categories.GET("/:id/products", productController.GetProductsByCategory) 
        categories.GET("/:id/subcategories", productController.GetSubCategoriesByCategory) // ✅ Changed :category_id to :id
        categories.PUT("/:id", productController.UpdateCategory)
		categories.DELETE("/:id", productController.DeleteCategory)

}
    // SubCategory routes
    subcategories := v1.Group("/subcategories")
    {
		subcategories.GET("", productController.ListSubCategories)
    
        subcategories.GET("/:id", productController.GetSubCategoryByID)                    
        subcategories.GET("/:id/sub-subcategories", productController.GetSubSubCategoriesBySubCategory) // ✅ Changed :subcategory_id to :id
        subcategories.GET("/:id/products", productController.GetProductsBySubCategory)
        subcategories.DELETE("/:id", productController.DeleteSubCategory)
    }

    // SubSubCategory routes
    subSubcategories := v1.Group("/sub-subcategories")
    {
        subSubcategories.GET("/:id", productController.GetSubSubCategoryByID)            
        subSubcategories.GET("/:id/products", productController.GetProductsBySubSubCategory)
        subSubcategories.DELETE("/:id", productController.DeleteSubSubCategory)
    }

    // Product routes
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