# Research: Go-Study2 学习闭环与测验体系

**Date**: 2025-12-11  
**Branch**: 011-learning-progress-quiz  
**Spec**: specs/011-learning-progress-quiz/spec.md

## Findings

1) Decision: 进度上报失败采用指数退避+抖动，最多 5 次，页面卸载前强制同步  
Rationale: 减少网络抖动时的风暴，并最大化数据保留，符合 FR-007。  
Alternatives considered: 固定间隔重试 3 次（失败后丢数据），前端队列批量补发（增加复杂度且当前规模不需要）。

2) Decision: API 响应统一使用 `Response{code,message,data}`，并保持 CLI/HTTP 内容源一致  
Rationale: 符合宪章 XXII-XXV 及现有后端规范，便于前端拦截器与错误处理一致。  
Alternatives considered: 按接口自定义响应结构（增加前端适配成本、违背一致性）。

3) Decision: SQLite 作为唯一存储，进度/测验写操作使用事务 + 最后更新时间戳防回退  
Rationale: 现有依赖与部署简单，数据规模小；时间戳防回退满足并发窗口需求。  
Alternatives considered: 引入外部 DB（超出当前规模与 YAGNI）；乐观锁字段（时间戳即可满足）。

4) Decision: 性能与容量预期——进度/测验 API p95 < 300ms，日活 <100 并发；题库约 400-500 题  
Rationale: 依据 PRD 章节量与单机 SQLite 能力，提前设定可验证的性能门槛。  
Alternatives considered: 不设目标（缺乏验收标准），或过高目标（与当前架构不符）。

