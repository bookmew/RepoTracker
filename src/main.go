package main

import (
	"log"
	"githubsource-demo/src/controller"
	"githubsource-demo/src/repository"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := repository.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	repoRepository := repository.NewGithubRepoRepository(db)
	
	r := gin.Default()

	controller.RegisterRoutes(r, repoRepository)

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}