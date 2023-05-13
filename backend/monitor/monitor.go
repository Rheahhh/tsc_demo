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
	IsInBlacklist bool   `json:"isInBlacklist"`
}

func BrowserHistory(blacklist []string) ([]BrowserHistoryRecord, error) {
	// Get the browser history
	history, err := getBrowserHistory()
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
	for _, url := range history {
		record := BrowserHistoryRecord{
			URL: url,
		}

		// If the record's URL is in the blacklist, mark it as such
		if _, ok := blacklistMap[url]; ok {
			record.IsInBlacklist = true
		}

		records = append(records, record)
	}

	return records, nil
}

func getBrowserHistory() ([]string, error) {
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

	rows, err := db.Query("SELECT url FROM urls")
	if err != nil {
		return nil, fmt.Errorf("cannot query database: %w", err)
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, fmt.Errorf("cannot scan row: %w", err)
		}
		urls = append(urls, url)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return urls, nil
}
