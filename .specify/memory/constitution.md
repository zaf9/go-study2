<!--
Sync Impact Report:
- Version change: 1.0.0 -> 1.1.0
- Added sections:
  - Principle XX: Hierarchical Chapter Learning Structure
- Removed sections: none
- Modified sections: none
- Impact: 新增 Principle 规范了 Go 语言规范章节学习内容的组织结构,
  要求每个章节创建独立 package,每个子章节对应一个 .go 文件,
  支持多层次的子 package 结构,确保学习内容的层次化和规范化组织。
- Templates requiring updates: none
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

### Principle XVII: Hierarchical Menu Navigation
Interactive applications MUST support hierarchical menu structures to improve user experience and content organization.
Menu systems MUST:
* Support multi-level navigation with clear entry and exit points
* Pass I/O streams to enable interactive sub-menus
* Provide intuitive navigation controls (e.g., 'q' to return to parent menu)
* Handle invalid input gracefully with clear error messages
* Maintain consistent numbering schemes (starting from 0) across menu levels

### Principle XVIII: Dual Learning Mode Support
All new Go learning chapter specifications MUST support both learning modes:
* Command-line interactive mode (CLI) for terminal-based learning
* HTTP request mode for web-based access
This ensures consistent accessibility across different user preferences and integration scenarios.

### Principle XIX: Documentation Synchronization
After completing the development of any new specification, the `README.md` file MUST be updated to reflect the changes.
Updates MUST include relevant sections such as features, usage instructions, project structure, and roadmap status.

### Principle XX: Hierarchical Chapter Learning Structure
Go 语言规范章节的学习内容必须遵循层次化的包结构组织:
* 每个规范章节必须创建独立的 package(如 `constants`、`lexical_elements`)
* 每个子章节必须对应一个独立的 .go 文件,文件名使用小写下划线命名(如 `boolean.go`、`integer.go`)
* 当子章节包含更深层次的子章节时,必须在当前章节 package 下创建子 package,子子章节同样每个对应一个 .go 文件
* 每个章节 package 必须包含 `README.md` 说明文档和主入口文件(如 `constants.go`)
* 每个 .go 文件必须包含详细的示例代码和中文说明,便于学习者理解和实践
* 层次结构必须清晰反映 Go 语言规范的章节组织,便于学习者按规范顺序学习


## Governance

Compliance with this constitution is mandatory for all contributions. Pull requests and code reviews must verify that these principles are upheld. Any deviation requires explicit justification and approval.

**Version**: 1.1.0 | **Ratified**: 2025-12-02 | **Last Amended**: 2025-12-09