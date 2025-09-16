package handler

import (
	"context"
	"net/http"

	"github.com/chandrasitinjak/integrate-pokeapi/internal/services"
	"github.com/gin-gonic/gin"
)

type PokemonHandler struct {
	service services.PokemonService
}

func NewPokemonHandler(s services.PokemonService) *PokemonHandler {
	return &PokemonHandler{service: s}
}

func (h *PokemonHandler) Sync(c *gin.Context) {
	if err := h.service.Sync(context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "sync success"})
}

func (h *PokemonHandler) GetAll(c *gin.Context) {
	data, err := h.service.GetAll(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
