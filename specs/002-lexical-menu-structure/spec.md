# Feature Specification: Lexical Menu Structure

<!--
  NOTE: Per the constitution, all user-facing documentation and code comments
  must be in Chinese.
-->

**Feature Branch**: `002-lexical-menu-structure`  
**Created**: 2025-12-03  
**Status**: Draft  
**Input**: User description: "优化首页的main.go的选择菜单，加入二级菜单。 在选择0 Lexical elements之后，为lexical_elements下的每个学功能go增加一个二级菜单选择，从0开始，顺序增加。最后添加一个q返回一级菜单"

## User Scenarios & Testing *(mandatory)*

<!--
  NOTE: Per the constitution, all features must achieve at least 80% unit test coverage.
  Ensure acceptance scenarios are comprehensive enough to meet this requirement.
-->

### User Story 1 - Access Lexical Elements Sub-menu (Priority: P1)

As a user, I want to see a detailed list of lexical element topics when I select the "Lexical elements" option, so that I can choose a specific topic to study.

**Why this priority**: This is the core functionality requested to improve navigation and usability.

**Independent Test**: Can be tested by running the application, selecting option '0', and verifying the output displays the secondary menu.

**Acceptance Scenarios**:

1. **Given** the application is running and showing the main menu, **When** I enter "0", **Then** the system should display the "Lexical Elements" secondary menu with options 0-10 (corresponding to the topics) and 'q'.
2. **Given** I am in the secondary menu, **When** I enter "q", **Then** the system should return to the main menu.

---

### User Story 2 - Execute Specific Topic (Priority: P1)

As a user, I want to select a specific topic from the secondary menu (e.g., "0. Comments") and see the output for that topic only, so that I can focus on one concept at a time.

**Why this priority**: Essential for the "study" aspect of the tool.

**Independent Test**: Can be tested by selecting a specific option in the sub-menu and checking the output.

**Acceptance Scenarios**:

1. **Given** I am in the "Lexical Elements" secondary menu, **When** I enter "0" (Comments), **Then** the system should display the content for "Comments" and then show the secondary menu again (or return to it).
2. **Given** I am in the "Lexical Elements" secondary menu, **When** I enter an invalid option, **Then** the system should display an error message and show the menu again.

### Edge Cases

- **Invalid Input**: User enters non-numeric characters or numbers outside the valid range (0-10). System should display an error and re-prompt.
- **Empty Input**: User presses Enter without typing anything. System should ignore or re-prompt.
- **Whitespace**: User enters valid numbers with leading/trailing whitespace. System should trim and accept the input.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The Main Menu System MUST support passing I/O streams to menu actions to enable interactive sub-menus.
- **FR-002**: The Lexical Elements Module MUST implement an interactive menu loop similar to the main app.
- **FR-003**: The "Lexical Elements" sub-menu MUST list the following options (order may vary but must be sequential starting from 0):
    - Comments
    - Tokens
    - Semicolons
    - Identifiers
    - Keywords
    - Operators
    - Integers
    - Floats
    - Imaginary literals
    - Runes
    - Strings
- **FR-004**: The sub-menu MUST include a 'q' option to return to the main menu.
- **FR-005**: The sub-menu MUST handle invalid input gracefully.

### Key Entities

- **MenuItem**: Updated to support interactive actions.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: User can successfully navigate to the sub-menu and back to the main menu.
- **SC-002**: User can execute each of the 11 lexical element topics individually.
- **SC-003**: Unit test coverage for the modified components MUST be at least 80%.
