-- +goose Up
-- +goose StatementBegin

-- Таблица заявок
CREATE TYPE request_status AS ENUM ('Новая', 'Принята', 'Отклонена', 'Опубликована');

CREATE TABLE requests
(
    id                SERIAL PRIMARY KEY,                      -- Уникальный идентификатор объекта
    title             VARCHAR(200)   NOT NULL,                 -- Название объекта
    description       TEXT           NOT NULL,                 -- Описание объекта
    email             VARCHAR(100)   NOT NULL,                 -- Почта пользователя, который отправил заявку
    telegram_username VARCHAR(100)   NOT NULL,                 -- Имя пользователя в Telegram, который отправил заявку
    archive_url       VARCHAR(200)   NOT NULL,                 -- Ссылка на архив с сайтом в бакете
    status            request_status NOT NULL DEFAULT 'Новая', -- Статус заявки
    admin_comment     TEXT,                                    -- Комментарий админа, при отклонение заявки
    created_at        TIMESTAMP      NOT NULL DEFAULT NOW(),   -- Дата создания записи
    updated_at        TIMESTAMP      NOT NULL DEFAULT NOW()    -- Дата последнего обновления записи
);

CREATE INDEX ind_requests_id ON requests (id);
CREATE INDEX ind_requests_email ON requests (email);
CREATE INDEX ind_requests_telegram_username ON requests (telegram_username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
