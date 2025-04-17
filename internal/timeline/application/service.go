package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nicodelara/microblogging-uala/internal/common"
	"github.com/nicodelara/microblogging-uala/internal/timeline/domain"
	"github.com/nicodelara/microblogging-uala/internal/timeline/domain/ports"
)

type TweetView struct {
	TweetID   string `json:"tweetId"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

type timelineService struct {
	repo        ports.TimelineRepository
	userChecker common.UserChecker
	cache       ports.CacheRepository
}

func NewTimelineService(
	repo ports.TimelineRepository,
	userChecker common.UserChecker,
	cache ports.CacheRepository,
) ports.TimelineService {
	return &timelineService{
		repo:        repo,
		userChecker: userChecker,
		cache:       cache,
	}
}

func (s *timelineService) GetTimeline(ctx context.Context, username string, offset, limit int) (*domain.Timeline, error) {
	// Verificar que el usuario existe
	_, err := s.userChecker.GetUser(username)
	if err != nil {
		return nil, errors.New("user does not exist")
	}

	// Obtener la lista de followings
	followings, err := s.userChecker.GetFollowings(ctx, username)
	if err != nil {
		return nil, err
	}
	if len(followings) == 0 {
		return domain.NewTimeline(username), nil
	}

	// Generar clave de cache con offset y limit
	cacheKey := fmt.Sprintf("timeline:%s:offset=%d:limit=%d", username, offset, limit)
	cached, err := s.cache.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var cachedTweets []domain.Tweet
		if err := json.Unmarshal([]byte(cached), &cachedTweets); err == nil {
			timeline := domain.NewTimeline(username)
			for _, tweet := range cachedTweets {
				if err := timeline.AddTweet(tweet); err != nil {
					return nil, err
				}
			}
			return timeline, nil
		}
	}

	// Obtener tweets de los usuarios seguidos
	tweets, err := s.repo.GetTweetsForUsers(ctx, followings, offset, limit)
	if err != nil {
		return nil, err
	}

	if len(tweets) == 0 {
		return domain.NewTimeline(username), nil
	}

	// Crear timeline
	timeline := domain.NewTimeline(username)
	for _, tweet := range tweets {
		if err := timeline.AddTweet(tweet); err != nil {
			return nil, err
		}
	}

	// Guardar en cache
	tweetsJSON, err := json.Marshal(tweets)
	if err == nil {
		_ = s.cache.Set(ctx, cacheKey, string(tweetsJSON))
	}

	return timeline, nil
}
