package order

import (
	"context"
	"time"

	"github.com/pkg/errors"

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

type OrderRepository struct {
	db DB
}

func (r *OrderRepository) Create(ctx context.Context, order *Order) (uint64, error) {
	var id uint64
	err := r.db.ExecQueryRow(ctx, `INSERT INTO orders(status, user_id) VALUES ('Создан', $1) RETURNING id`, order.Status, order.UserID).Scan(&id)
	return id, errors.Wrapf(err, "error creating order: %v", order)
}

func (r *OrderRepository) ReadCurrentUserArrangingOrderLazy(ctx context.Context, userID uint64) (*Order, error) {
	var o Order
	err := r.db.Get(ctx, &o, "SELECT id, total, status, is_arranged, user_id, created_at, updated_at FROM orders WHERE user_id = $1 AND is_finished = false", userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, OrderNotFoundErr
	}
	return &o, nil
}

func (r *OrderRepository) ReadCurrentUserArrangedOrders(ctx context.Context, userID uint64) ([]*Order, error) {
	orders := make([]*Order, 0)
	err := r.db.Select(ctx, &orders,
		"SELECT id, total, status, is_arranged, user_id, created_at, updated_at FROM orders WHERE user_id = $1", userID)
	return orders, errors.Wrapf(err, "error getting orders of user with id: %v", userID)
}

func (r *OrderRepository) Update(ctx context.Context, order *Order) (bool, error) {
	order.UpdatedAt = time.Now().UTC()
	result, err := r.db.Exec(ctx,
		"UPDATE orders SET is_finished = $1, updated_at = $2 WHERE id = $3",
		order.IsArranged, order.UpdatedAt, order.ID)
	return result.RowsAffected() > 0, errors.Wrapf(err, "error updating order: %v", order)
}
