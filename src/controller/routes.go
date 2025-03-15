package controller

import (
	"github-repo-tracker/src/handler"
	"github-repo-tracker/src/repository"
	"github-repo-tracker/src/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, repo repository.RepoRepository) {
	repoService := service.NewRepoService(repo)

	repoHandler := handler.NewRepoHandler(repoService)

	api := r.Group("/api")
	{
		repositories := api.Group("/repositories")
		{
			repositories.GET("", repoHandler.GetAllRepositories)
			repositories.GET("/:id", repoHandler.GetRepositoryByID)
			repositories.GET("/name/:fullName", repoHandler.GetRepositoryByFullName)
			repositories.POST("", repoHandler.CreateRepository)
			repositories.PUT("/:id/stats", repoHandler.UpdateRepositoryStats)
			repositories.POST("/:fullName/fetch", repoHandler.FetchAndUpdateRepositoryStats)
		}
	}

	repoService.StartPeriodicStatsFetcher(r.Context(), 1*60*60*1000*1000*1000) // 1 hour in nanoseconds
} 