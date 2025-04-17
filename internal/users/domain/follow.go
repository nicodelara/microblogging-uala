package domain

import (
	"time"

	"github.com/google/uuid"
)

// Follow es una entidad del dominio que representa una relación de seguimiento entre usuarios
type Follow struct {
	ID        string    `bson:"_id" json:"id"`
	Username  string    `bson:"username" json:"username"`
	Following string    `bson:"following" json:"following"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

// NewFollow crea una nueva instancia de Follow
func NewFollow(username, following string) *Follow {
	return &Follow{
		ID:        uuid.New().String(),
		Username:  username,
		Following: following,
		CreatedAt: time.Now(),
	}
}
