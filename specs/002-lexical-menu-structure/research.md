# Research: Lexical Menu Structure

**Feature**: 002-lexical-menu-structure  
**Date**: 2025-12-04  
**Status**: Complete

## Overview

This feature requires implementing a hierarchical menu system for the Go learning tool. Since the project already has a working menu pattern in `main.go`, this research focuses on adapting that pattern for sub-menus.

## Research Tasks

### 1. Menu Pattern Analysis

**Question**: How should we pass I/O streams to enable interactive sub-menus?

**Decision**: Modify `MenuItem.Action` from `func()` to `func(io.Reader, io.Writer, io.Writer)` to accept stdin, stdout, stderr.

**Rationale**: 
- Maintains consistency with the existing `App` struct pattern
- Enables sub-menus to be interactive and testable with mock I/O
- Simple signature change, minimal refactoring required
- Follows Go's explicit dependency injection pattern

**Alternatives Considered**:
- Global I/O variables: Rejected due to poor testability and non-idiomatic Go
- Context-based I/O: Rejected as overly complex for this use case
- Closure capturing App's I/O: Rejected as it couples menu items to App instance

### 2. Sub-Menu Implementation Pattern

**Question**: Should the sub-menu be a separate type or reuse the existing menu pattern?

**Decision**: Implement a simple menu loop function in `lexical_elements` package that mirrors `App.Run()` logic.

**Rationale**:
- Reuses proven pattern from main menu
- Keeps code simple and beginner-friendly (Principle I)
- No need for abstraction or new types (YAGNI - Principle V)
- Easy to test with table-driven tests

**Alternatives Considered**:
- Generic Menu type: Rejected as premature abstraction (YAGNI violation)
- Recursive menu system: Rejected as overly complex for 2-level requirement
- State machine: Rejected as over-engineered for simple navigation

### 3. Topic Display Functions

**Question**: Should individual topic display functions (DisplayComments, DisplayTokens, etc.) be modified?

**Decision**: No changes needed. They already write to stdout via fmt package.

**Rationale**:
- Current functions work correctly
- No I/O redirection needed for testing (can be tested separately)
- Minimizes scope of changes (YAGNI - Principle V)

**Alternatives Considered**:
- Refactor all Display* functions to accept I/O: Rejected as unnecessary for this feature

### 4. Menu Numbering Scheme

**Question**: How should the 11 topics be numbered in the sub-menu?

**Decision**: Sequential numbering from 0-10, matching the order in the current `Display()` function.

**Rationale**:
- Matches user requirement ("从0开始，顺序增加")
- Consistent with main menu numbering (starts at 0)
- Predictable and easy to remember

**Order**:
0. Comments
1. Tokens
2. Semicolons
3. Identifiers
4. Keywords
5. Operators
6. Integers
7. Floats
8. Imaginary
9. Runes
10. Strings

### 5. Error Handling

**Question**: How should invalid input be handled in the sub-menu?

**Decision**: Display error message and re-prompt, same as main menu.

**Rationale**:
- Consistent with existing main menu behavior
- User-friendly (doesn't exit on error)
- Handles edge cases: invalid numbers, non-numeric input, whitespace

## Technical Decisions Summary

| Decision | Choice | Impact |
|----------|--------|--------|
| MenuItem signature | `func(io.Reader, io.Writer, io.Writer)` | Refactor main.go and lexical_elements.go |
| Sub-menu pattern | Simple loop function | New function in lexical_elements.go |
| Display functions | No changes | Zero impact |
| Numbering | 0-10 sequential | Documentation only |
| Error handling | Re-prompt on error | Consistent UX |

## Dependencies

**No external dependencies required**. All functionality can be implemented using:
- Go standard library (`io`, `bufio`, `fmt`, `strings`, `sort`)
- Existing project structure

## Testing Strategy

- **Unit tests**: Table-driven tests with mock I/O (bytes.Buffer)
- **Coverage target**: 80%+ (per Principle VI)
- **Test scenarios**:
  - Valid menu selections (0-10)
  - 'q' to return to main menu
  - Invalid input (out of range, non-numeric, empty)
  - Whitespace handling

## Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Breaking existing tests | Medium | High | Run full test suite after refactoring |
| I/O signature mismatch | Low | Medium | Use compiler to catch signature errors |
| Test coverage drop | Low | Medium | Write tests before implementation |

## Conclusion

All technical unknowns have been resolved. The implementation is straightforward and aligns with all constitution principles. No external research or dependencies required. Ready to proceed to Phase 1 (Design).
