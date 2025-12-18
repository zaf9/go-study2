# Technology Research: Quiz UX Enhancement

## 1. Label Stability Strategy
**Status**: ✅ RESOLVED  
**Decision**: Frontend Dynamic Mapping based on Index.
**Rationale**:
- Backend shuffles content strictly (Option A, B, C, D text is randomized).
- Backend does NOT include "A." prefixes in the text.
- Frontend iterates the array: Index 0 always renders as "A", Index 1 as "B".
- This ensures visual stability: Labels are always A-D in order, while content varies.
- **Verification**: `String.fromCharCode(65 + index)` is reliable for 0-25 range.

## 2. Persistence Layer
**Status**: ✅ RESOLVED  
**Decision**: SQLite + `gdb` (GoFrame ORM).
**Details**:
- `QuizSession` (Parent): `session_id`, `user_id`, `score`, `started_at`, `completed_at`.
- `QuizAttempt` (Child): `session_id`, `question_id`, `user_answer`, `is_correct`.
- **Transaction Pattern**: Use `gdb.DB.Transaction(ctx, func(ctx, tx) error { ... })` to ensure session and attempts are saved atomically.
- **Reference**: `backend/internal/infra/repository/quiz_repo.go` `SaveAttempts` method already demonstrates batch insert pattern.

## 3. GoFrame ORM Transaction Best Practice
**Status**: ✅ RESOLVED  
**Question**: 如何在 GoFrame gdb 中实现 QuizSession + QuizAttempt 的原子性保存？

**Solution**: 使用 `g.DB().Transaction()` 方法封装事务操作。

```go
// 最佳实践示例
func (r *QuizRepository) SaveSessionWithAttempts(ctx context.Context, session *entity.QuizSession, attempts []*entity.QuizAttempt) error {
    return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
        // 1. 插入 QuizSession
        _, err := tx.Model("quiz_session").Ctx(ctx).Data(session).Insert()
        if err != nil {
            return err
        }
        
        // 2. 批量插入 QuizAttempt
        if len(attempts) > 0 {
            _, err = tx.Model("quiz_attempt").Ctx(ctx).Data(attempts).Insert()
            if err != nil {
                return err
            }
        }
        
        return nil
    })
}
```

**Key Points**:
- `g.DB().Transaction()` 自动处理 `Begin/Commit/Rollback`
- 事务内所有操作使用同一个 `tx` 对象
- 任何返回 `error` 都会触发自动回滚
- Context 传递确保链路追踪正常工作

## 4. Loading State (Skeleton)
**Status**: ✅ RESOLVED  
**Question**: 如何配置 AntD Skeleton 以匹配 QuizQuestionCard 布局，实现平滑加载？

**Solution**: 创建专用 `QuizSkeletonLoader.tsx` 组件，精确模拟 QuizQuestionCard 结构。

```tsx
// QuizSkeletonLoader.tsx 最佳实践
import { Card, Skeleton, Space } from 'antd';

export const QuizSkeletonLoader: React.FC = () => {
  return (
    <Card style={{ marginBottom: 16 }}>
      {/* 题号 + 题型标签 */}
      <Space style={{ marginBottom: 12 }}>
        <Skeleton.Button active size="small" style={{ width: 60 }} />
        <Skeleton.Button active size="small" style={{ width: 80 }} />
      </Space>
      
      {/* 题干文本 (2-3 行) */}
      <Skeleton 
        active 
        paragraph={{ rows: 2, width: ['100%', '80%'] }} 
        title={false} 
      />
      
      {/* 选项列表 (4 个选项) */}
      <div style={{ marginTop: 16 }}>
        {[0, 1, 2, 3].map((i) => (
          <Skeleton.Input 
            key={i}
            active 
            block 
            style={{ height: 44, marginBottom: 12, borderRadius: 8 }} 
          />
        ))}
      </div>
    </Card>
  );
};

// 使用示例：加载多道题目
export const QuizSkeletonList: React.FC<{ count?: number }> = ({ count = 3 }) => (
  <>
    {Array.from({ length: count }).map((_, i) => (
      <QuizSkeletonLoader key={i} />
    ))}
  </>
);
```

**Key Points**:
- 骨架屏高度/间距与实际组件保持一致，避免布局抖动
- 使用 `paragraph.width` 数组模拟不等长文本
- 选项使用 `Skeleton.Input` 而非 `Skeleton.Button`，更接近真实选项样式
- 默认加载 3 道题目的骨架屏，可通过 `count` 参数调整

## 5. Review Mode Architecture
**Status**: ✅ RESOLVED  
**Decision**: Reused `QuizViewer` with `mode` prop.
- `mode="interactive"`: Enables selection, disables "Explanation".
- `mode="review"`: Disables selection, shows "Explanation" + "Your Answer" vs "Correct Answer" indicators.
- **Data Source**:
  - Interactive: Fetches from `/api/v1/quiz/{topic}/{chapter}`.
  - Review: Fetches from `/api/v1/quiz/history/{sessionId}`.

## 6. Security Validation
**Status**: ✅ RESOLVED  
**Requirement**: Verify answer length matches question count.
- **Backend Check**: `len(submit.Answers) == len(quiz.Questions)`.
- **Frontend Check**: Modal warning if `Object.keys(answers).length < total`.
