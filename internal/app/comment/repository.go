package comment

import (
	"context"
	"log"
	"time"

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

type CommentRepository struct {
	db DB
}

func NewRepository(db DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (r *CommentRepository) Create(ctx context.Context, comment *Comment) (uint64, error) {
	var id uint64
	err := r.db.ExecQueryRow(ctx, `INSERT INTO comments(message, product_id, user_id) VALUES ($1, $2, $3) RETURNING product_id`,
		comment.Message, comment.ProductID, comment.User.ID).Scan(&id)
	return id, errors.Wrapf(err, "error creating comment: %v", comment)
}

func (r *CommentRepository) ReadByProductID(ctx context.Context, productID uint64) ([]*Comment, error) {
	comments := make([]*Comment, 0)
	err := r.db.Select(ctx, &comments,
		`SELECT comments.message, comments.user_id AS "user.id", users.email as "user.email", comments.updated_at 
				FROM comments 
				JOIN users ON users.id = comments.user_id
				WHERE product_id = $1`, productID)
	return comments, errors.Wrapf(err, "error getting comments for product with id: %d", productID)
}

func (r *CommentRepository) ReadByUserAndProductID(ctx context.Context, userID, productID uint64) (*Comment, error) {
	var c Comment
	err := r.db.Get(ctx, &c, `SELECT comments.message, comments.user_id AS "user.id", users.email as "user.email", comments.updated_at 
				FROM comments 
				JOIN users ON users.id = comments.user_id
				WHERE product_id = $1 AND user_id = $2`, productID, userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, CommentNotFoundErr
	}
	return &c, nil
}

func (r *CommentRepository) Update(ctx context.Context, comment *Comment) (bool, error) {
	comment.UpdatedAt = time.Now().UTC()
	log.Println("COMMENT MESSAGE REPOSITORY LEVEL:", comment.Message)
	result, err := r.db.Exec(ctx,
		"UPDATE comments SET message = $1, updated_at = $2 WHERE user_id = $3 AND product_id = $4",
		comment.Message, comment.UpdatedAt, comment.User.ID, comment.ProductID)
	log.Println("IS UPDATED:", result.RowsAffected())
	return result.RowsAffected() > 0, errors.Wrapf(err, "error updating comment: %v", comment)
}
