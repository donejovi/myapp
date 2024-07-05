package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"myapp/database"
	"myapp/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoginSuccess(t *testing.T) {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}
	database.DB = db

	// Auto-migrate model
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Setup Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Endpoint handler
	router.POST("/login", Login)

	// Mock request JSON
	jsonStr := `{"phone_number": "08123456789", "pin": "123456"}`

	// Create mock user
	mockUser := models.User{
		PhoneNumber: "08123456789",
		PIN:         "123456",
	}
	err = db.Create(&mockUser).Error
	if err != nil {
		t.Fatalf("Failed to create mock user: %v", err)
	}

	// Perform request
	req, err := http.NewRequest("POST", "/login", strings.NewReader(jsonStr))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse actual JSON response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check status and result
	assert.Equal(t, "SUCCESS", response["status"])

	result := response["result"].(map[string]interface{})
	assert.NotNil(t, result["access_token"])
	assert.NotNil(t, result["refresh_token"])

	// Clean up mock data
	err = db.Delete(&mockUser).Error
	if err != nil {
		t.Fatalf("Failed to delete mock user: %v", err)
	}
}
