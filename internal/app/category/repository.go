package category

import (
	"context"
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

type CategoryRepository struct {
	db DB
}

func NewRepository(db DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, category *Category) (uint64, error) {
	var id uint64
	err := r.db.ExecQueryRow(ctx, `INSERT INTO categories(name) VALUES ($1) RETURNING id`, category.Name).Scan(&id)
	return id, errors.Wrapf(err, "error creating category: %v", category)
}

func (r *CategoryRepository) ReadAll(ctx context.Context) ([]*Category, error) {
	categories := make([]*Category, 0)
	err := r.db.Select(ctx, &categories,
		"SELECT id, name, created_at, updated_at FROM categories")
	return categories, errors.Wrap(err, "error getting categories")
}

func (r *CategoryRepository) Read(ctx context.Context, id uint64) (*Category, error) {
	var c Category
	err := r.db.Get(ctx, &c, "SELECT id,name,created_at,updated_at FROM categories WHERE id = $1", id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, CategoryNotFoundErr
	}
	return &c, nil
}

func (r *CategoryRepository) Update(ctx context.Context, category *Category) (bool, error) {
	category.UpdatedAt = time.Now().UTC()
	result, err := r.db.Exec(ctx,
		"UPDATE categories SET name = $1, updated_at = $2 WHERE id = $3",
		category.Name, category.UpdatedAt, category.ID)
	return result.RowsAffected() > 0, errors.Wrapf(err, "error updating category: %v", category)
}

func (r *CategoryRepository) Delete(ctx context.Context, id uint64) (bool, error) {
	result, err := r.db.Exec(ctx, "DELETE FROM categories WHERE id = $1", id)
	return result.RowsAffected() > 0, errors.Wrapf(err, "error deleting category with id: %d", id)
}
