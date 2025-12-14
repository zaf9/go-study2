# Quickstart: 题库使用快速指南

**Branch**: `013-quiz-question-bank` | **Date**: 2025-12-14

## 概述

本指南帮助开发者和内容维护者快速了解测验题库系统的使用方法，包括如何添加/修改题目、触发重新加载、排查验证错误等常见操作。

## 题库文件组织

### 目录结构

```text
backend/quiz_data/
├── README.md                         # 题库文件说明
├── lexical_elements/                 # 词法元素主题（11章节）
│   ├── comments.yaml
│   ├── tokens.yaml
│   ├── semicolons.yaml
│   ├── identifiers.yaml
│   ├── keywords.yaml
│   ├── operators.yaml
│   ├── integers.yaml
│   ├── floats.yaml
│   ├── imaginary.yaml
│   ├── runes.yaml
│   └── strings.yaml
├── constants/                        # 常量主题（12章节）
│   ├── boolean.yaml
│   ├── rune.yaml
│   ├── integer.yaml
│   ├── floating_point.yaml
│   ├── complex.yaml
│   ├── string.yaml
│   ├── expressions.yaml
│   ├── typed_untyped.yaml
│   ├── conversions.yaml
│   ├── builtin_functions.yaml
│   ├── iota.yaml
│   └── implementation_restrictions.yaml
├── variables/                        # 变量主题（4章节）
│   ├── storage.yaml
│   ├── static.yaml
│   ├── dynamic.yaml
│   └── zero.yaml
└── types/                            # 类型主题（14章节）
    ├── boolean.yaml
    ├── numeric.yaml
    ├── string.yaml
    ├── array.yaml
    ├── slice.yaml
    ├── struct.yaml
    ├── pointer.yaml
    ├── function.yaml
    ├── interface_basic.yaml
    ├── interface_embedded.yaml
    ├── interface_general.yaml
    ├── interface_impl.yaml
    ├── map.yaml
    └── channel.yaml
```

### 文件命名规则

- 文件名：`{chapter}.yaml`（小写，下划线分隔）
- 主题目录：`{topic}/`（与代码中的主题标识一致）
- 示例：`constants/boolean.yaml`, `types/interface_basic.yaml`

## 快速开始：添加题目

### 步骤1：选择目标文件

根据题目所属主题和章节，找到对应的YAML文件。

**示例**：为"常量-布尔常量"章节添加题目
```bash
# 打开文件
backend/quiz_data/constants/boolean.yaml
```

### 步骤2：编写题目

按照YAML格式添加题目到`questions`数组末尾。

**单选题模板**：
```yaml
  - id: const-boolean-XXX           # 替换XXX为递增序号（如015）
    type: single
    difficulty: easy                # easy/medium/hard
    stem: "题干内容？"
    options:
      - "A: 选项A内容"
      - "B: 选项B内容"
      - "C: 选项C内容"
      - "D: 选项D内容"
    answer: "A"                     # 单个字母
    explanation: "答案解析：说明为何A正确，B/C/D为何错误。"
    topic: "constants"              # 与目录名一致
    chapter: "boolean"              # 与文件名一致
```

**多选题模板**：
```yaml
  - id: const-boolean-XXX
    type: multiple
    difficulty: medium
    stem: "题干内容？（多选）"      # 注意标注"（多选）"
    options:
      - "A: 选项A内容"
      - "B: 选项B内容"
      - "C: 选项C内容"
      - "D: 选项D内容"
    answer: "ACD"                   # 2-4个字母，升序排列
    explanation: "A正确因为...。B错误因为...。C正确因为...。D正确因为...。"
    topic: "constants"
    chapter: "boolean"
```

### 步骤3：验证格式

使用以下检查清单确保题目格式正确：

- [ ] ID符合格式：`{prefix}-{chapter}-{seq}`（如`const-boolean-015`）
- [ ] type为`single`或`multiple`
- [ ] difficulty为`easy`/`medium`/`hard`
- [ ] stem非空，10-500字符
- [ ] options数量符合要求（单选2-4，多选3-5）
- [ ] options格式为`"X: ..."`（X为A-E）
- [ ] answer格式正确（单选1字母，多选2-4字母升序）
- [ ] explanation非空，20-1000字符，中文
- [ ] topic和chapter与文件路径一致

### 步骤4：保存并重启服务

```bash
# 保存YAML文件
# Ctrl+S 或 :w（vim）

# 重启后端服务（题库在启动时加载）
cd backend
go run main.go
```

**重要提示**：题库更新需要重启服务才能生效，暂不支持热加载。

## 修改现有题目

### 步骤1：定位题目

在对应的YAML文件中搜索题目ID。

```bash
# 示例：查找ID为const-boolean-005的题目
grep "const-boolean-005" backend/quiz_data/constants/boolean.yaml
```

### 步骤2：修改内容

直接编辑YAML文件中的题目字段。

**常见修改场景**：
- 修正题干错别字：编辑`stem`字段
- 修改选项：编辑`options`数组
- 更正答案：编辑`answer`字段
- 完善解析：编辑`explanation`字段
- 调整难度：修改`difficulty`为`easy`/`medium`/`hard`

### 步骤3：验证并重启

参考"添加题目"的步骤3和步骤4。

## 删除题目

### 方法1：直接删除（推荐）

从YAML文件中删除整个题目对象（包括`- id: ...`到下一个`- id: ...`之间的所有行）。

### 方法2：注释（临时禁用）

在题目前添加`#`注释整个题目块（不推荐，YAML注释多行较麻烦）。

**注意**：删除题目后，ID不会自动重排，保持原有ID不变以避免历史数据混乱。

## 批量添加题目

### 场景：为新章节添加30-50题

1. **准备题目内容**：
   - 按照research.md中的知识点提纲准备题目
   - 使用AI辅助生成（参考research.md的Prompt模板）
   - 人工审核生成结果

2. **创建YAML文件**：
   ```yaml
   # backend/quiz_data/{topic}/{chapter}.yaml
   questions:
     - id: {prefix}-{chapter}-001
       # ... 第1题
     - id: {prefix}-{chapter}-002
       # ... 第2题
     # ... 共30-50题
   ```

3. **验证题目分布**：
   - 单选题：约50%（15-25题）
   - 多选题：约50%（15-25题）
   - 难度分布：easy 40%, medium 40%, hard 20%

4. **自动验证**：
   启动服务时会自动验证所有题目，错误会在日志中输出。

## 触发重新加载

### 当前方案：重启服务

```bash
# Linux/Mac
cd backend
./main   # 假设已编译

# 或使用go run
go run main.go

# Windows
cd backend
main.exe

# 或使用go run
go run main.go
```

### 未来可能支持（未实现）

- 热加载：修改YAML文件后自动重载
- 管理API：通过HTTP接口触发重载

## 验证错误排查

### 启动时验证失败

**症状**：服务启动失败，日志输出题库验证错误

**常见错误及解决方法**：

#### 错误1：文件不存在

```text
错误: 题库文件不存在: backend/quiz_data/constants/boolean.yaml
```

**解决**：
- 检查文件路径是否正确
- 确认文件已创建且命名正确
- 检查配置文件中的`quiz.dataPath`设置

#### 错误2：YAML语法错误

```text
错误: 解析YAML失败: yaml: line 15: mapping values are not allowed in this context
```

**解决**：
- 检查第15行附近的缩进（YAML严格要求2空格缩进）
- 检查引号是否配对
- 检查冒号后是否有空格
- 使用YAML在线验证工具检查格式

#### 错误3：必填字段缺失

```text
错误: 题目ID为空 (文件: constants/boolean.yaml, 索引: 5)
```

**解决**：
- 检查第6个题目（索引5，从0开始）是否缺少`id`字段
- 确保所有9个字段都存在

#### 错误4：枚举值错误

```text
错误: 无效的题型: Single, 仅支持 single/multiple (题目ID: const-boolean-010)
```

**解决**：
- 修改`type: "Single"`为`type: "single"`（小写）
- 检查`difficulty`是否为`easy`/`medium`/`hard`（小写）

#### 错误5：选项数量不符

```text
错误: 单选题选项数量必须在2-4之间，当前: 5 (题目ID: const-boolean-012)
```

**解决**：
- 单选题保留2-4个选项
- 多选题保留3-5个选项
- 删除多余选项或将题目改为多选题

#### 错误6：答案格式错误

```text
错误: 单选题答案必须为单个字母，当前: AB (题目ID: const-boolean-008)
```

**解决**：
- 单选题`answer`改为单个字母：`answer: "A"`
- 或将题目改为多选题：`type: "multiple"`

#### 错误7：答案不在选项范围内

```text
错误: 答案 E 不在选项范围内 (题目ID: const-boolean-020)
```

**解决**：
- 检查选项列表，确认有选项E
- 或修改答案为有效的选项字母（A-D）

#### 错误8：主题/章节不一致

```text
错误: 题目主题 lexical_elements 与文件路径 constants 不一致 (题目ID: const-boolean-003)
```

**解决**：
- 修改题目的`topic`字段为`"constants"`
- 或将题目移动到正确的主题目录

#### 错误9：ID重复

```text
错误: 题目ID重复: const-boolean-015 (索引 14 和 20)
```

**解决**：
- 将索引20的题目ID修改为未使用的序号（如`const-boolean-021`）

### 调试技巧

1. **查看完整错误日志**：
   ```bash
   go run main.go 2>&1 | tee startup.log
   ```

2. **逐个文件验证**：
   临时注释配置中的其他主题，仅加载一个主题的题库。

3. **使用YAML Lint工具**：
   ```bash
   # 在线工具
   https://www.yamllint.com/

   # 命令行工具
   yamllint backend/quiz_data/constants/boolean.yaml
   ```

4. **检查文件编码**：
   确保文件使用UTF-8编码（中文解析需要）。

## 配置文件说明

### 题库配置项

编辑`backend/configs/config.yaml`：

```yaml
quiz:
  # 题库文件根目录（相对于backend目录）
  dataPath: "quiz_data"
  
  # 每次测验抽取的题目数量
  questionCount:
    single: 4      # 单选题数量
    multiple: 4    # 多选题数量
  
  # 难度分布百分比（总和应为100）
  difficultyDistribution:
    easy: 40       # 简单题占比40%
    medium: 40     # 中等题占比40%
    hard: 20       # 困难题占比20%
  
  # 题库加载超时时间
  loadTimeout: 5s
  
  # 验证配置
  validation:
    strictMode: true   # 严格验证模式，建议保持true
    failFast: true     # 遇到首个错误即停止，建议保持true
```

### 修改抽题数量

**场景**：将每次测验改为抽取5道单选题和3道多选题

```yaml
quiz:
  questionCount:
    single: 5      # 修改为5
    multiple: 3    # 修改为3
```

保存后重启服务即可生效。

### 修改难度分布

**场景**：提高困难题比例到30%

```yaml
quiz:
  difficultyDistribution:
    easy: 35       # 35%
    medium: 35     # 35%
    hard: 30       # 30%（总和=100）
```

## 题库统计

### 查看题库概况

启动服务时，日志会输出题库加载统计：

```text
[INFO] 题库加载成功: 41个章节, 共2050题
[INFO] - lexical_elements: 11章节, 550题
[INFO] - constants: 12章节, 600题
[INFO] - variables: 4章节, 200题
[INFO] - types: 14章节, 700题
```

### 查看单个章节统计

**方法1：通过API查询**（如果实现了stats接口）

```bash
curl http://localhost:8080/api/v1/quiz/constants/boolean/stats
```

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 35,
    "byType": {
      "single": 18,
      "multiple": 17
    },
    "byDifficulty": {
      "easy": 14,
      "medium": 14,
      "hard": 7
    }
  }
}
```

**方法2：手动统计YAML文件**

```bash
# 统计题目总数
grep "^  - id:" backend/quiz_data/constants/boolean.yaml | wc -l

# 统计单选题数量
grep "type: single" backend/quiz_data/constants/boolean.yaml | wc -l

# 统计简单题数量
grep "difficulty: easy" backend/quiz_data/constants/boolean.yaml | wc -l
```

## 最佳实践

### 题目编写建议

1. **题干清晰**：
   - ✅ "Go语言中，iota的值在什么情况下会重置为0？"
   - ❌ "iota什么时候重置？"

2. **选项互斥**：
   - ✅ 选项间无重叠，彼此独立
   - ❌ "A: 大于0", "B: 大于等于1"（重叠）

3. **解析详细**：
   - ✅ 说明正确答案的原因 + 错误选项为何错误
   - ❌ "答案是A。"

4. **难度准确**：
   - easy：直接定义，无需推理
   - medium：需要理解和应用
   - hard：综合运用，边界情况

### 题目分布建议

每章节30-50题，建议分布：

| 题型 | 数量 | 难度分布 |
|------|------|----------|
| 单选题 | 15-25题 | easy:6-10, medium:6-10, hard:3-5 |
| 多选题 | 15-25题 | easy:6-10, medium:6-10, hard:3-5 |

### 版本控制建议

1. **提交前验证**：
   - 启动服务确认无验证错误
   - 检查题目质量

2. **提交信息格式**：
   ```bash
   git commit -m "feat(quiz): 新增constants/boolean章节15个题目

   - 单选题8题（easy:3, medium:3, hard:2）
   - 多选题7题（easy:3, medium:3, hard:1）
   - 覆盖知识点：布尔常量定义、运算、类型推断
   "
   ```

3. **批量修改**：
   - 使用分支开发，确认无误后合并
   - 大批量修改建议分章节提交

## 常见问题FAQ

### Q1: 如何知道题目ID的下一个序号？

**A**: 查看该章节YAML文件中最后一个题目的ID，序号+1。

```bash
# 示例：查看constants/boolean.yaml最后一题
tail -20 backend/quiz_data/constants/boolean.yaml | grep "id:"
# 输出: id: const-boolean-030
# 下一个题目使用: const-boolean-031
```

### Q2: 能否修改已有题目的ID？

**A**: 不建议。如果修改ID，可能影响已有的测验记录和统计数据。如需修改，确保没有依赖该ID的数据。

### Q3: 多选题答案字母顺序有要求吗？

**A**: 是的，必须升序排列。`"ACD"`正确，`"CAD"`错误。这是验证规则的要求，便于标准化比较。

### Q4: 如何快速生成大量题目？

**A**: 参考`research.md`中的AI辅助生成策略：
1. 人工编写5-8个核心题目
2. 使用AI生成变体（附带Prompt模板）
3. 人工审核和修正
4. 统一格式化

### Q5: 服务启动后如何测试抽题功能？

**A**: 调用开始测验API：

```bash
# 开始测验（获取随机题目）
curl http://localhost:8080/api/v1/quiz/constants/boolean/start

# 响应示例
{
  "code": 0,
  "message": "success",
  "data": {
    "topic": "constants",
    "chapter": "boolean",
    "questions": [...],  # 8个随机题目
    "total": 35,
    "selected": 8
  }
}
```

### Q6: 如何临时禁用某个章节的题库？

**A**: 重命名YAML文件（如`boolean.yaml.bak`），重启服务后该章节题库不会加载。

### Q7: 题库文件可以放在其他目录吗？

**A**: 可以，修改`config.yaml`中的`quiz.dataPath`为新路径（相对或绝对路径）。

### Q8: 如何检查题目解析是否包含中文？

**A**: 启动时验证器会自动检查，如果解析不包含中文会报错。

---

## 下一步

- **开始添加题目**：参考本指南的"快速开始"章节
- **查看数据模型**：阅读`data-model.md`了解详细结构
- **查看格式规范**：阅读`contracts/yaml-schema.md`了解字段约束
- **查看研究文档**：阅读`research.md`了解题目生成策略

---

**Quickstart Status**: ✅ Complete  
**Last Updated**: 2025-12-14
