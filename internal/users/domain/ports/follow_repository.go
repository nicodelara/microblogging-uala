package ports

import (
	"context"

	"github.com/nicodelara/uala-challenge/internal/users/domain"
)

type FollowRepository interface {
	GetFollowings(ctx context.Context, username string) ([]string, error)
	FollowUser(ctx context.Context, follow *domain.Follow) error
}
