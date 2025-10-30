CREATE TABLE urls (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    code VARCHAR(8) NOT NULL UNIQUE,
    -- Примерный лимит по инфе из https://stackoverflow.com/a/417184
    original_url VARCHAR(2048) NOT NULL UNIQUE
);
