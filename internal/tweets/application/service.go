package application

import (
	"context"

	"github.com/nicodelara/uala-challenge/internal/common"
	"github.com/nicodelara/uala-challenge/internal/tweets/domain"
	"github.com/nicodelara/uala-challenge/internal/tweets/domain/ports"
)

type tweetService struct {
	repo    ports.TweetRepository
	checker common.UserChecker
}

func NewTweetService(repo ports.TweetRepository, checker common.UserChecker) ports.TweetService {
	return &tweetService{
		repo:    repo,
		checker: checker,
	}
}

func (s *tweetService) CreateTweet(ctx context.Context, username, content string) (*domain.Tweet, error) {
	// Verificar que el usuario exista
	_, err := s.checker.GetUser(username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	tweet := domain.NewTweet(username, content)

	err = s.repo.SaveTweet(ctx, tweet)
	if err != nil {
		return nil, err
	}

	return tweet, nil
}
