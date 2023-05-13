package commentary

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB interface {
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	GetPool() *pgxpool.Pool
}

type CommentaryRepository struct {
	db DB
}

func (r *CommentaryRepository) Create() {

}

func (r *CommentaryRepository) ReadByProductID() {

}

func (r *CommentaryRepository) Update() {

}

func (r *CommentaryRepository) Delete() {

}