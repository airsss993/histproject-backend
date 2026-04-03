-- +goose Up
ALTER TABLE telegram_subscriptions RENAME TO bot_subscriptions;
ALTER TABLE bot_subscriptions RENAME COLUMN telegram_username TO username;

-- +goose Down
ALTER TABLE bot_subscriptions RENAME COLUMN username TO telegram_username;
ALTER TABLE bot_subscriptions RENAME TO telegram_subscriptions;
