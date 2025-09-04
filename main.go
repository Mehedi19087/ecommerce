package main

import (
	"ecommerce/config"
	"ecommerce/database"
	"ecommerce/internal/auth"
	"ecommerce/internal/cart"
	"ecommerce/internal/catalog"

	//"ecommerce/internal/health"
	"ecommerce/internal/order"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	gin.SetMode(gin.ReleaseMode)
	config.InitGoogleAuth()

	// Initialize repositories
	userRepo := auth.NewUserRepository(db)
	productRepo := catalog.NewProductRepository(db)
	cartRepo := cart.NewCartRepository(db)
	orderRepo := order.NewOrderRepository(db)

	// Initialize services
	userService := auth.NewUserService(userRepo)
	productService := catalog.NewProductService(productRepo)
	cartService := cart.NewCartService(cartRepo, productRepo)
	orderService := order.NewOrderService(orderRepo, cartService) // No db parameter

	// Initialize controllers
	userController := auth.NewUserController(userService)
	productController := catalog.NewProductController(productService)
	cartController := cart.NewCartController(cartService)
	orderController := order.NewOrderController(orderService)

	// Setup router and routes
	router := gin.Default()

	// Add request logging
	router.Use(gin.Logger())
	//router.Use(auth.LocationTrackingMiddleware(db))

	// Configure CORS
	router.Use(cors.New(cors.Config{
        AllowOrigins: []string{
		"https://admin-ecommarce.web.app",
        "https://alrizvan.com",      // 🔥 Your production domain
        "https://www.alrizvan.com",  // 🔥 With www subdomain
        "http://localhost:3000",          // 🔥 For local development
        "http://localhost:3001",          // 🔥 Alternative local port
    },
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:          12 * time.Hour,
    }))


	router.Use(func(c *gin.Context) {
		log.Printf("Request received: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
		log.Printf("Response sent: %d", c.Writer.Status())
	})

	auth.SetupAuthRoutes(router, userController)
	catalog.SetupCatalogRoutes(router, productController)
	cart.SetupCartRoutes(router, cartController)
	order.SetupOrderRoutes(router, orderController)

	//router.GET("/api/v1/visitor-division", health.VisitorDivision)

	// Start server - bind to all interfaces
	log.Println("Starting server on 0.0.0.0:8080")
	if err := router.Run("0.0.0.0:8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	//formating the full project - go fmt ./...
}
