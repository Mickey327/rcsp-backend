package orderItem

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

type OrderItemRepository struct {
	db DB
}

func NewRepository(db DB) *OrderItemRepository {
	return &OrderItemRepository{
		db: db,
	}
}

func (r *OrderItemRepository) Create(ctx context.Context, orderItem *OrderItem) (bool, error) {
	result, err := r.db.Exec(ctx, `INSERT INTO order_items(quantity, order_id, product_id) VALUES ($1, $2, $3)`, orderItem.Quantity, orderItem.OrderID, orderItem.Product.ID)
	return result.RowsAffected() > 0, errors.Wrapf(err, "error creating order item: %v", orderItem)
}

func (r *OrderItemRepository) Update(ctx context.Context, orderItem *OrderItem) (bool, error) {
	orderItem.UpdatedAt = time.Now().UTC()
	result, err := r.db.Exec(ctx,
		"UPDATE order_items SET quantity = quantity + $1, updated_at = $2 WHERE order_id = $3 AND product_id = $4",
		orderItem.Quantity, orderItem.UpdatedAt, orderItem.OrderID, orderItem.Product.ID)
	return result.RowsAffected() > 0, errors.Wrapf(err, "error updating order item: %v", orderItem)
}

func (r *OrderItemRepository) Delete(ctx context.Context, orderItem *OrderItem) (bool, error) {
	result, err := r.db.Exec(ctx, "DELETE FROM order_items WHERE order_id = $1 AND product_id = $2", orderItem.OrderID, orderItem.Product.ID)
	return result.RowsAffected() > 0, errors.Wrapf(err, "error deleting order item product_id: %d, order_id: %d", orderItem.Product.ID, orderItem.OrderID)
}

func (r *OrderItemRepository) ReadByOrderAndProductID(ctx context.Context, orderID, productID uint64) (*OrderItem, error) {
	var orderItem OrderItem
	err := r.db.Get(ctx, &orderItem,
		`
			SELECT order_id, quantity, created_at, updated_at, product_id as "product.id"
			FROM order_items
			WHERE order_id = $1 AND product_id = $2
				`, orderID, productID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, OrderItemNotFound
	}
	log.Println(orderItem.OrderID, orderItem.Quantity, orderItem.UpdatedAt.String(), orderItem.CreatedAt.String())
	return &orderItem, nil
}
