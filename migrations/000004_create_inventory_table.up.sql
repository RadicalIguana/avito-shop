CREATE TABLE inventory (
    user_id SERIAL REFERENCES users(id),
    item_name VARCHAR(255) REFERENCES merch(name),
    quantity INT NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, item_name)
);

-- Для быстрого доступа к инвенторю пользователя
CREATE INDEX inventory_user_id_idx ON inventory (user_id);

