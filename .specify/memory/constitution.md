<!--
Sync Impact Report:
- Version change: 2.0.0 -> 2.1.0 (新增多项治理原则以确保实现一致性)
- Added sections: 
  * Principle XIII: Fail-Fast Startup Strategy
  * Principle XIV: Observability and Structured Logging
  * Principle XV: Configuration Management
  * Principle XVI: Asynchronous and Concurrent Processing
  * Principle XVII: Graceful Shutdown
  * Principle XXXI: Standardized Learning Content Structure
  * Principle XXXII: Quiz and Assessment Standards
  * Principle XXXIII: Learning Progress Tracking
  * Principle XXXIV: Feature Independence and MVP Thinking
- Removed sections: none
- Modified sections: All subsequent principles renumbered (XIII->XVIII, XIV->XIX, etc.)
- Templates requiring updates:
  ✅ .specify/templates/plan-template.md
  ✅ .specify/templates/spec-template.md
  ✅ .specify/templates/tasks-template.md
  ✅ .specify/templates/checklist-template.md
  ⚠ .specify/templates/commands/ (目录缺失, 待确认是否需要补充)
- Follow-up TODOs: 
  * 更新所有现有 spec 文档以引用新的原则编号
  * 在下一个新 spec 开发时验证新原则的实践效果
  * 如需 commands 模板请创建目录并补齐相应文件
-->

# go-study2 Constitution

## Core Principles

### Principle I: Code Quality and Maintainability
Code must be clear, concise, and maintainable. Prioritize readability over premature optimization. Each module must have a single responsibility and be easy to understand and test.

### Principle II: Explicit Error Handling
All errors must be explicitly handled. Silent failures, ignored return values, and vague error messages are prohibited. Both frontend and backend must implement comprehensive error handling mechanisms.

### Principle III: Comprehensive Testing
All features must be accompanied by tests. Backend unit test coverage must be ≥80%, and frontend core business component test coverage must be ≥80%.

### Principle IV: Single Responsibility Principle
Each file, function, component, and package must bear a single, clear responsibility. Logic must be split when responsibilities diverge.

### Principle V: Consistent Documentation Standards
All code comments and documentation must be clear, organized, and logical. Backend must use Chinese; frontend may use Chinese or English but must maintain consistency within the project.

### Principle VI: YAGNI (You Ain't Gonna Need It)
Do not prematurely implement complex design patterns or features. Focus on providing the simplest solution that meets current requirements.

### Principle VII: Security First
Security is mandatory. Sensitive information must not be hardcoded, all inputs must be validated, API communication must use HTTPS, and authentication and authorization must be strictly enforced.

### Principle VIII: Predictable Project Structure
Directory layout, naming conventions, and initialization patterns must remain consistent and predictable throughout the project.

### Principle IX: Strict Dependency Discipline
Dependencies must only be introduced when absolutely necessary. External libraries should be minimal, stable, and widely adopted in their respective ecosystems.

### Principle X: Performance Optimization
Performance impacts must be considered. Backend avoids unnecessary database queries and memory allocations, frontend avoids unnecessary re-renders and large bundles.

### Principle XI: Documentation Synchronization
After completing development of any new specification, the `README.md` file must be updated to reflect the changes. Updates must include relevant sections such as features, usage instructions, project structure, and roadmap status.

### Principle XII: Git Workflow Standards
* Commit specification: Follow Conventional Commits specification
  - Format: `<type>(<scope>): <subject>`
  - Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`
  - Example: `feat(auth): add login page with form validation`

### Principle XIII: Fail-Fast Startup Strategy
System must validate critical resources and configurations at startup and refuse to start with clear error messages when issues are detected. This includes:
* Configuration file missing or invalid: refuse to start with detailed error diagnostics (file path, error location, error reason)
* Log directory write permissions insufficient: refuse to start with error message containing directory path and permission requirements
* Required configuration items missing: refuse to start with explanation of missing items
* Port already in use: detect and output clear error message before starting
* Avoid degraded silent operation - explicit failure is better than unclear behavior

### Principle XIV: Observability and Structured Logging
System must implement comprehensive logging infrastructure to support troubleshooting and performance analysis:
* Structured logging with configurable levels (DEBUG/INFO/WARN/ERROR) using YAML configuration
* Separate log instances for different concerns (app/access/error/slow query)
* Full-chain request tracing with TraceID propagation through Context
* Automatic TraceID generation when missing, with WARNING log indicating broken propagation chain
* Asynchronous log writing to avoid blocking business logic
* Log file retention policy (configurable, default 30 days) with automatic cleanup
* Log rotation by date and size (configurable, default max 100MB per file)
* Fallback strategy for log write failures (e.g., output to stderr), never silent failure
* Key operation logging points: content loading, menu navigation, error handling, database operations, slow queries

### Principle XV: Configuration Management
All system runtime parameters must be managed through configuration files with explicit validation:
* Use YAML format for configuration files, consistent with framework standards
* Required configuration items must have no default values (explicit configuration required)
* Support environment-specific configurations (dev/test/prod)
* Configuration validation at startup before resource initialization
* Environment variable support for deployment flexibility
* Configuration changes require restart (no hot-reload complexity unless explicitly needed)

### Principle XVI: Asynchronous and Concurrent Processing
System must handle asynchronous operations and concurrent requests appropriately:
* Asynchronous logging to prevent blocking business logic
* Thread-safe concurrent request handling for HTTP services
* Proper use of goroutines with error handling and panic recovery
* Resource pooling for concurrent operations (database connections, etc.)
* Performance benchmarks for concurrent scenarios (e.g., 1000 concurrent requests)

### Principle XVII: Graceful Shutdown
Applications must implement graceful shutdown mechanisms:
* Respond to system interrupt signals (SIGINT, SIGTERM)
* Complete in-flight requests before shutdown
* Clean up resources and close connections properly
* Timeout control for shutdown operations (max 30 seconds)
* Log shutdown events and final system state

---

## Backend Principles

### Backend Core Principles

#### Principle XVIII: Simplicity for Go Beginners
Code must be concise and easy to understand, suitable for Go beginners. The primary goal is readability and maintainability.

#### Principle XIX: Clear Layered Comments
Each logical layer of the application (such as controllers, services, repositories) must have clear comments explaining its specific responsibilities and functionality.

#### Principle XX: Backend Chinese Documentation Requirement
All code comments and user-facing documentation must be written in Chinese to ensure consistency and clarity for the target developers.

#### Principle XXI: Shallow Logic
Deeply nested logic (such as multiple nested if statements or loops) must be avoided. Prioritize guard clauses, early returns, and function decomposition to maintain a flat code structure.

#### Principle XXII: Consistent Developer Experience
The project must provide a consistent, beginner-friendly environment. Setup steps, development workflows, and documentation must minimize confusion and friction.

### Backend Project Standards

#### Principle XXIII: Standard Go Project Structure
The project must follow standard Go project layout. The root directory must contain `go.mod` and `go.sum`. Subdirectories must not contain `main.go`; only the root directory may define executable entry points.

#### Principle XXIV: Package-Level Documentation
Each package directory must contain a `README.md` describing:
* The package's purpose and functionality
* Detailed usage instructions

#### Principle XXV: Backend Code Quality Tools
The following tools must be executed regularly to ensure code quality:
* `go fmt` for formatting
* `go vet` for static analysis
* `golint` for style checking
* `go mod tidy` for dependency maintenance

#### Principle XXVI: Backend Testing Requirements
Each package must contain corresponding test files (`*_test.go`). Example functions (`ExampleXxx`) must be provided when appropriate to demonstrate usage.

#### Principle XXVII: Hierarchical Menu Navigation
Interactive applications must support hierarchical menu structures to improve user experience and content organization. The menu system must:
* Support multi-level navigation with clear entry and exit points
* Pass I/O streams to enable interactive submenus
* Provide intuitive navigation controls (such as 'q' to return to parent menu)
* Gracefully handle invalid input and display clear error messages
* Maintain consistent numbering schemes (starting from 0) between menu levels

#### Principle XXVIII: Dual Learning Mode Support
All new Go learning chapter specifications must support two learning modes:
* Command-line interactive mode (CLI) for terminal-based learning
* HTTP request mode for web-based access
This ensures consistent accessibility across different user preferences and integration scenarios.

#### Principle XXIX: Hierarchical Chapter Learning Structure
Learning content for Go language specification chapters must follow a hierarchical package structure:
* Each specification chapter must create an independent package (e.g., `constants`, `lexical_elements`)
* Each subsection must correspond to an independent .go file, named in lowercase with underscores (e.g., `boolean.go`, `integer.go`)
* When subsections contain deeper sub-subsections, a sub-package must be created under the current chapter package, with each sub-subsection also corresponding to a .go file
* Each chapter package must contain a `README.md` documentation file and a main entry file (e.g., `constants.go`)
* Each .go file must contain detailed example code and Chinese explanations to facilitate learner understanding and practice
* The hierarchical structure must clearly reflect the chapter organization of the Go language specification, facilitating learners to study in specification order

#### Principle XXX: HTTP/CLI Implementation Consistency
* CLI uniformly adopts `DisplayMenu(stdin, stdout, stderr)` interactive loop, numbering starts from 0 and increments, 'q' returns to parent level, uses mapping to bind numbers to subtopic display/loading functions, signature consistent with `main.App`'s `MenuItem.Action` can be integrated into the main menu
* CLI and HTTP share the same content source: new chapters refer to Lexical/Constants to reuse their respective `GetXContent`/display function output strings; new chapters must provide reusable content and quiz reading interfaces to avoid bifurcation between the two modes
* HTTP routes are fixed as `/api/v1/topic/{chapter}` (menu) and `/api/v1/topic/{chapter}/{subtopic}` (content), following `middleware.Format`'s `format=json|html` negotiation; menu responses use `Response{code,message,data.items}` (items structure same as `LexicalMenuItem`), content responses use `Response{code,message,data}` wrapping, HTML output through `getHtmlPage` with return link attached
* Error handling remains explicit: unknown subtopics return 404 with JSON/HTML prompts; missing content or quizzes refer to Variables' approach to return business error messages while maintaining stable response structure, no silent failures
* Topic registration consistency: new chapters need to be synchronously added to `/api/v1/topics` list and CLI main menu, path identifiers use English short names (snake_case) consistent with files/directories; if conflicts with existing conventions arise, this principle takes precedence

#### Principle XXXI: Standardized Learning Content Structure
All learning chapter content must follow a unified, predictable structure to ensure consistent learning experience:
* Each chapter must contain: Overview, Key Points, Detailed Explanation, Code Examples, Common Pitfalls (optional), Practical Recommendations
* Each chapter must include at least 2 runnable code examples (15-40 lines each)
* Code examples must include: title, inline Chinese comments, expected output
* Examples must follow official Go code style and be verified to compile and run correctly
* Content must be suitable for Go beginners with clear Chinese explanations

#### Principle XXXII: Quiz and Assessment Standards
Learning chapters must provide comprehensive assessment mechanisms:
* Quiz question count based on chapter complexity: basic chapters 5-8 questions, intermediate 8-12 questions, complex 12-15 questions
* Question types must include: single choice, multiple choice, true/false, code output prediction, error correction
* Pass threshold standardized at 60% correct rate
* Each question must provide detailed explanation referencing relevant specification rules
* Instant grading and feedback upon submission
* Quiz history tracking for progress monitoring
* Support for retaking quizzes to reinforce learning

#### Principle XXXIII: Learning Progress Tracking
System must track and persist learner progress to support continuous learning:
* Chapter status management: not_started, in_progress, completed, tested
* Progress metrics: read_duration (seconds, cumulative), scroll_progress (0-100%), last_position (pixels), quiz_score, quiz_passed
* Completion criteria: reading time ≥80% of estimated, scroll progress ≥90%, quiz passed (≥60% correct)
* Progress auto-save strategy: debounced updates every 10 seconds during reading, exponential backoff retry on failure (max 5 attempts), final sync before page unload
* Automatic scroll position restoration on return to chapter
* Overall progress calculation with topic-weighted formula
* Support for "Continue Learning" feature to resume at first incomplete chapter

#### Principle XXXIV: Feature Independence and MVP Thinking
New features must be designed with independent testability and incremental delivery:
* Each User Story must be independently testable and deliverable as a viable MVP
* Use P1/P2/P3 priority levels where P1 represents minimum viable functionality
* Features should be implementable in isolation without requiring other features to be complete
* Each feature module must have independent verification criteria
* Support for incremental rollout and testing

---

## Frontend Principles

### Frontend Core Principles

#### Principle XXXV: Type Safety First
Prioritize TypeScript type definitions to ensure type safety. Props, State, and API responses must all have explicit type definitions.

#### Principle XXXVI: Performance Optimization for Static Export
Leverage Next.js static export features to ensure optimal page loading performance. Reasonably use React.memo, useMemo, and useCallback to avoid unnecessary re-renders.

#### Principle XXXVII: Consistent UI/UX Standards
Strictly follow Ant Design design specifications to maintain interface consistency. Custom styles must be based on Ant Design's theme system.

#### Principle XXXVIII: Accessibility Standards
Follow WCAG 2.1 AA standards to ensure the application is accessible to all users. Use semantic HTML and appropriate ARIA attributes.

#### Principle XXXIX: Client-Side Security
All API calls must undergo appropriate error handling and data validation. Sensitive information must not be stored on the client side. Use HttpOnly Cookies to store authentication tokens. Guard against XSS attacks.

### Frontend Project Standards

#### Principle XL: Component Organization
* Use function components and React Hooks, avoid using Class components
* Component file structure: `components/[ComponentName]/index.tsx`
* Each component directory may contain: `index.tsx`, `styles.module.css`, `types.ts`
* Page components are placed in `app/` or `pages/` directory

#### Principle XLI: State Management Standards
* Local state prioritizes `useState`
* Cross-component shared state uses `useContext` + `useReducer`
* Avoid overuse of global state, maintain state proximity principle
* Context should be divided by functional domain (e.g., AuthContext, ThemeContext)

#### Principle XLII: API Integration Standards
* All Axios requests must use unified instance configuration
* Create `lib/api.ts` or `services/` directory to uniformly manage API calls
* Must implement request/response interceptors to handle common logic (token, error handling)
* API errors must have user-friendly prompts (using Ant Design's message or notification)

#### Principle XLIII: Styling Standards
* Prioritize CSS Modules for component-level style isolation
* Global styles only for reset, theme variables, and common utility classes
* Follow BEM naming convention (can be simplified in CSS Modules)
* Responsive design uses Ant Design's grid system
* Configure theme globally through ConfigProvider
* Do not directly modify internal styles of Ant Design components

#### Principle XLIV: Static Export Optimization
* Ensure all pages support static export, avoid using server-only features
* Use `next/image` component to optimize image loading
* Reasonably use `next/link` to implement client-side routing
* Avoid using `getServerSideProps` in components
* Use dynamic import `next/dynamic` to lazy load heavy components
* Import third-party libraries on-demand (such as Ant Design's tree-shaking)
* Highlight.js only imports needed language packages, avoid full import

#### Principle XLV: Frontend Testing Requirements
* Key business components must have unit tests (recommended Jest + React Testing Library)
* API service layer must have integration tests
* Test coverage target: core functionality > 80%
* Test user behavior rather than implementation details
* Use `data-testid` rather than CSS selectors to locate elements
* Mock external dependencies (API calls, third-party libraries)

#### Principle XLVI: Pull Request Requirements
* Must pass CI/CD checks (lint, test, build)

#### Principle XLVII: Code Review Checklist
- [ ] Code complies with ESLint and Prettier rules
- [ ] No remaining debug code such as console.log
- [ ] Component reusability and single responsibility
- [ ] Complete error handling (try-catch, error boundaries)
- [ ] Performance considerations (avoid unnecessary re-renders)
- [ ] Complete and accurate type definitions
- [ ] Styles adapted for mobile
- [ ] No hardcoded strings (use internationalization or constants)

#### Principle XLVIII: Frontend Code Quality Tools
The following tools must be executed regularly to ensure code quality:
* ESLint: Enforce code style and best practices
* Prettier: Unify code format
* TypeScript: Type checking
* Husky + lint-staged: Git hooks pre-commit checks

#### Principle XLIX: Frontend Documentation Requirements
* `README.md`: Project overview, installation, running, deployment guide
* `CONTRIBUTING.md`: Contribution guide and development process
* Components with complex logic must have JSDoc comments
* API service functions must comment parameters and return values
* Complex business logic needs inline comments for explanation
* Comments explain "why" rather than "what"

---

## Frontend-Backend Integration Standards

### Principle L: API Contract Consistency
Frontend and backend must use unified API contracts. Backend HTTP response format must be consistent with the format expected by frontend Axios interceptors.

### Principle LI: Development Environment Separation
Separate frontend and backend during development:
* Backend runs on independent port (e.g., `:8080`)
* Frontend development server runs on independent port (e.g., `:3000`)
* Frontend accesses backend API through proxy or CORS configuration

### Principle LII: Production Deployment Integration
Merged deployment in production environment:
* Next.js static export files (`out/`) must be hosted by the backend server
* Backend must configure static file service routes
* API routes must be clearly separated from static file routes (e.g., `/api/*` vs `/*`)
* Frontend build artifacts must be automatically integrated into the backend deployment package through CI/CD

### Principle LIII: Frontend-Backend Configuration Coordination
* Environment variable configuration must remain consistent between frontend and backend
* API base URL, ports, and other configurations must be configurable through environment variables
* Development/production environment configurations must be clearly separated

---

## Governance

Adherence to this constitution is mandatory for all contributions. Pull requests and code reviews must verify that these principles are maintained. Any deviations require explicit justification and approval.

**Version**: 2.1.0 | **Ratified**: 2025-12-10 | **Last Amended**: 2025-12-14