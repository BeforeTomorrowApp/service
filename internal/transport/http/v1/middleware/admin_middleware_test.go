package middleware

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestContext(authHeader string) (
	*gin.Context,
	*httptest.ResponseRecorder,
	*bytes.Buffer,
) {
	w := httptest.NewRecorder() // w stands for writer/recorder
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("GET", "/test", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	c.Request = req
	// capture log output
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	return c, w, &logBuffer
}

func teardownTest() {
	log.SetOutput(os.Stderr)
}

func TestAdminAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	c, w, logBuffer := setupTestContext("")
	defer teardownTest()

	middleware := AdminAuthMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.True(t, c.IsAborted())

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Authorization header required.", response["error"])

	logOutput := logBuffer.String()
	assert.Contains(t, logOutput, AUTH_LOG_PREFIX)
	assert.Contains(t, logOutput, "Invalid request")
	assert.Contains(t, logOutput, "Authorization header required.")
}

func TestAdminAuthMiddleware_InvalidFormat_NoBearerPrefix(t *testing.T) {
	c, w, logBuffer := setupTestContext("InvalidToken123")
	defer teardownTest()

	middleware := AdminAuthMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.True(t, c.IsAborted())

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Authorization header format must be Bearer")

	logOutput := logBuffer.String()
	assert.Contains(t, logOutput, AUTH_LOG_PREFIX)
	assert.Contains(t, logOutput, "Invalid header")
}

func TestAdminAuthMiddleware_InvalidHeader(t *testing.T) {
	c, w, logBuffer := setupTestContext("Basic token123")
	defer teardownTest()

	middleware := AdminAuthMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.True(t, c.IsAborted())

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Authorization header format must be Bearer")

	logOutput := logBuffer.String()
	assert.Contains(t, logOutput, AUTH_LOG_PREFIX)
	assert.Contains(t, logOutput, "Invalid header")
}

func TestAdminAuthMiddleware_ValidFormat(t *testing.T) {
	c, w, logBuffer := setupTestContext("Bearer valid_token")
	defer teardownTest()

	middleware := AdminAuthMiddleware()
	middleware(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.False(t, c.IsAborted())

	logOutput := logBuffer.String()
	assert.NotContains(t, logOutput, AUTH_LOG_PREFIX)
}
