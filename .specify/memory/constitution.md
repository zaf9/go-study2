<!--
Sync Impact Report:
- Version change: 1.1.0 -> 1.1.0 (无宪法内容变更,仅模板同步)
- Added sections: none
- Removed sections: none
- Modified sections: none
- Templates requiring updates:
  ✅ .specify/templates/plan-template.md
  ✅ .specify/templates/spec-template.md
  ✅ .specify/templates/tasks-template.md
  ✅ .specify/templates/checklist-template.md
  ⚠ .specify/templates/commands/ (目录不存在,暂无法核对)
- Follow-up TODOs: none
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

### Principle XXI: HTTP/CLI Implementation Consistency
* CLI 统一采用 `DisplayMenu(stdin, stdout, stderr)` 交互循环,编号从 0 递增、`q` 返回上级,使用映射将编号绑定到子主题展示/加载函数,签名与 `main.App` 的 `MenuItem.Action` 保持一致即可接入主菜单。
* CLI 与 HTTP 共用同一内容源: 新章节参考Lexical/Constants 复用各自 `GetXContent`/展示函数输出字符串; 新章节必须提供可复用的内容与测验读取接口,避免两种模式分叉。
* HTTP 路由固定为 `/api/v1/topic/{chapter}`(菜单) 与 `/api/v1/topic/{chapter}/{subtopic}`(内容),沿用 `middleware.Format` 的 `format=json|html` 协商; 菜单响应使用 `Response{code,message,data.items}`(items 结构同 `LexicalMenuItem`),内容响应使用 `Response{code,message,data}` 包装,HTML 通过 `getHtmlPage` 输出并附返回链接。
* 错误处理保持显式: 未知子主题返回 404 并给出 JSON/HTML 提示; 内容或测验缺失参考 Variables 的做法返回业务错误信息但保持响应结构稳定,不静默失败。
* 主题注册一致: 新章节需同步加入 `/api/v1/topics` 列表与 CLI 主菜单,路径标识使用与文件/目录一致的英文短名(snake_case); 若与既有约定冲突,以本原则为准。


## Governance

Compliance with this constitution is mandatory for all contributions. Pull requests and code reviews must verify that these principles are upheld. Any deviation requires explicit justification and approval.

**Version**: 1.1.0 | **Ratified**: 2025-12-02 | **Last Amended**: 2025-12-09