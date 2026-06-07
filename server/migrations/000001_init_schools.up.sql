-- 000001: 初始化学校表
CREATE TABLE IF NOT EXISTS schools (
    id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name         VARCHAR(100)    NOT NULL COMMENT '学校名称',
    code         VARCHAR(20)     NOT NULL COMMENT '学校代码',
    province     VARCHAR(20)     DEFAULT NULL COMMENT '所在省份',
    city         VARCHAR(20)     DEFAULT NULL COMMENT '所在城市',
    status       TINYINT         NOT NULL DEFAULT 1 COMMENT '1=启用 0=禁用',
    created_at   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='学校表';
