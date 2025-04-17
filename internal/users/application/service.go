package application

import (
	"context"
	"slices"

	"github.com/nicodelara/uala-challenge/internal/users/domain"
	"github.com/nicodelara/uala-challenge/internal/users/domain/ports"
)

type userService struct {
	userRepo   ports.UserRepository
	followRepo ports.FollowRepository
}

func NewUserService(userRepo ports.UserRepository, followRepo ports.FollowRepository) ports.UserService {
	return &userService{
		userRepo:   userRepo,
		followRepo: followRepo,
	}
}

func (s *userService) GetUser(ctx context.Context, username string) (*domain.User, error) {
	return s.userRepo.GetUser(ctx, username)
}

func (s *userService) CreateUser(ctx context.Context, username, email string) (*domain.User, error) {
	// Verificar que el username no exista
	existingUser, err := s.GetUser(ctx, username)
	if err == nil && existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Verificar que el email no exista
	existingUser, err = s.userRepo.GetUserByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	user := domain.NewUser(username, email)
	err = s.userRepo.SaveUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) FollowUser(ctx context.Context, username, followUsername string) (*domain.Follow, error) {
	// Verificar que el usuario que quiere seguir exista
	follower, err := s.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}
	if follower == nil {
		return nil, ErrFollowerNotFound
	}

	// Verificar que el usuario a seguir exista
	following, err := s.GetUser(ctx, followUsername)
	if err != nil {
		return nil, err
	}
	if following == nil {
		return nil, ErrFollowingNotFound
	}

	// Verificar que no esté siguiendo ya al usuario
	followings, err := s.followRepo.GetFollowings(ctx, username)
	if err != nil {
		return nil, err
	}
	if slices.Contains(followings, followUsername) {
		return nil, ErrAlreadyFollowing
	}

	follow := domain.NewFollow(username, followUsername)
	err = s.followRepo.FollowUser(ctx, follow)
	if err != nil {
		return nil, err
	}

	return follow, nil
}
