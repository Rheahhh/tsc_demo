package controllers

import (
	"net/http"
	"tsc_demo/backend/models"

	"github.com/gin-gonic/gin"
)

type AlertController struct {
	manager models.AlertManager
}

func NewAlertController(m models.AlertManager) *AlertController {
	return &AlertController{manager: m}
}

func (c *AlertController) GetAlerts(g *gin.Context) {
	alerts, err := c.manager.GetAlerts()
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, alerts)
}

func (c *AlertController) ReceiveBrowserHistory(g *gin.Context) {
	var browserHistoryInputs models.BrowserHistoryInputs
	if err := g.ShouldBindJSON(&browserHistoryInputs); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.manager.ReceiveBrowserHistory(browserHistoryInputs)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, gin.H{"status": "success"})
}
