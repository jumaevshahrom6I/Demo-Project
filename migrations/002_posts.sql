CREATE TABLE IF NOT EXISTS posts (
    id         BIGSERIAL   PRIMARY KEY,
    -- удалили пользователя — удаляются и его посты
    user_id    BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title      TEXT        NOT NULL,
    content    TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- лента: новые посты первыми
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts (created_at DESC);
-- посты одного автора
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts (user_id, created_at DESC);
