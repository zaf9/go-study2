# Phase 0: 技术调研与决策

**Feature**: 学习进度与测验体验优化  
**Created**: 2025-12-30

## 调研任务概览

本阶段解决Technical Context中标记的NEEDS CLARIFICATION项,并为Phase 1设计提供技术基础。

## 1. 章节数量统计方案

### 问题描述
主题列表页面显示的章节数量不准确,需要确定准确的数据源和统计方法。

### 调研结果

**现状分析**:
- 前端已有静态数据源: `frontend/lib/static-routes.ts` 定义了完整的主题和章节映射关系
- 数据结构: `topicChapters: Record<TopicKey, string[]>` 明确列出每个主题的章节列表
- 后端quiz_data目录与前端static-routes.ts保持同步

**决策**: 前端直接从`static-routes.ts`计算章节数量

**理由**:
1. static-routes.ts是单一数据源(Single Source of Truth),已被Next.js用于静态路由生成
2. 避免额外API调用,减少网络开销和延迟
3. 章节列表变化时,仅需更新一处代码
4. 计算逻辑简单: `topicChapters[topic].length`

**实现方式**:
```typescript
// frontend/lib/chapter-count.ts
import { topicChapters, TopicKey } from './static-routes';

export function getChapterCount(topic: TopicKey): number {
  return (topicChapters[topic] || []).length;
}

export function getAllChapterCounts(): Record<TopicKey, number> {
  return Object.fromEntries(
    Object.entries(topicChapters).map(([topic, chapters]) => [topic, chapters.length])
  );
}
```

**替代方案**:
- ❌ 后端API统计: 增加网络往返,且后端需维护与前端一致的章节列表
- ❌ 数据库存储章节数: 需要同步机制,增加复杂度

---

## 2. 章节学习状态计算逻辑

### 问题描述
如何在不大规模修改数据库schema的情况下,实现三状态(未开始/学习中/已完成)的章节追踪?

### 调研结果

**现状分析**:
- 数据库表`learning_progress`已有`status`字段(VARCHAR)
- 当前使用状态: `not_started`, `in_progress`, `done`, `tested`
- `LearningProgress`领域模型已定义状态常量

**决策**: 复用现有status字段,优化状态更新逻辑

**状态转换规则**:
1. **未开始(not_started)**: 默认状态,learning_progress表中无记录即为未开始
2. **学习中(in_progress)**: 用户访问章节内容页时创建记录,设置status=in_progress
3. **已完成(completed)**: 用户完成章节测验且通过(quiz_passed=true)时,更新status=completed

**实现逻辑**:
```go
// backend/internal/domain/progress/service.go
func (s *Service) UpdateChapterStatus(ctx context.Context, userID int64, topic, chapter string) error {
    existing, err := s.repo.Get(ctx, userID, topic, chapter)
    
    if err != nil || existing == nil {
        // 首次访问,创建记录
        return s.repo.CreateOrUpdate(ctx, &LearningProgress{
            UserID: userID,
            Topic: topic,
            Chapter: chapter,
            Status: StatusInProgress,
            FirstVisitAt: time.Now(),
            LastVisitAt: time.Now(),
        })
    }
    
    // 已有记录,更新访问时间
    existing.LastVisitAt = time.Now()
    return s.repo.CreateOrUpdate(ctx, existing)
}

func (s *Service) CompleteChapter(ctx context.Context, userID int64, topic, chapter string, quizScore int, passed bool) error {
    existing, err := s.repo.Get(ctx, userID, topic, chapter)
    if err != nil {
        return err
    }
    
    existing.QuizScore = quizScore
    existing.QuizPassed = passed
    if passed {
        existing.Status = StatusCompleted
        now := time.Now()
        existing.CompletedAt = &now
    }
    
    return s.repo.CreateOrUpdate(ctx, existing)
}
```

**前端状态判断**:
```typescript
// frontend/lib/chapter-status.ts
export type ChapterStatus = 'not_started' | 'in_progress' | 'completed';

export function getChapterStatus(progress: ChapterProgress | null): ChapterStatus {
  if (!progress) return 'not_started';
  if (progress.quizPassed) return 'completed';
  return 'in_progress';
}
```

**替代方案**:
- ❌ 新增chapter_status表: 过度设计,现有字段已足够
- ❌ 仅依赖quiz_passed判断: 无法区分"未开始"和"学习中"

---

## 3. 测验加载与提交流程修复

### 问题描述
当前测验功能存在加载失败和提交错误,需要定位根因并提出修复方案。

### 调研结果

**问题诊断**:

1. **Session机制缺陷**: 
   - 现象: 前端调用`GET /api/v1/quiz/:topic/:chapter`返回空session
   - 根因: 后端未正确初始化quiz_sessions表记录,或sessionId生成逻辑有误

2. **题目加载失败**:
   - 现象: quiz_data目录存在题目文件,但API返回questions为空
   - 根因: 后端加载题目时路径拼接错误,或YAML解析失败但未抛出明确错误

3. **提交重复问题**:
   - 现象: 用户快速点击提交按钮导致重复记录
   - 根因: 前端防重复提交逻辑不完善,后端缺少幂等性保证

**决策**: 分三步修复

**Step 1: 规范Session创建流程**
```go
// backend/internal/domain/quiz/service.go
func (s *Service) GetOrCreateSession(ctx context.Context, userID int64, topic, chapter string) (*QuizSession, error) {
    // 1. 尝试查找活跃session (24小时内)
    existing, err := s.repo.GetActiveSession(ctx, userID, topic, chapter)
    if err == nil && existing != nil {
        return existing, nil
    }
    
    // 2. 加载题目
    questions, err := s.loadQuestions(topic, chapter)
    if err != nil {
        return nil, fmt.Errorf("加载题目失败: %w", err)
    }
    
    if len(questions) == 0 {
        return nil, fmt.Errorf("章节 %s/%s 暂无测验题目", topic, chapter)
    }
    
    // 3. 创建新session
    sessionID := generateSessionID(userID, topic, chapter)
    session := &QuizSession{
        SessionID: sessionID,
        UserID: userID,
        Topic: topic,
        Chapter: chapter,
        Questions: questions,
        CreatedAt: time.Now(),
    }
    
    if err := s.repo.CreateSession(ctx, session); err != nil {
        return nil, fmt.Errorf("创建session失败: %w", err)
    }
    
    return session, nil
}

func generateSessionID(userID int64, topic, chapter string) string {
    ts := time.Now().Unix()
    data := fmt.Sprintf("%d-%s-%s-%d", userID, topic, chapter, ts)
    hash := md5.Sum([]byte(data))
    return hex.EncodeToString(hash[:])
}
```

**Step 2: 优化题目加载逻辑**
```go
func (s *Service) loadQuestions(topic, chapter string) ([]QuizQuestion, error) {
    // 1. 构建文件路径
    basePath := "quiz_data"
    filePath := filepath.Join(basePath, topic, chapter+".yaml")
    
    // 2. 检查文件是否存在
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return []QuizQuestion{}, nil // 无题目不是错误,返回空切片
    }
    
    // 3. 读取并解析YAML
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("读取题目文件失败 %s: %w", filePath, err)
    }
    
    var questions []QuizQuestion
    if err := yaml.Unmarshal(data, &questions); err != nil {
        return nil, fmt.Errorf("解析题目文件失败 %s: %w", filePath, err)
    }
    
    return questions, nil
}
```

**Step 3: 防重复提交机制**

前端:
```typescript
// frontend/hooks/useQuiz.ts (已有部分实现,需强化)
const isSubmittingRef = useRef(false);

const submit = async () => {
  if (isSubmittingRef.current || submitting) {
    console.warn('请勿重复提交');
    return null;
  }
  
  isSubmittingRef.current = true;
  setSubmitting(true);
  
  try {
    const result = await submitQuiz(payload);
    return result;
  } catch (error) {
    if (error.code === 409) {
      // 后端返回已提交错误
      throw new Error('该测验已提交,请勿重复提交');
    }
    throw error;
  } finally {
    // 成功后不重置,防止再次提交
    // isSubmittingRef.current = false;
  }
};
```

后端:
```go
func (s *Service) SubmitQuiz(ctx context.Context, sessionID string, answers []UserAnswer) (*QuizResult, error) {
    // 1. 检查session是否已提交
    session, err := s.repo.GetSessionByID(ctx, sessionID)
    if err != nil {
        return nil, err
    }
    
    if session.SubmittedAt != nil {
        return nil, &ConflictError{Message: "该测验已提交,请勿重复操作"}
    }
    
    // 2. 评分逻辑
    result := s.gradeAnswers(session.Questions, answers)
    
    // 3. 使用事务更新
    return result, s.repo.Transaction(ctx, func(tx *gdb.TX) error {
        // 标记session已提交
        if _, err := tx.Model("quiz_sessions").Where("session_id", sessionID).Update(g.Map{
            "submitted_at": time.Now(),
        }); err != nil {
            return err
        }
        
        // 保存答题记录
        for _, ans := range answers {
            attempt := &QuizAttempt{
                SessionID: sessionID,
                QuestionID: ans.QuestionID,
                UserChoice: ans.UserAnswer,
                IsCorrect: ans.IsCorrect,
                AttemptedAt: time.Now(),
            }
            if _, err := tx.Model("quiz_attempts").Insert(attempt); err != nil {
                return err
            }
        }
        
        // 更新章节进度
        return s.progressService.CompleteChapter(ctx, session.UserID, session.Topic, session.Chapter, result.Score, result.Passed)
    })
}
```

**替代方案**:
- ❌ 前端轮询检查session状态: 增加服务器负载,用户体验差
- ❌ 仅前端防重复: 不可靠,需后端配合幂等性保证

---

## 4. 前端状态管理一致性方案

### 问题描述
不同页面显示的进度数据不一致,需要统一数据来源和更新机制。

### 调研结果

**现状分析**:
- 前端已使用SWR进行数据缓存
- `useProgress` Hook封装进度查询逻辑
- 问题: 各页面独立调用API,缓存key不统一导致数据不同步

**决策**: 规范SWR缓存key和数据更新策略

**缓存key规范**:
```typescript
// frontend/services/progressService.ts
export const progressKeys = {
  overview: 'progress/overview',              // 全局进度概览
  topic: (topic: string) => ['progress/topic', topic],  // 单个主题进度
  chapter: (topic: string, chapter: string) => ['progress/chapter', topic, chapter], // 单个章节进度
};
```

**统一数据更新流程**:
```typescript
// 1. 用户进入章节
await updateProgress({ topic, chapter, status: 'in_progress' });
await mutate(progressKeys.topic(topic));  // 刷新主题进度
await mutate(progressKeys.overview);      // 刷新全局进度

// 2. 用户完成测验
await submitQuiz(payload);
await mutate(progressKeys.chapter(topic, chapter));  // 刷新章节状态
await mutate(progressKeys.topic(topic));
await mutate(progressKeys.overview);

// 3. 乐观更新策略
mutate(progressKeys.overview, (current) => {
  if (!current) return current;
  return {
    ...current,
    completedChapters: current.completedChapters + 1,
  };
}, false); // false表示不立即revalidate
```

**Context层聚合**:
```typescript
// frontend/contexts/LearningContext.tsx
export function LearningProvider({ children }) {
  const { data: overview, mutate: mutateOverview } = useSWR(progressKeys.overview, fetchProgressOverview);
  
  const refreshAll = useCallback(async () => {
    await Promise.all([
      mutateOverview(),
      mutate(() => true, undefined, { revalidate: true }), // 刷新所有SWR缓存
    ]);
  }, [mutateOverview]);
  
  return (
    <LearningContext.Provider value={{ overview, refreshAll }}>
      {children}
    </LearningContext.Provider>
  );
}
```

**替代方案**:
- ❌ Redux/Zustand全局状态: 过度设计,SWR已足够
- ❌ 每次都从后端拉取: 增加网络开销,用户体验差

---

## 5. 测验中心导航设计

### 问题描述
需要提供便捷的测验访问入口,优化用户导航路径。

### 调研结果

**决策**: 新增独立的"测验中心"页面

**路由设计**:
```
/quiz-center                    # 测验中心首页
/quiz-center?topic=constants    # 筛选特定主题的测验记录
```

**导航入口**:
1. **主导航栏**: 添加"测验中心"菜单项(与"主题"、"进度"并列)
2. **主题列表**: 每个主题卡片添加"查看测验"快速链接
3. **章节详情**: 测验完成后显示"查看历史记录"按钮

**页面功能**:
- 显示所有已完成的测验记录(按时间倒序)
- 支持按主题筛选
- 展示: 章节名称、得分、正确率、完成时间
- 点击记录可查看详细答题情况

**实现**:
```typescript
// frontend/app/(protected)/quiz-center/page.tsx
export default function QuizCenterPage({ searchParams }) {
  const topic = searchParams.topic;
  const { data, isLoading } = useQuizHistory(topic);
  
  return (
    <div>
      <PageHeader title="测验中心" />
      <TopicFilter current={topic} />
      <QuizHistoryList records={data} />
    </div>
  );
}
```

**替代方案**:
- ❌ 将测验记录整合到"进度"页面: 功能混杂,不符合单一职责原则
- ❌ 仅在章节页面显示测验入口: 缺少全局视图,不便于复习

---

## 技术选型总结

| 决策点 | 选择方案 | 关键理由 |
|--------|----------|----------|
| 章节计数数据源 | 前端static-routes.ts | 单一数据源,避免API调用 |
| 章节状态存储 | 复用learning_progress.status | 无需schema变更,向后兼容 |
| 测验Session管理 | 优化创建逻辑+幂等性保证 | 修复当前问题,保证数据一致性 |
| 前端状态管理 | SWR + 规范缓存key | 利用现有基础设施,最小改动 |
| 测验导航 | 新增独立测验中心页面 | 功能聚焦,用户路径清晰 |

所有决策均遵循YAGNI原则,优先利用现有技术栈,避免引入新依赖。
