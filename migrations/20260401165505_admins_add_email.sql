-- +goose Up

-- Добавляем email к таблице администраторов для возможности отправки уведомлений и связи с ними.
ALTER TABLE admins ADD COLUMN email VARCHAR(255) NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE admins DROP COLUMN IF EXISTS email;
