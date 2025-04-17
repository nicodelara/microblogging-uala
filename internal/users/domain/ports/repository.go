package ports

import (
	"context"

	"github.com/nicodelara/microblogging-uala/internal/users/domain"
)

// UserRepository define las operaciones que un repositorio de usuarios debe implementar
type UserRepository interface {
	// GetUser obtiene un usuario por su username
	GetUser(ctx context.Context, username string) (*domain.User, error)
	// GetUserByEmail obtiene un usuario por su email
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	// SaveUser guarda un nuevo usuario
	SaveUser(ctx context.Context, user *domain.User) error
}
