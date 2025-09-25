package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGreetHandler(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParam     string
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "greet with default name",
			queryParam:     "",
			expectedStatus: http.StatusOK,
			expectedMsg:    "Hello, World!",
		},
		{
			name:           "greet with custom name",
			queryParam:     "?name=Alice",
			expectedStatus: http.StatusOK,
			expectedMsg:    "Hello, Alice!",
		},
		{
			name:           "greet with empty name",
			queryParam:     "?name=",
			expectedStatus: http.StatusOK,
			expectedMsg:    "Error: Name cannot be empty",
		},
		{
			name:           "greet with whitespace name",
			queryParam:     "?name=" + "%20%20%20", // URL encoded spaces
			expectedStatus: http.StatusOK,
			expectedMsg:    "Error: Name cannot be empty",
		},
		{
			name:           "greet with special characters",
			queryParam:     "?name=José",
			expectedStatus: http.StatusOK,
			expectedMsg:    "Hello, José!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create handler
			handler := NewGreetHandler()

			// Create request
			req := httptest.NewRequest("GET", "/greet"+tt.queryParam, nil)
			w := httptest.NewRecorder()

			// Create Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call handler
			handler.Greet(c)

			// Check response
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response GreetResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedMsg, response.Message)
		})
	}
}

func TestGreetHandler_JSONResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewGreetHandler()
	req := httptest.NewRequest("GET", "/greet?name=Test", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.Greet(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response GreetResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Message)
}