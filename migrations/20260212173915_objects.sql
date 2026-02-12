-- +goose Up
-- +goose StatementBegin

-- Таблица всех исторических объектов, которые будут отображены на карте
CREATE TABLE objects
(
    id                SERIAL PRIMARY KEY,                                          -- Уникальный идентификатор объекта
    request_id        INTEGER                NOT NULL REFERENCES requests (id),    -- ID заявки, из которой создан объект
    title             VARCHAR(200)           NOT NULL,                             -- Название объекта
    description       TEXT                   NOT NULL,                             -- Описание объекта
    coordinates       GEOGRAPHY(Point, 4326) NOT NULL,                             -- Географические координаты
    event_date        DATE                   NOT NULL,                             -- Дата исторического события
    event_type_id     INTEGER                NOT NULL REFERENCES event_types (id), -- ID типа события
    site_url          VARCHAR(200)           NOT NULL,                             -- URL сайта пользователя
    preview_image_url VARCHAR(200),                                                -- URL скриншота сайта для отображения в карточке
    created_at        TIMESTAMP              NOT NULL DEFAULT NOW(),               -- Дата создания записи
    updated_at        TIMESTAMP              NOT NULL DEFAULT NOW()                -- Дата последнего обновления
);

CREATE INDEX ind_objects_event_type ON objects (event_type_id);
CREATE INDEX ind_objects_coordinates ON objects USING GIST (coordinates);
CREATE INDEX ind_objects_request_id ON objects (request_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS objects CASCADE;
-- +goose StatementEnd
