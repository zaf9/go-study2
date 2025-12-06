---
description: "Task list for Go Lexical Elements Learning Tool"
---

# Tasks: Go Lexical Elements Learning Tool

**Input**: Design documents from `specs/001-go-learn-lexical-elements/`
**Prerequisites**: plan.md (required), spec.md (required for user stories)

**Tests**: Per the constitution, features MUST have at least 80% unit test coverage. Test tasks are **MANDATORY**.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions.

## Path Conventions

The project structure is based on the GoFrame standard layout as defined in `plan.md`, with `main.go` at the project root.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and directory structure creation.

- [x] T001 [P] Create project directory `api/`
- [x] T002 [P] Create project directory `internal/`
- [x] T003 [P] Create project directory `manifest/`
- [x] T004 [P] Create project directory `resource/`
- [x] T005 [P] Create project directory `hack/`
- [x] T006 [P] Create project directory `utility/`
- [x] T007 Initialize Go module in project root: `go mod init go-study2`
- [x] T008 Add GoFrame dependency: `go get github.com/gogf/gf/v2`

---

## Phase 2: User Story 1 - Main Menu Navigation (Priority: P1) 🎯 MVP

**Goal**: As a Go learner, I want to run the application and see a menu of topics so that I can choose what to study.

**Independent Test**: Running `go run main.go` starts the application, displays the main menu, and entering 'q' exits the program correctly.

### Implementation for User Story 1

- [x] T009 [US1] Create the main application entry point `main.go`.
- [x] T010 [US1] Implement the main menu loop in `main.go` to display options and read user input.
- [x] T011 [US1] Implement the 'quit' functionality (exit on 'q') in `main.go`.
- [x] T012 [US1] Implement handling for invalid menu input in `main.go`, which should show an error and re-display the menu.

### Tests for User Story 1 (MANDATORY) ⚠️
> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T013 [P] [US1] Create test file `main_test.go`.
- [x] T014 [US1] In `main_test.go`, write a unit test to verify that input 'q' exits the application.
- [x] T015 [US1] In `main_test.go`, write a unit test to verify invalid input displays an error message.

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently. The application starts, shows a menu, and can be quit.

---

## Phase 3: User Story 2 - Topic Exploration (Priority: P1)

**Goal**: As a Go learner, I want to select the "Lexical Elements" topic from the menu to view the code examples and explanations for its sub-topics.

**Independent Test**: From the main menu, entering '0' triggers the functions within the 'lexical_elements' package, and the formatted content for all 11 sub-topics is printed to the console.

### Implementation for User Story 2

- [x] T016 [US2] Create package directory `internal/app/lexical_elements/`.
- [x] T017 [US2] In `main.go`, add logic to handle input '0' to call a master function in the `lexical_elements` package.
- [x] T018 [P] [US2] Create content file `internal/app/lexical_elements/comments.go`.
- [x] T019 [P] [US2] Create content file `internal/app/lexical_elements/tokens.go`.
- [x] T020 [P] [US2] Create content file `internal/app/lexical_elements/semicolons.go`.
- [x] T021 [P] [US2] Create content file `internal/app/lexical_elements/identifiers.go`.
- [x] T022 [P] [US2] Create content file `internal/app/lexical_elements/keywords.go`.
- [x] T023 [P] [US2] Create content file `internal/app/lexical_elements/operators.go`.
- [x] T024 [P] [US2] Create content file `internal/app/lexical_elements/integers.go`.
- [x] T025 [P] [US2] Create content file `internal/app/lexical_elements/floats.go`.
- [x] T026 [P] [US2] Create content file `internal/app/lexical_elements/imaginary.go`.
- [x] T027 [P] [US2] Create content file `internal/app/lexical_elements/runes.go`.
- [x] T028 [P] [US2] Create content file `internal/app/lexical_elements/strings.go`.
- [x] T029 [US2] In each of the 11 content files created above, add a placeholder public function (e.g., `DisplayComments()`) that prints its topic name.
- [x] T030 [US2] Create a master `Display()` function in a new file `internal/app/lexical_elements/lexical_elements.go` that calls all 11 sub-topic functions in order.
- [x] T031 [P] [US2] Populate `internal/app/lexical_elements/comments.go` with runnable examples and detailed explanations in Chinese.
- [x] T032 [P] [US2] Populate `internal/app/lexical_elements/tokens.go` with runnable examples and detailed explanations in Chinese.
- [x] T033 [P] [US2] Populate `internal/app/lexical_elements/semicolons.go` with runnable examples and detailed explanations in Chinese.
- [x] T034 [P] [US2] Populate `internal/app/lexical_elements/identifiers.go` with runnable examples and detailed explanations in Chinese.
- [x] T035 [P] [US2] Populate `internal/app/lexical_elements/keywords.go` with runnable examples and detailed explanations in Chinese.
- [x] T036 [P] [US2] Populate `internal/app/lexical_elements/operators.go` with runnable examples and detailed explanations in Chinese.
- [x] T037 [P] [US2] Populate `internal/app/lexical_elements/integers.go` with runnable examples and detailed explanations in Chinese.
- [x] T038 [P] [US2] Populate `internal/app/lexical_elements/floats.go` with runnable examples and detailed explanations in Chinese.
- [x] T039 [P] [US2] Populate `internal/app/lexical_elements/imaginary.go` with runnable examples and detailed explanations in Chinese.
- [x] T040 [P] [US2] Populate `internal/app/lexical_elements/runes.go` with runnable examples and detailed explanations in Chinese.
- [x] T041 [P] [US2] Populate `internal/app/lexical_elements/strings.go` with runnable examples and detailed explanations in Chinese.

### Tests for User Story 2 (MANDATORY) ⚠️
- [x] T042 [P] [US2] Create test file for master function `internal/app/lexical_elements/lexical_elements_test.go`.
- [x] T043 [US2] In `lexical_elements_test.go`, write a test to ensure the master `Display()` function executes without errors.
- [x] T044 [P] [US2] Create test file `internal/app/lexical_elements/comments_test.go`.
- [x] T045 [P] [US2] Create test file `internal/app/lexical_elements/tokens_test.go`.
- [x] T046 [P] [US2] Create test file `internal/app/lexical_elements/semicolons_test.go`.
- [x] T047 [P] [US2] Create test file `internal/app/lexical_elements/identifiers_test.go`.
- [x] T048 [P] [US2] Create test file `internal/app/lexical_elements/keywords_test.go`.
- [x] T049 [P] [US2] Create test file `internal/app/lexical_elements/operators_test.go`.
- [x] T050 [P] [US2] Create test file `internal/app/lexical_elements/integers_test.go`.
- [x] T051 [P] [US2] Create test file `internal/app/lexical_elements/floats_test.go`.
- [x] T052 [P] [US2] Create test file `internal/app/lexical_elements/imaginary_test.go`.
- [x] T053 [P] [US2] Create test file `internal/app/lexical_elements/runes_test.go`.
- [x] T054 [P] [US2] Create test file `internal/app/lexical_elements/strings_test.go`.
- [x] T055 [US2] In each of the 11 specific `_test.go` files (T044-T054), add a unit test to verify its corresponding display function runs without error and produces output.
- [x] T056 [US2] In `main_test.go`, add a test to verify that input '0' correctly calls the `lexical_elements` package's `Display` function.

**Checkpoint**: At this point, User Stories 1 AND 2 should both work. The "Lexical elements" topic can be selected and displays all content.

---

## Phase 4: User Story 3 - Extensible Menu (Priority: P2)

**Goal**: As a developer, I want the `main.go` file to be easily extensible so that new learning packages can be added to the menu with minimal effort.

**Independent Test**: A new menu item and corresponding function can be added by modifying a single data structure (e.g., a map) in `main.go`.

### Implementation for User Story 3

- [x] T057 [US3] Refactor the menu selection logic in `main.go` to use a map or similar data structure instead of a large switch/if-else block.
- [x] T058 [US3] Add code comments to `main.go` clearly explaining how to add a new learning module to the menu map.

### Tests for User Story 3 (MANDATORY) ⚠️

- [x] T059 [US3] Update tests in `main_test.go` to reflect the refactored menu logic and test its scalability.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories and ensure final quality.

- [x] T060 [P] Final Review: Verify all code comments and user-facing documentation are in Chinese.
- [x] T061 Run `go fmt ./...` to format all Go source files.
- [x] T062 Run `go vet ./...` to check for suspicious constructs in the code.
- [x] T063 Run `go test -cover ./...` and verify that overall test coverage is >= 80%.

---

## Dependencies & Execution Order

- **Phase 1 (Setup)** must complete before all other phases.
- **Phase 2 (US1)** can start after Phase 1.
- **Phase 3 (US2)** can start after Phase 1. Task T017 depends on T010 from US1.
- **Phase 4 (US3)** depends on the menu logic from Phase 2 (US1).
- **Phase 5 (Polish)** should be performed after all implementation and testing is complete.

### Parallel Opportunities

- Within Phase 1, all directory creation tasks (T001-T006) can run in parallel.
- Within Phase 3 (US2), the creation of content files (T018-T028), population of content files (T031-T041), and creation of test files (T044-T054) are highly parallelizable.
- Development on US1 and US2 can occur in parallel after Phase 1, with the exception of the integration points in `main.go`.

## Implementation Strategy

1.  **MVP First**: Complete Phase 1 and Phase 2 (User Story 1) to have a runnable application skeleton.
2.  **Incremental Delivery**: Complete Phase 3 (User Story 2) to add the core learning content.
3.  **Refinement**: Complete Phase 4 (User Story 3) to improve maintainability.
4.  **Finalize**: Complete Phase 5 (Polish) to ensure code quality and adherence to all requirements.