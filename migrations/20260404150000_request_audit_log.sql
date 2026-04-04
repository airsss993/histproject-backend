-- +goose Up
-- Таблица истории изменений заявок: фиксирует действия администраторов
CREATE TABLE
    request_audit_log (
        id SERIAL PRIMARY KEY, -- Уникальный идентификатор записи
        request_id INT NOT NULL REFERENCES requests (id) ON DELETE CASCADE, -- Ссылка на заявку, к которой относится запись
        admin_login VARCHAR(100) NOT NULL, -- Логин администратора, который совершил действие
        admin_role VARCHAR(50) NOT NULL, -- Роль администратора: 'super_admin' | 'admin'
        action VARCHAR(50) NOT NULL, -- Действие: 'Одобрена' | 'Отклонена'
        comment TEXT, -- Комментарий администратора (может быть пустым)
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW () -- Дата и время создания записи
    );

CREATE INDEX ind_request_audit_log_request_id ON request_audit_log (request_id);

-- +goose Down
DROP TABLE IF EXISTS request_audit_log;