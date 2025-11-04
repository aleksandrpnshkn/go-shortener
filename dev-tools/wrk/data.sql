-- тестовые данные для запуска wrk

INSERT INTO users VALUES (DEFAULT);

INSERT INTO urls (code, original_url, user_id)
SELECT 'tEsT' AS code, 
    'https://example.com' AS original_url, 
    id AS user_id
FROM users
LIMIT 1;
