package ports

import (
	"context"

	"github.com/nicodelara/microblogging-uala/internal/users/domain"
)

// UserService define las operaciones que un servicio de usuarios debe implementar
type UserService interface {
	// GetUser obtiene un usuario por su username
	GetUser(ctx context.Context, username string) (*domain.User, error)
	// CreateUser crea un nuevo usuario
	CreateUser(ctx context.Context, username, email string) (*domain.User, error)
	// FollowUser crea una relación de seguimiento entre usuarios
	FollowUser(ctx context.Context, username, followUsername string) (*domain.Follow, error)
}
