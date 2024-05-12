package models

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

func (m UserModel) GetForToken(db *gorm.DB, tokenScope, tokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	var user User

	result := db.Select("users.id, users.created_at, users.name, users.email, users.password_hash, users.activated, users.version").
		Joins("INNER JOIN tokens ON users.id = tokens.user_id").
		Where("tokens.hash = ? AND tokens.scope = ? AND tokens.expiry > ?", tokenHash[:], tokenScope, time.Now()).
		First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
