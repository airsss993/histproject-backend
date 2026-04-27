-- +goose Up
-- +goose StatementBegin
-- Переводим заявки со статусом 'Одобрена' обратно в 'На проверке'
UPDATE requests SET status = 'На проверке' WHERE status = 'Одобрена';
-- Пересоздаём тип без значения 'Одобрена'
ALTER TABLE requests ALTER COLUMN status TYPE VARCHAR(50);
DROP TYPE request_status;
CREATE TYPE request_status AS ENUM ('В обработке', 'Новая', 'На проверке', 'Отклонена', 'Опубликована');
ALTER TABLE requests ALTER COLUMN status TYPE request_status USING status::request_status;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE requests ALTER COLUMN status TYPE VARCHAR(50);
DROP TYPE request_status;
CREATE TYPE request_status AS ENUM ('В обработке', 'Новая', 'На проверке', 'Одобрена', 'Отклонена', 'Опубликована');
ALTER TABLE requests ALTER COLUMN status TYPE request_status USING status::request_status;
-- +goose StatementEnd
