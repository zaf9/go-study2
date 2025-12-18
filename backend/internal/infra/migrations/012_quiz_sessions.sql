-- 012_quiz_sessions.sql
-- 测验会话表：记录单次测验的汇总信息
BEGIN TRANSACTION;

DROP TABLE IF EXISTS quiz_sessions;

CREATE TABLE quiz_sessions (
    session_id TEXT PRIMARY KEY,           -- UUID 字符串作为主键
    user_id INTEGER NOT NULL,              -- 用户 ID
    topic TEXT NOT NULL,                   -- 主题 (如 lexical)
    chapter TEXT NOT NULL,                 -- 章节 (如 comment)
    total_questions INTEGER NOT NULL DEFAULT 0, -- 总题数
    correct_answers INTEGER NOT NULL DEFAULT 0, -- 正确数
    score INTEGER NOT NULL DEFAULT 0,      -- 百分制得分 (0-100)
    passed INTEGER NOT NULL DEFAULT 0,     -- 是否通过 (0, 1)
    started_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 开始时间
    completed_at DATETIME,                 -- 完成时间
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_quiz_sessions_user_completed ON quiz_sessions(user_id, completed_at DESC);

COMMIT;
