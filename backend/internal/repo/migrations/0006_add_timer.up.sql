ALTER TABLE games
    ADD COLUMN turn_deadline_at TIMESTAMPTZ;

CREATE INDEX games_turn_deadline_idx
    ON games (turn_deadline_at)
    WHERE status = 'active';