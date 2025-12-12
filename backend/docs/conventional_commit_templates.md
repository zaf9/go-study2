# Conventional Commit 模板（示例）

以下为推荐的 Git 提交信息模板，团队应遵循 Conventional Commits 规范：

示例 - 功能添加:

feat(logger): 支持按环境加载日志配置

描述: 在 `backend/configs` 中添加了 `logger` 区段, 支持多实例配置( access/error/slow )。

示例 - 修复:

fix(middleware): 修复 access log 中状态码格式化错误

描述: 使用 strconv.Itoa 代替错误的 rune->string 转换，并改为使用配置化时间格式。

示例 - 文档:

docs(logging): 增加日志最佳实践与故障排查文档

描述: 新增 `backend/docs/logging_best_practices.md` 与 `logging_troubleshooting.md`。

使用指南：
- commit 类型：feat | fix | docs | chore | refactor | test | perf | ci
- 范围（可选）：括号内写模块名，如 logger、middleware
- 主题行不超过 72 个字符
