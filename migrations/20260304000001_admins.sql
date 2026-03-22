-- +goose Up

-- Роли администраторов: super_admin (полный доступ), admin (обычный модератор)
CREATE TYPE admin_role AS ENUM ('super_admin', 'admin');

-- Таблица администраторов для авторизации в админ-панели
CREATE TABLE admins
(
    id            SERIAL PRIMARY KEY,                        -- Уникальный идентификатор администратора
    login         VARCHAR(100) NOT NULL UNIQUE,              -- Логин для входа
    password_hash VARCHAR(255) NOT NULL,                     -- Хеш пароля
    role          admin_role   NOT NULL DEFAULT 'admin',     -- Роль администратора
    created_at    TIMESTAMP    NOT NULL DEFAULT NOW()        -- Дата создания записи
);

-- Гарантирует единственность super_admin на уровне БД
CREATE UNIQUE INDEX unique_super_admin ON admins ((true)) WHERE role = 'super_admin';

-- +goose Down
DROP TABLE IF EXISTS admins;
DROP TYPE IF EXISTS admin_role;