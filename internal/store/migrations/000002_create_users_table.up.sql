CREATE TABLE users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY
);

ALTER TABLE urls
    -- Внешний ключ сделать не получилось, автотест ломается из-за того что дропает таблицы без миграций
    -- https://github.com/Yandex-Practicum/go-autotests/blob/fc542b82d9614ff2f450835e059263f71ee9af92/cmd/wipedb/main.go#L48
    ADD user_id BIGINT; 
