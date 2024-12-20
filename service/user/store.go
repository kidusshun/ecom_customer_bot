package user

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetUserByEmail(email string) (*User, error) {
	rows := s.db.QueryRow("SELECT * FROM users WHERE email = ?", email)

	u := new(User)
	u, err := ScanRowToUser(rows)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil

}

func ScanRowToUser(rows *sql.Row) (*User, error) {

	user := new(User)
	err := rows.Scan()

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) GetUserByID(id uuid.UUID) (*User, error) {
	rows := s.db.QueryRow("SELECT * FROM users WHERE id = ?", id)

	u := new(User)
	u, err := ScanRowToUser(rows)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return u, nil

}

func (s *Store) CreateUser(user User) error {
	_, err := s.db.Query("INSERT INTO users (name, email, profile_picture) VALUES (?,?,?)", user.Name, user.Email, user.ProfilePicture)
	if err != nil {
		return err
	}
	return nil
}
