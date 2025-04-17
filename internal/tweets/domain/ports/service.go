package ports

import (
	"context"

	"github.com/nicodelara/uala-challenge/internal/tweets/domain"
)

// TweetService define la interfaz para el servicio de tweets
type TweetService interface {
	// CreateTweet crea un nuevo tweet
	CreateTweet(ctx context.Context, username, content string) (*domain.Tweet, error)
}
