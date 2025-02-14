CREATE TABLE transfers (
    id SERIAL PRIMARY KEY,
    from_user SERIAL REFERENCES users(id),
    to_user SERIAL REFERENCES users(id),
    amount INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Для поиска исходящих переводов
CREATE INDEX transfers_from_user_idx ON transfers (from_user);

-- Для поиска входящих переводов
CREATE INDEX transfers_to_user_idx ON transfers (to_user);