package main

import (
	"fmt"
	"log"
	"os"

	"RepoTracker/src/handler"
	"RepoTracker/src/repository"
	"RepoTracker/src/service"
	"RepoTracker/src/util"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	// Connect to the database
	db := util.ConnectDB()
	defer db.Close()

	// Initialize repositories
	repoStatsRepo := repository.NewRepoStatsRepository(db)
	tokenRepo := repository.NewTokenRepository(db)

	// Initialize services
	repoStatsService := service.NewRepoStatsService(repoStatsRepo)
	tokenService := service.NewTokenService(tokenRepo)

	// Initialize handlers
	repoStatsHandler := handler.NewRepoStatsHandler(repoStatsService)
	tokenHandler := handler.NewTokenHandler(tokenService)

	// Set up Gin router
	router := gin.Default()

	// Register routes
	repoStatsHandler.RegisterRoutes(router)
	tokenHandler.RegisterRoutes(router)

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	// Start the server
	fmt.Printf("Server running on port %s...\n", port)
	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
