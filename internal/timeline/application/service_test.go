package application

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	timelinedomain "github.com/nicodelara/uala-challenge/internal/timeline/domain"
	usersdomain "github.com/nicodelara/uala-challenge/internal/users/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTimelineRepository struct {
	mock.Mock
}

func (m *mockTimelineRepository) GetTweetsForUsers(ctx context.Context, usernames []string, offset, limit int) ([]timelinedomain.Tweet, error) {
	args := m.Called(ctx, usernames, offset, limit)
	return args.Get(0).([]timelinedomain.Tweet), args.Error(1)
}

func (m *mockTimelineRepository) SaveTweet(ctx context.Context, tweet timelinedomain.Tweet) error {
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

type mockCacheRepository struct {
	mock.Mock
}

func (m *mockCacheRepository) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *mockCacheRepository) Set(ctx context.Context, key string, value string) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func TestTimelineService_GetTimeline(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name          string
		username      string
		offset        int
		limit         int
		checkerSetup  func(*mockUserChecker)
		repoSetup     func(*mockTimelineRepository)
		cacheSetup    func(*mockCacheRepository)
		expectedError error
	}{
		{
			name:     "successful retrieval from cache",
			username: "testuser",
			offset:   0,
			limit:    10,
			checkerSetup: func(m *mockUserChecker) {
				m.On("GetUser", "testuser").Return(&usersdomain.User{}, nil)
				m.On("GetFollowings", mock.Anything, "testuser").Return([]string{"user1", "user2"}, nil)
			},
			cacheSetup: func(m *mockCacheRepository) {
				tweets := []timelinedomain.Tweet{
					{ID: "1", Username: "user1", Content: "Tweet 1", CreatedAt: now},
					{ID: "2", Username: "user2", Content: "Tweet 2", CreatedAt: now},
				}
				tweetsJSON, _ := json.Marshal(tweets)
				m.On("Get", mock.Anything, "timeline:testuser:offset=0:limit=10").Return(string(tweetsJSON), nil)
			},
			expectedError: nil,
		},
		{
			name:     "successful retrieval from repository",
			username: "testuser",
			offset:   0,
			limit:    10,
			checkerSetup: func(m *mockUserChecker) {
				m.On("GetUser", "testuser").Return(&usersdomain.User{}, nil)
				m.On("GetFollowings", mock.Anything, "testuser").Return([]string{"user1", "user2"}, nil)
			},
			cacheSetup: func(m *mockCacheRepository) {
				m.On("Get", mock.Anything, "timeline:testuser:offset=0:limit=10").Return("", nil)
				tweets := []timelinedomain.Tweet{
					{ID: "1", Username: "user1", Content: "Tweet 1", CreatedAt: now},
					{ID: "2", Username: "user2", Content: "Tweet 2", CreatedAt: now},
				}
				tweetsJSON, _ := json.Marshal(tweets)
				m.On("Set", mock.Anything, "timeline:testuser:offset=0:limit=10", string(tweetsJSON)).Return(nil)
			},
			repoSetup: func(m *mockTimelineRepository) {
				tweets := []timelinedomain.Tweet{
					{ID: "1", Username: "user1", Content: "Tweet 1", CreatedAt: now},
					{ID: "2", Username: "user2", Content: "Tweet 2", CreatedAt: now},
				}
				m.On("GetTweetsForUsers", mock.Anything, []string{"user1", "user2"}, 0, 10).Return(tweets, nil)
			},
			expectedError: nil,
		},
		{
			name:     "user does not exist",
			username: "nonexistent",
			offset:   0,
			limit:    10,
			checkerSetup: func(m *mockUserChecker) {
				m.On("GetUser", "nonexistent").Return(nil, errors.New("user does not exist"))
			},
			expectedError: errors.New("user does not exist"),
		},
		{
			name:     "no followings",
			username: "testuser",
			offset:   0,
			limit:    10,
			checkerSetup: func(m *mockUserChecker) {
				m.On("GetUser", "testuser").Return(&usersdomain.User{}, nil)
				m.On("GetFollowings", mock.Anything, "testuser").Return([]string{}, nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := new(mockUserChecker)
			repo := new(mockTimelineRepository)
			cache := new(mockCacheRepository)
			tt.checkerSetup(checker)
			if tt.repoSetup != nil {
				tt.repoSetup(repo)
			}
			if tt.cacheSetup != nil {
				tt.cacheSetup(cache)
			}

			service := NewTimelineService(repo, checker, cache)
			timeline, err := service.GetTimeline(context.Background(), tt.username, tt.offset, tt.limit)

			if tt.expectedError != nil {
				assert.Nil(t, timeline)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NotNil(t, timeline)
				assert.NoError(t, err)
				assert.Equal(t, tt.username, timeline.Username)
			}
			checker.AssertExpectations(t)
			repo.AssertExpectations(t)
			cache.AssertExpectations(t)
		})
	}
}
