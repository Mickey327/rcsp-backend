package orderItem

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

type OrderItemRepository struct {
	db DB
}

func NewRepository(db DB) *OrderItemRepository {
	return &OrderItemRepository{
		db: db,
	}
}

func (r *OrderItemRepository) ReadByOrderID(ctx context.Context, orderID uint64) ([]*OrderItem, error) {
	orderItems := make([]*OrderItem, 0)
	err := r.db.Select(ctx, &orderItems,
		"SELECT order_id, product_id, quantity, created_at, updated_at FROM order_items WHERE order_id = $1", orderID)
	return orderItems, errors.Wrapf(err, "error getting order items of order with id: %v", orderID)
}

func (r *OrderItemRepository) ChangeQuantityProductToOrderByID(ctx context.Context, orderItem *OrderItem, quantityDelta int) (bool, error) {
	result, err := r.db.Exec(ctx, "CALL order_items_procedure($1,$2,$3);", orderItem.OrderID, orderItem.ProductID, quantityDelta)
	return result.RowsAffected() > 0, errors.Wrapf(err, "error changing product quantity: %d", orderItem.Quantity)
}
