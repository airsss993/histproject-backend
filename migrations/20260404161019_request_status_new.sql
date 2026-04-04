-- +goose Up
-- +goose StatementBegin
ALTER TYPE request_status ADD VALUE 'Новая' BEFORE 'На проверке';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Удаление значения из ENUM в PostgreSQL невозможно напрямую; выполняется пересоздание типа
ALTER TABLE requests ALTER COLUMN status TYPE VARCHAR(50);
DROP TYPE request_status;
CREATE TYPE request_status AS ENUM ('В обработке', 'На проверке', 'Отклонена', 'Опубликована');
ALTER TABLE requests ALTER COLUMN status TYPE request_status USING status::request_status;
-- +goose StatementEnd
