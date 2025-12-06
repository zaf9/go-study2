# Implementation Plan: Lexical Menu Structure

**Branch**: `002-lexical-menu-structure` | **Date**: 2025-12-04 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-lexical-menu-structure/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Add a hierarchical menu system to the Go learning tool, enabling users to navigate into a secondary menu when selecting "Lexical elements" (option 0). The secondary menu will display 11 individual topics (Comments, Tokens, Semicolons, Identifiers, Keywords, Operators, Integers, Floats, Imaginary, Runes, Strings) numbered 0-10, with 'q' to return to the main menu. This requires refactoring the MenuItem structure to support passing I/O streams and implementing an interactive menu loop in the lexical_elements package.

## Technical Context

**Language/Version**: Go 1.24.5  
**Primary Dependencies**: GoFrame v2.9.5 (already in project, minimal usage)  
**Storage**: N/A (static content only)  
**Testing**: Go standard testing (`go test`)  
**Target Platform**: Cross-platform CLI (Windows/Linux/macOS)
**Project Type**: Single project (CLI application)  
**Performance Goals**: Instant menu response (<10ms)  
**Constraints**: Must maintain 80%+ test coverage, beginner-friendly code  
**Scale/Scope**: 11 lexical element topics, 2-level menu hierarchy

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Principle I (Simplicity):** ✅ PASS - Refactoring MenuItem to accept I/O streams is straightforward. Menu loop pattern already exists in main.go and can be replicated.
- **Principle II (Comments):** ✅ PASS - Will add clear Chinese comments explaining menu navigation responsibilities.
- **Principle III (Language):** ✅ PASS - All new comments and documentation will be in Chinese.
- **Principle IV (Nesting):** ✅ PASS - Menu logic uses flat switch/map patterns with early returns, no deep nesting.
- **Principle V (YAGNI):** ✅ PASS - Only implementing the minimal changes needed: I/O passing and a simple sub-menu loop.
- **Principle VI (Testing):** ✅ PASS - Existing test infrastructure supports table-driven tests with mock I/O. Can achieve 80%+ coverage.
- **Principle XVII (Hierarchical Menu):** ✅ PASS - This feature directly implements the hierarchical menu navigation principle.

## Project Structure

### Documentation (this feature)

```text
specs/002-lexical-menu-structure/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (minimal - no external research needed)
├── data-model.md        # Phase 1 output (minimal - menu structure only)
├── quickstart.md        # Phase 1 output (user guide for navigation)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
# Single project structure (CLI application)
main.go                                    # Modified: MenuItem struct, NewApp signature
main_test.go                               # Modified: Add sub-menu navigation tests

internal/app/lexical_elements/
├── lexical_elements.go                    # Modified: Replace Display() with DisplayMenu()
├── lexical_elements_test.go               # Modified: Test menu navigation
├── comments.go                            # Unchanged
├── tokens.go                              # Unchanged
├── semicolons.go                          # Unchanged
├── identifiers.go                         # Unchanged
├── keywords.go                            # Unchanged
├── operators.go                           # Unchanged
├── integers.go                            # Unchanged
├── floats.go                              # Unchanged
├── imaginary.go                           # Unchanged
├── runes.go                               # Unchanged
└── strings.go                             # Unchanged

go.mod                                     # Unchanged
go.sum                                     # Unchanged
```

**Structure Decision**: Using the existing single project structure. All changes are confined to `main.go` and the `internal/app/lexical_elements/` package. No new directories or packages needed.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations detected. All principles are satisfied by the proposed design.

---

## Phase 0: Research - COMPLETE ✅

**Output**: [research.md](./research.md)

**Key Decisions**:
1. MenuItem signature changed to `func(io.Reader, io.Writer, io.Writer)`
2. Sub-menu implemented as simple loop function mirroring main menu pattern
3. Topic numbering: 0-10 sequential
4. No external dependencies required

**Status**: All technical unknowns resolved. No NEEDS CLARIFICATION items remain.

---

## Phase 1: Design & Contracts - COMPLETE ✅

**Outputs**:
- [data-model.md](./data-model.md) - Menu structure and state transitions
- [quickstart.md](./quickstart.md) - User guide in Chinese
- Agent context updated (GEMINI.md)

**Key Design Elements**:
1. **MenuItem Entity**: Modified to support I/O streams
2. **Menu Hierarchy**: Two-level structure (main → lexical elements)
3. **State Machines**: Documented navigation flow
4. **Validation Rules**: Input handling and error cases

**Status**: Design complete. Ready for task generation.

---

## Post-Design Constitution Re-Check ✅

*Re-evaluating all principles after completing Phase 1 design*

- **Principle I (Simplicity):** ✅ PASS - Design maintains simplicity. No complex abstractions introduced.
- **Principle II (Comments):** ✅ PASS - Documentation plan includes Chinese comments for all new code.
- **Principle III (Language):** ✅ PASS - quickstart.md written in Chinese, code comments will be Chinese.
- **Principle IV (Nesting):** ✅ PASS - State machines show flat control flow, no deep nesting.
- **Principle V (YAGNI):** ✅ PASS - No premature abstractions. Simple map-based menu, no generic types.
- **Principle VI (Testing):** ✅ PASS - Test strategy documented in research.md, table-driven tests planned.
- **Principle VII (Single Responsibility):** ✅ PASS - Clear separation: main.go handles main menu, lexical_elements.go handles sub-menu.
- **Principle VIII (Predictable Structure):** ✅ PASS - No new directories, follows existing patterns.
- **Principle IX (Dependencies):** ✅ PASS - Zero new dependencies, uses only Go stdlib.
- **Principle X (Error Handling):** ✅ PASS - Error handling documented in data-model.md validation rules.
- **Principle XI (Developer Experience):** ✅ PASS - quickstart.md provides clear user guidance.
- **Principle XVII (Hierarchical Menu):** ✅ PASS - Fully implements hierarchical navigation requirement.

**Final Assessment**: All constitution principles satisfied. Design approved for implementation.

---

## Next Steps

Run `/speckit.tasks` to generate the implementation task breakdown.
