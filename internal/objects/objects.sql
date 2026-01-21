CREATE TABLE objects
(
    id                SERIAL PRIMARY KEY,                            -- Уникальный идентификатор объекта
    submission_id     INTEGER,                                       -- ID заявки, из которой создан объект
    title             VARCHAR(200)           NOT NULL,               -- Название объекта
    description       TEXT,                                          -- Описание объекта
    coordinates       GEOGRAPHY(Point, 4326) NOT NULL,               -- Географические координаты
    event_date        DATE,                                          -- Дата исторического события
    event_type_id     INTEGER,                                       -- ID типа события
    preview_image_url VARCHAR(500),                                  -- URL объекта для отображения в карточке
    created_at        TIMESTAMP              NOT NULL DEFAULT NOW(), -- Дата создания записи
    updated_at        TIMESTAMP              NOT NULL DEFAULT NOW()  -- Дата последнего обновления
); -- Таблица всех исторических объектов, которые будут отображены на карте

CREATE INDEX ind_objects_event_type ON objects (event_type_id);
CREATE INDEX ind_objects_coordinates ON objects USING GIST (coordinates);
CREATE INDEX ind_objects_submission_id ON objects (submission_id);