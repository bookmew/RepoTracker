package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"RepoTracker/src/handler"
	"RepoTracker/src/repository"
	"RepoTracker/src/service"
	"RepoTracker/src/util"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	db := util.ConnectDB()
	defer db.Close()

	repoStatsRepo := repository.NewRepoStatsRepository(db)
	tokenRepo := repository.NewTokenRepository(db)

	repoStatsService := service.NewRepoStatsService(repoStatsRepo)
	tokenService := service.NewTokenService(tokenRepo)

	repoStatsHandler := handler.NewRepoStatsHandler(repoStatsService)
	tokenHandler := handler.NewTokenHandler(tokenService)

	fetchIntervalStr := os.Getenv("FETCH_INTERVAL_MINUTES")
	fetchInterval := 1 * time.Minute
	if fetchIntervalStr != "" {
		if minutes, err := time.ParseDuration(fetchIntervalStr + "m"); err == nil {
			fetchInterval = minutes
		}
	}
	
	scheduler := util.NewRepoStatsScheduler(repoStatsRepo, tokenRepo, fetchInterval)
	scheduler.Start()
	defer scheduler.Stop()

	router := gin.Default()

	repoStatsHandler.RegisterRoutes(router)
	tokenHandler.RegisterRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on port %s...\n", port)
	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
