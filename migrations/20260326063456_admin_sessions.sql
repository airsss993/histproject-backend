-- +goose Up

-- Таблица сессий администраторов для хранения refresh-токенов
CREATE TABLE admin_sessions
(
    id            SERIAL PRIMARY KEY,                                    -- Уникальный идентификатор сессии
    admin_id      INT         NOT NULL REFERENCES admins(id) ON DELETE CASCADE, -- ID администратора
    refresh_token TEXT        NOT NULL UNIQUE,                           -- Refresh-токен
    expires_at    TIMESTAMPTZ NOT NULL,                                  -- Дата истечения токена
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()                     -- Дата создания сессии
);

-- +goose Down
DROP TABLE IF EXISTS admin_sessions;
