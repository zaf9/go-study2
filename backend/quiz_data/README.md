# 题库（quiz_data）说明

此目录存放用于测验的 YAML 题库文件，按主题（topic）和章节（chapter）组织。

目录结构示例：

- quiz_data/
  - constants/
    - boolean.yaml
    - integer.yaml
    - ...
  - lexical_elements/
    - comments.yaml
    - tokens.yaml
  - variables/
    - storage.yaml
  - types/
    - string.yaml

每个 YAML 文件代表一个章节题库，要求包含 30-50 道题目，题型包括单选（single）与多选（multiple）两类。

YAML字段规范（每题）:

- id: 唯一标识（字符串或数字）
- type: "single" 或 "multiple"
- difficulty: "easy" | "medium" | "hard"
- question: 题干，中文为主，简洁明确
- options: 列表，按字母顺序或任意顺序均可（系统会在返回时打乱）
- answer: 单选为单个字母，多选为多个字母字符串（例如 "AC"）
- explanation: 中文解析，解释正确答案和关键点
- source: 可选，参考资料或出处
- tags: 可选，关键字数组（用于检索或统计）

题目创建流程（简要）：

1. 在对应主题目录下创建或编辑 YAML 文件（确保遵循字段规范）。
2. 运行本地验证脚本（见下）检查格式与唯一性。
3. 提交 PR 并在 PR 描述中列出人工审核项（语言、答案、解析完整性）。
4. 审核通过后，合并并触发 CI，服务启动时会加载并校验题库。

本地验证（建议步骤）:

- 使用项目内的验证单元测试：

```powershell
cd backend
go test ./internal/domain/quiz -run TestValidateYAML -v
```

- 或运行自定义脚本（如存在）：

```powershell
cd backend
go run ./internal/cmd/quiz_validator main --path ../quiz_data
```

人工审核清单（参考 T020）:

- 题目数量是否在 30-50 之间？
- 单选/多选比例是否合理？
- 难度分布是否大致符合要求？
- 题干是否通顺，中文表述是否准确？
- 选项是否唯一、无重复；答案是否正确？
- 解析是否包含关键点且为中文说明？

维护指南：

- 新增题目请保持字段完整并在 PR 中说明来源与验证方法。
- 若需批量修复，可先在分支中运行验证脚本，修复所有错误后再提交。

如需更多帮助，请参阅 `specs/013-quiz-question-bank` 中的设计文档与任务说明。
