package monitor

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"os/user"
	"path/filepath"
	"tsc_demo/backend/common"
)

// defaultBrowserHistoryGetter is the default BrowserHistoryGetter
var defaultBrowserHistoryGetter BrowserHistoryGetter

// defaultRawHistoryGetter is the default RawHistoryGetter
var defaultRawHistoryGetter RawHistoryGetter

// BrowserHistoryGetter interface
type BrowserHistoryGetter interface {
	BrowserHistory(blacklist []string) ([]common.BrowserHistoryRecord, error)
}

// RawHistoryGetter interface
type RawHistoryGetter interface {
	GetBrowserHistory() ([]common.BrowserHistoryRecord, error)
}

// SetBrowserHistoryGetter sets the default BrowserHistoryGetter
func SetBrowserHistoryGetter(getter BrowserHistoryGetter) {
	defaultBrowserHistoryGetter = getter
}

// NewRawHistoryGetter creates a new instance of RawHistoryGetter
func NewRawHistoryGetter() RawHistoryGetter {
	if defaultRawHistoryGetter == nil {
		defaultRawHistoryGetter = &rawHistoryGetterImpl{}
	}
	return defaultRawHistoryGetter
}

type rawHistoryGetterImpl struct{}

// NewBrowserHistoryGetter creates a new instance of BrowserHistoryGetter
func NewBrowserHistoryGetter() BrowserHistoryGetter {
	if defaultBrowserHistoryGetter == nil {
		defaultBrowserHistoryGetter = &browserHistoryGetterImpl{}
	}
	return defaultBrowserHistoryGetter
}

type browserHistoryGetterImpl struct{}

// SetRawHistoryGetter sets the default RawHistoryGetter
func SetRawHistoryGetter(getter RawHistoryGetter) {
	defaultRawHistoryGetter = getter
}

func (b *browserHistoryGetterImpl) BrowserHistory(blacklist []string) ([]common.BrowserHistoryRecord, error) {
	rawGetter := NewRawHistoryGetter()
	// Get the browser history
	histories, err := rawGetter.GetBrowserHistory()
	if err != nil {
		return nil, err
	}

	// Create a map for faster lookup of the blacklist
	blacklistMap := make(map[string]bool)
	for _, url := range blacklist {
		blacklistMap[url] = true
	}

	// Check each browser record against the blacklist
	var records []common.BrowserHistoryRecord
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

func (r *rawHistoryGetterImpl) GetBrowserHistory() ([]common.BrowserHistoryRecord, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("cannot get current user: %w", err)
	}

	// Edge browser history database path
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

	var records []common.BrowserHistoryRecord
	for rows.Next() {
		var record common.BrowserHistoryRecord
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
