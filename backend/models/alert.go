package models

import (
	_ "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"time"
)

type BrowserHistoryInputs []BrowserHistoryInput

type BrowserHistoryInput struct {
	ClientID  string `json:"client_id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	ViewCount int    `json:"view_count"`
	VisitTime string `json:"visit_time"`
}

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
	ReceiveBrowserHistory(inputs []BrowserHistoryInput) error
}

type defaultAlertManager struct{}

func NewAlertManager() AlertManager {
	return &defaultAlertManager{}
}

func (m *defaultAlertManager) GetAlerts() ([]Alert, error) {
	var alerts []Alert

	var myTimeStr string
	var createAtStr string

	// 获取当前的黑名单
	b := NewBlacklistManager()
	blacklist, err := b.GetBlacklist()
	if err != nil {
		return nil, err
	}

	// 将黑名单转换为 map 以便快速查找
	blacklistMap := make(map[string]bool)
	for _, url := range blacklist {
		blacklistMap[url] = true
	}

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

		// 如果此告警的URL在黑名单中，则添加到告警列表中
		if _, exists := blacklistMap[a.URL]; exists {
			alerts = append(alerts, a)
		}
	}

	return alerts, nil
}

func (m *defaultAlertManager) ReceiveBrowserHistory(inputs []BrowserHistoryInput) error {
	values := []interface{}{}
	valueStrings := []string{}
	for _, input := range inputs {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, STR_TO_DATE(?, '%Y-%m-%d %H:%i:%s'), NOW())")
		values = append(values, input.ClientID, input.Name, input.URL, input.ViewCount, input.VisitTime)
	}

	query := fmt.Sprintf(`
        INSERT INTO alert (client_id, name, url, view_count, view_timestamp, created_at)
        VALUES %s
        ON DUPLICATE KEY UPDATE
        view_count = view_count + VALUES(view_count),
        view_timestamp = VALUES(view_timestamp)
    `, strings.Join(valueStrings, ","))

	_, err := db.Exec(query, values...)
	if err != nil {
		fmt.Println("mysql err:", err)
	}
	return err
}
