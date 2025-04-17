package application

import (
	"context"
	"errors"
	"testing"

	tweetsdomain "github.com/nicodelara/uala-challenge/internal/tweets/domain"
	usersdomain "github.com/nicodelara/uala-challenge/internal/users/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTweetRepository struct {
	mock.Mock
}

func (m *mockTweetRepository) SaveTweet(ctx context.Context, tweet *tweetsdomain.Tweet) error {
	args := m.Called(ctx, tweet)
	return args.Error(0)
}

type mockUserChecker struct {
	mock.Mock
}

func (m *mockUserChecker) GetUser(username string) (*usersdomain.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usersdomain.User), args.Error(1)
}

func (m *mockUserChecker) GetFollowings(ctx context.Context, username string) ([]string, error) {
	args := m.Called(ctx, username)
	return args.Get(0).([]string), args.Error(1)
}

func TestTweetService_CreateTweet(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		content       string
		checkerSetup  func(*mockUserChecker)
		repoSetup     func(*mockTweetRepository)
		expectedError error
	}{
		{
			name:     "successful creation",
			username: "testuser",
			content:  "Hello, world!",
			checkerSetup: func(m *mockUserChecker) {
				m.On("GetUser", "testuser").Return(&usersdomain.User{}, nil)
			},
			repoSetup: func(m *mockTweetRepository) {
				m.On("SaveTweet", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "user does not exist",
			username: "nonexistent",
			content:  "Hello, world!",
			checkerSetup: func(m *mockUserChecker) {
				m.On("GetUser", "nonexistent").Return(nil, errors.New("user not found"))
			},
			expectedError: ErrUserNotFound,
		},
		{
			name:     "repository error",
			username: "testuser",
			content:  "Hello, world!",
			checkerSetup: func(m *mockUserChecker) {
				m.On("GetUser", "testuser").Return(&usersdomain.User{}, nil)
			},
			repoSetup: func(m *mockTweetRepository) {
				m.On("SaveTweet", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := new(mockUserChecker)
			repo := new(mockTweetRepository)
			tt.checkerSetup(checker)
			if tt.repoSetup != nil {
				tt.repoSetup(repo)
			}

			service := NewTweetService(repo, checker)
			tweet, err := service.CreateTweet(context.Background(), tt.username, tt.content)

			if tt.expectedError != nil {
				assert.Nil(t, tweet)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NotNil(t, tweet)
				assert.NoError(t, err)
				assert.Equal(t, tt.username, tweet.Username)
				assert.Equal(t, tt.content, tweet.Content)
			}
			checker.AssertExpectations(t)
			repo.AssertExpectations(t)
		})
	}
}
