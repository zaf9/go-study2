# Data Model: Go-Study2 前端 UI

## Entities

- 用户 (User)  
  - Fields: id, username(唯一, 3-50), passwordHash(bcrypt), createdAt, updatedAt  
  - Relationships: 拥有多条学习进度 LearningProgress 与测验记录 QuizRecord  
  - Validation Rules: username 需唯一且不可为空；passwordHash 为 bcrypt 哈希；时间字段自动维护

- 主题 (Topic)  
  - Fields: key(enum: lexical_elements|constants|variables|types), title, summary, chapterCount, order  
  - Relationships: 包含多个章节 Chapter；对应测验集合 QuizItem  
  - Validation Rules: key 必填且来自枚举；order、chapterCount 为非负整数

- 章节 (Chapter)  
  - Fields: id(slug), topicKey, title, summary, order, contentRef(API 路径或内容源标识)  
  - Relationships: 从属于 Topic；与 LearningProgress、QuizItem 关联  
  - Validation Rules: id 与 topicKey 组合唯一；order 非负；contentRef 非空

- 学习进度 (LearningProgress)  
  - Fields: id, userId, topicKey, chapterId, status(enum: not_started|in_progress|done), lastVisit(datetime), lastPosition(scroll offset 或段落锚点)  
  - Relationships: 关联 User 与 Chapter  
  - Validation Rules: (userId, topicKey, chapterId) 唯一；status 需在枚举内；lastVisit 自动更新

- 测验题目 (QuizItem)  
  - Fields: id, topicKey, chapterId(nullable), stem, options[{id,label}], answer(single|multiple), explanation, difficulty(tag)  
  - Relationships: 归属 Topic/Chapter；被 QuizRecord.answers 引用  
  - Validation Rules: 至少 2 个选项，答案唯一或多选需与 options 对齐；explanation 必填

- 测验记录 (QuizRecord)  
  - Fields: id, userId, topicKey, chapterId(nullable for 综合), score, total, answers(JSON), createdAt  
  - Relationships: 关联 User；引用 QuizItem 结果  
  - Validation Rules: score、total 为非负整数且 score ≤ total；answers 需包含题目与选项映射

- 会话令牌 (TokenPair)  
  - Fields: accessToken(JWT, 7 天), refreshToken(HttpOnly Cookie), expiresAt, userId  
  - Relationships: 与 User 一对一（当前活跃会话）  
  - Validation Rules: expiresAt 必填；refreshToken 仅服务端可读；accessToken 需签名与过期校验

- 静态路由清单 (StaticRouteManifest)  
  - Fields: topics[{key, chapters:[chapterId]}], buildTime(datetime)  
  - Relationships: 前端构建用于生成 `generateStaticParams`；需与 Topic/Chapter 数据一致  
  - Validation Rules: topics 不为空；chapters 与后端已知章节列表同步

## Derived Views

- 学习进度总览: 聚合 LearningProgress 按 topicKey 计算完成率与最近访问章节  
- 测验历史: 按 userId、topicKey 排序 QuizRecord（倒序时间）；支持筛选 topicKey  
- 继续学习入口: 基于 LearningProgress.lastVisit 与 status != done 确定跳转章节

