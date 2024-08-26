-- SELECT 'CREATE DATABASE auth'
-- WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'mydb')\gexec
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id serial,
    PRIMARY KEY(id)
);

CREATE TABLE sessions (
    user_id int references users(id),
    refresh_token bytea not null,
    ip text not null,
    expired_at date not null
);

CREATE UNIQUE INDEX idx_refresh_token on sessions (user_id);