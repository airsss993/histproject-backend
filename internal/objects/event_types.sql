-- Таблица с типами меток
CREATE TABLE histproject.event_types
(
    id          SERIAL PRIMARY KEY,    -- ID метки
    name        VARCHAR(50)  NOT NULL, -- Название типа метки
    description VARCHAR(200) NOT NULL  -- Описание типа метки
);

CREATE INDEX ind_event_types_id ON histproject.event_types (id);