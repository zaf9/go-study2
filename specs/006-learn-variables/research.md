# Research: Variables章节学习完成

## Findings

- Decision: 采用 Go 1.22 标准库实现内容与示例，保持零依赖交付。
  Rationale: 章节教学以概念与示例为主，标准库已足够，无需增加学习负担。
  Alternatives considered: 引入第三方CLI/HTTP框架；增加动态内容存储（放弃，超出范围且增加复杂度）。

- Decision: 同时提供 CLI 与 HTTP 两种访问模式，接口契约保持一致的内容与测验数据。
  Rationale: 遵循双学习模式原则，兼顾终端与网页访问；复用同一数据源降低维护成本。
  Alternatives considered: 单一 CLI 或单一 HTTP（放弃，违背双模式要求）。

- Decision: 按子主题拆分文件（变量基础、静态类型、动态类型、零值），并在包级 README 汇总。
  Rationale: 符合层次化章节结构，便于初学者按主题学习与检索。
  Alternatives considered: 单文件集中（放弃，可读性下降；不利于示例与测试的聚合）。

- Decision: 表驱动测试覆盖核心示例与测验题目，合约测试覆盖 CLI/HTTP 输出结构。
  Rationale: 确保>=80%覆盖率且验证双模式一致性。
  Alternatives considered: 仅手动验证或样例运行（放弃，无法满足测试原则）。