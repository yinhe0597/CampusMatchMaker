-- 000006: poll_options 添加 day_of_week 字段
ALTER TABLE poll_options ADD COLUMN day_of_week TINYINT NOT NULL DEFAULT 0 COMMENT '星期几(1-7)' AFTER slot_end_time;
