-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    username TEXT NOT NULL,
    email TEXT,
    password_hash TEXT,

    is_guest BOOLEAN NOT NULL DEFAULT false,
    role TEXT NOT NULL DEFAULT 'player',

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_seen_at TIMESTAMP
);

-- Индексы
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Уникальность
CREATE UNIQUE INDEX users_username_uq ON users (username);
CREATE UNIQUE INDEX users_email_uq ON users (email) WHERE email IS NOT NULL;

-- Проверка guest / registered
ALTER TABLE users
    ADD CONSTRAINT users_password_hash_guest_chk
        CHECK (
            (is_guest = TRUE  AND password_hash IS NULL) OR
            (is_guest = FALSE AND password_hash IS NOT NULL)
        );

-- Автообновление updated_at
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();