-- PostgreSQL init script for docker-compose
-- The official postgres image creates the DB defined by POSTGRES_DB before
-- running scripts in /docker-entrypoint-initdb.d/, so we connect directly
-- and create tables.

\c insight;

-- ------------------------------------------------------------
-- 清理旧表
-- ------------------------------------------------------------
DROP TABLE IF EXISTS user_fans;
DROP TABLE IF EXISTS user_follow;
DROP TABLE IF EXISTS user_stat;
DROP TABLE IF EXISTS user_base;

-- ------------------------------------------------------------
-- table user_base
-- ------------------------------------------------------------
CREATE TABLE user_base (
    id          BIGSERIAL     PRIMARY KEY,
    username    VARCHAR(64)   NOT NULL,
    password    VARCHAR(255)  NOT NULL,
    email       VARCHAR(255)  NOT NULL,
    avatar      VARCHAR(255)  NOT NULL DEFAULT '/images/avatars/avatar-1.png',
    role        VARCHAR(32)   NOT NULL DEFAULT 'user',
    created_at  TIMESTAMP     NULL,
    updated_at  TIMESTAMP     NULL,
    CONSTRAINT uniq_username UNIQUE (username),
    CONSTRAINT uniq_email    UNIQUE (email)
);

INSERT INTO user_base (id, username, password, email, avatar, role, created_at, updated_at) VALUES
  (1, 'admin', '$2a$10$WhJY.MCtsp5kmnyl/UAdQuWbbMzxvmLCPeDhcpxyL84lYey829/ym', 'admin@insight.com', '/images/avatars/avatar-1.png', 'admin', NOW(), NOW());
SELECT setval(pg_get_serial_sequence('user_base', 'id'), (SELECT MAX(id) FROM user_base));
