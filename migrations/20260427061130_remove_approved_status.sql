-- +goose Up
-- +goose StatementBegin
UPDATE requests SET status = 'На проверке' WHERE status = 'Одобрена';
ALTER TABLE requests ALTER COLUMN status DROP DEFAULT;
ALTER TABLE requests ALTER COLUMN status TYPE VARCHAR(50);
DROP TYPE request_status;
CREATE TYPE request_status AS ENUM ('В обработке', 'Новая', 'На проверке', 'Отклонена', 'Опубликована');
ALTER TABLE requests ALTER COLUMN status TYPE request_status USING status::request_status;
ALTER TABLE requests ALTER COLUMN status SET DEFAULT 'В обработке';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE requests ALTER COLUMN status DROP DEFAULT;
ALTER TABLE requests ALTER COLUMN status TYPE VARCHAR(50);
DROP TYPE request_status;
CREATE TYPE request_status AS ENUM ('В обработке', 'Новая', 'На проверке', 'Одобрена', 'Отклонена', 'Опубликована');
ALTER TABLE requests ALTER COLUMN status TYPE request_status USING status::request_status;
ALTER TABLE requests ALTER COLUMN status SET DEFAULT 'В обработке';
-- +goose StatementEnd
