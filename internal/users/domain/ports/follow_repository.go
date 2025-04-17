package ports

import (
	"context"

	"github.com/nicodelara/microblogging-uala/internal/users/domain"
)

type FollowRepository interface {
	GetFollowings(ctx context.Context, username string) ([]string, error)
	FollowUser(ctx context.Context, follow *domain.Follow) error
}
