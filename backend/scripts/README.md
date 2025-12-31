# 测验题库工具集

本目录包含用于管理和验证 Go 学习测验题库的实用工具脚本。

## 工具列表

### 1. generate_quiz_yaml.go - 题目模板生成器

**用途**: 为新章节生成 YAML 题库骨架文件，包含元数据和题目占位符。

**使用方法**:
```bash
go run generate_quiz_yaml.go <topic> <chapter> <question_count>
```

**示例**:
```bash
# 为 lexical_elements/keywords 章节生成包含 40 道题的模板
go run generate_quiz_yaml.go lexical_elements keywords 40
```

**输出**: 在 `backend/quiz_data/<topic>/<chapter>.yaml` 创建模板文件。

---

### 2. validate_quiz_yaml.go - 格式验证器

**用途**: 100% 自动化检查 YAML 题库格式规范，确保无语法错误。

**使用方法**:
```bash
# 验证单个文件
go run validate_quiz_yaml.go <yaml_file_path>

# 验证整个目录
go run validate_quiz_yaml.go <directory_path>
```

**示例**:
```bash
# 验证单个章节
go run validate_quiz_yaml.go ../quiz_data/lexical_elements/keywords.yaml

# 验证所有 lexical_elements 章节
go run validate_quiz_yaml.go ../quiz_data/lexical_elements/
```

**检查项**:
- ✅ YAML 语法正确性
- ✅ 题目 ID 格式 (如 `keywords-001`)
- ✅ 题目类型枚举 (`single_choice`, `multiple_choice`, `code_output`, `code_fix`)
- ✅ 难度枚举 (`easy`, `medium`, `hard`)
- ✅ 选项数量 (单选4个, 多选4-6个)
- ✅ 答案格式正确性

**输出**: 验证通过或详细错误报告（行号 + 错误类型）。

---

### 3. quiz_quality_check.go - 内容质量检查器

**用途**: 80%+ 自动化内容质量验证，检查题目语义合理性和 Go 规范对齐。

**使用方法**:
```bash
go run quiz_quality_check.go <yaml_file_or_directory>
```

**示例**:
```bash
# 检查单个章节内容质量
go run quiz_quality_check.go ../quiz_data/lexical_elements/keywords.yaml

# 批量检查所有章节
go run quiz_quality_check.go ../quiz_data/
```

**检查项**:
- ✅ 解析 (explanation) 长度 (≥50 字符, 推荐 ≥100)
- ✅ Go 规范引用 (包含章节编号或关键词)
- ✅ 代码可执行性 (code_output/code_fix 题型)
- ✅ 重复题目检测 (题干相似度 >80%)
- ✅ 难度分布合理性 (easy 50-70%, medium 20-40%, hard 10-20%)

**输出**: 质量评分 (0-100) + 改进建议列表。

---

### 4. enhance_quiz_yaml.go - 题目增强工具

**用途**: 自动优化现有题目质量（补充解析、格式化代码、添加规范引用）。

**使用方法**:
```bash
go run enhance_quiz_yaml.go <yaml_file_path>
```

**示例**:
```bash
go run enhance_quiz_yaml.go ../quiz_data/lexical_elements/keywords.yaml
```

**功能**:
- 自动格式化代码片段 (使用 `gofmt`)
- 补充缺失的难度标签 (基于题目复杂度推断)
- 建议 Go 规范章节引用
- 标准化选项格式

**输出**: 在原文件同目录生成 `<chapter>_enhanced.yaml`。

---

### 5. check_quiz_progress.sh - 进度跟踪脚本

**用途**: 统计题库完成情况，显示各章节状态和总体进度。

**使用方法**:
```bash
./check_quiz_progress.sh
```

**输出示例**:
```
==== 题库完成进度报告 ====
总章节数: 41
已完成: 19 (46.3%)
待补全: 22 (53.7%)
总题目数: 875 / 目标 1200-2000

--- 按主题分类 ---
[lexical_elements] 9/9 完成 (100%)  ✅
[constants]        2/6 完成 (33.3%)  🔶
[variables]        0/5 完成 (0%)     ❌
[types]           8/21 完成 (38.1%)  🔶
```

---

### 6. stress_client.go - 性能压力测试

**用途**: 模拟并发用户访问测验 API，验证性能指标。

**使用方法**:
```bash
go run stress_client.go -url http://localhost:8080/api/v1/quiz -users 1000 -duration 60s
```

**参数**:
- `-url`: 测验 API 基础地址
- `-users`: 并发用户数 (默认 100)
- `-duration`: 测试时长 (默认 30s)

**输出**:
```
压力测试结果:
总请求数: 15234
成功率: 99.8%
平均响应时间: 45ms
P95 响应时间: 120ms
P99 响应时间: 280ms
```

---

## 工作流程示例

### 新章节创建完整流程

```bash
# 1. 生成模板
go run generate_quiz_yaml.go types slice 45

# 2. 手动编辑题目内容 (使用编辑器)
vim ../quiz_data/types/slice.yaml

# 3. 格式验证
go run validate_quiz_yaml.go ../quiz_data/types/slice.yaml

# 4. 内容质量检查
go run quiz_quality_check.go ../quiz_data/types/slice.yaml

# 5. (可选) 自动增强
go run enhance_quiz_yaml.go ../quiz_data/types/slice.yaml

# 6. 检查整体进度
./check_quiz_progress.sh
```

---

## 批量验证所有章节

```bash
# 验证所有已完成章节的格式
go run validate_quiz_yaml.go ../quiz_data/

# 检查所有章节的内容质量
go run quiz_quality_check.go ../quiz_data/

# 如果发现问题，针对性修复后重新验证
```

---

## 开发环境要求

- **Go 版本**: 1.21+ (支持泛型和最新标准库)
- **依赖包**:
  - `gopkg.in/yaml.v3` (YAML 解析)
  - `github.com/stretchr/testify` (测试框架)

---

## 注意事项

1. **编码格式**: 所有 YAML 文件使用 UTF-8 编码
2. **命名规范**: 章节文件名使用小写 + 下划线 (如 `type_assertions.yaml`)
3. **ID 格式**: 题目 ID 格式为 `<chapter>-<序号>` (如 `slice-001`, `slice-002`)
4. **备份**: 运行增强工具前建议先备份原文件
5. **并发限制**: 批量验证大量文件时注意内存使用 (建议单次处理 <100 个文件)

---

## 故障排查

**问题**: `validate_quiz_yaml.go` 报错 "无法解析 YAML"
- **解决**: 检查 YAML 缩进 (必须使用空格，不能用 Tab)，验证 YAML 语法工具: https://www.yamllint.com/

**问题**: `quiz_quality_check.go` 质量评分过低
- **解决**: 查看详细报告，重点优化：
  1. 补充解析内容 (≥100 字符)
  2. 添加 Go 规范引用 (如 "参见 Go 语言规范 §3.1")
  3. 确保代码示例可执行

**问题**: `check_quiz_progress.sh` 在 Windows 不可用
- **解决**: 使用 PowerShell 版本或 Git Bash 运行

---

## 贡献指南

新增工具脚本请遵循：
1. 使用 Go 编写 (便于跨平台)
2. 提供 `--help` 参数说明用法
3. 输出清晰的错误信息 (包含文件名和行号)
4. 更新本 README 文档
5. 添加单元测试 (coverage ≥80%)

---

最后更新: 2026-01-01
维护者: Go 学习平台开发团队
