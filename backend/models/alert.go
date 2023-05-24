package models

import (
	_ "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Alert struct {
	ID            int       `json:"id"`
	ClientID      string    `json:"client_id"`
	Name          string    `json:"name"`
	URL           string    `json:"url"`
	ViewCount     int       `json:"view_count"`
	ViewTimestamp time.Time `json:"view_timestamp"`
	CreatedAt     time.Time `json:"created_at"`
}

type AlertManager interface {
	GetAlerts() ([]Alert, error)
	ReceiveBrowserHistory(clientID, name, url string, viewCount int, visitTime time.Time) error
}

type defaultAlertManager struct{}

func NewAlertManager() AlertManager {
	return &defaultAlertManager{}
}

func (m *defaultAlertManager) GetAlerts() ([]Alert, error) {
	var alerts []Alert

	var myTimeStr string
	var createAtStr string

	rows, err := db.Query("SELECT id, client_id, name, url, view_count, view_timestamp, created_at FROM alert")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a Alert
		err := rows.Scan(&a.ID, &a.ClientID, &a.Name, &a.URL, &a.ViewCount, &myTimeStr, &createAtStr)
		if err != nil {
			return nil, err
		}
		a.ViewTimestamp, err = time.Parse("2006-01-02 15:04:05", myTimeStr)
		if err != nil {
			return nil, err
		}
		a.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createAtStr)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, a)
	}

	return alerts, nil
}

func (m *defaultAlertManager) ReceiveBrowserHistory(clientID, name, url string, viewCount int, visitTime time.Time) error {

	_, err := db.Exec(`
		INSERT INTO alert (client_id, name, url, view_count, view_timestamp, created_at)
		VALUES (?, ?, ?, ?, ?, NOW())
		ON DUPLICATE KEY UPDATE
		view_count = view_count + VALUES(view_count),
		view_timestamp = VALUES(view_timestamp)
	`, clientID, name, url, viewCount, visitTime)

	return err
}
