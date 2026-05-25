-- PostgreSQL schema for insight
-- Note: create the database (insight) before running this script:
--   CREATE DATABASE insight;
-- Then connect to it:
--   \c insight
-- and run the rest of this file.

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
COMMENT ON TABLE  user_fans                IS '用户粉丝表';
COMMENT ON COLUMN user_fans.user_id        IS '用户id';
COMMENT ON COLUMN user_fans.follower_uid   IS '粉丝的uid';
COMMENT ON COLUMN user_fans.status         IS '状态 1:已关注 0:取消关注';

INSERT INTO user_fans (id, user_id, follower_uid, status, created_at, updated_at) VALUES
  (1, 2,  1, 1, '2020-05-23 00:12:30', NULL),
  (2, 4,  1, 1, '2020-05-23 00:23:10', NULL),
  (3, 12, 1, 1, '2020-05-23 00:25:48', '2020-05-23 00:27:03'),
  (5, 13, 1, 1, '2020-05-29 12:50:54', NULL);
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
COMMENT ON TABLE  user_follow              IS '用户关注表';
COMMENT ON COLUMN user_follow.user_id      IS '发起关注的人';
COMMENT ON COLUMN user_follow.followed_uid IS '被关注用户的uid';
COMMENT ON COLUMN user_follow.status       IS '关注状态 1:已关注 0:取消关注';

INSERT INTO user_follow (id, user_id, followed_uid, status, created_at, updated_at) VALUES
  (1, 1, 2,  1, '2020-05-23 00:12:30', NULL),
  (2, 1, 4,  1, '2020-05-23 00:23:10', NULL),
  (3, 1, 12, 1, '2020-05-23 00:25:48', '2020-05-23 00:27:03'),
  (5, 1, 13, 1, '2020-05-29 12:50:54', NULL);
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
COMMENT ON TABLE  user_stat                 IS '用户统计表';
COMMENT ON COLUMN user_stat.user_id         IS '用户id';
COMMENT ON COLUMN user_stat.follow_count    IS '关注数';
COMMENT ON COLUMN user_stat.follower_count  IS '粉丝数';
COMMENT ON COLUMN user_stat.status          IS '状态  1:正常';

INSERT INTO user_stat (id, user_id, follow_count, follower_count, status, created_at, updated_at) VALUES
  (1,  1,  3, 0, 1, '2020-05-23 00:12:30', '2020-05-29 12:50:54'),
  (2,  2,  0, 0, 1, '2020-05-23 00:12:30', '2020-05-23 00:20:09'),
  (8,  4,  0, 1, 1, '2020-05-23 00:23:10', NULL),
  (10, 12, 0, 1, 1, '2020-05-23 00:25:48', '2020-05-23 00:27:03'),
  (16, 13, 0, 1, 1, '2020-05-29 12:50:54', NULL);
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
COMMENT ON TABLE  user_base            IS '用户表';
COMMENT ON COLUMN user_base.avatar     IS '头像';
COMMENT ON COLUMN user_base.phone      IS '手机号';
COMMENT ON COLUMN user_base.email      IS '邮箱';
COMMENT ON COLUMN user_base.sex        IS '性别 0:未知 1:男 2:女';

INSERT INTO user_base (id, username, password, avatar, phone, email, sex, deleted_at, created_at, updated_at) VALUES
  (1,  'test-name', '$2a$10$WhJY.MCtsp5kmnyl/UAdQuWbbMzxvmLCPeDhcpxyL84lYey829/ym', '/uploads/avatar.jpg', 13010102020, '123@cc.com',   1, NULL, '2020-02-09 10:23:33', '2020-05-09 10:23:33'),
  (2,  'admin',     '$2a$10$WhJY.MCtsp5kmnyl/UAdQuWbbMzxvmLCPeDhcpxyL84lYey829/ym', '',                   1,           '1234@cc.com',  0, NULL, '2020-05-20 22:42:18', '2020-05-20 22:42:18'),
  (4,  'admin2',    '$2a$10$Dps9oN3Oe3ZDMACih3DCGeTvR.jW/I8WD1NqapCJ6Vq3PzjnusI9i', '',                   0,           '12345@cc.com', 0, NULL, '2020-05-20 22:43:21', '2020-05-20 22:43:21'),
  (12, 'user001',   '123456',                                                       '',                   13810002000, '',             0, NULL, '2020-01-01 00:00:00', '2020-01-01 00:00:00'),
  (13, 'user002',   '123456',                                                       '',                   13810002001, '',             0, NULL, '2020-01-01 00:00:00', '2020-01-01 00:00:00');
SELECT setval(pg_get_serial_sequence('user_base', 'id'), (SELECT MAX(id) FROM user_base));
