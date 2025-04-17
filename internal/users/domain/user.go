package domain

import (
	"time"

	"github.com/google/uuid"
)

// User es una entidad del dominio que representa un usuario
type User struct {
	ID        string    `bson:"_id" json:"id"`
	Username  string    `bson:"username" json:"username"`
	Email     string    `bson:"email" json:"email"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

// NewUser crea una nueva instancia de User
func NewUser(username, email string) *User {
	return &User{
		ID:        uuid.New().String(),
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
	}
}
