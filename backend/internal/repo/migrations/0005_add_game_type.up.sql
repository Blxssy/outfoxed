-- +goose Up
ALTER TABLE games
    ADD COLUMN visibility TEXT NOT NULL DEFAULT 'public'
        CHECK (visibility IN ('public', 'private'));

ALTER TABLE games
    ADD COLUMN join_code CHAR(6);

ALTER TABLE games
    ADD COLUMN title TEXT NOT NULL DEFAULT 'Комната';

CREATE UNIQUE INDEX games_join_code_uq
    ON games (join_code)
    WHERE join_code IS NOT NULL;

CREATE INDEX games_visibility_status_idx
    ON games (visibility, status);

CREATE INDEX games_join_code_idx
    ON games (join_code);