-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role_name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS categories(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS companies(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS products(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    price BIGINT NOT NULL,
    stock BIGINT NOT NULL,
    image TEXT NOT NULL,
    category_id BIGINT NOT NULL REFERENCES categories(id) ON UPDATE CASCADE,
    company_id BIGINT NOT NULL REFERENCES companies(id) ON UPDATE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS orders(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    total BIGINT NOT NULL DEFAULT 0,
    status TEXT NOT NULL,
    is_arranged BOOLEAN NOT NULL DEFAULT FALSE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON UPDATE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS order_items(
    quantity BIGINT NOT NULL,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON UPDATE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id) ON UPDATE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY(order_id, product_id)
);
CREATE TABLE IF NOT EXISTS reviews(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    message TEXT NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON UPDATE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id) ON UPDATE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION update_total_price() RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        UPDATE orders
        SET total = COALESCE((SELECT sum(quantity*price) as total
                            from (SELECT * FROM orders) as something
                                     INNER JOIN order_items oi on something.id = oi.order_id
                                     INNER JOIN products p on oi.product_id = p.id
                                     INNER JOIN users u on orders.user_id = u.id
                            WHERE something.id = old.order_id
                            GROUP BY user_id),0)
        WHERE old.order_id = orders.id;
        RETURN old;
    ELSIF (TG_OP = 'UPDATE') OR (TG_OP = 'INSERT') THEN
        UPDATE orders
        SET total = COALESCE((SELECT sum(quantity*price) as total
                            from (SELECT * FROM orders) as something
                                     INNER JOIN order_items oi on something.id = oi.order_id
                                     INNER JOIN products p on oi.product_id = p.id
                                     INNER JOIN users u on orders.user_id = u.id
                            WHERE something.id = new.order_id
                            GROUP BY user_id),0)
        WHERE new.order_id = orders.id;
        RETURN new;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE PROCEDURE order_items_procedure(orders_id BIGINT, products_id BIGINT, quantity_delta BIGINT)
    LANGUAGE SQL
    AS $$
MERGE INTO order_items w
USING (VALUES (orders_id,
               products_id,
               quantity_delta)) s
    ON (s.column1 = w.order_id AND s.column2 = w.product_id)
    WHEN NOT MATCHED AND s.column3 > 0 THEN
        INSERT VALUES(s.column1, s.column2, s.column3)
    WHEN MATCHED AND w.quantity + s.column3 > 0 THEN
        UPDATE SET quantity = w.quantity + s.column3
    WHEN MATCHED THEN
        DELETE;
$$;

CREATE TRIGGER update_order_total_price
    AFTER INSERT OR UPDATE OR DELETE ON order_items
        FOR EACH ROW EXECUTE FUNCTION update_total_price();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_order_total_price ON order_items;
DROP PROCEDURE IF EXISTS order_items_procedure(order_id BIGINT, product_id BIGINT, quantity_delta BIGINT);
DROP TABLE order_items;
DROP TABLE orders;
DROP TABLE reviews;
DROP TABLE products;
DROP TABLE users;
DROP TABLE categories;
DROP TABLE companies;
-- +goose StatementEnd
