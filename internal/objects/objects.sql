CREATE EXTENSION IF NOT EXISTS postgis;

-- Таблица всех исторических объектов, которые будут отображены на карте
CREATE TABLE histproject.objects
(
    id                SERIAL PRIMARY KEY,                                          -- Уникальный идентификатор объекта
    request_id        INTEGER                NOT NULL,                             -- ID заявки, из которой создан объект
    title             VARCHAR(200)           NOT NULL,                             -- Название объекта
    description       TEXT                   NOT NULL,                             -- Описание объекта
    coordinates       GEOGRAPHY(Point, 4326) NOT NULL,                             -- Географические координаты
    event_date        DATE                   NOT NULL,                             -- Дата исторического события
    event_type_id     INTEGER                NOT NULL REFERENCES event_types (id), -- ID типа события
    preview_image_url VARCHAR(200),                                                -- URL объекта для отображения в карточке
    created_at        TIMESTAMP              NOT NULL DEFAULT NOW(),               -- Дата создания записи
    updated_at        TIMESTAMP              NOT NULL DEFAULT NOW()                -- Дата последнего обновления
);

CREATE INDEX ind_objects_event_type ON histproject.objects (event_type_id);
CREATE INDEX ind_objects_coordinates ON histproject.objects USING GIST (coordinates);
CREATE INDEX ind_objects_request_id ON histproject.objects (request_id);