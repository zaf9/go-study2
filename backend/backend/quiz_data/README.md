# 题库数据目录 (占位生成)

本目录由脚本 backend/scripts/generate_quiz_yaml.go 生成，包含按主题分组的章节 YAML 文件。

注意：此脚本生成的是占位题目，供人工审核和替换为高质量题目。生成的题目示例中包含中文占位题干和解析，请在发布前进行人工校验并完善题目内容。

生成规则：
- 每个章节文件为 YAML 格式，根节点为 'questions'。
- 每题包含字段：id, type, difficulty, stem, options, answer, explanation, topic, chapter
- 请参照 specs/013-quiz-question-bank/data-model.md 中的数据模型与约束进行最终校验。

生成命令：

    go run backend/scripts/generate_quiz_yaml.go

