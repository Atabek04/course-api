package models

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	//Password  password  `gorm:"foreignKey:UserID" json:"-"`
	PasswordHash string `gorm:"column:password_hash"`
	Activated    bool   `json:"activated"`
	Version      int    `json:"-"`
}

type password struct {
	UserID    int64
	Plaintext *string
	Hash      []byte
}

func (u *User) SetPassword(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

type UserModel struct {
	DB *sql.DB
}

func (User) TableName() string {
	return "users"
}
