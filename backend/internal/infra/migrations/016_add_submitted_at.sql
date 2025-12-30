-- 添加submitted_at字段到quiz_sessions表
-- 用于跟踪测验提交时间,防止重复提交

-- ALTER TABLE quiz_sessions ADD COLUMN submitted_at DATETIME;
-- 注意: submitted_at字段已存在于表结构中,无需再次添加

-- 回填已有session的submitted_at值
-- 对于已提交的测验(有completed_at的),使用completed_at作为submitted_at
UPDATE quiz_sessions 
SET submitted_at = completed_at 
WHERE completed_at IS NOT NULL AND submitted_at IS NULL;
