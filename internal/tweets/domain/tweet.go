package domain

import (
	"time"

	"github.com/google/uuid"
)

// Tweet es una entidad del dominio que representa un tweet
type Tweet struct {
	ID        string    `bson:"_id" json:"id"`
	Username  string    `bson:"username" json:"username"`
	Content   string    `bson:"content" json:"content"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

// NewTweet crea una nueva instancia de Tweet
func NewTweet(username, content string) *Tweet {
	return &Tweet{
		ID:        uuid.New().String(),
		Username:  username,
		Content:   content,
		CreatedAt: time.Now(),
	}
}
