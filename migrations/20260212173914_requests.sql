-- +goose Up
-- +goose StatementBegin

-- Таблица заявок
CREATE TYPE request_status AS ENUM ('Новая', 'На модерации', 'Отклонена', 'Опубликована');

CREATE TABLE requests
(
    id                SERIAL PRIMARY KEY,                                  -- Уникальный идентификатор объекта
    title             VARCHAR(200)   NOT NULL,                             -- Название объекта
    description       TEXT           NOT NULL,                             -- Описание объекта
    event_date        DATE           NOT NULL,                             -- Дата исторического события
    event_type_id     INTEGER        NOT NULL REFERENCES event_types (id), -- ID типа события
    email             VARCHAR(70)   NOT NULL,                             -- Почта пользователя, который отправил заявку
    telegram_username VARCHAR(100)   NOT NULL,                             -- Имя пользователя в Telegram, который отправил заявку
    archive_id        VARCHAR(200)   NOT NULL,                             -- ID архива с сайтом в бакете
    site_url          VARCHAR(200),                                        -- URL развернутого сайта для просмотра
    screenshot_url    VARCHAR(200),                                        -- Скриншот главной страницы
    status            request_status NOT NULL DEFAULT 'Новая',             -- Статус заявки
    admin_comment     TEXT,                                                -- Комментарий админа, при отклонение заявки
    created_at        TIMESTAMP      NOT NULL DEFAULT NOW(),               -- Дата создания записи
    updated_at        TIMESTAMP      NOT NULL DEFAULT NOW()                -- Дата последнего обновления записи
);

CREATE INDEX ind_requests_id ON requests (id);
CREATE INDEX ind_requests_email ON requests (email);
CREATE INDEX ind_requests_telegram_username ON requests (telegram_username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS requests CASCADE;
DROP TYPE IF EXISTS request_status CASCADE;
-- +goose StatementEnd
