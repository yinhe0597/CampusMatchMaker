-- 000005: 初始化投票相关表
CREATE TABLE IF NOT EXISTS polls (
    id                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    creator_user_id   BIGINT UNSIGNED NOT NULL COMMENT '创建者',
    title             VARCHAR(200)    NOT NULL COMMENT '投票标题',
    description       TEXT            DEFAULT NULL COMMENT '投票说明',
    scope_type        VARCHAR(20)     NOT NULL COMMENT '范围：class/group',
    scope_id          BIGINT UNSIGNED NOT NULL COMMENT '范围ID',
    status            VARCHAR(20)     NOT NULL DEFAULT 'draft' COMMENT 'draft/open/closed/finalized',
    deadline          DATETIME        DEFAULT NULL COMMENT '截止时间',
    min_participants  INT             DEFAULT 2 COMMENT '最少参与人数',
    final_option_id   BIGINT UNSIGNED DEFAULT NULL COMMENT '最终确定的选项',
    created_at        DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    closed_at         DATETIME        DEFAULT NULL COMMENT '关闭时间',
    PRIMARY KEY (id),
    KEY idx_creator (creator_user_id),
    KEY idx_scope (scope_type, scope_id),
    KEY idx_status (status),
    CONSTRAINT fk_poll_creator FOREIGN KEY (creator_user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='投票表';

CREATE TABLE IF NOT EXISTS poll_options (
    id                  BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    poll_id             BIGINT UNSIGNED NOT NULL COMMENT '所属投票',
    slot_date           DATE            NOT NULL COMMENT '时段日期',
    slot_start_time     TIME            NOT NULL COMMENT '开始时间',
    slot_end_time       TIME            NOT NULL COMMENT '结束时间',
    is_recommended      TINYINT         NOT NULL DEFAULT 0 COMMENT '是否引擎推荐',
    recommendation_rate DECIMAL(5,4)    DEFAULT NULL COMMENT '推荐参与率(0.0000~1.0000)',
    sort_order          INT             NOT NULL DEFAULT 0 COMMENT '排序序号',
    created_at          DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_poll_sort (poll_id, sort_order),
    CONSTRAINT fk_option_poll FOREIGN KEY (poll_id) REFERENCES polls(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='投票选项表';

CREATE TABLE IF NOT EXISTS poll_votes (
    id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    poll_id        BIGINT UNSIGNED NOT NULL COMMENT '所属投票',
    option_id      BIGINT UNSIGNED NOT NULL COMMENT '选择的选项',
    voter_user_id  BIGINT UNSIGNED NOT NULL COMMENT '投票人',
    choice         VARCHAR(10)     NOT NULL COMMENT 'yes/no/maybe',
    voted_at       DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_poll_option_voter (poll_id, option_id, voter_user_id),
    KEY idx_poll_voter (poll_id, voter_user_id),
    CONSTRAINT fk_vote_poll FOREIGN KEY (poll_id) REFERENCES polls(id),
    CONSTRAINT fk_vote_option FOREIGN KEY (option_id) REFERENCES poll_options(id),
    CONSTRAINT fk_vote_voter FOREIGN KEY (voter_user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='投票记录表';

CREATE TABLE IF NOT EXISTS poll_results (
    id                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    poll_id            BIGINT UNSIGNED NOT NULL COMMENT '所属投票',
    option_id          BIGINT UNSIGNED NOT NULL COMMENT '对应选项',
    yes_count          INT             NOT NULL DEFAULT 0,
    no_count           INT             NOT NULL DEFAULT 0,
    maybe_count        INT             NOT NULL DEFAULT 0,
    total_votes        INT             NOT NULL DEFAULT 0,
    participation_rate DECIMAL(5,4)    NOT NULL DEFAULT 0 COMMENT '参与率',
    calculated_at      DATETIME        NOT NULL COMMENT '计算时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_poll_option (poll_id, option_id),
    CONSTRAINT fk_result_poll FOREIGN KEY (poll_id) REFERENCES polls(id),
    CONSTRAINT fk_result_option FOREIGN KEY (option_id) REFERENCES poll_options(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='投票结果汇总表';
