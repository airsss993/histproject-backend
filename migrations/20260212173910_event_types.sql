-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS postgis;

-- Таблица с типами меток
CREATE TABLE event_types
(
    id          SERIAL PRIMARY KEY,    -- ID метки
    name        VARCHAR(50)  NOT NULL, -- Название типа метки
    description VARCHAR(200) NOT NULL  -- Описание типа метки
);

CREATE INDEX ind_event_types_id ON event_types (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS event_types CASCADE;
-- +goose StatementEnd
