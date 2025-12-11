# 进度模块说明

本模块负责章节级学习进度的领域定义与校验。核心内容：

- `LearningProgress`：记录用户在特定主题/章节的状态、阅读时长、滚动进度、测验结果与时间戳。
- 状态枚举：`not_started`、`in_progress`、`completed`、`tested`，禁止回退。
- 仓储接口 `ProgressRepository`：提供创建/更新、按用户与主题查询的抽象，具体实现位于 `internal/infra/repository`。

设计要求：字段与数据库迁移 `011_learning_progress_quiz.sql` 对齐，时间与数值校验在仓储/服务层显式处理，注释保持中文。***

