# Feature Specification: Go-Study2 日志系统重构

<!--
  NOTE: Per the constitution, all user-facing documentation and code comments
  must be in Chinese.
-->

**Feature Branch**: `012-logging-system`  
**Created**: 2025-12-12  
**Status**: Draft  
**Input**: User description: "Go-Study2项目日志系统重构需求"

## Constitution Guardrails

- 注释与用户文档需清晰且后端全中文(Principle V/XV)。
- 方案需保持可维护性与单一职责,避免过度设计并保持浅层逻辑(Principle I/IV/VI/XVI)。
- 明确错误处理,无静默失败(Principle II)。
- 规划测试覆盖率≥80%,各包具备 *_test.go 与示例; 前端核心组件同样达标(Principle III/XXI/XXXVI)。
- 目录/职责可预测且遵循标准 Go 布局,仅根目录 main, go.mod/go.sum 完整,各包需 README 说明(Principle VIII/XVIII/XIX)。
- 依赖最小且必要(Principle IX)。
- 安全优先: 输入校验、鉴权、HTTPS、敏感信息保护(Principle VII)。
- 如涉及章节/菜单/主题,需同时支持 CLI 与 HTTP,共享内容源,菜单导航与路由/响应格式一致且显式错误(Principle XXII/XXIII/XXV)。
- Go 规范章节需按章节->子章节->子包层次组织,文件命名与示例齐备(Principle XXIV)。
- 完成后需同步更新 README 等文档(Principle XI)。

## Clarifications

### Session 2025-12-12

- Q: When the application cannot write to the log directory (permissions, disk full, etc.), what should be the fallback behavior? → A: 应用启动失败,返回明确错误信息并退出 (Fail fast - application refuses to start)
- Q: When a middleware or service layer fails to propagate the Context (and thus loses the TraceID), what should happen? → A: 生成新的 TraceID 并记录警告,说明 Context 传递中断 (Generate new TraceID, log warning about broken propagation)
- Q: How long should log files be retained before automatic deletion or archival? → A: 日志保留天数可配置,默认保留30天 (Configurable retention period, default 30 days)
- Q: Should there be access control for viewing/querying logs, or are logs accessible to all authenticated users? → A: 本期不实现应用层访问控制,依赖文件系统权限 (No application-layer access control in this phase, rely on filesystem permissions)
- Q: How should the system handle invalid or missing log configuration? → A: 配置文件缺失或无效时拒绝启动,返回详细错误 (Refuse to start if config is missing/invalid, return detailed error)

## User Scenarios & Testing *(mandatory)*

<!--
  NOTE: Per the constitution, all features must achieve at least 80% unit test coverage.
  Ensure acceptance scenarios are comprehensive enough to meet this requirement.
-->

### User Story 1 - 统一日志配置管理 (Priority: P1)

开发者需要通过配置文件统一管理项目中所有日志实例的行为,包括日志级别、输出格式、文件路径、分割策略等,而无需在代码中硬编码这些配置。

**Why this priority**: 这是日志系统的基础设施,所有其他日志功能都依赖于正确的配置。没有统一的配置管理,后续的日志记录、追踪、监控都无法正常工作。

**Independent Test**: 可以通过修改配置文件中的日志级别、输出路径等参数,然后启动应用程序,验证日志是否按照配置正确输出到指定位置和格式。

**Acceptance Scenarios**:

1. **Given** 配置文件中设置日志级别为 INFO, **When** 应用启动并记录 DEBUG 和 INFO 级别日志, **Then** 只有 INFO 及以上级别的日志被输出
2. **Given** 配置文件中设置日志输出格式为 JSON, **When** 应用记录日志, **Then** 日志以结构化 JSON 格式输出
3. **Given** 配置文件中设置日志文件路径为 `./logs/app.log`, **When** 应用启动并记录日志, **Then** 日志文件在指定路径创建
4. **Given** 配置文件中设置日志按日期分割, **When** 跨天运行应用, **Then** 自动创建新的日期命名的日志文件
5. **Given** 配置文件中定义多个日志实例(app/access/error), **When** 应用启动, **Then** 每个实例按各自配置独立工作

---

### User Story 2 - HTTP 请求全链路追踪 (Priority: P2)

开发者希望在处理 HTTP 请求时,能够通过 TraceID 追踪整个请求的生命周期,包括请求进入、业务处理、数据库操作、响应返回等各个环节的日志,便于问题排查和性能分析。

**Why this priority**: 这是生产环境问题排查的关键能力。在分布式系统或高并发场景下,没有请求追踪会导致日志混乱,无法定位问题。

**Independent Test**: 可以通过发送 HTTP 请求,然后在日志文件中搜索该请求的 TraceID,验证是否能够找到该请求从进入到返回的完整日志链路。

**Acceptance Scenarios**:

1. **Given** HTTP 服务已启动, **When** 客户端发送请求, **Then** 系统自动为该请求生成唯一 TraceID
2. **Given** 请求携带 TraceID, **When** 请求经过中间件、控制器、服务层, **Then** 所有日志都包含该 TraceID
3. **Given** 请求处理过程中发生错误, **When** 查看日志, **Then** 可以通过 TraceID 快速定位错误发生的具体位置
4. **Given** 请求完成处理, **When** 查看访问日志, **Then** 记录包含请求方法、路径、状态码、耗时、TraceID 等信息
5. **Given** 客户端请求头包含自定义 TraceID, **When** 服务处理请求, **Then** 使用客户端提供的 TraceID 而非生成新的

---

### User Story 3 - 关键操作日志埋点 (Priority: P3)

开发者需要在关键业务操作点(如学习内容加载、菜单导航、错误处理)记录结构化日志,包含操作类型、参数、结果、耗时等信息,用于业务监控和性能优化。

**Why this priority**: 这是业务可观测性的基础。虽然不如配置和追踪紧急,但对于理解系统运行状态和优化性能至关重要。

**Independent Test**: 可以通过执行特定业务操作(如访问学习内容),然后检查日志文件,验证是否记录了操作的详细信息和性能指标。

**Acceptance Scenarios**:

1. **Given** 用户请求学习内容, **When** 系统加载 Markdown 文件, **Then** 记录文件路径、加载耗时、是否成功等信息
2. **Given** 用户导航菜单, **When** 系统处理菜单选择, **Then** 记录菜单路径、选择的选项、处理结果
3. **Given** 系统发生错误, **When** 错误被捕获, **Then** 记录错误类型、错误消息、堆栈信息、上下文参数
4. **Given** 数据库查询执行, **When** 查询耗时超过阈值, **Then** 记录慢查询日志,包含 SQL、参数、耗时
5. **Given** 应用启动或关闭, **When** 生命周期事件发生, **Then** 记录启动/关闭时间、配置信息、资源状态

---

### User Story 4 - 日志查询与分析支持 (Priority: P4)

运维人员需要能够方便地查询和分析日志,包括按时间范围、日志级别、TraceID、关键字等条件过滤日志,以便快速定位问题和生成运营报告。

**Why this priority**: 这是日志系统的增强功能。虽然重要,但可以在基础日志功能完善后再实现,初期可以使用基本的文本搜索工具。

**Independent Test**: 可以通过生成一批测试日志,然后使用日志查询功能按不同条件(时间、级别、TraceID)查询,验证返回结果的准确性。

**Acceptance Scenarios**:

1. **Given** 日志文件包含多天的日志, **When** 查询指定日期范围的日志, **Then** 只返回该时间范围内的日志记录
2. **Given** 日志包含多个级别, **When** 查询 ERROR 级别日志, **Then** 只返回 ERROR 及以上级别的日志
3. **Given** 日志包含多个请求的 TraceID, **When** 查询特定 TraceID, **Then** 返回该请求的完整日志链路
4. **Given** 日志文件较大, **When** 执行查询, **Then** 查询响应时间在可接受范围内(如 5 秒内)
5. **Given** 日志为 JSON 格式, **When** 查询特定字段值, **Then** 能够精确匹配 JSON 字段内容

---

### Edge Cases

- **日志文件权限不足**: 应用启动时必须验证日志目录的写入权限。如果无权限写入,应用必须启动失败并返回明确的错误信息(包括目录路径和权限要求),不允许降级运行
- **日志文件磁盘空间不足**: 当磁盘空间耗尽时,日志写入失败是否会影响业务逻辑?如何保证业务不中断?
- **高并发日志写入**: 在高并发场景下,日志写入是否会成为性能瓶颈?是否需要异步写入或缓冲机制?
- **日志文件分割期间**: 在日志文件正在分割时(如跨天),新的日志记录如何处理?是否会丢失或写入错误文件?
- **TraceID 传递中断**: 当某个中间层未正确传递 Context 导致 TraceID 丢失时,系统必须自动生成新的 TraceID 并记录 WARNING 级别日志,说明 Context 传递链路中断的位置(文件名和行号),以便开发者修复代码
- **配置文件热更新**: 是否支持在不重启应用的情况下更新日志配置?如果支持,如何保证配置一致性?
- **日志内容敏感信息**: 如何防止敏感信息(如密码、Token)被记录到日志中?是否需要自动脱敏机制?
- **日志文件过大**: 当单个日志文件超过预期大小时,是否自动分割?分割策略如何配置?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 系统必须支持通过 YAML 配置文件管理所有日志实例的配置,包括日志级别、输出格式、文件路径、分割策略
- **FR-002**: 系统必须支持至少三种日志实例类型:应用日志(app)、访问日志(access)、错误日志(error)
- **FR-003**: 系统必须支持 JSON 结构化日志格式和纯文本格式,可通过配置切换
- **FR-004**: 系统必须支持按日期和文件大小两种日志分割策略
- **FR-005**: 系统必须为每个 HTTP 请求自动生成或传递 TraceID,并在整个请求生命周期中保持一致
- **FR-006**: 系统必须提供 HTTP 日志中间件,自动记录请求方法、路径、状态码、耗时、TraceID
- **FR-007**: 系统必须支持链式调用方式记录日志,如 `g.Log().Line().File().Cat("module").Print()`
- **FR-008**: 系统必须在关键业务操作点记录日志,包括:学习内容加载、菜单导航、错误处理、数据库操作
- **FR-009**: 系统必须支持慢查询监控,当数据库查询耗时超过配置阈值时自动记录慢查询日志
- **FR-010**: 系统必须支持异步日志写入,避免日志操作阻塞业务逻辑
- **FR-011**: 系统必须在日志写入失败时有明确的降级策略(如输出到标准错误流),不能静默失败
- **FR-012**: 系统必须支持日志级别动态过滤,只输出配置级别及以上的日志
- **FR-013**: 系统必须在每条日志中包含时间戳、日志级别、文件位置(文件名和行号)、日志内容
- **FR-014**: 系统必须支持自定义日志分类(Category),用于区分不同模块或业务领域的日志
- **FR-015**: 系统必须提供统一的日志记录接口,封装 GoFrame 的 glog 组件,便于全局使用
- **FR-016**: 系统必须在应用启动时验证日志目录的写入权限,如果验证失败则拒绝启动并返回包含目录路径和权限要求的明确错误信息
- **FR-017**: 系统必须在检测到 Context 传递中断(TraceID 丢失)时自动生成新的 TraceID,并记录 WARNING 级别日志说明中断位置
- **FR-018**: 系统必须支持可配置的日志文件保留策略,自动删除超过保留期限的日志文件,默认保留期限为 30 天
- **FR-019**: 系统必须在应用启动时验证日志配置文件的存在性和有效性,如果配置文件缺失或格式无效则拒绝启动并返回详细的错误诊断信息(包括配置文件路径、错误位置、错误原因)

### Key Entities

- **日志配置(LogConfig)**: 表示日志系统的配置信息,包括日志级别、输出路径、格式、分割策略、是否异步等属性
- **日志实例(LogInstance)**: 表示一个具体的日志记录器,如 app、access、error,每个实例有独立的配置
- **日志记录(LogRecord)**: 表示一条日志记录,包含时间戳、级别、TraceID、分类、文件位置、消息内容等属性
- **TraceID**: 表示请求追踪标识,用于关联同一请求的所有日志记录,在 HTTP 请求的 Context 中传递
- **日志中间件(LogMiddleware)**: 表示 HTTP 请求的日志拦截器,负责生成 TraceID、记录请求日志、传递 Context

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 开发者能够在 5 分钟内通过修改配置文件完成日志级别、格式、路径的调整,无需修改代码
- **SC-002**: 系统在高并发场景下(1000 并发请求)日志写入不会导致请求响应时间相比无日志场景增加超过 10% (基准测试:无日志配置下的平均响应时间作为 100% 基准)
- **SC-003**: 运维人员能够在 30 秒内通过 TraceID 定位到某个请求的完整日志链路
- **SC-004**: 日志文件能够按配置自动分割,单个日志文件大小不超过配置的阈值(如 100MB)
- **SC-005**: 所有关键业务操作(学习内容加载、菜单导航、错误处理)的日志记录覆盖率达到 100%
- **SC-006**: 慢查询日志能够准确记录所有耗时超过阈值(如 1 秒)的数据库操作
- **SC-007**: 日志系统的单元测试覆盖率达到 80% 以上
- **SC-008**: 在日志写入失败的情况下,业务逻辑不受影响,错误信息能够通过备用渠道(如标准错误流)输出
- **SC-009**: 日志查询功能能够在 5 秒内从 1GB 的日志文件中检索出指定 TraceID 的所有日志记录 (单条件精确匹配查询)
- **SC-010**: 开发团队对日志系统的易用性满意度达到 90% 以上(通过内部调研)

## Assumptions

- 假设项目已经使用 GoFrame v2.9.5 框架,无需升级或更换框架
- 假设日志文件存储在本地文件系统,不涉及远程日志收集系统(如 ELK、Loki)
- 假设日志配置文件使用 YAML 格式,与 GoFrame 的标准配置格式保持一致
- 假设日志文件的默认存储路径为 `./logs/`,可通过配置修改
- 假设日志级别默认为 INFO,可通过配置调整为 DEBUG、WARNING、ERROR 等
- 假设慢查询阈值默认为 1 秒,可通过配置调整
- 假设 TraceID 使用 UUID v4 格式生成,长度为 36 字符
- 假设日志文件按日期分割时,文件名格式为 `app-2025-12-12.log`
- 假设日志文件按大小分割时,文件名格式为 `app.log.1`, `app.log.2` 等
- 假设异步日志写入使用 GoFrame 内置的异步机制,无需引入第三方队列
- 假设敏感信息脱敏功能在本期不实现,由开发者在记录日志时自行处理
- 假设日志文件保留期限默认为 30 天,可通过配置文件调整
- 假设日志访问控制在本期不实现应用层权限管理,依赖操作系统的文件系统权限控制
- 假设日志查询功能在本期仅提供基础的文本搜索能力,不实现复杂的查询语法
