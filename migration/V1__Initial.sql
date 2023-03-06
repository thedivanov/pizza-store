CREATE TABLE IF NOT EXISTS items(
    id serial PRIMARY KEY,
    name VARCHAR(50),
    amount numeric,
    currency VARCHAR(5),
    comment VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS orders(
    id serial PRIMARY KEY,
    status VARCHAR(30) DEFAULT 'new',
    delivery_address VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS orders_items(
    order_id int REFERENCES orders(id),
	item_id int REFERENCES items(id)
);

CREATE TABLE IF NOT EXISTS kitchen_orders(
    id serial PRIMARY KEY,
    status VARCHAR(30) DEFAULT 'new',
    order_id int REFERENCES orders(id)
);

CREATE TABLE IF NOT EXISTS delivery_orders(
    id serial PRIMARY KEY,
    status VARCHAR(30) DEFAULT 'new',
    order_id int REFERENCES orders(id)
);