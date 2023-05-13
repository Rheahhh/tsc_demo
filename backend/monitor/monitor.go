package monitor

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"os/user"
	"path/filepath"
)

// BrowserHistoryRecord is a struct to hold a browser history record
type BrowserHistoryRecord struct {
	URL           string `json:"url"`
	Title         string `json:"title"`
	VisitCount    int    `json:"visit_count"`
	LastVisitTime int64  `json:"last_visit_time"`
	IsInBlacklist bool   `json:"isInBlacklist"`
}

func BrowserHistory(blacklist []string) ([]BrowserHistoryRecord, error) {
	// Get the browser history
	histories, err := getBrowserHistory()
	if err != nil {
		return nil, err
	}

	// Create a map for faster lookup of the blacklist
	blacklistMap := make(map[string]bool)
	for _, url := range blacklist {
		blacklistMap[url] = true
	}

	// Check each browser record against the blacklist
	var records []BrowserHistoryRecord
	for _, history := range histories {

		// If the record's URL is in the blacklist, mark it as such
		if _, ok := blacklistMap[history.URL]; ok {
			history.IsInBlacklist = true
		} else {
			history.IsInBlacklist = false
		}
		records = append(records, history)
	}
	return records, nil
}

func getBrowserHistory() ([]BrowserHistoryRecord, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("cannot get current user: %w", err)
	}

	// Edge 浏览器历史记录数据库路径
	historyDB := filepath.Join(usr.HomeDir, "AppData", "Local", "Microsoft", "Edge", "User Data", "Default", "History")

	if _, err := os.Stat(historyDB); os.IsNotExist(err) {
		return nil, fmt.Errorf("History database does not exist: %w", err)
	}

	db, err := sql.Open("sqlite3", historyDB)
	if err != nil {
		return nil, fmt.Errorf("cannot open database: %w", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT url, title, visit_count, last_visit_time FROM urls")
	if err != nil {
		return nil, fmt.Errorf("cannot query database: %w", err)
	}
	defer rows.Close()

	var records []BrowserHistoryRecord
	for rows.Next() {
		var record BrowserHistoryRecord
		if err := rows.Scan(&record.URL, &record.Title, &record.VisitCount, &record.LastVisitTime); err != nil {
			return nil, fmt.Errorf("cannot scan row: %w", err)
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return records, nil
}
