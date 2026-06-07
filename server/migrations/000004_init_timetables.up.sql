-- 000004: 初始化课表相关表
CREATE TABLE IF NOT EXISTS class_timetables (
    id                  BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    class_id            BIGINT UNSIGNED NOT NULL COMMENT '所属班级',
    day_of_week         TINYINT         NOT NULL COMMENT '周几：1=周一...7=周日',
    period_start        TINYINT         NOT NULL COMMENT '第几节开始(1-12)',
    period_end          TINYINT         NOT NULL COMMENT '第几节结束(1-12)',
    course_name         VARCHAR(100)    NOT NULL COMMENT '课程名称',
    teacher             VARCHAR(50)     DEFAULT NULL COMMENT '授课教师',
    room                VARCHAR(50)     DEFAULT NULL COMMENT '教室',
    contributor_user_id BIGINT UNSIGNED DEFAULT NULL COMMENT '录入人',
    version             INT             NOT NULL DEFAULT 1 COMMENT '版本号',
    status              TINYINT         NOT NULL DEFAULT 1 COMMENT '1=有效 0=删除 2=已纠错替换',
    created_at          DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_class_day (class_id, day_of_week),
    KEY idx_class_status (class_id, status),
    CONSTRAINT fk_ct_class FOREIGN KEY (class_id) REFERENCES classes(id),
    CONSTRAINT fk_ct_contributor FOREIGN KEY (contributor_user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='班级公共课表';

CREATE TABLE IF NOT EXISTS personal_timetables (
    id                      BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id                 BIGINT UNSIGNED NOT NULL COMMENT '所属用户',
    class_id                BIGINT UNSIGNED DEFAULT NULL COMMENT '关联班级',
    day_of_week             TINYINT         NOT NULL COMMENT '周几',
    period_start            TINYINT         NOT NULL COMMENT '第几节开始',
    period_end              TINYINT         NOT NULL COMMENT '第几节结束',
    course_name             VARCHAR(100)    NOT NULL COMMENT '课程名称',
    source                  VARCHAR(20)     NOT NULL COMMENT 'inherited/personal',
    ref_class_timetable_id BIGINT UNSIGNED DEFAULT NULL COMMENT '继承来源公共课表ID',
    is_overridden           TINYINT         NOT NULL DEFAULT 0 COMMENT '是否被覆盖',
    created_at              DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at              DATETIME        DEFAULT NULL,
    PRIMARY KEY (id),
    KEY idx_user_day (user_id, day_of_week),
    KEY idx_user_class (user_id, class_id),
    KEY idx_ref_ct (ref_class_timetable_id),
    CONSTRAINT fk_pt_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_pt_class FOREIGN KEY (class_id) REFERENCES classes(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='个人课表';

CREATE TABLE IF NOT EXISTS timetable_corrections (
    id                     BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    class_timetable_id     BIGINT UNSIGNED NOT NULL COMMENT '关联公共课表条目',
    reporter_user_id       BIGINT UNSIGNED NOT NULL COMMENT '举报人',
    correction_type        VARCHAR(20)     NOT NULL COMMENT 'error/missing',
    description            TEXT            DEFAULT NULL COMMENT '纠错描述',
    suggested_course_name  VARCHAR(100)    DEFAULT NULL COMMENT '建议课程名',
    suggested_period_start TINYINT         DEFAULT NULL COMMENT '建议开始节次',
    suggested_period_end   TINYINT         DEFAULT NULL COMMENT '建议结束节次',
    status                 TINYINT         NOT NULL DEFAULT 0 COMMENT '0=待审核 1=已采纳 2=已驳回',
    reviewed_by            BIGINT UNSIGNED DEFAULT NULL COMMENT '审核人',
    reviewed_at            DATETIME        DEFAULT NULL COMMENT '审核时间',
    created_at             DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    resolved_at            DATETIME        DEFAULT NULL COMMENT '解决时间',
    PRIMARY KEY (id),
    KEY idx_ct_id (class_timetable_id),
    KEY idx_status (status),
    CONSTRAINT fk_tc_ct FOREIGN KEY (class_timetable_id) REFERENCES class_timetables(id),
    CONSTRAINT fk_tc_reporter FOREIGN KEY (reporter_user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='课表纠错表';
