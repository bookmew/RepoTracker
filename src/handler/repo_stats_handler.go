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

func (h *RepoStatsHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/repo-stats", h.GetAllRepoStats)
		api.GET("/repo-stats/:mint", h.GetRepoStatsByMint)
	}
}
