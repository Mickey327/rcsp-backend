package user

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type DB interface {
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	GetPool() *pgxpool.Pool
}

type UserRepository struct {
	db DB
}

func NewRepository(db DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) Register(ctx context.Context, user *User) (uint64, error) {
	var id uint64
	err := u.db.ExecQueryRow(ctx, `INSERT INTO users(email, password, role_name) VALUES ($1, $2, $3) RETURNING id`, user.Email, user.Password, user.Role).Scan(&id)
	return id, errors.Wrapf(err, "error registering user: %v", user)
}

func (u *UserRepository) GetByEmail(ctx context.Context, user *User) (*User, error) {
	var dbUser User
	err := u.db.Get(ctx, &dbUser, "SELECT id, email, password, role_name FROM users WHERE email = $1", user.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.Wrapf(err, "user with such email not found: %v", user.Email)
	}
	return &dbUser, nil
}
