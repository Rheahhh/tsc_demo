package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMonitorHandler(t *testing.T) {
	// Create a new instance of your router
	router := gin.Default()
	router.POST("/monitor", func(c *gin.Context) { /* your handler code here */ })

	// Test case 1: The JSON data format is correct
	requestData := RequestData{Blacklist: []string{"http://example.com"}}
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		t.Fatalf("failed to marshal request data: %v", err)
	}
	req, err := http.NewRequest("POST", "/monitor", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("failed to create new HTTP request: %v", err)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}
