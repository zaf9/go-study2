# Data Model: Go-Study2 学习闭环与测验体系

**Branch**: 011-learning-progress-quiz  
**Spec**: specs/011-learning-progress-quiz/spec.md  
**Source**: backend SQLite（GF ORM）

## Entities

### LearningProgress
- Fields: id(PK), user_id, topic, chapter, status(not_started|in_progress|completed|tested), read_duration(int, sec, >=0,累加), scroll_progress(int 0-100,覆盖), last_position(int>=0,覆盖), quiz_score(int 0-100), quiz_passed(bool), first_visit_at, last_visit_at, completed_at, created_at, updated_at  
- Unique: (user_id, topic, chapter)  
- Index: user_id, (user_id,status), (user_id,last_visit_at)  
- Rules: 状态仅随最新 updated_at 前进；完成条件=时长≥80%预估 & 滚动≥90% & quiz_passed==true。离开/卸载做最终同步。

### QuizQuestion
- Fields: id(PK), topic, chapter, type(single|multiple|truefalse|code_output|code_correction), difficulty(easy|medium|hard), question(markdown), options(json array), correct_answers(json array index), explanation(text), code_snippet(text|null), created_at, updated_at  
- Index: (topic,chapter), difficulty  
- Rules: 题目归属固定章节；题型/难度决定抽题比例。

### QuizSession
- Fields: id(PK), user_id, topic, chapter, total_questions, correct_answers, score(0-100), passed(bool), started_at, completed_at  
- Index: (user_id,completed_at)  
- Rules: 每次测验创建新 session；提交后写入结果并更新进度状态（passed→completed，否则 tested）。

### QuizAttempt
- Fields: id(PK), user_id, topic, chapter, question_id(FK QuizQuestion), user_answers(json array), is_correct(bool), attempted_at  
- Index: user_id, question_id  
- Rules: 关联 session 通过时间关联（可额外存 session_id if needed in impl）；用于详细解析与历史。

## Relationships
- LearningProgress 1:1 (user, topic, chapter)  
- QuizSession N:1 QuizQuestion (via attempts)  
- QuizAttempt N:1 QuizQuestion，N:1 (user,topic,chapter)

## Validation & Defaults
- status default not_started；read_duration/scroll_progress/last_position default 0；quiz_passed default false；timestamps default CURRENT_TIMESTAMP。  
- scroll_progress clamp 0-100；read_duration 只增不减；status 不得从 completed 回退（除非删除记录）。  
- quiz_score 与 correct_answers 对齐；multiple/code_correction 判分需考虑部分正确。

## Derived Metrics
- 主题进度 = completed_chapters / total_chapters ×100%。  
- 整体进度 = Σ(主题进度 × 权重)。  
- 学习天数/总时长：基于 first_visit_at/last_visit_at 及 read_duration 聚合。

