package handler

import (
	"net/http"
	"strconv"

	"github-repo-tracker/src/model"
	"github-repo-tracker/src/service"

	"github.com/gin-gonic/gin"
)

type RepoHandler struct {
	service *service.RepoService
}

func NewRepoHandler(service *service.RepoService) *RepoHandler {
	return &RepoHandler{service: service}
}

func (h *RepoHandler) GetAllRepositories(c *gin.Context) {
	repos, err := h.service.GetAllRepositories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, repos)
}

func (h *RepoHandler) GetRepositoryByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	repo, err := h.service.GetRepositoryByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	c.JSON(http.StatusOK, repo)
}

func (h *RepoHandler) GetRepositoryByFullName(c *gin.Context) {
	fullName := c.Param("fullName")
	if fullName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Repository full name is required"})
		return
	}

	repo, err := h.service.GetRepositoryByFullName(c.Request.Context(), fullName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	c.JSON(http.StatusOK, repo)
}

func (h *RepoHandler) CreateRepository(c *gin.Context) {
	var repo model.Repository
	if err := c.ShouldBindJSON(&repo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if repo.Name == "" || repo.Owner == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Repository name and owner are required"})
		return
	}

	if err := h.service.CreateRepository(c.Request.Context(), &repo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, repo)
}

func (h *RepoHandler) UpdateRepositoryStats(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	var stats model.RepositoryStats
	if err := c.ShouldBindJSON(&stats); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateRepositoryStats(c.Request.Context(), id, &stats); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Repository statistics updated successfully"})
}

func (h *RepoHandler) FetchAndUpdateRepositoryStats(c *gin.Context) {
	fullName := c.Param("fullName")
	if fullName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Repository full name is required"})
		return
	}

	if err := h.service.FetchAndUpdateRepositoryStats(c.Request.Context(), fullName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Repository statistics fetched and updated successfully"})
} 