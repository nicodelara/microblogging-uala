package application

import (
	"context"
	"testing"

	"github.com/nicodelara/microblogging-uala/internal/users/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) GetUser(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepository) SaveUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

type mockFollowRepository struct {
	mock.Mock
}

func (m *mockFollowRepository) GetFollowings(ctx context.Context, username string) ([]string, error) {
	args := m.Called(ctx, username)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockFollowRepository) FollowUser(ctx context.Context, follow *domain.Follow) error {
	args := m.Called(ctx, follow)
	return args.Error(0)
}

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		email         string
		userRepoSetup func(*mockUserRepository)
		expectedError error
	}{
		{
			name:     "successful creation",
			username: "testuser",
			email:    "test@example.com",
			userRepoSetup: func(m *mockUserRepository) {
				m.On("GetUser", mock.Anything, "testuser").Return(nil, nil)
				m.On("GetUserByEmail", mock.Anything, "test@example.com").Return(nil, nil)
				m.On("SaveUser", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "username already exists",
			username: "existinguser",
			email:    "test@example.com",
			userRepoSetup: func(m *mockUserRepository) {
				m.On("GetUser", mock.Anything, "existinguser").Return(&domain.User{}, nil)
			},
			expectedError: ErrUserAlreadyExists,
		},
		{
			name:     "email already exists",
			username: "testuser",
			email:    "existing@example.com",
			userRepoSetup: func(m *mockUserRepository) {
				m.On("GetUser", mock.Anything, "testuser").Return(nil, nil)
				m.On("GetUserByEmail", mock.Anything, "existing@example.com").Return(&domain.User{}, nil)
			},
			expectedError: ErrEmailAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(mockUserRepository)
			followRepo := new(mockFollowRepository)
			tt.userRepoSetup(userRepo)

			service := NewUserService(userRepo, followRepo)
			user, err := service.CreateUser(context.Background(), tt.username, tt.email)

			if tt.expectedError != nil {
				assert.Nil(t, user)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NotNil(t, user)
				assert.NoError(t, err)
				assert.Equal(t, tt.username, user.Username)
				assert.Equal(t, tt.email, user.Email)
			}
			userRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_FollowUser(t *testing.T) {
	tests := []struct {
		name            string
		username        string
		followUsername  string
		userRepoSetup   func(*mockUserRepository)
		followRepoSetup func(*mockFollowRepository)
		expectedError   error
	}{
		{
			name:           "successful follow",
			username:       "follower",
			followUsername: "following",
			userRepoSetup: func(m *mockUserRepository) {
				m.On("GetUser", mock.Anything, "follower").Return(&domain.User{}, nil)
				m.On("GetUser", mock.Anything, "following").Return(&domain.User{}, nil)
			},
			followRepoSetup: func(m *mockFollowRepository) {
				m.On("GetFollowings", mock.Anything, "follower").Return([]string{}, nil)
				m.On("FollowUser", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:           "follower does not exist",
			username:       "nonexistent",
			followUsername: "following",
			userRepoSetup: func(m *mockUserRepository) {
				m.On("GetUser", mock.Anything, "nonexistent").Return(nil, nil)
			},
			followRepoSetup: func(m *mockFollowRepository) {},
			expectedError:   ErrFollowerNotFound,
		},
		{
			name:           "following does not exist",
			username:       "follower",
			followUsername: "nonexistent",
			userRepoSetup: func(m *mockUserRepository) {
				m.On("GetUser", mock.Anything, "follower").Return(&domain.User{}, nil)
				m.On("GetUser", mock.Anything, "nonexistent").Return(nil, nil)
			},
			followRepoSetup: func(m *mockFollowRepository) {},
			expectedError:   ErrFollowingNotFound,
		},
		{
			name:           "already following",
			username:       "follower",
			followUsername: "following",
			userRepoSetup: func(m *mockUserRepository) {
				m.On("GetUser", mock.Anything, "follower").Return(&domain.User{}, nil)
				m.On("GetUser", mock.Anything, "following").Return(&domain.User{}, nil)
			},
			followRepoSetup: func(m *mockFollowRepository) {
				m.On("GetFollowings", mock.Anything, "follower").Return([]string{"following"}, nil)
			},
			expectedError: ErrAlreadyFollowing,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(mockUserRepository)
			followRepo := new(mockFollowRepository)
			tt.userRepoSetup(userRepo)
			tt.followRepoSetup(followRepo)

			service := NewUserService(userRepo, followRepo)
			follow, err := service.FollowUser(context.Background(), tt.username, tt.followUsername)

			if tt.expectedError != nil {
				assert.Nil(t, follow)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NotNil(t, follow)
				assert.NoError(t, err)
				assert.Equal(t, tt.username, follow.Username)
				assert.Equal(t, tt.followUsername, follow.Following)
			}
			userRepo.AssertExpectations(t)
			followRepo.AssertExpectations(t)
		})
	}
}
