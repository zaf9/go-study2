# Research: 补全题库内容实现完整的学习测验系统

## Phase 0: Research & Technical Decisions

### 1. Go 1.24 语言规范章节结构分析

**Decision**: 基于 Go 1.24 官方语言规范组织题库章节结构

**Rationale**: 
- Go 官方规范是最权威的知识来源,确保题目内容准确性
- 规范已经提供了完整的章节层次结构,可直接映射到题库组织
- 保证学习路径与 Go 语言官方文档一致,降低学习曲线

**章节映射 (41个章节)**:

#### lexical_elements (词法元素 - 9个子章节)
1. comments (注释) - **已完成, 35题**
2. keywords (关键字)
3. identifiers (标识符)
4. tokens (令牌)
5. semicolons (分号)
6. integers (整数字面量)
7. floats (浮点数字面量)
8. imaginary (虚数字面量)
9. runes (rune字面量)
10. strings (字符串字面量)

#### constants (常量 - 6个子章节)
11. boolean (布尔常量)
12. rune (rune常量)
13. integer (整型常量)
14. floating_point (浮点常量)
15. complex (复数常量)
16. string (字符串常量)

#### variables (变量 - 5个子章节)
17. variable_declarations (变量声明)
18. short_declarations (短变量声明)
19. blank_identifier (空白标识符)
20. type_conversions (类型转换)
21. zero_value (零值)

#### types (类型 - 20个子章节)
22. boolean (布尔类型)
23. numeric (数值类型)
24. string (字符串类型)
25. array (数组类型)
26. slice (切片类型)
27. struct (结构体类型)
28. pointer (指针类型)
29. function (函数类型)
30. interface (接口类型)
31. map (map类型)
32. channel (channel类型)
33. type_definitions (类型定义)
34. type_aliases (类型别名)
35. type_parameters (类型参数)
36. type_constraints (类型约束)
37. type_inference (类型推断)
38. type_unification (类型统一)
39. underlying_types (底层类型)
40. type_identity (类型同一性)
41. method_sets (方法集)

**Alternatives Considered**:
- ❌ 自定义章节结构: 可能与官方文档脱节,增加学习者心智负担
- ❌ 仅按难度分级: 缺乏系统性,不符合规范学习路径
- ✅ **官方规范章节映射**: 权威、系统、易于维护

---

### 2. 题目生成策略与质量保证

**Decision**: 采用"规范引用+代码实践+场景应用"三层题目模型

**Rationale**:
- Go 规范文档提供了丰富的示例代码和边界情况说明
- 实践表明单纯的概念题无法有效检验学习者的实际掌握程度
- 代码输出预测题和错误修正题能显著提升测验的实用价值

**题目分类与比例** (每章节30-50题):

| 题型 | 比例 | 示例 | 难度分布 |
|------|------|------|----------|
| 单选题 | 40% | "以下哪种注释方式用于单行注释?" | 简单60% 中等30% 困难10% |
| 多选题 | 30% | "关于多行注释,正确的说法是?(多选)" | 简单40% 中等40% 困难20% |
| 代码输出题 | 20% | "下列代码输出是什么?" | 中等50% 困难50% |
| 错误修正题 | 10% | "指出代码中的错误并修正" | 中等40% 困难60% |

**质量保证机制**:

1. **100% 格式验证 (自动化)**:
   ```yaml
   validation:
     - id: 必须唯一且符合 {chapter}-{topic}-{序号} 格式
     - type: 必须是 [single, multiple, code_output, code_fix] 之一
     - difficulty: 必须是 [easy, medium, hard] 之一
     - stem: 非空字符串,长度 20-500 字符
     - options: 单选4个选项,多选4-6个选项
     - answer: 单选单字母,多选字母组合(如 "AC")
     - explanation: 非空字符串,必须引用 Go 规范章节
     - topic: 必须匹配已定义的 topic 列表
     - chapter: 必须匹配已定义的 chapter 列表
   ```

2. **80%+ 内容验证 (自动化+人工抽查)**:
   ```yaml
   automated_checks:
     - 引用检查: explanation 必须包含"Go 规范"/"官方文档"等关键词
     - 代码验证: code_output 题型的代码必须能通过 go run 执行
     - 答案一致性: 正确答案必须在 options 中存在
     - 难度校准: easy 题目 explanation < 100字, hard 题目 explanation > 200字
   
   manual_spot_check:
     - 每章节随机抽查20%题目验证内容准确性
     - 优先抽查 medium/hard 难度题目
     - 检查 explanation 是否引用正确的规范章节
   ```

3. **去重与多样性检查**:
   ```bash
   # 使用脚本检查题目相似度
   check-quiz-diversity:
     - 同一章节题目的 stem 不允许完全重复
     - 代码题的代码片段至少有30%差异
     - explanation 不允许复制粘贴相同内容
   ```

**Alternatives Considered**:
- ❌ 纯人工审核: 1200-2000题目工作量巨大,不可持续
- ❌ AI 生成无校验: 题目质量无法保证,可能出现错误知识点
- ✅ **自动化验证 + 抽查**: 平衡效率和质量,符合 80%+ 自动化要求

---

### 3. 存储方案与性能优化

**Decision**: 继续使用 YAML 文件存储,增加索引和缓存机制

**Rationale**:
- 现有系统已采用 YAML 存储 (comments.yaml 验证可行)
- YAML 可读性强,便于人工审核和版本控制
- 题库量级 (1200-2000题) 不足以构成性能瓶颈

**存储结构**:
```yaml
# backend/quiz_data/{topic}/{chapter}.yaml
metadata:
  topic: lexical_elements
  chapter: keywords
  total_questions: 45
  last_updated: "2026-01-15"
  difficulty_distribution:
    easy: 27    # 60%
    medium: 14  # 31%
    hard: 4     # 9%

questions:
  - id: lexical-keywords-001
    type: single
    difficulty: easy
    stem: "以下哪个不是 Go 语言的关键字?"
    options:
      - "A: func"
      - "B: class"
      - "C: interface"
      - "D: package"
    answer: B
    explanation: "Go 语言没有 class 关键字,使用 struct 和 interface 实现面向对象。参见 Go 规范 Keywords 章节。"
    topic: lexical_elements
    chapter: keywords
    tags: ["keywords", "basic"]
```

**性能优化策略**:

1. **索引文件生成** (构建时):
   ```yaml
   # backend/quiz_data/index.yaml
   topics:
     lexical_elements:
       chapters:
         keywords:
           file: lexical_elements/keywords.yaml
           count: 45
           difficulties: {easy: 27, medium: 14, hard: 4}
         identifiers:
           file: lexical_elements/identifiers.yaml
           count: 40
           difficulties: {easy: 24, medium: 12, hard: 4}
   ```

2. **缓存机制** (运行时):
   ```go
   // 使用 sync.Map 缓存已加载的章节题库
   var quizCache sync.Map

   func LoadChapterQuiz(topic, chapter string) (*Quiz, error) {
       cacheKey := fmt.Sprintf("%s/%s", topic, chapter)
       if cached, ok := quizCache.Load(cacheKey); ok {
           return cached.(*Quiz), nil
       }
       
       quiz, err := loadFromYAML(topic, chapter)
       if err != nil {
           return nil, err
       }
       quizCache.Store(cacheKey, quiz)
       return quiz, nil
   }
   ```

3. **随机题目抽取算法**:
   ```go
   // 使用 Fisher-Yates 洗牌算法确保随机性
   func SelectRandomQuestions(questions []Question, count int, singleRatio float64) []Question {
       // 1. 按类型分组
       singles, multiples := groupByType(questions)
       
       // 2. 计算各类型数量
       singleCount := int(float64(count) * singleRatio)
       multipleCount := count - singleCount
       
       // 3. 随机抽取
       selectedSingles := shuffle(singles)[:singleCount]
       selectedMultiples := shuffle(multiples)[:multipleCount]
       
       // 4. 合并并打乱顺序
       return shuffle(append(selectedSingles, selectedMultiples...))
   }
   ```

**Alternatives Considered**:
- ❌ SQLite 数据库: 增加部署复杂度,现阶段数据量不需要数据库
- ❌ JSON 存储: 可读性不如 YAML,且现有系统已使用 YAML
- ✅ **YAML + 缓存**: 简单可靠,符合项目 "YAGNI" 原则

---

### 4. 题目生成工具链

**Decision**: 开发半自动化题目生成和校验工具

**Rationale**:
- 纯手工编写 1200+ 题目耗时过长
- 完全自动化生成质量难以保证
- 半自动化工具可以提供模板和校验,提高效率

**工具链组成**:

1. **题目模板生成器** (`scripts/generate_quiz_template.go`):
   ```go
   // 输入: 章节名称, 题目数量
   // 输出: 包含元数据的 YAML 模板
   func GenerateTemplate(topic, chapter string, count int) error {
       template := QuizTemplate{
           Metadata: Metadata{
               Topic: topic,
               Chapter: chapter,
               TotalQuestions: count,
               LastUpdated: time.Now().Format("2006-01-02"),
           },
           Questions: make([]QuestionTemplate, count),
       }
       
       // 自动分配题型和难度
       for i := range template.Questions {
           template.Questions[i] = assignTypeAndDifficulty(i, count)
       }
       
       return saveYAML(template, fmt.Sprintf("quiz_data/%s/%s.yaml", topic, chapter))
   }
   ```

2. **格式验证器** (`scripts/validate_quiz_yaml.go`):
   ```go
   func ValidateQuiz(filePath string) []ValidationError {
       quiz, err := loadYAML(filePath)
       if err != nil {
           return []ValidationError{{Type: "parse_error", Message: err.Error()}}
       }
       
       var errors []ValidationError
       
       // 验证每个题目
       for i, q := range quiz.Questions {
           errors = append(errors, validateQuestion(q, i)...)
       }
       
       // 验证难度分布
       if dist := calculateDifficultyDistribution(quiz); !isValidDistribution(dist) {
           errors = append(errors, ValidationError{
               Type: "distribution_error",
               Message: fmt.Sprintf("难度分布不符合要求: %+v", dist),
           })
       }
       
       return errors
   }
   ```

3. **批量质量检查工具** (`scripts/quiz_quality_check.go`):
   ```go
   func RunQualityChecks(dir string) QualityReport {
       var report QualityReport
       
       // 遍历所有 YAML 文件
       filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
           if !strings.HasSuffix(path, ".yaml") {
               return nil
           }
           
           // 格式检查
           report.FormatErrors = append(report.FormatErrors, ValidateQuiz(path)...)
           
           // 内容检查
           report.ContentWarnings = append(report.ContentWarnings, checkContent(path)...)
           
           // 去重检查
           report.Duplicates = append(report.Duplicates, checkDuplicates(path)...)
           
           return nil
       })
       
       return report
   }
   ```

4. **进度跟踪脚本** (`scripts/check_quiz_progress.sh`):
   ```bash
   #!/bin/bash
   # 统计各章节题目完成情况
   
   echo "题库完成进度:"
   echo "===================="
   
   total_chapters=41
   completed=0
   total_questions=0
   
   for yaml in backend/quiz_data/**/*.yaml; do
       if [ -f "$yaml" ]; then
           count=$(grep -c "^  - id:" "$yaml")
           chapter=$(basename "$yaml" .yaml)
           echo "$chapter: $count 题"
           completed=$((completed + 1))
           total_questions=$((total_questions + count))
       fi
   done
   
   echo "===================="
   echo "已完成章节: $completed / $total_chapters"
   echo "已编写题目: $total_questions / 1200 (目标)"
   echo "完成率: $(echo "scale=2; $completed * 100 / $total_chapters" | bc)%"
   ```

**Alternatives Considered**:
- ❌ 完全手工编写: 效率低,易出错
- ❌ GPT 批量生成: 质量不可控,需大量人工审核
- ✅ **半自动化工具链**: 平衡效率和质量,提供可复用的验证机制

---

### 5. 测验系统实现方案

**Decision**: 复用现有 quiz 模块,扩展支持多章节和历史记录

**Rationale**:
- 现有代码已实现基础的题目加载和测验逻辑
- 避免重复造轮子,符合 DRY 原则
- 扩展现有模块比重写更可控

**核心模块扩展**:

1. **QuizService 接口扩展**:
   ```go
   // backend/internal/domain/quiz.go
   type QuizService interface {
       // 现有方法
       LoadChapterQuiz(topic, chapter string) (*Quiz, error)
       SelectRandomQuestions(quiz *Quiz, count int) []Question
       
       // 新增方法
       GetAllChapters() []ChapterInfo
       RecordQuizAttempt(userID string, attempt QuizAttempt) error
       GetUserHistory(userID string, topic, chapter string) ([]QuizAttempt, error)
       GetUserProgress(userID string) ([]ChapterProgress, error)
   }
   
   type ChapterInfo struct {
       Topic           string `json:"topic"`
       Chapter         string `json:"chapter"`
       TotalQuestions  int    `json:"total_questions"`
       AvailableForQuiz bool  `json:"available_for_quiz"`
   }
   
   type QuizAttempt struct {
       ID          string    `json:"id"`
       UserID      string    `json:"user_id"`
       Topic       string    `json:"topic"`
       Chapter     string    `json:"chapter"`
       Questions   []Question `json:"questions"`
       UserAnswers map[string]string `json:"user_answers"`
       Score       float64   `json:"score"`
       CreatedAt   time.Time `json:"created_at"`
   }
   
   type ChapterProgress struct {
       Topic        string  `json:"topic"`
       Chapter      string  `json:"chapter"`
       BestScore    float64 `json:"best_score"`
       AttemptCount int     `json:"attempt_count"`
       LastAttempt  time.Time `json:"last_attempt"`
   }
   ```

2. **存储层实现**:
   ```go
   // backend/internal/infra/quiz_repository.go
   type QuizRepository interface {
       SaveAttempt(attempt QuizAttempt) error
       GetAttemptsByUser(userID string) ([]QuizAttempt, error)
       GetAttemptsByChapter(userID, topic, chapter string) ([]QuizAttempt, error)
   }
   
   // 使用 JSON 文件存储历史记录 (暂不引入数据库)
   type FileQuizRepository struct {
       dataDir string
   }
   
   func (r *FileQuizRepository) SaveAttempt(attempt QuizAttempt) error {
       userDir := filepath.Join(r.dataDir, "attempts", attempt.UserID)
       os.MkdirAll(userDir, 0755)
       
       filename := fmt.Sprintf("%s_%s_%s.json", 
           attempt.Topic, attempt.Chapter, attempt.CreatedAt.Format("20060102_150405"))
       
       data, err := json.MarshalIndent(attempt, "", "  ")
       if err != nil {
           return err
       }
       
       return os.WriteFile(filepath.Join(userDir, filename), data, 0644)
   }
   ```

3. **HTTP API 路由设计**:
   ```go
   // backend/internal/interfaces/quiz_handler.go
   func RegisterQuizRoutes(r *gin.RouterGroup) {
       r.GET("/chapters", GetChaptersHandler)                  // 获取所有章节列表
       r.GET("/quiz/:topic/:chapter", StartQuizHandler)        // 开始测验
       r.POST("/quiz/:topic/:chapter", SubmitQuizHandler)      // 提交测验
       r.GET("/history/:topic/:chapter", GetHistoryHandler)    // 获取历史记录
       r.GET("/progress", GetProgressHandler)                  // 获取进度概览
   }
   
   // 示例响应
   // GET /api/v1/quiz/chapters
   {
       "code": 200,
       "message": "success",
       "data": [
           {
               "topic": "lexical_elements",
               "chapter": "comments",
               "total_questions": 35,
               "available_for_quiz": true
           },
           {
               "topic": "lexical_elements",
               "chapter": "keywords",
               "total_questions": 0,
               "available_for_quiz": false
           }
       ]
   }
   ```

**Alternatives Considered**:
- ❌ 重写整个测验系统: 浪费现有代码,增加风险
- ❌ 引入 PostgreSQL/MySQL: 过度设计,当前数据量用文件存储足够
- ✅ **扩展现有模块 + 文件存储**: 最小化变更,符合 YAGNI 原则

---

## 关键技术决策总结

| 决策点 | 选择方案 | 主要原因 |
|--------|----------|----------|
| 章节结构 | 基于 Go 1.24 官方规范 | 权威性、系统性、易维护 |
| 题目质量保证 | 100%格式验证 + 80%内容验证 | 平衡自动化效率和质量 |
| 存储方案 | YAML + 索引 + 缓存 | 复用现有方案,符合 YAGNI |
| 工具链 | 半自动化生成+验证工具 | 提高效率,保证质量 |
| 测验系统 | 扩展现有 quiz 模块 | 最小化变更,降低风险 |

---

## 待解决的技术问题

1. **题目内容生成效率**: 如何在保证质量的前提下加速 1200+ 题目的编写?
   - **建议**: 优先完成高频章节 (lexical_elements, types),允许分阶段上线

2. **题目难度校准**: 如何确保 easy/medium/hard 的区分度?
   - **建议**: 通过实际用户测试数据调整难度标签,初期可保守估计

3. **多语言支持预留**: 是否需要考虑未来的英文版本?
   - **建议**: 当前仅支持中文,如需扩展可通过 i18n 键值对重构

4. **性能压测**: 1000并发用户同时进行测验是否会有性能瓶颈?
   - **建议**: Phase 2 实施后进行压测,必要时引入 Redis 缓存

