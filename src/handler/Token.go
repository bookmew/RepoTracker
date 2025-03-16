package handler

import (
	"net/http"

	"RepoTracker/src/service"

	"github.com/gin-gonic/gin"
)

type TokenHandler struct {
	service *service.TokenService
}

func NewTokenHandler(service *service.TokenService) *TokenHandler {
	return &TokenHandler{service: service}
}

func (h *TokenHandler) GetAllTokens(c *gin.Context) {
	tokens, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *TokenHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/tokens", h.GetAllTokens)
	}
}
