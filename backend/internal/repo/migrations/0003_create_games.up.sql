CREATE TABLE games (
                       id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       status        TEXT NOT NULL CHECK (status IN ('waiting', 'active', 'finished')),
                       state_json    JSONB NOT NULL,
                       version       INTEGER NOT NULL DEFAULT 1,
                       fox_escape_at INTEGER NOT NULL DEFAULT 15,
                       culprit_id    INTEGER NOT NULL DEFAULT 0,
                       created_by    UUID REFERENCES users (id),
                       created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
                       updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX games_status_idx ON games (status);
CREATE INDEX games_created_by_idx ON games (created_by);

CREATE TABLE game_players (
                              game_id   UUID NOT NULL REFERENCES games (id) ON DELETE CASCADE,
                              user_id   UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
                              seat      SMALLINT NOT NULL CHECK (seat >= 0 AND seat < 4),
                              joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                              PRIMARY KEY (game_id, user_id),
                              CONSTRAINT game_players_game_seat_uq UNIQUE (game_id, seat)
);

CREATE INDEX game_players_user_id_idx ON game_players (user_id);
CREATE INDEX game_players_game_id_idx ON game_players (game_id);