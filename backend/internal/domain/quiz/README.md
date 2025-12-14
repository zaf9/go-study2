% quiz 包说明

此包负责题库的加载、验证、抽题与在内存中的索引。

目录: `backend/internal/domain/quiz`

主要职责：
- loader.go: 从 YAML 文件加载题库
- validator.go: 校验题目结构与业务规则
- selector.go: 根据配置执行随机抽题（难度分布、题型数量）
- repository.go: 内存索引与查询

开发说明：
- 所有日志与注释应使用中文
- 单元测试覆盖率目标 ≥ 80%
# 测验模块说明

本模块定义章节测验的领域模型：

- `QuizQuestion`：题库题目，包含题型、难度、选项、答案与解析。
- `QuizSession`：一次测验的结果摘要（得分、通过与起止时间）。
- `QuizAttempt`：单题作答记录，支持复查与解析。
- 仓储接口 `QuizRepository`：抽题、创建会话、保存作答、查询历史，具体实现位于 `internal/infra/repository`。

设计要求：题型/难度常量与合约一致，表结构与迁移 `011_learning_progress_quiz.sql` 对齐，注释保持中文，错误处理显式。***

