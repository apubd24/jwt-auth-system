package routes

import (
	"jwt-auth-backend/handlers"
	"jwt-auth-backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	/*
		For production, set AllowOrigins to your actual domain (e.g., https://yourfrontend.com). You can also use an environment variable:

			frontendURL := os.Getenv("FRONTEND_URL") // e.g., https://myapp.com
			r.Use(cors.New(cors.Config{
		    AllowOrigins: []string{frontendURL},
		    // ...
		}))
	*/

	// CORS for React (adjust to your frontend URL)
	r.Use(cors.New(cors.Config{
		// AllowOrigins: []string{"*"},
		AllowOrigins:     []string{"http://172.17.18.188:3000"}, // React dev server OR Front end server if its diffrent replaceIPwith localhost
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // if you use cookies/auth headers
	}))

	// Public routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// Protected routes (require JWT)
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// Token validation
		api.GET("/validate", handlers.ValidateToken)

		// Device routes – read for all authenticated
		api.GET("/devices", handlers.GetAllDevices)

		// Admin-only device mutations
		adminDevices := api.Group("/devices")
		adminDevices.Use(middleware.RequireRole("admin"))
		{
			adminDevices.POST("/", handlers.CreateDevice)
			adminDevices.GET("/:id", handlers.GetDeviceByID)
			adminDevices.PUT("/:id", handlers.UpdateDevice)
			adminDevices.DELETE("/:id", handlers.DeleteDevice)
		}

		// Device Products – read for all authenticated
		api.GET("/products", handlers.GetAllProducts)

		// adminProducts := api.Group("/products")
		// adminProducts.POST("", handlers.CreateProduct) // 👈 route: /api/products/

		// Admin-only Product mutations
		adminProducts := api.Group("/products")
		adminProducts.Use(middleware.RequireRole("admin"))
		{
			adminProducts.POST("", handlers.CreateProduct)
		}

		// User management – only admin
		adminUsers := api.Group("/users")
		adminUsers.Use(middleware.RequireRole("admin"))
		{
			adminUsers.GET("/", handlers.GetAllUsers)
			adminUsers.GET("/:id", handlers.GetUserByID)
			adminUsers.POST("", handlers.CreateUser)
			adminUsers.PUT("/:id", handlers.UpdateUser)
			adminUsers.DELETE("/:id", handlers.DeleteUser)
			adminUsers.PUT("/:id/password", handlers.ChangePassword)
		}
	}
	return r
}
