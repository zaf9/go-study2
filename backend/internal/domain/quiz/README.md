# 测验模块说明

本模块定义章节测验的领域模型：

- `QuizQuestion`：题库题目，包含题型、难度、选项、答案与解析。
- `QuizSession`：一次测验的结果摘要（得分、通过与起止时间）。
- `QuizAttempt`：单题作答记录，支持复查与解析。
- 仓储接口 `QuizRepository`：抽题、创建会话、保存作答、查询历史，具体实现位于 `internal/infra/repository`。

设计要求：题型/难度常量与合约一致，表结构与迁移 `011_learning_progress_quiz.sql` 对齐，注释保持中文，错误处理显式。***

