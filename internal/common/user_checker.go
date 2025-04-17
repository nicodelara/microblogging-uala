package common

import (
	"context"

	"github.com/nicodelara/uala-challenge/internal/users/domain"
)

// UserChecker define la interfaz para verificar la existencia de un usuario
// y obtener, además, la lista de usuarios a los que sigue.
type UserChecker interface {
	GetUser(username string) (*domain.User, error)
	GetFollowings(ctx context.Context, username string) ([]string, error)
}
