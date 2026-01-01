# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**go-study2** is a comprehensive Go language learning platform with dual-mode architecture (CLI + HTTP) and a modern Next.js frontend. It provides interactive learning modules for Go lexical elements, constants, variables, and types with integrated quizzes and progress tracking.

**Current Branch**: `017-complete-quiz-bank`
**Tech Stack**: Go 1.24.5 + GoFrame v2.9.6 (backend), Next.js 14 + Ant Design 5 (frontend)

## Common Commands

### Backend Development

```bash
# CLI mode (default)
cd backend
go run main.go

# HTTP daemon mode (port 8080)
go run main.go -d

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test suite
go test -v ./tests/unit/...

# Build for production
./build.bat  # Windows (format → vet → test → build)
# Or manually:
go fmt ./...
go vet ./...
go test -cover ./...
go build -o ../bin/gostudy.exe main.go
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Development server (port 3000, proxies to backend at 8080)
npm run dev

# Build for production (static export to out/)
npm run build
npm run export

# Run tests
npm test -- --coverage
npm test -- --watch
```

### Production Deployment

```bash
# 1. Build frontend
cd frontend && npm install && npm run build && npm run export

# 2. Build backend
cd ../backend
go test ./...
go build -o ../bin/gostudy main.go

# 3. Run (serves frontend + API on port 8080)
cd ../bin
./gostudy -d
```

**Environment Variables**:
- `JWT_SECRET`: Required for JWT signing (≥32 chars)
- `NEXT_PUBLIC_API_URL`: Backend API base URL (default: http://localhost:8080)

## Architecture

### Backend - Clean Architecture / DDD

**Layer Separation**:
1. **Domain Layer** (`internal/domain/`): Business entities and interfaces
   - `user/`: User entity with auth service
   - `progress/`: Learning progress tracking
   - `quiz/`: Quiz management and scoring

2. **Application Layer** (`internal/app/`): Use cases and coordination
   - `http_server/handler/`: Request handlers for auth/progress/quiz
   - `http_server/middleware/`: Auth, logging, CORS, recovery
   - `http_server/router.go`: API route definitions

3. **Infrastructure Layer** (`internal/infrastructure/`): External concerns
   - `database/`: SQLite connection with WAL mode, migrations
   - `repository/`: Repository implementations (User/Progress/Quiz)

**Key Patterns**:
- **Repository Pattern**: Domain defines interfaces, infrastructure provides implementations
- **Middleware Pipeline**: CORS → PanicRecovery → AccessLog → Auth → Format
- **Hub Pattern** (`internal/websocket/`): WebSocket connection management with broadcast
- **Dependency Injection**: Services receive repository interfaces

### Frontend - Next.js App Router

**Route Organization**:
- `(auth)/`: Public routes (login, register)
- `(protected)/`: Authenticated routes (dashboard, topics, quiz, progress, profile)
- File-based routing with `layout.tsx` for shared UI
- Static export enabled (`output: 'export'`)

**State Management**:
- **Server State**: SWR for API caching and revalidation
- **Client State**: React Context (`AuthProvider`, `WebSocketProvider`)
- **URL State**: Search params for filters

**Key Patterns**:
- **Custom Hooks**: `useAuth`, `useProgress`, `useQuiz`, `useScrollPosition`
- **Provider Pattern**: Global auth and WebSocket state
- **Composition**: Reusable quiz components with controlled inputs
- **Error Boundaries**: Graceful degradation

## Critical Files

**Backend Entry Points**:
- `backend/main.go` - CLI/HTTP mode entry point
- `backend/internal/app/http_server/router.go` - API route definitions
- `backend/configs/config.yaml` - Server configuration

**Frontend Entry Points**:
- `frontend/app/layout.tsx` - Root layout with providers
- `frontend/lib/api.ts` - Axios client with token refresh
- `frontend/lib/websocket.ts` - WebSocket client with auto-reconnect

**Domain Logic**:
- `backend/internal/domain/user/` - User entity and auth service
- `backend/internal/domain/progress/` - Progress tracking service
- `backend/internal/domain/quiz/` - Quiz management and scoring

## API Design

### RESTful Endpoints

**Base URL**: `http://localhost:8080/api/v1`

**Authentication**:
- `POST /auth/register` - User registration
- `POST /auth/login` - Login (returns access token + HttpOnly refresh token)
- `POST /auth/refresh` - Refresh access token
- `GET /auth/profile` - Get current user
- `POST /auth/logout` - Logout

**Learning**:
- `GET /topics` - List all topics
- `GET /topic/{topic}/{chapter}` - Get chapter content
- `GET /progress` - Get all progress (requires auth)
- `POST /progress` - Save progress (requires auth)

**Quiz**:
- `GET /quiz/{topic}/{chapter}` - Start quiz (requires auth)
- `POST /quiz/submit` - Submit answers and score (requires auth)
- `GET /quiz/history` - Quiz history with filters (requires auth)

**WebSocket**:
- `GET /ws/dashboard` - Real-time updates for progress and quiz completion

### Response Format

Unified JSON responses:
```json
{
  "code": 20000,
  "message": "success",
  "data": { ... }
}
```

Error codes:
- `20000`: Success
- `40001`: Invalid request
- `40100`: Unauthorized
- `40300`: Forbidden
- `50000`: Internal error

## Key Data Flows

### Authentication Flow
1. User posts credentials to `/auth/login`
2. Backend validates with bcrypt, generates JWT access token (7d expiry) + HttpOnly refresh token
3. Frontend stores access token in localStorage, refresh token in HttpOnly cookie
4. Axios interceptor adds `Authorization: Bearer` header to all requests
5. On 401, interceptor calls `/auth/refresh`, retries original request
6. If refresh fails, redirect to login

### Learning Progress Flow
1. User navigates to `/topics/{topic}/{chapter}`
2. Chapter content loads via SWR from `/topic/{topic}/{chapter}`
3. Frontend tracks scroll position and time spent
4. Debounced POST to `/progress` every 10 seconds
5. Backend calculates completion % based on: `time_spent / (base_time * completion_fraction)`
6. Progress saved to DB, broadcast via WebSocket
7. Dashboard updates in real-time

### Quiz Flow
1. User starts quiz: GET `/quiz/{topic}/{chapter}`
2. Backend loads from `quiz_data/{topic}/{chapter}.md` (YAML format)
3. Random selection: 4 single-choice + 4 multiple-choice questions
4. Difficulty distribution: 40% easy, 40% medium, 20% hard
5. Backend returns session ID and shuffled questions
6. User answers and submits to `/quiz/submit`
7. Backend scores: `(correct / total) * 100`, pass threshold = 60%
8. Results saved to `quiz_sessions` and `quiz_answers` tables
9. WebSocket broadcasts completion to Dashboard

## WebSocket Real-time Updates

**Backend** (`internal/websocket/hub.go`):
- Maintains active connections map
- Broadcasts on events: `progress_updated`, `quiz_completed`
- Clients authenticated via query token

**Frontend** (`lib/websocket.ts`):
- `WebSocketProvider` wraps app
- Auto-reconnect with exponential backoff (1s → 30s max)
- Dashboard subscribes to progress/quiz events
- Updates UI in real-time without refresh

## Testing Strategy

**Backend Tests** (`backend/tests/`):
- `unit/`: Component tests (handlers, services)
- `integration/`: API endpoint tests with test DB
- `contract/`: Interface contract verification
- Target coverage: ≥80%
- Tools: `go test`, `testify/assert`

**Frontend Tests** (`frontend/__tests__/`):
- `app/`: Page component tests
- `components/`: UI component tests
- `hooks/`: Custom hook tests
- `lib/`: Utility function tests
- Tools: Jest, React Testing Library, SWR devtools

## Configuration Management

**Backend Config** (`backend/configs/config.yaml`):
```yaml
http:
  port: 8080
database:
  path: "./data/gostudy.db"
  pragmas:
    - "journal_mode=WAL"
    - "busy_timeout=5000"
jwt:
  secret: "${JWT_SECRET}"
  accessTokenExpiry: 604800  # 7 days
static:
  enabled: true
  path: "../frontend/out"
  spaFallback: true
progress:
  readCharsPerSec: 5.0
  quiz:
    dataPath: quiz_data
    questionCount:
      single: 4
      multiple: 4
```

**Frontend Config**:
- `NEXT_PUBLIC_API_URL`: Backend API base URL
- `next.config.ts`: Rewrites `/api/v1/*` to backend in dev mode, static export for production

## Branch Strategy

- `master`: Production-ready code
- Feature branches: `017-complete-quiz-bank`, `009-frontend-ui`, etc.
- Commit conventions: `feat:`, `fix:`, `docs:`, `test:`, `refactor:`

## Quality Gates

Before committing:
1. Backend: `go fmt ./...` → `go vet ./...` → `go test -cover ./...`
2. Frontend: `npm run lint` → `npm test -- --coverage`
3. Both: Build must succeed

## Default Credentials

- **Admin**: `admin` / `GoStudy@123`
- First login requires password change
- All other users must register

## Development Workflow

1. Ensure `./build.bat` exists and run it first (dependency check + compile)
2. Start backend: `cd backend && go run main.go -d`
3. Start frontend: `cd frontend && npm run dev`
4. Access at `http://localhost:3000/` (dev) or `http://localhost:8080/` (production)
5. Run tests before committing
6. Use conventional commit messages

## Deployment Notes

- Single-port deployment: Backend serves static files from `frontend/out/`
- SPA fallback: Non-API routes return `index.html`
- SQLite auto-migration on startup
- Seed default admin if not exists
- Quiz data imported from YAML on first run
- `JWT_SECRET` environment variable required in production

## Important Architecture Constraints

1. **Repository Pattern**: Always inject interfaces, not implementations
2. **Middleware Order**: CORS → PanicRecovery → AccessLog → Auth → Format
3. **WebSocket Auth**: Via query parameter token, auto-reconnect with exponential backoff
4. **Progress Tracking**: Debounced 10s, algorithm based on content length and difficulty
5. **Quiz Questions**: Loaded from YAML, shuffled per session, difficulty-weighted
6. **Static Export**: Frontend uses `output: 'export'`, pre-generates all routes
7. **Token Refresh**: Automatic via Axios interceptor, silent to user
8. **Error Handling**: Unified `{code, message, data}` format across all APIs

## Related Documentation

- `README.md` - Project overview and quick start
- `docs/DEPLOYMENT.md` - Deployment guide
- `docs/API.md` - Detailed API documentation
- `.cursor/rules/specify-rules.mdc` - Auto-generated development guidelines
- `specs/` - Feature specifications and implementation plans
