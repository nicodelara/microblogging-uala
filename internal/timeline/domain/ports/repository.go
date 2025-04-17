package ports

import (
	"context"

	"github.com/nicodelara/uala-challenge/internal/timeline/domain"
)

// TimelineRepository define la interfaz para el repositorio de timeline
type TimelineRepository interface {
	// GetTweetsForUsers obtiene los tweets de una lista de usuarios
	GetTweetsForUsers(ctx context.Context, usernames []string, offset, limit int) ([]domain.Tweet, error)
}

// UserRepository define la interfaz para el repositorio de usuarios
type UserRepository interface {
	// GetFollowedUsers obtiene la lista de usuarios seguidos
	GetFollowedUsers(ctx context.Context, username string) ([]string, error)
}

// CacheRepository define la interfaz para el caché
type CacheRepository interface {
	// Get obtiene un valor del caché
	Get(ctx context.Context, key string) (string, error)

	// Set guarda un valor en el caché
	Set(ctx context.Context, key string, value string) error
}
