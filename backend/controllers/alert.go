package controllers

import (
	"net/http"
	"time"
	"tsc_demo/backend/models"

	"github.com/gin-gonic/gin"
)

type BrowserHistoryInput struct {
	ClientID  string    `json:"client_id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	ViewCount int       `json:"view_count"`
	VisitTime time.Time `json:"visit_time"`
}

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
	var browserHistoryInput BrowserHistoryInput
	if err := g.ShouldBindJSON(&browserHistoryInput); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.manager.ReceiveBrowserHistory(browserHistoryInput.ClientID, browserHistoryInput.Name, browserHistoryInput.URL, browserHistoryInput.ViewCount, browserHistoryInput.VisitTime)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{"status": "success"})
}
