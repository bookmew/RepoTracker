package handler

import (
	"net/http"

	"RepoTracker/src/service"

	"github.com/gin-gonic/gin"
)

type RepoStatsHandler struct {
	service *service.RepoStatsService
}

func NewRepoStatsHandler(service *service.RepoStatsService) *RepoStatsHandler {
	return &RepoStatsHandler{service: service}
}

func (h *RepoStatsHandler) GetRepoStatsByMint(c *gin.Context) {
	mint := c.Param("mint")
	if mint == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "mint parameter is required"})
		return
	}

	stats, err := h.service.GetByMint(mint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *RepoStatsHandler) GetAllRepoStats(c *gin.Context) {
	stats, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *RepoStatsHandler) FetchRepoStats(c *gin.Context) {
	repoURL := c.Query("repo_url")
	if repoURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "repo_url is required"})
		return
	}
	
	err := h.service.SaveStats(repoURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Repository statistics fetched successfully"})
}

func (h *RepoStatsHandler) GetLatestAllRepoStats(c *gin.Context) {
	stats, err := h.service.GetLatestAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *RepoStatsHandler) GetLatestRepoStatsByMint(c *gin.Context) {
	mint := c.Param("mint")
	if mint == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "mint parameter is required"})
		return
	}

	stats, err := h.service.GetLatestByMint(mint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *RepoStatsHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/repo-stats/history", h.GetAllRepoStats)
		api.GET("/repo-stats/history/:mint", h.GetRepoStatsByMint)
		api.GET("/fetch-repo-stats", h.FetchRepoStats)
		api.GET("/repo-stats/latest", h.GetLatestAllRepoStats)
		api.GET("/repo-stats/latest/:mint", h.GetLatestRepoStatsByMint)
	}
}
