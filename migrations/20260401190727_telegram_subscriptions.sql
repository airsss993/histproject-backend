-- +goose Up

-- Подписки пользователей на Telegram-бота: username → chatId
CREATE TABLE telegram_subscriptions (
    id                SERIAL PRIMARY KEY,                        -- Уникальный идентификатор подписки
    telegram_username VARCHAR(100) NOT NULL UNIQUE,              -- Telegram username пользователя (@username)
    chat_id           BIGINT       NOT NULL,                     -- Telegram chat ID для отправки личных сообщений
    created_at        TIMESTAMP    NOT NULL DEFAULT NOW()        -- Дата подписки (первый /start боту)
);

-- +goose Down
DROP TABLE IF EXISTS telegram_subscriptions;
