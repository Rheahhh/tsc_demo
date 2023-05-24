package main

import (
	"tsc_demo/backend/controllers"
	"tsc_demo/backend/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(cors.Default())

	bm := models.NewBlacklistManager()
	bc := controllers.NewBlacklistController(bm)

	am := models.NewAlertManager()
	ac := controllers.NewAlertController(am)

	r.GET("/blacklist", bc.GetBlacklist)
	r.POST("/blacklist", bc.ManageBlacklist)

	r.POST("/browser-history", ac.ReceiveBrowserHistory)
	r.GET("/alerts", ac.GetAlerts)

	r.Run(":8080")
}
