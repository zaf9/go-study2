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
All code comments and user-facing documentation MUST be written in Chinese to ensure consistency and clarity for the target developers.

### Principle IV: Shallow Logic
Deeply nested logic (e.g., multiple nested if-statements or loops) MUST be avoided. Prefer guard clauses, early returns, and function decomposition to maintain a flat code structure.

### Principle V: YAGNI (You Ain't Gonna Need It)
Do not implement complex design patterns or functionality prematurely. Focus on delivering the simplest solution that meets the current requirements.

### Principle VI: Comprehensive Testing
All features MUST be accompanied by unit tests. The total unit test coverage for the project MUST be maintained at 80% or higher.

### Principle VII: Single Responsibility Enforcement
Each file, function, and package MUST serve a single, clear responsibility.
Split logic when responsibilities diverge.

### Principle VIII: Predictable Project Structure
Directory layouts, naming conventions, and initialization patterns MUST be consistent and predictable throughout the project.

### Principle IX: Strict Dependency Discipline
Dependencies MUST be introduced only when absolutely necessary.
External libraries SHOULD be minimal, stable, and widely adopted within the Go ecosystem.

### Principle X: Explicit Error Handling
All errors MUST be handled explicitly.
Silent failures, ignored return values, and ambiguous error messages are prohibited.

### Principle XI: Consistent Developer Experience
The project MUST provide a consistent, beginner-friendly environment.
Setup steps, development workflow, and documentation MUST minimize confusion and friction.

## Project Standards

### Principle XII: Standard Go Project Structure
The project MUST follow a standard Go project layout.
The root directory MUST contain `go.mod` and `go.sum`.
Subdirectories MUST NOT contain `main.go`; only the root may define an executable entry.

### Principle XIII: Clear Chinese Comments
All comment content MUST be:
* Clear
* Organized
* Logical
  and MUST be written fully in Chinese.

### Principle XIV: Package-Level Documentation
Each package directory MUST contain a `README.md` describing:
* The package’s purpose and functionality
* Detailed usage instructions

### Principle XV: Code Quality Enforcement
The following tools MUST be executed regularly to ensure code quality:
* `go fmt` for formatting
* `go vet` for static analysis
* `golint` for style compliance
* `go mod tidy` for dependency maintenance

### Principle XVI: Testing Requirements
Each package MUST contain corresponding test files (`*_test.go`).
Example functions (`ExampleXxx`) MUST be provided when appropriate to demonstrate usage.


## Governance

Compliance with this constitution is mandatory for all contributions. Pull requests and code reviews must verify that these principles are upheld. Any deviation requires explicit justification and approval.

**Version**: 1.0.0 | **Ratified**: 2025-12-02 | **Last Amended**: 2025-12-02