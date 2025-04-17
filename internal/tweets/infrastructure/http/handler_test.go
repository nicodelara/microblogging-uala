package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nicodelara/uala-challenge/internal/tweets/application"
	"github.com/nicodelara/uala-challenge/internal/tweets/domain"
	"github.com/nicodelara/uala-challenge/internal/tweets/domain/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTweetService struct {
	mock.Mock
}

func (m *mockTweetService) CreateTweet(ctx context.Context, username, content string) (*domain.Tweet, error) {
	args := m.Called(ctx, username, content)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tweet), args.Error(1)
}

func setupRouter(service ports.TweetService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := NewTweetHandler(service)
	router.POST("/tweets", handler.CreateTweet)
	return router
}

func TestTweetHandler_CreateTweet(t *testing.T) {
	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		requestBody    map[string]string
		serviceSetup   func(*mockTweetService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful creation",
			requestBody: map[string]string{
				"username": "testuser",
				"content":  "Hello, world!",
			},
			serviceSetup: func(m *mockTweetService) {
				tweet := &domain.Tweet{
					ID:        "123",
					Username:  "testuser",
					Content:   "Hello, world!",
					CreatedAt: createdAt,
				}
				m.On("CreateTweet", mock.Anything, "testuser", "Hello, world!").Return(tweet, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":"123","username":"testuser","content":"Hello, world!","createdAt":"2024-01-01T00:00:00Z"}`,
		},
		{
			name: "missing required fields",
			requestBody: map[string]string{
				"username": "testuser",
			},
			serviceSetup: func(m *mockTweetService) {
				// No se espera llamada al servicio
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Key: 'createTweetRequest.Content' Error:Field validation for 'Content' failed on the 'required' tag"}`,
		},
		{
			name: "user does not exist",
			requestBody: map[string]string{
				"username": "nonexistent",
				"content":  "Hello, world!",
			},
			serviceSetup: func(m *mockTweetService) {
				m.On("CreateTweet", mock.Anything, "nonexistent", "Hello, world!").Return(nil, application.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"user does not exist"}`,
		},
		{
			name: "content too long",
			requestBody: map[string]string{
				"username": "testuser",
				"content":  "This is a very long tweet that exceeds the maximum allowed length of 280 characters. This is a very long tweet that exceeds the maximum allowed length of 280 characters. This is a very long tweet that exceeds the maximum allowed length of 280 characters. This is a very long tweet that exceeds the maximum allowed length of 280 characters.",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"content exceeds 280 characters"}`,
		},
		{
			name: "internal server error",
			requestBody: map[string]string{
				"username": "testuser",
				"content":  "Hello, world!",
			},
			serviceSetup: func(m *mockTweetService) {
				m.On("CreateTweet", mock.Anything, "testuser", "Hello, world!").Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"assert.AnError general error for testing"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := new(mockTweetService)
			if tt.serviceSetup != nil {
				tt.serviceSetup(service)
			}

			router := setupRouter(service)
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/tweets", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
			service.AssertExpectations(t)
		})
	}
}
