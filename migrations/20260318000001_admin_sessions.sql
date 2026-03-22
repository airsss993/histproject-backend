-- +goose Up
CREATE TABLE admin_sessions
(
    id            SERIAL PRIMARY KEY,
    admin_id      INT         NOT NULL REFERENCES admins(id) ON DELETE CASCADE,
    refresh_token TEXT        NOT NULL UNIQUE,
    expires_at    TIMESTAMPTZ NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS admin_sessions;
