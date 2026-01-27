-- +goose Up
CREATE TABLE users (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username       TEXT,
    email          TEXT,
    password_hash  TEXT,
    is_guest       BOOLEAN NOT NULL DEFAULT FALSE,
    wins           INTEGER NOT NULL DEFAULT 0,
    losses         INTEGER NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_seen_at   TIMESTAMPTZ
);

-- Уникальность username, но разрешаем NULL (для гостей).
CREATE UNIQUE INDEX users_username_uq ON users (username) WHERE username IS NOT NULL;

-- Уникальность email, но разрешаем NULL.
CREATE UNIQUE INDEX users_email_uq ON users (email) WHERE email IS NOT NULL;

-- Базовые проверки
ALTER TABLE users
    ADD CONSTRAINT users_password_hash_guest_chk
        CHECK (
            (is_guest = TRUE  AND password_hash IS NULL) OR
            (is_guest = FALSE AND password_hash IS NOT NULL)
            );

-- +goose Down
DROP TABLE IF EXISTS users;
