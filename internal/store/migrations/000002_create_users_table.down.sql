ALTER TABLE urls
DROP CONSTRAINT urls_id_pkey,
DROP user_id CASCADE,
DROP is_deleted;

DROP TABLE IF EXISTS users;
