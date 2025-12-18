# QuickStart Guide: Quiz UX Enhancement

## Prerequisites
- Backend running: `cd backend && go run main.go`
- Frontend running: `cd frontend && npm run dev`

## New Features
1. **Interactive Quiz**: Stable A-D labels.
2. **Review Mode**: View history/details at `/quiz/history`.
3. **Persistence**: Scores saved to SQLite.

## How to Test
1. **Take a Quiz**: 
   - Navigate to any topic (e.g., `/topics/lexical/comment`).
   - Click "Start Quiz".
   - Verify Loading Skeleton appears effectively.
   - Verify Labels are A, B, C, D regardless of content.
   - Click Submit -> Confirm Modal should appear.
2. **Check History**:
   - Go to `/quiz/history`.
   - Click "Review" on the recent attempt.
   - Verify "Your Answer" and "Correct Answer" tags.
