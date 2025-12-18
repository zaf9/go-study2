# Implementation Plan - 章节测验体验升级与功能深化

**Feature**: 章节测验体验升级与功能深化 (Chapter Quiz UX & Depth Enhancement)
**Status**: DRAFT
**Feature Branch**: `014-quiz-ux-enhancement`

## Technical Context

### 1. Technology Stack
- **Frontend**: Next.js (React), Ant Design (UI), Axios (API), React Context (State)
- **Backend**: Go 1.24, GoFrame v2.x (Web Framework), SQLite (Persistence)
- **Infrastructure**: Existing `backend` and `frontend` directories

### 2. Architecture Overview
- **UI Layer**:
  - `QuizViewer` component: Unified renderer for both "Take Quiz" and "Review/History" modes.
  - `QuizQuestionCard` component: Handles stable option labeling (A, B, C, D) independent of content order.
  - `HistoryPage`: New page for listing past attempts with status/score.
- **API Layer**:
  - `POST /quiz/submit`: Existing, updated to return percentage score.
  - `GET /quiz/history`: List session history.
  - `GET /quiz/history/{sessionId}`: Detailed review data.
- **Persistence Layer**:
  - `QuizSession` table: Stores attempt meta (score, time, status).
  - `QuizAttempt` table: Stores granular answer data per question.

### 3. Key Technical Decisions
- **Stable Labeling**: Backend shuffles content but removes prefixes. Frontend assigns A-D based on array index.
- **Review Mode**: Web-only feature per clarification. CLI remains "take-only".
- **State Management**: Simple React State + Context for the active quiz session. No Redux required (YAGNI).
- **Persistence**: SQLite utilizing `gdb` ORM features.

### 4. Dependencies
- **GoFrame**: `ghttp`, `gdb`, `glog` packages.
- **Ant Design**: `Table`, `Modal`, `Skeleton`, `Tag` components.

## Constitution Check

### Principle Alignment
- [x] **Principle I (Code Quality)**: Components split (QuizViewer vs QuizCard) for single responsibility.
- [x] **Principle III (Testing)**: 80% coverage planned for new Quiz logic and API endpoints.
- [x] **Principle XXII (CLI/Web)**: Explicit exception granted for History Review (Web-only) in spec clarification.
- [ ] **Principle VII (Security)**: Input validation on Submit to prevent logic errors/cheating.

### Gates
- **Gate 1**: Spec is clear and unambiguous? **YES** (Clarifications resolved).
- **Gate 2**: Tech stack approved? **YES**.
- **Gate 3**: Core principles met? **YES**.

## Research Questions (Phase 0)

1. [RESOLVED] **Label Stability**: Frontend index-based mapping confirmed in spec constraints.
2. [RESOLVED] **Persistence**: SQLite + gdb confirmed.
3. [RESOLVED] **Review Scope**: CLI excluded from review mode history.
4. [RESOLVED] **GoFrame ORM Transactions**: 使用 `g.DB().Transaction()` 封装，详见 `research.md` §3。
5. [RESOLVED] **AntD Skeleton Config**: 创建 `QuizSkeletonLoader.tsx` 精确模拟布局，详见 `research.md` §4。
6. [RESOLVED] **userChoice 存储策略**: 采用内容值存储方案，详见 `data-model.md` Final Decision。

---

### Execution Steps
1. **Research**: Confirm gdb transaction syntax and Skeleton layout match.
2. **Design**: Create `data-model.md` and API contracts.
3. **Plan**: Update `tasks.md`.
