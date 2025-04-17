package common

import (
	"context"
	"errors"
	"time"

	"github.com/nicodelara/microblogging-uala/internal/users/domain"
	userPorts "github.com/nicodelara/microblogging-uala/internal/users/domain/ports"
)

// User representa un usuario en el sistema
type User struct {
	Username string
}

// UserCheckerAdapter adapta un UserRepository y FollowRepository para implementar la interfaz UserChecker
type UserCheckerAdapter struct {
	userRepo   userPorts.UserRepository
	followRepo userPorts.FollowRepository
}

// NewUserCheckerAdapter crea una nueva instancia de UserCheckerAdapter
func NewUserCheckerAdapter(userRepo userPorts.UserRepository, followRepo userPorts.FollowRepository) *UserCheckerAdapter {
	return &UserCheckerAdapter{
		userRepo:   userRepo,
		followRepo: followRepo,
	}
}

// GetUser obtiene un usuario por su nombre de usuario
func (a *UserCheckerAdapter) GetUser(username string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := a.userRepo.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user does not exist")
	}

	return user, nil
}

// GetFollowings obtiene los usuarios seguidos por un usuario
func (a *UserCheckerAdapter) GetFollowings(ctx context.Context, username string) ([]string, error) {
	return a.followRepo.GetFollowings(ctx, username)
}
