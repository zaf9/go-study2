BEGIN TRANSACTION;

-- 学习进度表：记录章节级进度、阅读行为与测验结果
CREATE TABLE IF NOT EXISTS learning_progress (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    topic TEXT NOT NULL,
    chapter TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'not_started' CHECK(status IN ('not_started','in_progress','completed','tested')),
    read_duration INTEGER NOT NULL DEFAULT 0 CHECK(read_duration >= 0),
    scroll_progress INTEGER NOT NULL DEFAULT 0 CHECK(scroll_progress >= 0 AND scroll_progress <= 100),
    last_position INTEGER NOT NULL DEFAULT 0 CHECK(last_position >= 0),
    quiz_score INTEGER,
    quiz_passed INTEGER NOT NULL DEFAULT 0 CHECK(quiz_passed IN (0,1)),
    first_visit_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_visit_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, topic, chapter)
);
CREATE INDEX IF NOT EXISTS idx_learning_progress_user ON learning_progress(user_id);
CREATE INDEX IF NOT EXISTS idx_learning_progress_user_status ON learning_progress(user_id, status);
CREATE INDEX IF NOT EXISTS idx_learning_progress_user_last_visit ON learning_progress(user_id, last_visit_at DESC);

-- 题库表：按章节存储题目与解析
CREATE TABLE IF NOT EXISTS quiz_questions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    topic TEXT NOT NULL,
    chapter TEXT NOT NULL,
    type TEXT NOT NULL CHECK(type IN ('single','multiple','truefalse','code_output','code_correction')),
    difficulty TEXT NOT NULL CHECK(difficulty IN ('easy','medium','hard')),
    question TEXT NOT NULL,
    options TEXT NOT NULL,
    correct_answers TEXT NOT NULL,
    explanation TEXT NOT NULL,
    code_snippet TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_quiz_questions_topic_chapter ON quiz_questions(topic, chapter);
CREATE INDEX IF NOT EXISTS idx_quiz_questions_difficulty ON quiz_questions(difficulty);

-- 测验会话表：一次测验的得分与通过状态
CREATE TABLE IF NOT EXISTS quiz_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL UNIQUE,
    user_id INTEGER NOT NULL,
    topic TEXT NOT NULL,
    chapter TEXT NOT NULL,
    total_questions INTEGER NOT NULL DEFAULT 0 CHECK(total_questions >= 0),
    correct_answers INTEGER NOT NULL DEFAULT 0 CHECK(correct_answers >= 0),
    score INTEGER NOT NULL DEFAULT 0 CHECK(score >= 0 AND score <= 100),
    passed INTEGER NOT NULL DEFAULT 0 CHECK(passed IN (0,1)),
    started_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_quiz_sessions_user_completed ON quiz_sessions(user_id, completed_at DESC);

-- 测验作答表：存储用户对单题的回答
CREATE TABLE IF NOT EXISTS quiz_attempts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    topic TEXT NOT NULL,
    chapter TEXT NOT NULL,
    question_id INTEGER NOT NULL,
    user_answers TEXT NOT NULL,
    is_correct INTEGER NOT NULL DEFAULT 0 CHECK(is_correct IN (0,1)),
    attempted_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES quiz_sessions(session_id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES quiz_questions(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_quiz_attempts_user ON quiz_attempts(user_id);
CREATE INDEX IF NOT EXISTS idx_quiz_attempts_question ON quiz_attempts(question_id);

COMMIT;

