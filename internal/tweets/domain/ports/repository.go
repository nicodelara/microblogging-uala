package ports

import (
	"context"

	"github.com/nicodelara/microblogging-uala/internal/tweets/domain"
)

// TweetRepository define la interfaz para el repositorio de tweets
type TweetRepository interface {
	// SaveTweet guarda un tweet en el repositorio
	SaveTweet(ctx context.Context, tweet *domain.Tweet) error
}
