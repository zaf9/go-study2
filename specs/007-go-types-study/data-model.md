# Data Model: Go 类型章节学习方案

## Entities

- 类型主题 (TypeConcept)  
  - Fields: id, category(basic|numeric|string|array|slice|struct|pointer|function|interface|map|channel|type_set), title, summary, goVersion, rules[], keywords[], printableOutline  
  - Relationships: 关联多个示例 ExampleCase、测验 QuizItem、索引 ReferenceIndex  
  - Validation Rules: id、title 必填；category 需在预设枚举内；rules/keywords 非空以满足检索

- 类型规则 (TypeRule)  
  - Fields: ruleId, conceptId, ruleType(identity|comparability|recursion_limit|direction|key_constraint|type_set), description, references(specSection, goVersion), severity(info|warning|invalid)  
  - Relationships: 被 ExampleCase 与 QuizItem 引用；与 TypeConcept 多对一  
  - Validation Rules: ruleId、ruleType 必填；references.specSection 非空

- 示例/反例 (ExampleCase)  
  - Fields: id, conceptId, title, code, expectedOutput, isValid(bool), ruleRef(ruleId), notes[]  
  - Relationships: 归属 TypeConcept；引用 TypeRule 说明合规或违规原因  
  - Validation Rules: conceptId、title 必填；code/expectedOutput 不为空；ruleRef 必填以满足 FR-007

- 测验项 (QuizItem)  
  - Fields: id, conceptId, stem, options[], answer, explanation, ruleRef(ruleId), difficulty(tag)  
  - Relationships: 归属 TypeConcept；与 TypeRule 关联提供解析  
  - Validation Rules: 至少 2 个选项且唯一正确答案；explanation、ruleRef 必填；difficulty 取固定标签

- 检索索引 (ReferenceIndex)  
  - Fields: keyword, conceptId, summary, positiveExampleId, negativeExampleId, anchors(route|cliPath)  
  - Relationships: 指向 TypeConcept 与 ExampleCase 的正反例  
  - Validation Rules: keyword、conceptId 必填；anchors 至少包含 HTTP 路径或 CLI 选择编号

- 学习进度 (LearningProgress)  
  - Fields: userId, completedConcepts[], lastVisited(anchor), quizScores[{conceptId, score, timestamp}]  
  - Relationships: 记录用户与 TypeConcept 的完成度  
  - Validation Rules: userId 必填；completedConcepts 去重；score 范围 0-100

