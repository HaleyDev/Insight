-- PostgreSQL init script for docker-compose
-- The official postgres image creates the DB defined by POSTGRES_DB before
-- running scripts in /docker-entrypoint-initdb.d/, so we connect directly
-- and create tables.

\c insight;

-- ------------------------------------------------------------
-- table user_fans
-- ------------------------------------------------------------
DROP TABLE IF EXISTS user_fans;
CREATE TABLE user_fans (
    id            BIGSERIAL    PRIMARY KEY,
    user_id       BIGINT       NOT NULL DEFAULT 0,
    follower_uid  BIGINT       NOT NULL DEFAULT 0,
    status        SMALLINT     NOT NULL DEFAULT 0,
    created_at    TIMESTAMP    NULL,
    updated_at    TIMESTAMP    NULL,
    CONSTRAINT idx_uid_fid UNIQUE (user_id, follower_uid)
);
CREATE INDEX idx_user_fans_status_uid ON user_fans (status, user_id);

INSERT INTO user_fans (id, user_id, follower_uid, status, created_at, updated_at) VALUES
  (1, 2, 1, 1, '2020-05-23 00:12:30', NULL),
  (2, 3, 1, 1, '2020-05-23 00:23:10', NULL);
SELECT setval(pg_get_serial_sequence('user_fans', 'id'), (SELECT MAX(id) FROM user_fans));


-- ------------------------------------------------------------
-- table user_follow
-- ------------------------------------------------------------
DROP TABLE IF EXISTS user_follow;
CREATE TABLE user_follow (
    id            BIGSERIAL    PRIMARY KEY,
    user_id       BIGINT       NOT NULL DEFAULT 0,
    followed_uid  BIGINT       NOT NULL DEFAULT 0,
    status        SMALLINT     NOT NULL DEFAULT 0,
    created_at    TIMESTAMP    NULL,
    updated_at    TIMESTAMP    NULL,
    CONSTRAINT uniq_uid_fuid UNIQUE (user_id, followed_uid)
);
CREATE INDEX idx_user_follow_status_uid ON user_follow (status, user_id);

INSERT INTO user_follow (id, user_id, followed_uid, status, created_at, updated_at) VALUES
  (1, 1, 2, 1, '2020-05-23 00:12:30', NULL),
  (2, 1, 3, 1, '2020-05-23 00:23:10', NULL);
SELECT setval(pg_get_serial_sequence('user_follow', 'id'), (SELECT MAX(id) FROM user_follow));


-- ------------------------------------------------------------
-- table user_stat
-- ------------------------------------------------------------
DROP TABLE IF EXISTS user_stat;
CREATE TABLE user_stat (
    id              BIGSERIAL    PRIMARY KEY,
    user_id         BIGINT       NOT NULL DEFAULT 0,
    follow_count    BIGINT       NOT NULL DEFAULT 0,
    follower_count  BIGINT       NOT NULL DEFAULT 0,
    status          SMALLINT     NOT NULL DEFAULT 1,
    created_at      TIMESTAMP    NULL,
    updated_at      TIMESTAMP    NULL,
    CONSTRAINT uniq_uid UNIQUE (user_id)
);
CREATE INDEX idx_user_stat_status ON user_stat (status);

INSERT INTO user_stat (id, user_id, follow_count, follower_count, status, created_at, updated_at) VALUES
  (1, 1, 3, 0, 1, '2020-05-23 00:12:30', '2020-05-29 12:50:54'),
  (2, 2, 0, 0, 1, '2020-05-23 00:12:30', '2020-05-23 00:20:09'),
  (8, 3, 0, 1, 1, '2020-05-23 00:23:10', NULL);
SELECT setval(pg_get_serial_sequence('user_stat', 'id'), (SELECT MAX(id) FROM user_stat));


-- ------------------------------------------------------------
-- table user_base
-- ------------------------------------------------------------
DROP TABLE IF EXISTS user_base;
CREATE TABLE user_base (
    id          BIGSERIAL     PRIMARY KEY,
    username    VARCHAR(255)  NOT NULL DEFAULT '',
    password    VARCHAR(60)   NOT NULL DEFAULT '',
    avatar      VARCHAR(255)  NOT NULL DEFAULT '',
    phone       BIGINT        NOT NULL DEFAULT 0,
    email       VARCHAR(255)  NOT NULL DEFAULT '',
    sex         SMALLINT      NOT NULL DEFAULT 0,
    deleted_at  TIMESTAMP     NULL,
    created_at  TIMESTAMP     NULL,
    updated_at  TIMESTAMP     NULL,
    CONSTRAINT uniq_username UNIQUE (username),
    CONSTRAINT uniq_phone    UNIQUE (phone)
);

INSERT INTO user_base (id, username, password, avatar, phone, email, sex, deleted_at, created_at, updated_at) VALUES
  (1, 'test-name', '$2a$10$WhJY.MCtsp5kmnyl/UAdQuWbbMzxvmLCPeDhcpxyL84lYey829/ym', '/uploads/avatar.jpg', 13010102020, '123@cc.com',   1, NULL, '2020-02-09 10:23:33', '2020-05-09 10:23:33'),
  (2, 'admin',     '$2a$10$WhJY.MCtsp5kmnyl/UAdQuWbbMzxvmLCPeDhcpxyL84lYey829/ym', '13010102021',         1,           '1234@cc.com',  0, NULL, '2020-05-20 22:42:18', '2020-05-20 22:42:18'),
  (3, 'admin2',    '$2a$10$Dps9oN3Oe3ZDMACih3DCGeTvR.jW/I8WD1NqapCJ6Vq3PzjnusI9i', '13010102022',         0,           '12345@cc.com', 0, NULL, '2020-05-20 22:43:21', '2020-05-20 22:43:21');
SELECT setval(pg_get_serial_sequence('user_base', 'id'), (SELECT MAX(id) FROM user_base));
