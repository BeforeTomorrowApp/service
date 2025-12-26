package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"north-post/service/internal/domain/v1/models"
	"north-post/service/internal/services"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPromptService for testing
type MockPromptService struct {
	mock.Mock
}

func (m *MockPromptService) GetSystemAddressGenerationPrompt(ctx context.Context, input services.GetSystemAddressGenerationPromptInput) (*services.GetSystemAddressGenerationPromptOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*services.GetSystemAddressGenerationPromptOutput), args.Error(1)
}

func TestPromptHandler_GetSystemAddressGenerationPrompt(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name           string
		language       string
		mockOutput     *services.GetSystemAddressGenerationPromptOutput
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "success",
			language:       "en",
			mockOutput:     &services.GetSystemAddressGenerationPromptOutput{Prompt: "test prompt"},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data":"test prompt"}`,
		},
		{
			name:           "service error",
			language:       "en",
			mockOutput:     nil,
			mockError:      errors.New("service error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"service error"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockPromptService)
			logger := slog.Default()
			handler := NewPromptHandler(mockService, logger)
			mockService.On(
				"GetSystemAddressGenerationPrompt",
				mock.Anything,
				services.GetSystemAddressGenerationPromptInput{Language: models.Language(tt.language)},
			).Return(tt.mockOutput, tt.mockError)
			r := gin.New()
			r.GET("/admin/prompt/system/address", handler.GetSystemAddressGenerationPrompt)
			req := httptest.NewRequest("GET", "/admin/prompt/system/address?language="+tt.language, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
			mockService.AssertExpectations(t)
		})
	}
}
