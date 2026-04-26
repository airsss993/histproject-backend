-- +goose Up
-- +goose StatementBegin

-- Снимаем FK-ограничения и удаляем таблицу event_types
ALTER TABLE requests DROP CONSTRAINT requests_event_type_id_fkey;
ALTER TABLE objects DROP CONSTRAINT objects_event_type_id_fkey;
DROP TABLE IF EXISTS event_types CASCADE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

CREATE TABLE event_types
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(50)  NOT NULL,
    description VARCHAR(200) NOT NULL
);

CREATE INDEX ind_event_types_id ON event_types (id);

ALTER TABLE requests ADD CONSTRAINT requests_event_type_id_fkey FOREIGN KEY (event_type_id) REFERENCES event_types (id);
ALTER TABLE objects ADD CONSTRAINT objects_event_type_id_fkey FOREIGN KEY (event_type_id) REFERENCES event_types (id);

-- +goose StatementEnd
