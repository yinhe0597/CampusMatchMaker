-- 000002: 初始化用户表 + 认证表
CREATE TABLE IF NOT EXISTS users (
    id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    student_id     VARCHAR(256)    NOT NULL COMMENT '学号（SHA-256哈希存储）',
    nickname       VARCHAR(50)     NOT NULL COMMENT '昵称',
    avatar_url     VARCHAR(500)    DEFAULT NULL COMMENT '头像地址',
    phone          VARCHAR(128)    DEFAULT NULL COMMENT '手机号（加密存储）',
    password_hash  VARCHAR(255)    NOT NULL COMMENT '密码哈希',
    school_id      BIGINT UNSIGNED DEFAULT NULL COMMENT '所属学校',
    status         TINYINT         NOT NULL DEFAULT 1 COMMENT '1=正常 0=禁用',
    privacy_level  TINYINT         NOT NULL DEFAULT 1 COMMENT '隐私级别：1=最小暴露',
    created_at     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at     DATETIME        DEFAULT NULL COMMENT '软删除时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_student_id (student_id),
    KEY idx_school_id (school_id),
    KEY idx_deleted_at (deleted_at),
    CONSTRAINT fk_users_school FOREIGN KEY (school_id) REFERENCES schools(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

CREATE TABLE IF NOT EXISTS student_auth (
    id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id       BIGINT UNSIGNED NOT NULL COMMENT '关联用户',
    student_id    VARCHAR(256)    NOT NULL COMMENT '学号（明文存储，供审核用）',
    school_id     BIGINT UNSIGNED NOT NULL COMMENT '所属学校',
    auth_method   VARCHAR(20)     NOT NULL COMMENT '认证方式：edu_email/manual',
    auth_status   TINYINT         NOT NULL DEFAULT 0 COMMENT '0=待验证 1=已认证 2=失败',
    verified_at   DATETIME        DEFAULT NULL COMMENT '认证通过时间',
    created_at    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_user_id (user_id),
    KEY idx_school_student (school_id, student_id),
    CONSTRAINT fk_auth_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_auth_school FOREIGN KEY (school_id) REFERENCES schools(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='学号认证表';
