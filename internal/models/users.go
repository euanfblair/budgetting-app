package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type User struct {
	UserID       int
	FirstName    string
	Surname      string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(firstName, surname, email string, password []byte) error {
	stmt := `INSERT INTO users (first_name, surname, email, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := m.DB.Exec(stmt, firstName, surname, email, password, time.Now(), time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) ExistingEmail(email string) bool {
	stmt := `SELECT * FROM users WHERE  email = $1`
	rows, err := m.DB.Query(stmt, email)
	if err != nil {
		return true
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

func (m *UserModel) Login(email string) (int, string) {
	var storedHash string
	var id int
	stmt := `SELECT user_id, password_hash  FROM users WHERE email = $1`
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &storedHash)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ""
	}
	if err != nil {
		fmt.Println("Database error:", err)
		return 0, ""
	}

	return id, storedHash
}

func (m *UserModel) GetCurrentUser(userid int) []string {
	var first_name, surname, email string
	stmt := `SELECT first_name, surname, email FROM users WHERE  user_id = $1`
	row := m.DB.QueryRow(stmt, userid)
	err := row.Scan(&first_name, &surname, &email)
	if err != nil {
		fmt.Println(err.Error())
	}
	return []string{
		first_name,
		surname,
		email,
	}
}

func (m *UserModel) GetCurrentPassword(userid int) string {
	var password string
	stmt := `SELECT password_hash FROM users WHERE  user_id = $1`
	row := m.DB.QueryRow(stmt, userid)
	err := row.Scan(&password)
	if err != nil {
		return ""
	}
	return password
}

func (m *UserModel) UpdatePassword(userid int, password []byte) error {
	stmt := `UPDATE users SET password_hash = $1, updated_at = $2 WHERE user_id = $3`
	_, err := m.DB.Exec(stmt, password, time.Now(), userid)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) DeleteUser(userid int) error {
	stmt := `DELETE FROM users WHERE user_id = $1`
	_, err := m.DB.Exec(stmt, userid)
	if err != nil {
		return err
	}
	return nil
}
