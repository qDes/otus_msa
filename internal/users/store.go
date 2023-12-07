package users

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5"
)

type Store interface {
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int64) error
	GetUser(ctx context.Context, id int64) (*User, error)
}

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	Id        int64  `json:"id"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Username  string `json:"username"`
}
type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) (*PostgresStore, error) {
	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) CreateUser(ctx context.Context, user *User) error {
	query := `INSERT INTO users (id, email, first_name, last_name, phone, username) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := s.db.QueryRowContext(ctx, query, user.Id, user.Email, user.FirstName, user.LastName, user.Phone, user.Username).Scan(&user.Id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) UpdateUser(ctx context.Context, user *User) error {
	query := `UPDATE users SET email = $1, first_name = $2, last_name = $3, phone = $4, username = $5 WHERE id = $6`
	_, err := s.db.ExecContext(ctx, query, user.Email, user.FirstName, user.LastName, user.Phone, user.Username, user.Id)
	return err
}

func (s *PostgresStore) DeleteUser(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

func (s *PostgresStore) GetUser(ctx context.Context, id int64) (*User, error) {
	user := &User{}
	query := `SELECT id, email, first_name, last_name, phone, username FROM users WHERE id = $1`
	err := s.db.QueryRowContext(ctx, query, id).Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.Phone, &user.Username)
	if err != nil {
		return nil, err
	}
	return user, nil
}
