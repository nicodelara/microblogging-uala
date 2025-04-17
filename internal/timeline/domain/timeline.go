package domain

import (
	"errors"
	"time"
)

// Tweet es una entidad del dominio que representa un tweet
type Tweet struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

// Timeline es el agregado raíz que representa el timeline de un usuario
type Timeline struct {
	Username string  `json:"username"`
	Tweets   []Tweet `json:"tweets"`
}

// NewTimeline crea un nuevo timeline para un usuario
func NewTimeline(username string) *Timeline {
	return &Timeline{
		Username: username,
		Tweets:   make([]Tweet, 0),
	}
}

// AddTweet añade un tweet al timeline siguiendo las reglas de negocio
func (t *Timeline) AddTweet(tweet Tweet) error {
	if tweet.Username == "" {
		return errors.New("tweet debe tener un usuario")
	}
	if len(tweet.Content) == 0 {
		return errors.New("tweet no puede estar vacío")
	}

	t.Tweets = append(t.Tweets, tweet)
	return nil
}
