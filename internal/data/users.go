package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type User struct {
	ID            int64     `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Password      password  `json:"-"`
	Activated     bool      `json:"activated"`
	Version       int       `json:"-"`
	Role 	      string 	`json:"role"`
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(user *User) error {
	query := `
        INSERT INTO users (name, email, password_hash, activated, role)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, version`

	args := []interface{}{user.Name, user.Email, user.Password.hash, user.Activated,user.Role}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
        SELECT id, created_at, name, email, password_hash, activated, version, role
        FROM users
        WHERE email = $1`

	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
		&user.Role,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}


func (m UserModel) GetByID(id int64) (*User, error) {
	query := `
        SELECT id, created_at, name, email, password_hash, activated, version, role
        FROM users
        WHERE id = $1`

	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
		&user.Role,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (m UserModel) Count() (int, error) {
	query := "SELECT COUNT(*) FROM users"

	var count int
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}
