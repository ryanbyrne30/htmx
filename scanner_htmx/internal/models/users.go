package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(name, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := "INSERT INTO users (name, email, hashed_password, created) VALUES(?, ?, ?, DATETIME('now', 'utc'))"
	_, err = u.DB.Exec(stmt, name, email, hash)
	
	if err != nil {
		var sqliteError sqlite3.Error

		if errors.As(err, &sqliteError) {
			if errors.Is(sqliteError.Code, sqlite3.ErrConstraint) && strings.Contains(sqliteError.Error(), "users.email") {
				return ErrDuplicateEmail
			}
		}

		return err	
	}

	return nil
}

func (u *UserModel) Authenticate(email, password string) (int, error) {
	stmt := "SELECT id, hashed_password FROM users WHERE email=?"
	user := &User{}
	err := u.DB.QueryRow(stmt, email).Scan(&user.ID, &user.HashedPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, ErrNoRecord
		} else {
			return -1, err
		}
	}

	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
	if err != nil {
		return -1, ErrInvalidCredentials
	}

	return user.ID, nil
}

func (u *UserModel) Exists(id int) (bool, error) {
	return false, nil
}