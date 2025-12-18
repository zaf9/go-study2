-- 013_quiz_attempts.sql
-- 测验作答记录表：记录单次测验中每一道题的作答情况
BEGIN TRANSACTION;

DROP TABLE IF EXISTS quiz_attempts;

CREATE TABLE quiz_attempts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,              -- 关联的会话 ID
    question_id TEXT NOT NULL,             -- 题目 ID (可能是 UUID 或内容哈希)
    user_choice TEXT NOT NULL,             -- 用户的选择 (内容字符串)
    is_correct INTEGER NOT NULL DEFAULT 0, -- 是否正确 (0, 1)
    attempted_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 作答时间
    FOREIGN KEY (session_id) REFERENCES quiz_sessions(session_id) ON DELETE CASCADE
);

CREATE INDEX idx_quiz_attempts_session ON quiz_attempts(session_id);
CREATE INDEX idx_quiz_attempts_question ON quiz_attempts(question_id);

COMMIT;
