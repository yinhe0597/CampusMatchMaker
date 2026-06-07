-- 000003: 初始化班级表 + 成员表
CREATE TABLE IF NOT EXISTS classes (
    id                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    school_id         BIGINT UNSIGNED NOT NULL COMMENT '所属学校',
    grade             VARCHAR(20)     NOT NULL COMMENT '年级',
    department        VARCHAR(50)     DEFAULT NULL COMMENT '院系',
    name              VARCHAR(100)    NOT NULL COMMENT '班级名称',
    code              VARCHAR(50)     DEFAULT NULL COMMENT '班级唯一标识',
    creator_user_id   BIGINT UNSIGNED NOT NULL COMMENT '创建者',
    invite_code       VARCHAR(10)     DEFAULT NULL COMMENT '邀请码',
    timetable_status  TINYINT         NOT NULL DEFAULT 0 COMMENT '课表状态：0=未录入 1=已录入',
    created_at        DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_code (code),
    UNIQUE KEY uk_invite_code (invite_code),
    KEY idx_school_grade (school_id, grade),
    CONSTRAINT fk_class_school FOREIGN KEY (school_id) REFERENCES schools(id),
    CONSTRAINT fk_class_creator FOREIGN KEY (creator_user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='班级表';

CREATE TABLE IF NOT EXISTS class_members (
    id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    class_id   BIGINT UNSIGNED NOT NULL COMMENT '所属班级',
    user_id    BIGINT UNSIGNED NOT NULL COMMENT '用户',
    role       VARCHAR(20)     NOT NULL DEFAULT 'member' COMMENT 'owner/admin/member',
    status     TINYINT         NOT NULL DEFAULT 1 COMMENT '1=正常 0=已退出',
    joined_at  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_class_user (class_id, user_id),
    KEY idx_user_id (user_id),
    CONSTRAINT fk_member_class FOREIGN KEY (class_id) REFERENCES classes(id),
    CONSTRAINT fk_member_user FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='班级成员表';
