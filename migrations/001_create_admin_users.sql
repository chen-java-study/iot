-- 001_create_admin_users.sql
-- 管理员用户表

CREATE TABLE IF NOT EXISTS admin_users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    real_name VARCHAR(100),
    phone VARCHAR(20),
    email VARCHAR(100),
    status SMALLINT DEFAULT 1, -- 1:启用 0:禁用
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_admin_username ON admin_users(username);
CREATE INDEX idx_admin_status ON admin_users(status);

-- 插入默认管理员账号 (密码: admin123)
INSERT INTO admin_users (username, password_hash, real_name, status) VALUES
('admin', 'admin123', '系统管理员', 1);

COMMENT ON TABLE admin_users IS '管理员用户表';
COMMENT ON COLUMN admin_users.status IS '状态: 1=启用 0=禁用';
