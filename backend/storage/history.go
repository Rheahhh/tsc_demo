package storage

import (
	"database/sql"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"tsc_demo/backend/common"

	_ "github.com/mattn/go-sqlite3"
)

// RawHistoryGetter interface
type RawHistoryGetter interface {
	GetBrowserHistory() ([]common.BrowserHistoryRecord, error)
}

// defaultRawHistoryGetter is the default RawHistoryGetter
var defaultRawHistoryGetter RawHistoryGetter

// NewRawHistoryGetter creates a new instance of RawHistoryGetter
func NewRawHistoryGetter() RawHistoryGetter {
	if defaultRawHistoryGetter == nil {
		defaultRawHistoryGetter = &rawHistoryGetterImpl{}
	}
	return defaultRawHistoryGetter
}

// SetRawHistoryGetter sets the default RawHistoryGetter
func SetRawHistoryGetter(getter RawHistoryGetter) {
	defaultRawHistoryGetter = getter
}

type rawHistoryGetterImpl struct{}

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

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_busy_timeout=2000", historyDB))
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
