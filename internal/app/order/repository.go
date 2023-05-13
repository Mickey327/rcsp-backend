package order

import (
	"context"
	"time"

	"github.com/Mickey327/rcsp-backend/internal/app/orderItem"
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

func NewRepository(db DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Create(ctx context.Context, userID uint64) (uint64, error) {
	var id uint64
	err := r.db.ExecQueryRow(ctx, `INSERT INTO orders(status, user_id) VALUES ('Создан', $1) RETURNING id`, userID).Scan(&id)
	return id, errors.Wrapf(err, "error creating order for user with id: %v", userID)
}

func (r *OrderRepository) Update(ctx context.Context, order *Order) (bool, error) {
	order.UpdatedAt = time.Now().UTC()
	result, err := r.db.Exec(ctx,
		"UPDATE orders SET is_arranged = $1, updated_at = $2 WHERE id = $3",
		order.IsArranged, order.UpdatedAt, order.ID)
	return result.RowsAffected() > 0, errors.Wrapf(err, "error updating order: %v", order)
}

func (r *OrderRepository) ReadCurrentUserArrangingOrderLazy(ctx context.Context, userID uint64) (*Order, error) {
	var o Order

	err := r.db.Get(ctx, &o, "SELECT id, total, status, is_arranged, user_id, created_at, updated_at FROM orders WHERE user_id = $1 AND is_arranged = false", userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, OrderNotFoundErr
	}

	var count uint64
	err = r.db.Get(ctx, &count, `
			SELECT COUNT(*) 
			FROM order_items 
			WHERE order_id = (SELECT orders.id FROM orders WHERE orders.user_id = $1 AND is_arranged = false)`, userID)
	if err != nil {
		return nil, err
	}
	o.Count = count

	return &o, nil
}

func (r *OrderRepository) ReadCurrentUserArrangingOrderEager(ctx context.Context, userID uint64) (*Order, error) {
	var o Order
	err := r.db.Get(ctx, &o, `
		SELECT orders.id, orders.total, orders.status, orders.is_arranged, orders.user_id, orders.created_at, orders.updated_at
		FROM orders
		WHERE user_id = $1 AND is_arranged = false
		`, userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, OrderNotFoundErr
	}

	orderItems := make([]*orderItem.OrderItem, 0)
	err = r.db.Select(ctx, &orderItems, ` 
		SELECT 
		    order_items.quantity, order_items.order_id, order_items.created_at, order_items.updated_at,
		    p.id as "product.id", p.name as "product.name", p.description as "product.description", p.price as "product.price", p.stock as "product.stock",
        	p.image as "product.image", p.created_at as "product.created_at", p.updated_at as "product.updated_at",
        	c.id as "product.category.id", c.name as "product.category.name", c.updated_at as "product.category.updated_at", c.created_at as "product.category.created_at",
       		c2.id as "product.company.id", c2.name as "product.company.name", c2.updated_at as "product.company.updated_at", c2.created_at as "product.company.created_at"
		FROM order_items
			JOIN products p on p.id = order_items.product_id
			JOIN categories c on p.category_id = c.id
			JOIN companies c2 on p.company_id = c2.id
		WHERE order_items.order_id = $1
			`, o.ID)

	o.OrderItems = orderItems
	o.Count = uint64(len(orderItems))

	if err != nil {
		return nil, err
	}

	return &o, nil
}

func (r *OrderRepository) GetUserIDByNotArrangedOrderID(ctx context.Context, orderID uint64) (uint64, error) {
	var userID uint64

	err := r.db.Get(ctx, &userID, `
		SELECT orders.user_id
		FROM orders
		WHERE orders.id = $1 AND is_arranged = false
		`, orderID)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, OrderNotFoundErr
	}

	return userID, nil
}
