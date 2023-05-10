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
    name TEXT NOT NULL,
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
    is_finished BOOLEAN NOT NULL DEFAULT FALSE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON UPDATE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS order_items(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    quantity BIGINT NOT NULL,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON UPDATE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id) ON UPDATE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS reviews(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    message TEXT NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON UPDATE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id) ON UPDATE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE order_items;
DROP TABLE orders;
DROP TABLE reviews;
DROP TABLE products;
DROP TABLE users;
DROP TABLE categories;
DROP TABLE companies;
-- +goose StatementEnd
