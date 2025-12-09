# Research: Go 类型章节学习方案

## Findings

- Decision: 采用 Go 1.24.5 标准库+既有 GoFrame HTTP 服务，类型章节内容静态内置（无外部存储/依赖）。  
  Rationale: 章节以规则与示例为主，静态表即可满足 CLI/HTTP/打印需求，避免引入数据库或新框架。  
  Alternatives considered: 使用全文检索库或外部存储（放弃，超出范围且增加部署复杂度）。

- Decision: 在 `src/learning/types` 内以子主题拆分文件，统一的内容/题库结构体供 CLI 与 HTTP 复用。  
  Rationale: 符合层次化章节结构（Principle XX），同时保证双模式输出一致性，便于表驱动测试。  
  Alternatives considered: 单文件集中或模式分叉存储（放弃，降低可读性且易造成内容漂移）。

- Decision: 构建轻量关键词索引（ReferenceIndex），按类型名/规则关键词映射到摘要与正反例，支持 15 秒内检索。  
  Rationale: 满足 FR-004/SC-003，静态映射足够本地检索且无外部依赖。  
  Alternatives considered: 运行时全文扫描或引入搜索服务（放弃，性能与复杂度不划算）。

- Decision: 测验采用表驱动数据结构（题干/选项/答案/解析），在内容加载时同时返回，可支持 CLI 评分与 HTTP 提交接口。  
  Rationale: 满足 FR-003 对评分与解析要求，并便于契约测试与重做逻辑。  
  Alternatives considered: 仅展示解析不做评分（放弃，不满足需求）；动态生成题目（放弃，增加不确定性）。

- Decision: 为非法递归、不可比较类型、接口类型集等高风险规则编写专门反例与测试断言。  
  Rationale: 满足 FR-007 边界清单要求，并通过单元/契约测试防止内容遗漏。  
  Alternatives considered: 仅文档描述无测试（放弃，缺乏回归保障）。

