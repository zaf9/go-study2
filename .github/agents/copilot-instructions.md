# go-study2 Development Guidelines

Auto-generated from all feature plans. Last updated: 2025-12-14

## Active Technologies
- Go 1.21+ (013-quiz-question-bank)
- gopkg.in/yaml.v3 (YAML解析库，用于题库文件加载) (013-quiz-question-bank)
- YAML外部文件存储（quiz_data/{topic}/{chapter}.yaml），启动时加载到内存 (013-quiz-question-bank)
- 题库管理：随机抽题、难度分级（easy/medium/hard）、题型分类（single/multiple） (013-quiz-question-bank)
- 深度验证 + Fail-Fast启动策略（题库结构完整性检查） (013-quiz-question-bank)
- Go 1.24.5 + github.com/gogf/gf/v2 (GoFrame 框架) (005-https-protocol-support)

## Project Structure

```text
src/
tests/
```

## Commands

# Add commands for Go 1.24.5

## Code Style

Go 1.24.5: Follow standard conventions

## Recent Changes
- 013-quiz-question-bank: Added Go 1.21+

- 005-https-protocol-support: Added Go 1.24.5 + github.com/gogf/gf/v2 (GoFrame 框架)

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
