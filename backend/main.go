// main.go
package main

import (
	"net/http"
	"tsc_demo/backend/monitor" // 你需要替换为你的实际包路径

	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3" // SQLite3 驱动
)

type RequestData struct {
	Blacklist []string `json:"blacklist"`
}

func main() {
	//temp()
	run()

}

//localhost:8080/monitor
func run() {
	r := gin.Default()

	r.Use(cors.Default())

	r.POST("/monitor", func(c *gin.Context) {
		var data RequestData

		// Bind the JSON data from the request body to the variable
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Println("blacklist:  ", data.Blacklist)
		// Now you can use data.Blacklist

		browserHistoryGetter := monitor.NewBrowserHistoryGetter()
		result, err := browserHistoryGetter.BrowserHistory(data.Blacklist)

		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"result": result,
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
