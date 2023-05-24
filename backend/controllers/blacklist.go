package controllers

import (
	"net/http"
	"tsc_demo/backend/models"

	"github.com/gin-gonic/gin"
)

// BlacklistInput represents the input for the manage blacklist API.
type BlacklistInput struct {
	Action string `json:"action"`
	Url    string `json:"url"`
}

type BlacklistController struct {
	manager models.BlacklistManager
}

func NewBlacklistController(m models.BlacklistManager) *BlacklistController {
	return &BlacklistController{manager: m}
}

func (c *BlacklistController) GetBlacklist(g *gin.Context) {
	blacklist, err := c.manager.GetBlacklist()
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, blacklist)
}

func (c *BlacklistController) ManageBlacklist(g *gin.Context) {
	var input BlacklistInput
	if err := g.ShouldBindJSON(&input); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.manager.ManageBlacklist(input.Action, input.Url)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{"status": "success"})
}
