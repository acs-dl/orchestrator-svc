-- +migrate Up

ALTER TABLE requests
    ALTER COLUMN status TYPE VARCHAR(255);

DROP TYPE IF EXISTS request_status_enum CASCADE;
CREATE TYPE request_status_enum AS ENUM (
    'pending',
    'in progress',
    'success',
    'invited',
    'not found',
    'failed'
    );

ALTER TABLE requests
    ALTER COLUMN status TYPE request_status_enum
        USING (status::request_status_enum);

-- +migrate Down
