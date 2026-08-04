-- 1. Пользователи
CREATE TYPE user_role AS ENUM ('customer', 'admin');

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'customer',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. Категории
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL
);

-- 3. Товары (с полнотекстовым поиском)
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    category_id INT NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    price NUMERIC(12, 2) NOT NULL CHECK (price >= 0),
    tsv TSVECTOR,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_products_created_at_id ON products (created_at DESC, id DESC);
CREATE INDEX idx_products_tsv ON products USING GIN(tsv);

-- Триггер для автообновления tsvector
CREATE OR REPLACE FUNCTION products_trigger_tsv() RETURNS trigger AS $$
BEGIN
  new.tsv := setweight(to_tsvector('english', coalesce(new.title, '')), 'A') ||
             setweight(to_tsvector('english', coalesce(new.description, '')), 'B');
  RETURN new;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_products_tsv BEFORE INSERT OR UPDATE
ON products FOR EACH ROW EXECUTE FUNCTION products_trigger_tsv();

-- 4. Остатки на складе
CREATE TABLE inventories (
    product_id BIGINT PRIMARY KEY REFERENCES products(id) ON DELETE CASCADE,
    stock INT NOT NULL CHECK (stock >= 0),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Тестовые данные для мгновенного поиска
INSERT INTO categories (id, name, slug) VALUES (1, 'Electronics', 'electronics');

INSERT INTO products (id, category_id, title, slug, description, price) VALUES
(1, 1, 'Pro Wireless Gaming Mouse', 'pro-wireless-gaming-mouse', 'Ultra lightweight gaming mouse with high precision sensor', 89.99),
(2, 1, 'Mechanical Keyboard RGB', 'mechanical-keyboard-rgb', 'Tactile mechanical keyboard with customizable RGB backlight', 129.99),
(3, 1, 'Noise Canceling Headphones', 'noise-canceling-headphones', 'Over-ear Bluetooth headphones with active noise cancellation', 199.99);

INSERT INTO inventories (product_id, stock) VALUES
(1, 50),
(2, 20),
(3, 0);

-- 5. Заказы
CREATE TYPE order_status AS ENUM ('pending', 'processing', 'completed', 'cancelled');

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL, -- Без FOREIGN KEY на users для простоты тестирования без созданного user
    status order_status NOT NULL DEFAULT 'pending',
    total_amount NUMERIC(12, 2) NOT NULL CHECK (total_amount >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);

-- 6. Позиции заказа
CREATE TABLE IF NOT EXISTS order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id),
    quantity INT NOT NULL CHECK (quantity > 0),
    unit_price NUMERIC(12, 2) NOT NULL CHECK (unit_price >= 0)
);

CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
