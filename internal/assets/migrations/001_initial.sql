-- +migrate Up

CREATE TYPE request_status_enum AS ENUM ('created', 'pending', 'finished', 'failed');

CREATE TABLE IF NOT EXISTS modules (
    name TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    topic TEXT NOT NULL,
    prefix TEXT NOT NULL,
    link TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS requests (
    id UUID PRIMARY KEY,
    from_user_id BIGINT NOT NULL,
    to_user_id BIGINT NOT NULL,
    payload JSONB NOT NULL,
    status request_status_enum NOT NULL DEFAULT 'created',
    module_name TEXT NOT NULL,
    error TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_modules_name FOREIGN KEY (module_name) REFERENCES modules (name) ON DELETE CASCADE
);

-- +migrate Down

DROP TABLE IF EXISTS requests;
DROP TABLE IF EXISTS modules;
DROP TYPE IF EXISTS request_status_enum;