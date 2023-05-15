package common

// BrowserHistoryRecord is a struct to hold a browser history record
type BrowserHistoryRecord struct {
	URL           string `json:"url"`
	Title         string `json:"title"`
	VisitCount    int    `json:"visit_count"`
	LastVisitTime int64  `json:"last_visit_time"`
	IsInBlacklist bool   `json:"is_in_blacklist"`
}
