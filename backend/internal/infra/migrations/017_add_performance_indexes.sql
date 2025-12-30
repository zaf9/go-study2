-- 性能优化：为进度查询添加索引
-- Feature: 016-progress-quiz-ux Phase 8

-- 优化GetByUserAndTopic查询
-- 该索引已存在于 011_learning_progress_quiz.sql 中:
-- CREATE INDEX IF NOT EXISTS idx_learning_progress_user ON learning_progress(user_id);

-- 添加复合索引优化按主题查询进度
CREATE INDEX IF NOT EXISTS idx_learning_progress_user_topic 
ON learning_progress(user_id, topic);

-- 优化最近访问章节查询
-- 该索引已存在:
-- CREATE INDEX IF NOT EXISTS idx_learning_progress_user_last_visit ON learning_progress(user_id, last_visit_at DESC);

-- 优化测验历史查询
CREATE INDEX IF NOT EXISTS idx_quiz_sessions_user_created 
ON quiz_sessions(user_id, created_at DESC);

-- 优化测验答题记录查询
CREATE INDEX IF NOT EXISTS idx_quiz_attempts_session 
ON quiz_attempts(session_id);
