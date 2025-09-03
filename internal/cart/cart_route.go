package cart

import (
	"github.com/gin-gonic/gin"
	//"ecommerce/internal/auth"
)

func SetupCartRoutes(
	router *gin.Engine,
	cartController *CartController) {

	v1 := router.Group("/api/v1")

	// Cart routes - these should be protected by auth middleware in production
	cart := v1.Group("/cart")
	// cart.Use(auth.AuthMiddleware()) // Add this line
	// {
	cart.GET("", cartController.GetCart)
	cart.POST("/items", cartController.AddItemToCart)
	cart.PUT("/items/:id", cartController.UpdateCartItem)
	cart.DELETE("/items/:id", cartController.RemoveCartItem)
	cart.DELETE("", cartController.ClearCart)
}
