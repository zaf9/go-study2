# Feature Specification: HTTP学习模式

<!--
  NOTE: Per the constitution, all user-facing documentation and code comments
  must be in Chinese.
-->

**Feature Branch**: `003-http-learning-mode`  
**Created**: 2025-12-04  
**Status**: Draft  
**Input**: User description: "为当前的命令行交互式学习内容增加通过http请求的方式进行学习同样的内容。通过不同的启动方式来进行不同的交互式学习：1. 不带参数默认启动，则是命令行交互式学习 2. 加上-d 或者--deamon参数启动则是http请求的交互式学习。无论是那种学习方式，同样的章节返回的内容应该一致"

## Clarifications

### Session 2025-12-04

- Q: When HTTP clients request chapter content, what response format should the system provide? → A: Query parameter (?format=json or ?format=html), default returns JSON
- Q: What should be the default listening address when HTTP mode starts without explicit configuration? → A: Configurable via config file, must be explicitly set (no hardcoded default)
- Q: What logging and observability approach should the HTTP service implement? → A: Structured logging with configurable levels (DEBUG/INFO/WARN/ERROR)
- Q: What file format should be used for the configuration file? → A: YAML format
- Q: When HTTP requests fail (404, 500, etc.), what should the error response contain? → A: Match request format - return error in same format as successful response (JSON/HTML)

## User Scenarios & Testing *(mandatory)*

<!--
  NOTE: Per the constitution, all features must achieve at least 80% unit test coverage.
  Ensure acceptance scenarios are comprehensive enough to meet this requirement.
-->

### User Story 1 - 命令行交互式学习（默认模式） (Priority: P1)

作为Go语言学习者，我希望能够通过命令行交互方式学习词法元素内容，以便在终端环境中快速查看和学习相关知识点。

**Why this priority**: 这是现有功能的延续，是系统的基础模式，必须保持稳定运行。作为P1优先级，它代表了最小可用产品(MVP)的核心功能。

**Independent Test**: 可以通过直接运行程序（不带任何参数）来独立测试，验证命令行菜单显示、用户输入处理和内容展示是否正常工作。

**Acceptance Scenarios**:

1. **Given** 用户在终端中，**When** 执行程序不带任何参数（如 `./go-study2`），**Then** 系统启动命令行交互模式，显示主菜单选项
2. **Given** 系统处于命令行交互模式，**When** 用户选择某个词法元素章节（如"标识符"），**Then** 系统在终端中显示该章节的完整学习内容
3. **Given** 系统显示了某个章节内容，**When** 用户选择返回或继续浏览其他章节，**Then** 系统正确响应用户操作并保持会话状态
4. **Given** 用户在命令行交互模式中，**When** 用户输入退出命令，**Then** 系统正常退出程序

---

### User Story 2 - HTTP服务模式学习 (Priority: P2)

作为Go语言学习者，我希望能够通过HTTP请求方式访问学习内容，以便在Web浏览器或通过API客户端进行学习，实现更灵活的学习方式。

**Why this priority**: 这是新增的核心功能，提供了更现代化的访问方式，支持远程访问和集成到其他系统中。作为P2优先级，它在P1基础上扩展了系统的可用性。

**Independent Test**: 可以通过使用 `-d` 或 `--daemon` 参数启动程序，然后使用curl或浏览器访问HTTP端点来独立测试，验证HTTP服务是否正确启动并能返回学习内容。

**Acceptance Scenarios**:

1. **Given** 用户在终端中，**When** 执行程序带 `-d` 参数（如 `./go-study2 -d`），**Then** 系统启动HTTP服务模式，监听指定端口并输出服务地址信息
2. **Given** HTTP服务已启动，**When** 用户通过浏览器或HTTP客户端访问根路径（如 `http://localhost:8080/`），**Then** 系统返回可用章节列表或欢迎页面
3. **Given** HTTP服务已启动，**When** 用户请求特定章节内容（如 `http://localhost:8080/lexical/identifiers`），**Then** 系统默认返回JSON格式的章节学习内容
4. **Given** HTTP服务已启动，**When** 用户请求特定章节并指定格式参数（如 `http://localhost:8080/lexical/identifiers?format=html`），**Then** 系统返回HTML格式的章节学习内容
5. **Given** HTTP服务已启动，**When** 用户请求特定章节并指定格式参数（如 `http://localhost:8080/lexical/identifiers?format=json`），**Then** 系统返回JSON格式的章节学习内容
6. **Given** HTTP服务已启动，**When** 用户通过 `--daemon` 参数启动，**Then** 系统行为与 `-d` 参数完全一致
7. **Given** HTTP服务正在运行，**When** 用户发送停止信号（如Ctrl+C），**Then** 系统优雅关闭HTTP服务并清理资源

---

### User Story 3 - 内容一致性保证 (Priority: P1)

作为Go语言学习者，无论我使用命令行模式还是HTTP模式，我期望看到的同一章节内容完全一致，以确保学习体验的连贯性。

**Why this priority**: 内容一致性是用户体验的核心要求，直接影响学习质量。作为P1优先级，它确保了两种模式下的数据完整性。

**Independent Test**: 可以通过分别在两种模式下请求相同章节，然后比较返回内容来独立测试，验证内容源是否统一且格式化逻辑一致。

**Acceptance Scenarios**:

1. **Given** 系统同时支持命令行和HTTP两种模式，**When** 在命令行模式下查看"标识符"章节，**Then** 显示的内容与HTTP模式下请求相同章节返回的内容（去除格式差异后）完全一致
2. **Given** 学习内容数据源已更新，**When** 分别通过两种模式访问更新后的章节，**Then** 两种模式都能获取到最新内容
3. **Given** 某个章节包含代码示例、说明文字和练习题，**When** 通过不同模式访问该章节，**Then** 所有内容元素（代码、文字、练习）在两种模式下都完整呈现

---

### User Story 4 - HTTP服务配置灵活性 (Priority: P3)

作为系统管理员或高级用户，我希望能够配置HTTP服务的监听端口和地址，以便适应不同的部署环境和安全要求。

**Why this priority**: 配置灵活性提升了系统的适应性，但不是核心功能。作为P3优先级，它是在基础功能稳定后的增强特性。

**Independent Test**: 可以通过使用不同的配置参数（如 `-d --port=9090`）启动服务，然后验证服务是否在指定端口监听来独立测试。

**Acceptance Scenarios**:

1. **Given** 用户需要在特定端口运行服务，**When** 使用 `-d --port=9090` 启动程序，**Then** HTTP服务在9090端口监听
2. **Given** 用户需要绑定特定IP地址,**When** 使用 `-d --host=127.0.0.1` 启动程序,**Then** HTTP服务仅在本地回环地址监听
3. **Given** 用户未在配置文件中设置监听地址和端口,**When** 使用 `-d` 启动程序,**Then** 系统提示缺少必需配置并拒绝启动
4. **Given** 配置文件中已设置监听地址和端口,**When** 使用 `-d` 启动程序,**Then** HTTP服务按配置文件中的设置启动

---

### Edge Cases

- **并发访问**: 当多个HTTP客户端同时请求不同章节内容时，系统如何处理并发请求？是否会出现资源竞争或数据不一致？
- **无效章节请求**: 当用户通过HTTP请求不存在的章节路径时，系统应返回404状态码，并根据请求格式返回JSON或HTML格式的错误信息
- **启动参数冲突**: 如果用户同时提供了多个互斥的启动参数（虽然当前只有两种模式），系统应如何处理？
- **端口占用**: 当HTTP模式启动时，如果指定端口已被占用，系统应如何提示用户并处理？
- **内容加载失败**: 如果某个章节的内容文件损坏或缺失，两种模式应如何一致地处理这种错误情况？
- **长时间运行**: HTTP服务模式下，如果服务长时间运行（如数天），系统必须通过日志监控确保不会出现内存泄漏或性能退化
- **特殊字符处理**: 章节内容中包含特殊字符（如HTML标签、JSON特殊字符）时，HTTP响应必须根据返回格式正确转义（JSON格式需转义引号和控制字符，HTML格式需转义HTML实体）
- **空内容章节**: 如果某个章节当前没有内容（占位章节），两种模式应如何一致地展示？

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 系统必须支持无参数启动，此时进入命令行交互式学习模式
- **FR-002**: 系统必须支持 `-d` 参数启动，此时进入HTTP服务模式
- **FR-003**: 系统必须支持 `--daemon` 参数启动，其行为与 `-d` 参数完全相同
- **FR-004**: 命令行模式必须保持现有的所有交互功能不变（菜单显示、用户输入、内容展示）
- **FR-005**: HTTP模式必须提供RESTful API端点，用于获取章节列表和具体章节内容
- **FR-006**: HTTP模式必须在启动时输出服务监听地址和端口信息到控制台
- **FR-007**: 系统必须使用统一的内容数据源，确保两种模式访问相同的学习内容
- **FR-008**: 系统必须为HTTP模式提供优雅关闭机制，响应系统中断信号（如SIGINT、SIGTERM）
- **FR-009**: HTTP响应必须支持通过查询参数?format指定响应格式（json或html），默认返回JSON格式，并包含适当的Content-Type头（application/json或text/html）
- **FR-010**: 系统必须对HTTP请求中的无效章节路径返回适当的HTTP错误状态码（如404），错误响应格式必须与请求格式一致（如请求?format=json则返回JSON错误，请求?format=html则返回HTML错误页面）
- **FR-011**: 系统必须在HTTP模式下处理并发请求，保证线程安全
- **FR-012**: 命令行模式和HTTP模式返回的同一章节内容（核心数据）必须完全一致
- **FR-013**: 系统必须通过配置文件提供HTTP服务监听端口配置项,该配置项为必填项,无默认值
- **FR-014**: 系统必须通过配置文件提供HTTP服务监听地址配置项,该配置项为必填项,无默认值
- **FR-015**: HTTP模式启动时，如果端口被占用，系统必须输出清晰的错误信息并退出
- **FR-016**: 系统在HTTP模式启动时,如果配置文件中缺少监听地址或端口配置,必须输出清晰的错误信息并拒绝启动
- **FR-017**: 系统必须实现结构化日志记录，支持可配置的日志级别（DEBUG、INFO、WARN、ERROR）
- **FR-018**: HTTP模式必须记录关键操作日志，包括服务启动/关闭、请求处理（请求路径、响应状态码、处理时间）、错误事件
- **FR-019**: 系统必须使用YAML格式的配置文件存储HTTP服务配置


### Key Entities *(include if feature involves data)*

- **学习章节 (LearningChapter)**: 代表一个独立的学习主题（如"标识符"、"关键字"等），包含章节ID、标题、内容文本、代码示例等属性
- **内容提供者 (ContentProvider)**: 负责从统一数据源加载和提供章节内容，被命令行模式和HTTP模式共同使用
- **HTTP路由 (HTTPRoute)**: 定义HTTP API的端点路径与章节内容的映射关系
- **服务配置 (ServiceConfig)**: 包含HTTP服务的运行参数,如监听端口、地址、超时设置、日志级别等,这些配置必须通过配置文件显式指定
- **配置文件 (ConfigFile)**: 存储系统运行时配置的YAML格式文件,包含HTTP服务的必需配置项(监听地址、端口等)和可选配置项(日志级别等)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 用户能够在不带参数启动程序后，在3秒内看到命令行交互菜单
- **SC-002**: 用户能够使用 `-d` 或 `--daemon` 参数启动HTTP服务，服务在5秒内完成启动并输出监听地址
- **SC-003**: HTTP服务能够在1秒内响应章节内容请求（对于正常大小的章节内容，<100KB）
- **SC-004**: 通过自动化测试验证，命令行模式和HTTP模式返回的同一章节内容一致性达到100%（去除格式差异）
- **SC-005**: HTTP服务能够同时处理至少50个并发请求而不出现错误或明显性能下降（响应时间增加不超过50%）
- **SC-006**: 系统在接收到中断信号后，能够在2秒内完成优雅关闭（清理资源、关闭连接）
- **SC-007**: 单元测试覆盖率达到80%以上，包括命令行模式、HTTP模式和内容提供者的核心逻辑
- **SC-008**: 用户在使用HTTP模式时，对于无效章节请求能够收到明确的错误提示（HTTP 404状态码和与请求格式一致的错误消息）
- **SC-009**: 系统能够在端口被占用时，在启动阶段检测到并向用户输出清晰的错误信息，避免静默失败
- **SC-010**: 90%的用户能够在首次尝试时成功使用两种模式访问学习内容（基于用户测试反馈）
