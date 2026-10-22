CREATE TABLE IF NOT EXISTS roles (
    id          BIGSERIAL   PRIMARY KEY,
    name        TEXT        NOT NULL,
    code        TEXT        NOT NULL UNIQUE,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL   PRIMARY KEY,
    email         TEXT        NOT NULL UNIQUE,
    password_hash TEXT        NOT NULL,
    name          TEXT        NOT NULL,
    role_id       BIGINT      NOT NULL REFERENCES roles(id),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Самостоятельно зарегистрироваться можно только как user (проверяет service),
-- остальные роли назначает администратор.
INSERT INTO roles (name, code, description) VALUES
    ('Пользователь', 'user', 'Создаёт, редактирует и удаляет свои посты'),
    ('Администратор', 'admin', 'Полный доступ: пользователи, роли, настройки'),
    ('Аудитор', 'auditor', 'Только чтение: просмотр данных и истории действий для проверок'),
    ('Модератор', 'moderator', 'Проверяет и удаляет посты, нарушающие правила'),
    ('Поддержка', 'support', 'Отвечает на обращения пользователей, видит их профили'),
    ('Работодатель', 'employer', 'Публикует вакансии и рассматривает отклики'),
    ('Соискатель', 'applicant', 'Ведёт резюме и откликается на вакансии')
ON CONFLICT (code) DO NOTHING;
