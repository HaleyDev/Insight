-- PostgreSQL schema for insight
-- Note: create the database (insight) before running this script:
--   CREATE DATABASE insight;
-- Then connect to it:
--   \c insight
-- and run the rest of this file.

-- ------------------------------------------------------------
-- 清理旧表（示例代码遗留）
-- ------------------------------------------------------------
DROP TABLE IF EXISTS user_fans;
DROP TABLE IF EXISTS user_follow;
DROP TABLE IF EXISTS user_stat;
DROP TABLE IF EXISTS user_base;

-- ------------------------------------------------------------
-- table user_base
-- 字段与前端 UserInfo 对齐：id / username / email / avatar / role
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
COMMENT ON TABLE  user_base            IS '用户表';
COMMENT ON COLUMN user_base.username   IS '用户名';
COMMENT ON COLUMN user_base.password   IS '密码（bcrypt）';
COMMENT ON COLUMN user_base.email      IS '邮箱';
COMMENT ON COLUMN user_base.avatar     IS '头像 URL';
COMMENT ON COLUMN user_base.role       IS '角色：admin / user';

-- 默认管理员账号：admin@insight.com / 123456 （bcrypt cost=10）
INSERT INTO user_base (id, username, password, email, avatar, role, created_at, updated_at) VALUES
  (1, 'admin', '$2a$10$WhJY.MCtsp5kmnyl/UAdQuWbbMzxvmLCPeDhcpxyL84lYey829/ym', 'admin@insight.com', '/images/avatars/avatar-1.png', 'admin', NOW(), NOW());
SELECT setval(pg_get_serial_sequence('user_base', 'id'), (SELECT MAX(id) FROM user_base));
