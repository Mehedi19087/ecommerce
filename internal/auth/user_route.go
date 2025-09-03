package auth

import (
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, userController *UserController) {
	v1 := router.Group("/api/v1")

	// Public routes
	auth := v1.Group("/auth")
	{
		auth.GET("/auth/google/login",userController.GoogleLogin)
        auth.GET("/auth/google/callback",userController.GoogleCallBack)
		router.GET("/api/v1/admin/visitors/city", GetVisitorCountByCity)
	}

	// Protected routes
	protected := v1.Group("")
	protected.Use(JWTAuthMiddleware())
	{
		protected.POST("/auth/logout", userController.Logout) // ✅ NEW logout endpoint
		protected.GET("/profile", userController.GetProfile)
		protected.PUT("/profile", userController.UpdateProfile)

		addresses := protected.Group("/addresses")
		{
			addresses.GET("", userController.GetAddresses)
			addresses.POST("", userController.CreateAddress)
			addresses.PUT("/:id", userController.UpdateAddress)
			addresses.DELETE("/:id", userController.DeleteAddress)
		}
	}
}
