package models

type Blacklist struct {
	URL       string `json:"url"`
	IsDeleted bool   `json:"is_deleted"`
}

type BlacklistManager interface {
	GetBlacklist() ([]string, error)
	ManageBlacklist(action string, url string) error
}

type defaultBlacklistManager struct{}

func NewBlacklistManager() BlacklistManager {
	return &defaultBlacklistManager{}
}

func (m *defaultBlacklistManager) GetBlacklist() ([]string, error) {
	var blacklist []string

	rows, err := db.Query("SELECT url FROM blacklist WHERE is_deleted = FALSE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var url string
		err := rows.Scan(&url)
		if err != nil {
			return nil, err
		}
		blacklist = append(blacklist, url)
	}

	return blacklist, nil
}

func (m *defaultBlacklistManager) ManageBlacklist(action string, url string) error {
	var err error

	switch action {
	case "add":
		_, err = db.Exec("INSERT INTO blacklist (url, is_deleted) VALUES (?, FALSE) ON DUPLICATE KEY UPDATE is_deleted=FALSE", url)
	case "delete":
		_, err = db.Exec("UPDATE blacklist SET is_deleted = TRUE WHERE url = ?", url)
	}

	return err
}
