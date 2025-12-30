# 快速开始: 学习进度与测验体验优化

**Feature**: 016-progress-quiz-ux  
**Status**: 设计完成,待实施  
**Estimated Time**: 3-5天

---

## 功能概述

本功能修复学习进度追踪和测验系统的五大核心问题:

1. ✅ **准确章节数量**: 主题列表显示真实的章节总数
2. ✅ **三状态追踪**: 清晰区分未开始/学习中/已完成章节
3. ✅ **可用测验**: 修复加载失败和提交错误,确保测验功能正常工作
4. ✅ **一致数据**: 所有页面显示相同的进度统计
5. ✅ **测验导航**: 新增测验中心,便捷访问测验记录

---

## 前置条件

### 开发环境
- Go 1.24+
- Node.js 18+
- SQLite3
- Git

### 已有基础设施
- 后端: GoFrame v2, Gin framework
- 前端: Next.js 14, React 18, SWR, Ant Design 5
- 数据库表: `learning_progress`, `quiz_sessions`, `quiz_attempts`
- 认证: JWT Bearer Token

---

## 快速体验(用户视角)

### 场景1: 查看准确的章节数量

1. **访问主题列表**: 导航至 `/topics`
2. **观察章节数**: 每个主题卡片显示章节总数,例如"章节数:12"
3. **验证准确性**: 
   - Constants主题 → 显示12章节
   - Lexical Elements主题 → 显示11章节
   - 数量与`frontend/lib/static-routes.ts`中定义一致

**预期效果**: 章节数量100%准确,不受用户学习进度影响

---

### 场景2: 查看章节学习状态

1. **进入主题详情**: 点击某个主题,例如"Constants"
2. **观察状态标识**: 章节列表显示三种状态
   - 🔘 **未开始**(灰色): 从未访问过的章节
   - 📖 **学习中**(蓝色): 已访问但未完成测验的章节
   - ✅ **已完成**(绿色): 通过测验的章节
3. **访问章节**: 点击某个"未开始"章节
4. **返回列表**: 状态自动更新为"学习中"

**预期效果**: 状态实时更新,清晰反映学习进度

---

### 场景3: 完成章节测验

1. **进入章节详情**: 访问 `/topics/constants/boolean`
2. **点击测验按钮**: "开始测验"或"进入测验"
3. **答题**: 选择答案(单选/多选)
4. **提交**: 点击"提交测验"
5. **查看结果**: 
   - 显示总分、正确率
   - 每题显示用户答案和正确答案
   - 得分≥60%显示"通过",章节状态更新为"已完成"

**预期效果**: 测验流程顺畅,无加载失败或提交错误

---

### 场景4: 验证进度数据一致性

1. **主题列表页**: 记录显示的已完成章节数,例如"3/12"
2. **访问进度页**: 导航至 `/progress`
3. **对比数据**: 进度页显示相同的"3/12已完成"
4. **完成新测验**: 完成一个新章节的测验
5. **返回主题列表**: 自动更新为"4/12"
6. **刷新进度页**: 同样显示"4/12"

**预期效果**: 所有页面数据实时同步,无延迟或不一致

---

### 场景5: 使用测验中心

1. **访问测验中心**: 点击主导航的"测验中心"或访问 `/quiz-center`
2. **查看历史记录**: 显示所有已完成的测验(按时间倒序)
   - 章节名称
   - 得分和正确率
   - 完成时间
   - 通过/未通过标识
3. **筛选主题**: 使用筛选器查看特定主题的测验
4. **查看详情**: 点击某条记录,查看答题详情

**预期效果**: 测验记录一目了然,便于复习和追踪

---

## 开发指南

### Phase 1: 后端实现(2天)

#### 步骤1.1: 添加数据库迁移

```bash
cd backend
# 创建迁移文件
cat > internal/infra/migrations/016_add_submitted_at.sql << 'EOF'
ALTER TABLE quiz_sessions ADD COLUMN submitted_at DATETIME;

UPDATE quiz_sessions
SET submitted_at = (
    SELECT MIN(attempted_at)
    FROM quiz_attempts
    WHERE quiz_attempts.session_id = quiz_sessions.session_id
)
WHERE id IN (
    SELECT DISTINCT qs.id
    FROM quiz_sessions qs
    INNER JOIN quiz_attempts qa ON qs.session_id = qa.session_id
);
EOF
```

#### 步骤1.2: 优化进度Service

**文件**: `backend/internal/domain/progress/service.go`

关键方法:
- `GetOverviewWithChapterCounts()`: 聚合所有主题的进度,计算章节总数
- `GetTopicProgressWithStatus()`: 返回主题的章节列表,填充三状态
- `CompleteChapter()`: 测验通过时更新状态为completed

**测试**:
```bash
cd backend
go test ./internal/domain/progress/... -v
```

#### 步骤1.3: 修复测验Service

**文件**: `backend/internal/domain/quiz/service.go`

关键方法:
- `GetOrCreateSession()`: 规范session创建逻辑,复用活跃会话
- `loadQuestions()`: 优化题目加载,明确错误处理
- `SubmitQuiz()`: 添加幂等性检查,使用事务提交

**测试**:
```bash
go test ./internal/domain/quiz/... -v
```

#### 步骤1.4: 更新Handler层

**文件**: 
- `backend/internal/interfaces/progress_handler.go`
- `backend/internal/interfaces/quiz_handler.go`

新增/修改端点:
- `GET /api/v1/progress/overview`
- `GET /api/v1/progress/topic/:topic`
- `GET /api/v1/quiz/:topic/:chapter`
- `POST /api/v1/quiz/submit`

**测试**:
```bash
go test ./internal/interfaces/... -v
```

---

### Phase 2: 前端实现(2天)

#### 步骤2.1: 修复章节数量显示

**文件**: `frontend/components/learning/TopicCard.tsx`

```typescript
import { getChapterCount } from '@/lib/chapter-count';

export default function TopicCard({ topic }) {
  const chapterCount = getChapterCount(topic.key);
  
  return (
    <Card>
      <Text>章节数:{chapterCount}</Text>
    </Card>
  );
}
```

**测试**:
```bash
cd frontend
npm test -- TopicCard.test.tsx
```

#### 步骤2.2: 实现三状态显示

**新增组件**: `frontend/components/learning/ChapterStatusBadge.tsx`

```typescript
export default function ChapterStatusBadge({ status }: { status: ChapterStatus }) {
  const config = {
    not_started: { color: 'default', text: '未开始', icon: <BookOutlined /> },
    in_progress: { color: 'processing', text: '学习中', icon: <ReadOutlined /> },
    completed: { color: 'success', text: '已完成', icon: <CheckCircleOutlined /> },
  };
  
  const { color, text, icon } = config[status];
  
  return <Badge status={color} text={text} icon={icon} />;
}
```

**修改**: `frontend/app/(protected)/topics/[topic]/page.tsx`

**测试**:
```bash
npm test -- ChapterStatusBadge.test.tsx
```

#### 步骤2.3: 优化测验Hook

**文件**: `frontend/hooks/useQuiz.ts`

修复:
- 强化防重复提交逻辑(useRef + state双重锁定)
- 添加错误处理(409 Conflict → 提示已提交)
- 成功提交后不重置锁定状态

**测试**:
```bash
npm test -- useQuiz.test.ts
```

#### 步骤2.4: 统一进度数据

**文件**: `frontend/hooks/useProgress.ts`

规范SWR缓存key:
```typescript
export const progressKeys = {
  overview: 'progress/overview',
  topic: (topic: string) => ['progress/topic', topic],
  chapter: (topic: string, chapter: string) => ['progress/chapter', topic, chapter],
};
```

测验提交后刷新所有相关缓存:
```typescript
await mutate(progressKeys.chapter(topic, chapter));
await mutate(progressKeys.topic(topic));
await mutate(progressKeys.overview);
```

**测试**:
```bash
npm test -- useProgress.test.ts
```

#### 步骤2.5: 新增测验中心

**新增页面**: `frontend/app/(protected)/quiz-center/page.tsx`

**新增组件**:
- `frontend/components/quiz/QuizCenter.tsx`: 容器组件
- `frontend/components/quiz/QuizHistoryCard.tsx`: 历史记录卡片

**路由**: `/quiz-center?topic=constants`

**测试**:
```bash
npm test -- QuizCenter.test.tsx
```

---

### Phase 3: 集成测试(1天)

#### 端到端测试场景

1. **章节数量准确性**
   ```bash
   # 访问主题列表,验证章节数
   curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:8080/api/v1/progress/overview
   # 检查data.topics[*].totalChapters是否准确
   ```

2. **三状态流程**
   ```bash
   # 1. 验证初始状态(未开始)
   # 2. 访问章节 → 状态变为学习中
   curl -X POST -H "Authorization: Bearer $TOKEN" \
     -d '{"topic":"constants","chapter":"boolean"}' \
     http://localhost:8080/api/v1/progress
   # 3. 完成测验 → 状态变为已完成
   ```

3. **测验完整流程**
   ```bash
   # 1. 获取会话
   curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:8080/api/v1/quiz/constants/boolean
   # 2. 提交答案
   curl -X POST -H "Authorization: Bearer $TOKEN" \
     -d '{"sessionId":"...","answers":[...]}' \
     http://localhost:8080/api/v1/quiz/submit
   # 3. 验证防重复提交(再次提交应返回409)
   ```

4. **数据一致性**
   ```bash
   # 完成测验后,对比多个端点的数据
   curl http://localhost:8080/api/v1/progress/overview
   curl http://localhost:8080/api/v1/progress/topic/constants
   # 验证completedChapters数量一致
   ```

---

## 验收标准

### 功能验收

- [ ] 主题列表显示的章节数量100%准确(与static-routes.ts一致)
- [ ] 所有章节明确显示三种状态(未开始/学习中/已完成)
- [ ] 测验加载成功率≥95%(有题目的章节)
- [ ] 测验提交成功率≥98%(无重复提交错误)
- [ ] 跨页面进度数据一致性100%(同一时刻查询结果相同)
- [ ] 测验中心可访问,显示所有历史记录
- [ ] 从主页到测验中心≤2次点击

### 性能验收

- [ ] 主题列表页面加载时间≤2秒
- [ ] 测验提交响应时间≤1秒
- [ ] 前端状态更新延迟≤500ms

### 测试覆盖率

- [ ] 后端单元测试覆盖率≥80%
- [ ] 前端核心组件测试覆盖率≥80%
- [ ] 端到端测试覆盖所有关键用户场景

---

## 常见问题

### Q1: 章节数量如何保持同步?

**A**: 前端`static-routes.ts`是单一数据源(SSOT),后端不维护章节列表。新增章节时:
1. 在`quiz_data/:topic/`添加YAML文件
2. 更新`frontend/lib/static-routes.ts`的`topicChapters`映射
3. 重新构建前端

### Q2: 如何判断章节是"未开始"状态?

**A**: 后端查询`learning_progress`表无记录即为未开始。前端:
1. 调用`GET /api/v1/progress/topic/:topic`
2. 对比返回的chapters列表与`topicChapters[topic]`
3. 缺失的章节填充默认值`{status: 'not_started'}`

### Q3: 测验重复提交如何防止?

**A**: 双重保护:
- **前端**: `useRef`立即锁定 + 提交成功后不重置
- **后端**: 检查`quiz_sessions.submitted_at`字段,非NULL返回409

### Q4: 进度数据如何保证一致性?

**A**: SWR缓存策略:
1. 统一使用`progressKeys`定义缓存key
2. 更新进度后,刷新所有相关缓存:`mutate(progressKeys.overview)`
3. 测验提交成功时,后端自动更新进度,前端无需额外调用

### Q5: 如何处理无测验的章节?

**A**: 
- 后端: `loadQuestions()`返回空切片(非错误)
- API: 返回`questions: []`,明确告知前端无题目
- 前端: 显示"该章节暂无测验"提示,隐藏测验按钮

---

## 下一步

完成本功能后,运行:

```bash
# 更新项目文档
npm run docs:update

# 提交变更
git add .
git commit -m "feat(progress-quiz): 修复学习进度与测验体验问题 (#016)"
git push origin 016-progress-quiz-ux

# 创建PR
gh pr create --title "Feature 016: 学习进度与测验体验优化" \
  --body "修复章节计数、三状态追踪、测验功能等5个核心问题"
```

**后续规划**:
- Feature 017: 学习时长统计和知识点掌握度分析
- Feature 018: 学习提醒和个性化推荐

---

## 参考文档

- [Feature Spec](./spec.md)
- [Data Model](./data-model.md)
- [Progress API Contract](./contracts/progress-api.md)
- [Quiz API Contract](./contracts/quiz-api.md)
- [Constitution](./../.specify/memory/constitution.md)
