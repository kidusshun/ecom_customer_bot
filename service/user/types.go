package user

import (
	"github.com/google/uuid"
)

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id uuid.UUID) (*User, error)
	CreateUser(User) error
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Name 	string    `json:"name"`
	Email	 string    `json:"email"`
	ProfilePicture string `json:"profile_picture"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}
