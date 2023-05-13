// main.go
package main

import (
	"database/sql"
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
		result, err := monitor.BrowserHistory(data.Blacklist)
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

func temp() {
	// 你的 Edge 浏览器历史记录数据库的路径。注意，你需要替换为你自己的路径。
	dbPath := "C:/Users/zhengyujia/AppData/Local/Microsoft/Edge/User Data/Default/History"

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT url, title, visit_count, last_visit_time FROM urls")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var url string
		var title string
		var visitCount int
		var lastVisitTime int
		if err := rows.Scan(&url, &title, &visitCount, &lastVisitTime); err != nil {
			panic(err)
		}
		fmt.Printf("URL: %s, Title: %s, Visit Count: %d, Last Visit Time: %d\n", url, title, visitCount, lastVisitTime)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}
}
