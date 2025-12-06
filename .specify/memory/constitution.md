<!--
Sync Impact Report:
- Version change: none -> 1.0.0
- Added sections:
  - Principle I: Simplicity and Clarity
  - Principle II: Clear Layered Comments
  - Principle III: Chinese Language Documentation
  - Principle IV: Shallow Logic
  - Principle V: YAGNI (You Ain't Gonna Need It)
  - Principle VI: Comprehensive Testing
- Removed sections:
  - [PRINCIPLE_1_NAME] to [PRINCIPLE_5_NAME]
  - [SECTION_2_NAME]
  - [SECTION_3_NAME]
- Templates requiring updates:
  - ✅ .specify/templates/plan-template.md
  - ✅ .specify/templates/spec-template.md
  - ✅ .specify/templates/tasks-template.md
  - ✅ .specify/templates/checklist-template.md
-->
# go-study2 Constitution

## Core Principles

### Principle I: Simplicity and Clarity
Code MUST be concise and simple, suitable for Go beginners. The primary goal is readability and maintainability over clever or overly optimized solutions.

### Principle II: Clear Layered Comments
Every logical layer of the application (e.g., controllers, services, repositories) MUST have clear comments explaining its specific responsibilities and function.

### Principle III: Chinese Language Documentation
All code comments, variable names, and user-facing documentation MUST be written in Chinese to ensure consistency and clarity for the target developers.

### Principle IV: Shallow Logic
Deeply nested logic (e.g., multiple nested if-statements or loops) MUST be avoided. Prefer guard clauses, early returns, and function decomposition to maintain a flat code structure.

### Principle V: YAGNI (You Ain't Gonna Need It)
Do not implement complex design patterns or functionality prematurely. Focus on delivering the simplest solution that meets the current requirements.

### Principle VI: Comprehensive Testing
All features MUST be accompanied by unit tests. The total unit test coverage for the project MUST be maintained at 80% or higher.

## Governance

Compliance with this constitution is mandatory for all contributions. Pull requests and code reviews must verify that these principles are upheld. Any deviation requires explicit justification and approval.

**Version**: 1.0.0 | **Ratified**: 2025-12-02 | **Last Amended**: 2025-12-02