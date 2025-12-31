# Data Model: 补全题库内容实现完整的学习测验系统

## 核心实体

### 1. Question (题目)

```go
// 题目实体 - 表示单个测验题目
type Question struct {
    ID          string   `yaml:"id" json:"id"`                  // 唯一标识符,格式: {topic}-{chapter}-{序号}
    Type        QuestionType `yaml:"type" json:"type"`          // 题型: single, multiple, code_output, code_fix
    Difficulty  Difficulty `yaml:"difficulty" json:"difficulty"` // 难度: easy, medium, hard
    Stem        string   `yaml:"stem" json:"stem"`              // 题干(问题描述)
    Options     []string `yaml:"options" json:"options"`        // 选项列表
    Answer      string   `yaml:"answer" json:"answer"`          // 正确答案
    Explanation string   `yaml:"explanation" json:"explanation"` // 详细解析
    Topic       string   `yaml:"topic" json:"topic"`            // 所属主题
    Chapter     string   `yaml:"chapter" json:"chapter"`        // 所属章节
    Tags        []string `yaml:"tags,omitempty" json:"tags,omitempty"` // 可选标签
}

// 题型枚举
type QuestionType string

const (
    QuestionTypeSingle      QuestionType = "single"       // 单选题
    QuestionTypeMultiple    QuestionType = "multiple"     // 多选题
    QuestionTypeCodeOutput  QuestionType = "code_output"  // 代码输出题
    QuestionTypeCodeFix     QuestionType = "code_fix"     // 错误修正题
)

// 难度枚举
type Difficulty string

const (
    DifficultyEasy   Difficulty = "easy"   // 简单
    DifficultyMedium Difficulty = "medium" // 中等
    DifficultyHard   Difficulty = "hard"   // 困难
)
```

**字段说明**:
- `ID`: 全局唯一,格式为 `{topic}-{chapter}-{3位数字}`,例如 `lexical-comments-001`
- `Type`: 决定题目的交互方式和评分逻辑
- `Stem`: 题干长度限制 20-500 字符
- `Options`: 单选题固定4个选项,多选题4-6个选项
- `Answer`: 单选题为单字母(如 "B"),多选题为字母组合(如 "AC")
- `Explanation`: 必须引用 Go 规范章节,长度 ≥50 字符

**验证规则**:
```go
func (q *Question) Validate() error {
    // ID 格式验证
    if !regexp.MustCompile(`^[a-z_]+-[a-z_]+-\d{3}$`).MatchString(q.ID) {
        return fmt.Errorf("invalid ID format: %s", q.ID)
    }
    
    // 题型验证
    validTypes := map[QuestionType]bool{
        QuestionTypeSingle: true, QuestionTypeMultiple: true,
        QuestionTypeCodeOutput: true, QuestionTypeCodeFix: true,
    }
    if !validTypes[q.Type] {
        return fmt.Errorf("invalid question type: %s", q.Type)
    }
    
    // 选项数量验证
    if q.Type == QuestionTypeSingle && len(q.Options) != 4 {
        return errors.New("single choice question must have exactly 4 options")
    }
    if q.Type == QuestionTypeMultiple && (len(q.Options) < 4 || len(q.Options) > 6) {
        return errors.New("multiple choice question must have 4-6 options")
    }
    
    // 答案验证
    if q.Answer == "" {
        return errors.New("answer cannot be empty")
    }
    
    // 解析验证 (必须引用 Go 规范)
    if !strings.Contains(q.Explanation, "Go 规范") && !strings.Contains(q.Explanation, "官方文档") {
        return errors.New("explanation must reference Go specification")
    }
    
    return nil
}
```

---

### 2. ChapterQuizBank (章节题库)

```go
// 章节题库 - 包含某个章节的所有题目
type ChapterQuizBank struct {
    Metadata  QuizMetadata `yaml:"metadata" json:"metadata"`    // 元数据
    Questions []Question   `yaml:"questions" json:"questions"`  // 题目列表
}

// 题库元数据
type QuizMetadata struct {
    Topic                   string             `yaml:"topic" json:"topic"`
    Chapter                 string             `yaml:"chapter" json:"chapter"`
    TotalQuestions          int                `yaml:"total_questions" json:"total_questions"`
    LastUpdated             string             `yaml:"last_updated" json:"last_updated"` // YYYY-MM-DD
    DifficultyDistribution  DifficultyDist     `yaml:"difficulty_distribution" json:"difficulty_distribution"`
    TypeDistribution        TypeDist           `yaml:"type_distribution,omitempty" json:"type_distribution,omitempty"`
}

// 难度分布
type DifficultyDist struct {
    Easy   int `yaml:"easy" json:"easy"`
    Medium int `yaml:"medium" json:"medium"`
    Hard   int `yaml:"hard" json:"hard"`
}

// 题型分布
type TypeDist struct {
    Single     int `yaml:"single" json:"single"`
    Multiple   int `yaml:"multiple" json:"multiple"`
    CodeOutput int `yaml:"code_output,omitempty" json:"code_output,omitempty"`
    CodeFix    int `yaml:"code_fix,omitempty" json:"code_fix,omitempty"`
}
```

**使用示例**:
```go
// 加载章节题库
func LoadChapterQuizBank(topic, chapter string) (*ChapterQuizBank, error) {
    filePath := filepath.Join("quiz_data", topic, fmt.Sprintf("%s.yaml", chapter))
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read quiz file: %w", err)
    }
    
    var bank ChapterQuizBank
    if err := yaml.Unmarshal(data, &bank); err != nil {
        return nil, fmt.Errorf("failed to parse quiz YAML: %w", err)
    }
    
    // 验证题库完整性
    if err := bank.Validate(); err != nil {
        return nil, fmt.Errorf("quiz bank validation failed: %w", err)
    }
    
    return &bank, nil
}

func (b *ChapterQuizBank) Validate() error {
    // 验证元数据一致性
    if len(b.Questions) != b.Metadata.TotalQuestions {
        return fmt.Errorf("question count mismatch: got %d, want %d", 
            len(b.Questions), b.Metadata.TotalQuestions)
    }
    
    // 验证难度分布
    actualDist := b.CalculateDifficultyDistribution()
    if actualDist != b.Metadata.DifficultyDistribution {
        return fmt.Errorf("difficulty distribution mismatch")
    }
    
    // 验证每个题目
    for i, q := range b.Questions {
        if err := q.Validate(); err != nil {
            return fmt.Errorf("question %d validation failed: %w", i+1, err)
        }
    }
    
    return nil
}

func (b *ChapterQuizBank) CalculateDifficultyDistribution() DifficultyDist {
    dist := DifficultyDist{}
    for _, q := range b.Questions {
        switch q.Difficulty {
        case DifficultyEasy:
            dist.Easy++
        case DifficultyMedium:
            dist.Medium++
        case DifficultyHard:
            dist.Hard++
        }
    }
    return dist
}
```

---

### 3. QuizAttempt (测验记录)

```go
// 测验记录 - 表示用户的一次测验尝试
type QuizAttempt struct {
    ID          string              `json:"id"`           // 唯一标识符
    UserID      string              `json:"user_id"`      // 用户ID
    Topic       string              `json:"topic"`        // 主题
    Chapter     string              `json:"chapter"`      // 章节
    Questions   []Question          `json:"questions"`    // 抽取的题目列表
    UserAnswers map[string]string   `json:"user_answers"` // 用户答案 {question_id: answer}
    Score       float64             `json:"score"`        // 分数 (0-100)
    CorrectCount int                `json:"correct_count"` // 答对题数
    TotalCount  int                 `json:"total_count"`  // 总题数
    CreatedAt   time.Time           `json:"created_at"`   // 创建时间
    CompletedAt *time.Time          `json:"completed_at,omitempty"` // 完成时间
    TimeSpent   int                 `json:"time_spent"`   // 耗时(秒)
}

// 题目结果
type QuestionResult struct {
    QuestionID  string `json:"question_id"`
    UserAnswer  string `json:"user_answer"`
    CorrectAnswer string `json:"correct_answer"`
    IsCorrect   bool   `json:"is_correct"`
    Explanation string `json:"explanation"`
}
```

**评分逻辑**:
```go
func (a *QuizAttempt) CalculateScore() {
    correct := 0
    for _, q := range a.Questions {
        userAnswer := a.UserAnswers[q.ID]
        if userAnswer == q.Answer {
            correct++
        }
    }
    
    a.CorrectCount = correct
    a.TotalCount = len(a.Questions)
    a.Score = float64(correct) / float64(a.TotalCount) * 100
}

func (a *QuizAttempt) GetResults() []QuestionResult {
    results := make([]QuestionResult, len(a.Questions))
    for i, q := range a.Questions {
        userAnswer := a.UserAnswers[q.ID]
        results[i] = QuestionResult{
            QuestionID:    q.ID,
            UserAnswer:    userAnswer,
            CorrectAnswer: q.Answer,
            IsCorrect:     userAnswer == q.Answer,
            Explanation:   q.Explanation,
        }
    }
    return results
}
```

**存储结构** (JSON 文件):
```json
// data/attempts/{user_id}/{topic}_{chapter}_{timestamp}.json
{
  "id": "attempt-20260115-143052-abc123",
  "user_id": "user123",
  "topic": "lexical_elements",
  "chapter": "comments",
  "questions": [...],
  "user_answers": {
    "lexical-comments-001": "B",
    "lexical-comments-005": "AC"
  },
  "score": 85.0,
  "correct_count": 8,
  "total_count": 10,
  "created_at": "2026-01-15T14:30:52Z",
  "completed_at": "2026-01-15T14:35:30Z",
  "time_spent": 278
}
```

---

### 4. ChapterProgress (章节进度)

```go
// 章节进度 - 表示用户在某个章节的学习进度
type ChapterProgress struct {
    UserID       string    `json:"user_id"`
    Topic        string    `json:"topic"`
    Chapter      string    `json:"chapter"`
    Status       LearningStatus `json:"status"`      // 学习状态
    BestScore    float64   `json:"best_score"`      // 最高分数
    AttemptCount int       `json:"attempt_count"`   // 测验次数
    LastAttempt  time.Time `json:"last_attempt"`    // 最后一次测验时间
    FirstAttempt time.Time `json:"first_attempt"`   // 首次测验时间
}

// 学习状态枚举
type LearningStatus string

const (
    StatusNotStarted LearningStatus = "not_started" // 未开始
    StatusInProgress LearningStatus = "in_progress" // 学习中
    StatusCompleted  LearningStatus = "completed"   // 已完成 (最高分 ≥ 60)
    StatusMastered   LearningStatus = "mastered"    // 精通 (最高分 ≥ 90)
)
```

**进度计算逻辑**:
```go
func CalculateProgress(userID string) ([]ChapterProgress, error) {
    attempts, err := loadAllAttempts(userID)
    if err != nil {
        return nil, err
    }
    
    // 按章节分组
    progressMap := make(map[string]*ChapterProgress)
    
    for _, attempt := range attempts {
        key := fmt.Sprintf("%s/%s", attempt.Topic, attempt.Chapter)
        
        if progress, exists := progressMap[key]; exists {
            // 更新已有进度
            progress.AttemptCount++
            if attempt.Score > progress.BestScore {
                progress.BestScore = attempt.Score
            }
            if attempt.CreatedAt.After(progress.LastAttempt) {
                progress.LastAttempt = attempt.CreatedAt
            }
        } else {
            // 创建新进度
            progressMap[key] = &ChapterProgress{
                UserID:       userID,
                Topic:        attempt.Topic,
                Chapter:      attempt.Chapter,
                BestScore:    attempt.Score,
                AttemptCount: 1,
                FirstAttempt: attempt.CreatedAt,
                LastAttempt:  attempt.CreatedAt,
            }
        }
    }
    
    // 更新状态
    for _, progress := range progressMap {
        progress.Status = determineStatus(progress.BestScore, progress.AttemptCount)
    }
    
    return convertMapToSlice(progressMap), nil
}

func determineStatus(bestScore float64, attemptCount int) LearningStatus {
    if attemptCount == 0 {
        return StatusNotStarted
    }
    if bestScore >= 90 {
        return StatusMastered
    }
    if bestScore >= 60 {
        return StatusCompleted
    }
    return StatusInProgress
}
```

---

## 数据关系图

```
┌─────────────────┐
│ ChapterQuizBank │ (YAML 文件存储)
├─────────────────┤
│ - metadata      │
│ - questions[]   │
└────────┬────────┘
         │
         │ 1:N
         │
         ▼
  ┌──────────────┐
  │  Question    │
  ├──────────────┤
  │ - id         │
  │ - type       │
  │ - difficulty │
  │ - stem       │
  │ - options[]  │
  │ - answer     │
  │ - explanation│
  └──────┬───────┘
         │
         │ N:M (随机抽取)
         │
         ▼
  ┌──────────────────┐
  │  QuizAttempt     │ (JSON 文件存储)
  ├──────────────────┤
  │ - id             │
  │ - user_id        │
  │ - topic          │
  │ - chapter        │
  │ - questions[]    │◄─────┐
  │ - user_answers{} │      │
  │ - score          │      │
  │ - created_at     │      │
  └────────┬─────────┘      │
           │                │
           │ N:1            │ 聚合计算
           │                │
           ▼                │
    ┌─────────────────┐    │
    │ChapterProgress  │────┘
    ├─────────────────┤
    │ - user_id       │
    │ - topic         │
    │ - chapter       │
    │ - status        │
    │ - best_score    │
    │ - attempt_count │
    └─────────────────┘
```

---

## 数据验证规范

### 1. 题库 YAML 文件验证

```yaml
# 必需字段验证
required_fields:
  metadata:
    - topic
    - chapter
    - total_questions
    - last_updated
    - difficulty_distribution
  question:
    - id
    - type
    - difficulty
    - stem
    - options
    - answer
    - explanation
    - topic
    - chapter

# 格式验证
format_validation:
  id: "^[a-z_]+-[a-z_]+-\\d{3}$"
  last_updated: "^\\d{4}-\\d{2}-\\d{2}$"
  type: ["single", "multiple", "code_output", "code_fix"]
  difficulty: ["easy", "medium", "hard"]

# 内容验证
content_validation:
  stem_length: [20, 500]
  explanation_min_length: 50
  single_options_count: 4
  multiple_options_count: [4, 6]
  answer_single_format: "^[A-Z]$"
  answer_multiple_format: "^[A-Z]{2,6}$"

# 分布验证
distribution_validation:
  difficulty:
    easy: [0.5, 0.7]    # 50%-70%
    medium: [0.25, 0.4] # 25%-40%
    hard: [0.05, 0.15]  # 5%-15%
  type:
    single: [0.35, 0.45]   # 35%-45%
    multiple: [0.25, 0.35] # 25%-35%
```

### 2. 测验记录 JSON 验证

```go
func (a *QuizAttempt) Validate() error {
    // 基础字段验证
    if a.UserID == "" {
        return errors.New("user_id is required")
    }
    if a.Topic == "" || a.Chapter == "" {
        return errors.New("topic and chapter are required")
    }
    
    // 题目数量验证
    if len(a.Questions) < 6 || len(a.Questions) > 10 {
        return errors.New("quiz must have 6-10 questions")
    }
    
    // 答案完整性验证
    for _, q := range a.Questions {
        if _, exists := a.UserAnswers[q.ID]; !exists {
            return fmt.Errorf("missing answer for question: %s", q.ID)
        }
    }
    
    // 分数验证
    if a.Score < 0 || a.Score > 100 {
        return errors.New("score must be between 0 and 100")
    }
    
    return nil
}
```

---

## 性能优化策略

### 1. 缓存机制

```go
// 题库缓存 (减少 YAML 文件读取)
var quizBankCache = &sync.Map{}

func GetChapterQuizBankCached(topic, chapter string) (*ChapterQuizBank, error) {
    cacheKey := fmt.Sprintf("%s/%s", topic, chapter)
    
    if cached, ok := quizBankCache.Load(cacheKey); ok {
        return cached.(*ChapterQuizBank), nil
    }
    
    bank, err := LoadChapterQuizBank(topic, chapter)
    if err != nil {
        return nil, err
    }
    
    quizBankCache.Store(cacheKey, bank)
    return bank, nil
}

// 缓存失效策略
func InvalidateCache(topic, chapter string) {
    cacheKey := fmt.Sprintf("%s/%s", topic, chapter)
    quizBankCache.Delete(cacheKey)
}
```

### 2. 索引文件

```yaml
# backend/quiz_data/index.yaml
# 避免遍历文件系统,快速定位章节
index:
  lexical_elements:
    comments:
      file: lexical_elements/comments.yaml
      question_count: 35
      available: true
      last_updated: "2026-01-15"
    keywords:
      file: lexical_elements/keywords.yaml
      question_count: 0
      available: false
      last_updated: null
```

### 3. 批量加载优化

```go
// 并发加载多个章节题库
func LoadMultipleChapters(chapters []ChapterKey) (map[string]*ChapterQuizBank, error) {
    results := make(map[string]*ChapterQuizBank)
    errChan := make(chan error, len(chapters))
    var wg sync.WaitGroup
    
    for _, key := range chapters {
        wg.Add(1)
        go func(k ChapterKey) {
            defer wg.Done()
            bank, err := GetChapterQuizBankCached(k.Topic, k.Chapter)
            if err != nil {
                errChan <- err
                return
            }
            results[k.String()] = bank
        }(key)
    }
    
    wg.Wait()
    close(errChan)
    
    if err := <-errChan; err != nil {
        return nil, err
    }
    
    return results, nil
}
```

---

## 数据迁移与备份

### 备份策略

```bash
# 定期备份测验记录
backup_quiz_data.sh:
  - 每日备份所有用户测验记录到 data/backups/YYYYMMDD/
  - 保留最近7天的每日备份
  - 每周一进行完整备份并压缩
  - 备份文件命名格式: quiz_attempts_YYYYMMDD_HHmmss.tar.gz
```

### 数据清理策略

```go
// 清理超过180天的测验记录
func CleanOldAttempts(retentionDays int) error {
    cutoff := time.Now().AddDate(0, 0, -retentionDays)
    
    return filepath.Walk("data/attempts", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if !info.IsDir() && strings.HasSuffix(path, ".json") {
            if info.ModTime().Before(cutoff) {
                log.Printf("Deleting old attempt: %s", path)
                return os.Remove(path)
            }
        }
        
        return nil
    })
}
```

