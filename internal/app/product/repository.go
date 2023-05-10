package product

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

type ProductRepository struct {
	db DB
}

func NewRepository(db DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, product *Product) (uint64, error) {
	var id uint64
	err := r.db.ExecQueryRow(ctx, `INSERT INTO products(name, description, price, stock, image, category_id, company_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		product.Name, product.Description, product.Price, product.Stock, product.Image, product.CategoryID, product.CompanyID).Scan(&id)
	return id, errors.Wrapf(err, "error creating product: %v", product)
}

func (r *ProductRepository) Read(ctx context.Context, id uint64) (*Product, error) {
	var p Product
	err := r.db.Get(ctx, &p, "SELECT id, name, description, price, stock, image, category_id, company_id, created_at, updated_at FROM products WHERE id = $1", id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ProductNotFoundErr
	}
	return &p, nil
}

func (r *ProductRepository) ReadAll(ctx context.Context) ([]*Product, error) {
	products := make([]*Product, 0)
	err := r.db.Select(ctx, &products,
		"SELECT id, name, description, price, stock, image, category_id, company_id, created_at, updated_at FROM products")
	return products, errors.Wrap(err, "error getting products")
}

func (r *ProductRepository) ReadByCategoryID(ctx context.Context, categoryID uint64) ([]*Product, error) {
	products := make([]*Product, 0)
	err := r.db.Select(ctx, &products,
		"SELECT id, name, description, price, stock, image, category_id, company_id, created_at, updated_at FROM products WHERE category_id = $1",
		categoryID)
	return products, errors.Wrapf(err, "error getting products by category id: %d", categoryID)
}

func (r *ProductRepository) ReadByCompanyID(ctx context.Context, companyID uint64) ([]*Product, error) {
	products := make([]*Product, 0)
	err := r.db.Select(ctx, &products,
		"SELECT id, name, description, price, stock, image, category_id, company_id, created_at, updated_at FROM products WHERE company_id = $1",
		companyID)
	return products, errors.Wrapf(err, "error getting products by category id: %d", companyID)
}

func (r *ProductRepository) ReadByCompanyIDAndCategoryID(ctx context.Context, companyID, categoryID uint64) ([]*Product, error) {
	products := make([]*Product, 0)
	err := r.db.Select(ctx, &products,
		"SELECT id, name, description, price, stock, image, category_id, company_id, created_at, updated_at FROM products WHERE company_id = $1 AND category_id = $2",
		companyID, categoryID)
	return products, errors.Wrapf(err, "error getting products by company id and category id: %d; %d", companyID, categoryID)
}

func (r *ProductRepository) Update(ctx context.Context, product *Product) (bool, error) {
	product.UpdatedAt = time.Now().UTC()
	result, err := r.db.Exec(ctx,
		"UPDATE products SET name = $1, description = $2, price = $3, stock = $4, image = $5, category_id = $6, company_id = $7, updated_at = $8 WHERE id = $9",
		product.Name, product.Description, product.Price, product.Stock, product.Image, product.CategoryID, product.CompanyID, product.UpdatedAt, product.ID)
	return result.RowsAffected() > 0, errors.Wrapf(err, "error updating product: %v", product)
}

func (r *ProductRepository) Delete(ctx context.Context, id uint64) (bool, error) {
	result, err := r.db.Exec(ctx, "DELETE FROM products WHERE id = $1", id)
	return result.RowsAffected() > 0, errors.Wrapf(err, "error deleting product with id: %d", id)
}
