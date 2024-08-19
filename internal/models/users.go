package models

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	UserID       int
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(username, email string, password []byte) error {
	stmt := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`

	_, err := m.DB.Exec(stmt, username, email, password)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) Validate(username, email string) bool {
	stmt := `SELECT * FROM users WHERE username = $1 OR email = $2`

	rows, err := m.DB.Query(stmt, username, email)
	if err != nil {
		return false
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}
	return false

}

func (m *UserModel) Login(username, inputPassword string) (int, error) {
	var storedHash string
	var id int
	stmt := `SELECT user_id, password_hash  FROM users WHERE username = $1 OR email = $1`
	row := m.DB.QueryRow(stmt, username)
	err := row.Scan(&id, &storedHash)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	if err != nil {
		fmt.Println("Database error:", err)
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(inputPassword))
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return id, nil
}

func (m *UserModel) GetUserProfile(id int) []string {
	userData := make([]string, 2)

	stmt := `SELECT username, email  FROM users WHERE user_id = $1`
	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(&userData[0], &userData[1])
	if err != nil {
		fmt.Println("Database error:", err)
		return nil
	}
	return userData
}

func (m *UserModel) UpdatePassword(id int, newPassword []byte) error {
	stmt := `UPDATE users SET password_hash = $1, updated_at = $2 where user_id = $3`
	_, err := m.DB.Exec(stmt, newPassword, time.Now(), id)
	return err
}

func (m *UserModel) CheckPasswordMatches(id int, oldPasswordConfirm string) error {
	var oldPassword string
	stmt := `SELECT password_hash  FROM users WHERE user_id = $1`
	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(&oldPassword)
	if err != nil {
		fmt.Println("Database error:", err)
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(oldPassword), []byte(oldPasswordConfirm))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return err
	} else if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) DeleteUser(id int) error {
	stmt := `DELETE FROM users WHERE user_id = $1`
	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}
