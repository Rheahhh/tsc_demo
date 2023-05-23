package monitor

import (
	"tsc_demo/backend/common"
	"tsc_demo/backend/storage"
)

// defaultBrowserHistoryGetter is the default BrowserHistoryGetter
var defaultBrowserHistoryGetter BrowserHistoryGetter

// BrowserHistoryGetter interface
type BrowserHistoryGetter interface {
	BrowserHistory(blacklist []string) ([]common.BrowserHistoryRecord, error)
}

// SetBrowserHistoryGetter sets the default BrowserHistoryGetter
func SetBrowserHistoryGetter(getter BrowserHistoryGetter) {
	defaultBrowserHistoryGetter = getter
}

// NewBrowserHistoryGetter creates a new instance of BrowserHistoryGetter
func NewBrowserHistoryGetter() BrowserHistoryGetter {
	if defaultBrowserHistoryGetter == nil {
		defaultBrowserHistoryGetter = &browserHistoryGetterImpl{}
	}
	return defaultBrowserHistoryGetter
}

type browserHistoryGetterImpl struct{}

func (b *browserHistoryGetterImpl) BrowserHistory(blacklist []string) ([]common.BrowserHistoryRecord, error) {
	rawGetter := storage.NewRawHistoryGetter()
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
