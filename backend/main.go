package main

import (
	"jwt-auth-backend/database"
	"jwt-auth-backend/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	database.Connect()
	r := routes.SetupRouter()
	r.Static("/uploads", "./uploads") //For logo upload
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	log.Printf("Server running on port %s", port)
	r.Run(":" + port)
}
