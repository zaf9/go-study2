# Tasks: Lexical Menu Structure

**Input**: Design documents from `/specs/002-lexical-menu-structure/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, quickstart.md

**Tests**: Per the constitution, features MUST have at least 80% unit test coverage. Test tasks are **MANDATORY**.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

## Path Conventions

This is a single Go project with the following structure:
- Root: `main.go`, `main_test.go`
- Package: `internal/app/lexical_elements/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Verify existing project structure and prepare for refactoring

- [X] T001 Verify Go 1.24.5 environment and dependencies in go.mod
- [X] T002 [P] Review existing main.go structure and MenuItem implementation
- [X] T003 [P] Review existing lexical_elements package structure

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core refactoring that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T004 Refactor MenuItem struct in main.go to change Action signature from `func()` to `func(io.Reader, io.Writer, io.Writer)`
- [X] T005 Update NewApp function in main.go to pass I/O streams to menu actions
- [X] T006 Update App.Run() method in main.go to call menu actions with I/O streams (app.stdin, app.stdout, app.stderr)

**Checkpoint**: Foundation ready - MenuItem now supports interactive sub-menus

---

## Phase 3: User Story 1 - Access Lexical Elements Sub-menu (Priority: P1) 🎯 MVP

**Goal**: Enable users to see a detailed list of lexical element topics when selecting "Lexical elements" option

**Independent Test**: Run the application, select option '0', verify the secondary menu displays with options 0-10 and 'q'

### Tests for User Story 1 (MANDATORY) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T007 [P] [US1] Add test case in main_test.go for navigating to sub-menu (input "0", verify sub-menu display)
- [X] T008 [P] [US1] Add test case in main_test.go for returning from sub-menu (input "0" then "q", verify return to main menu)
- [X] T009 [P] [US1] Create test file internal/app/lexical_elements/lexical_elements_test.go with sub-menu display tests

### Implementation for User Story 1

- [X] T010 [US1] Create DisplayMenu function in internal/app/lexical_elements/lexical_elements.go with signature `func(stdin io.Reader, stdout, stderr io.Writer)`
- [X] T011 [US1] Implement sub-menu map structure in DisplayMenu with 11 topics (0-10) and descriptions in Chinese
- [X] T012 [US1] Implement menu display loop in DisplayMenu showing all options 0-10 and 'q'
- [X] T013 [US1] Add input reading logic in DisplayMenu using bufio.Reader
- [X] T014 [US1] Implement 'q' handler in DisplayMenu to return to main menu
- [X] T015 [US1] Update main.go menu["0"].Action to call lexical_elements.DisplayMenu instead of Display

**Checkpoint**: At this point, User Story 1 should be fully functional - users can navigate into and out of the sub-menu

---

## Phase 4: User Story 2 - Execute Specific Topic (Priority: P1)

**Goal**: Enable users to select a specific topic from the secondary menu and see the output for that topic only

**Independent Test**: Navigate to sub-menu, select option '0' (Comments), verify Comments content displays and menu re-appears

### Tests for User Story 2 (MANDATORY) ⚠️

- [X] T016 [P] [US2] Add test case in lexical_elements_test.go for selecting topic "0" (Comments) and verifying DisplayComments is called
- [X] T017 [P] [US2] Add test case in lexical_elements_test.go for selecting topic "5" (Operators) and verifying DisplayOperators is called
- [X] T018 [P] [US2] Add test case in lexical_elements_test.go for invalid input (e.g., "99", "abc") and verifying error message
- [X] T019 [P] [US2] Add test case in lexical_elements_test.go for empty input and verifying re-prompt behavior
- [X] T020 [P] [US2] Add test case in lexical_elements_test.go for whitespace handling (e.g., " 3 ") and verifying correct topic execution

### Implementation for User Story 2

- [X] T021 [US2] Map sub-menu option "0" to DisplayComments in DisplayMenu function
- [X] T022 [US2] Map sub-menu option "1" to DisplayTokens in DisplayMenu function
- [X] T023 [US2] Map sub-menu option "2" to DisplaySemicolons in DisplayMenu function
- [X] T024 [US2] Map sub-menu option "3" to DisplayIdentifiers in DisplayMenu function
- [X] T025 [US2] Map sub-menu option "4" to DisplayKeywords in DisplayMenu function
- [X] T026 [US2] Map sub-menu option "5" to DisplayOperators in DisplayMenu function
- [X] T027 [US2] Map sub-menu option "6" to DisplayIntegers in DisplayMenu function
- [X] T028 [US2] Map sub-menu option "7" to DisplayFloats in DisplayMenu function
- [X] T029 [US2] Map sub-menu option "8" to DisplayImaginary in DisplayMenu function
- [X] T030 [US2] Map sub-menu option "9" to DisplayRunes in DisplayMenu function
- [X] T031 [US2] Map sub-menu option "10" to DisplayStrings in DisplayMenu function
- [X] T032 [US2] Implement input validation in DisplayMenu (check range 0-10, handle non-numeric)
- [X] T033 [US2] Implement error handling in DisplayMenu for invalid input with Chinese error message
- [X] T034 [US2] Implement whitespace trimming in DisplayMenu using strings.TrimSpace
- [X] T035 [US2] Add menu loop logic to re-display sub-menu after topic execution

**Checkpoint**: All user stories should now be independently functional - complete navigation and topic execution

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories and ensure code quality

- [X] T036 [P] Add comprehensive Chinese comments to MenuItem struct in main.go explaining I/O stream parameters
- [X] T037 [P] Add comprehensive Chinese comments to DisplayMenu function in lexical_elements.go explaining menu navigation flow
- [X] T038 [P] Add Chinese comments to all sub-menu mapping logic explaining topic numbering scheme
- [X] T039 [P] Run `go test ./...` and verify 80%+ test coverage for main.go
- [X] T040 [P] Run `go test ./...` and verify 80%+ test coverage for internal/app/lexical_elements/
- [X] T041 [P] Run `go fmt ./...` to format all code
- [X] T042 [P] Run `go vet ./...` to check for code issues
- [X] T043 [P] Validate quickstart.md instructions by manually testing all navigation scenarios
- [X] T044 [P] Update main.go comments to reflect new MenuItem signature for future developers
- [X] T045 Remove old Display() function from lexical_elements.go (now replaced by DisplayMenu)
- [X] T046 [P] Run full integration test: main menu → sub-menu → all 11 topics → return → quit

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User Story 1 must complete before User Story 2 (US2 depends on sub-menu structure from US1)
- **Polish (Phase 5)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - Creates sub-menu infrastructure
- **User Story 2 (P1)**: Depends on User Story 1 completion - Adds topic execution to existing sub-menu

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Sub-menu structure (US1) before topic mapping (US2)
- Input validation after basic functionality
- Comments and formatting after implementation

### Parallel Opportunities

- All Setup tasks (T001-T003) marked [P] can run in parallel
- All test tasks within a user story marked [P] can run in parallel
- All topic mapping tasks (T021-T031) marked [P] can run in parallel within US2
- All Polish tasks marked [P] can run in parallel

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: "Add test case in main_test.go for navigating to sub-menu"
Task: "Add test case in main_test.go for returning from sub-menu"
Task: "Create test file internal/app/lexical_elements/lexical_elements_test.go"
```

## Parallel Example: User Story 2

```bash
# Launch all topic mapping tasks together (T021-T031):
Task: "Map sub-menu option '0' to DisplayComments"
Task: "Map sub-menu option '1' to DisplayTokens"
Task: "Map sub-menu option '2' to DisplaySemicolons"
# ... (all 11 mappings can be done in parallel)
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T003)
2. Complete Phase 2: Foundational (T004-T006) - CRITICAL
3. Complete Phase 3: User Story 1 (T007-T015)
4. **STOP and VALIDATE**: Test sub-menu navigation independently
5. Demo: Show users can navigate into and out of sub-menu

### Full Feature Delivery

1. Complete Setup + Foundational → MenuItem refactored, ready for sub-menus
2. Add User Story 1 → Test independently → Sub-menu displays correctly
3. Add User Story 2 → Test independently → All topics executable
4. Complete Polish → Code quality, coverage, documentation verified
5. Final validation using quickstart.md

### Sequential Execution (Single Developer)

1. T001-T003: Verify environment (5 min)
2. T004-T006: Refactor MenuItem (30 min)
3. T007-T009: Write US1 tests (20 min)
4. T010-T015: Implement US1 sub-menu (45 min)
5. **Checkpoint**: Test US1 independently
6. T016-T020: Write US2 tests (25 min)
7. T021-T035: Implement US2 topic execution (60 min)
8. **Checkpoint**: Test US2 independently
9. T036-T046: Polish and validate (30 min)

**Total Estimated Time**: ~3.5 hours

---

## Notes

- [P] tasks = different files, no dependencies, can run in parallel
- [US1] and [US2] labels map tasks to specific user stories for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing (TDD approach)
- Commit after each logical group of tasks
- Stop at checkpoints to validate story independently
- All user-facing text and code comments MUST be in Chinese (Principle III)
- Maintain 80%+ test coverage (Principle VI)
- Keep code simple and beginner-friendly (Principle I)
