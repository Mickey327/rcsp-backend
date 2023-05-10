package company

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

type CompanyRepository struct {
	db DB
}

func NewRepository(db DB) *CompanyRepository {
	return &CompanyRepository{
		db: db,
	}
}

func (r *CompanyRepository) Create(ctx context.Context, company *Company) (uint64, error) {
	var id uint64
	err := r.db.ExecQueryRow(ctx, `INSERT INTO companies(name) VALUES ($1) RETURNING id`, company.Name).Scan(&id)
	return id, errors.Wrapf(err, "error creating company: %v", company)
}

func (r *CompanyRepository) Read(ctx context.Context, id uint64) (*Company, error) {
	var c Company
	err := r.db.Get(ctx, &c, "SELECT id,name,created_at,updated_at FROM companies WHERE id = $1", id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, CompanyNotFoundErr
	}
	return &c, nil
}

func (r *CompanyRepository) ReadAll(ctx context.Context) ([]*Company, error) {
	companies := make([]*Company, 0)
	err := r.db.Select(ctx, &companies,
		"SELECT id, name, created_at, updated_at FROM companies")
	return companies, errors.Wrap(err, "error getting companies")
}

func (r *CompanyRepository) Update(ctx context.Context, company *Company) (bool, error) {
	company.UpdatedAt = time.Now().UTC()
	result, err := r.db.Exec(ctx,
		"UPDATE companies SET name = $1, updated_at = $2 WHERE id = $3",
		company.Name, company.UpdatedAt, company.ID)
	return result.RowsAffected() > 0, errors.Wrapf(err, "error updating company: %v", company)
}

func (r *CompanyRepository) Delete(ctx context.Context, id uint64) (bool, error) {
	result, err := r.db.Exec(ctx, "DELETE FROM companies WHERE id = $1", id)
	return result.RowsAffected() > 0, errors.Wrapf(err, "error deleting company with id: %d", id)
}
