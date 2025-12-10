package types

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Topic 表示 Types 章节的子主题。
type Topic string

const (
	TopicBoolean           Topic = "boolean"
	TopicNumeric           Topic = "numeric"
	TopicString            Topic = "string"
	TopicArray             Topic = "array"
	TopicSlice             Topic = "slice"
	TopicStruct            Topic = "struct"
	TopicPointer           Topic = "pointer"
	TopicFunction          Topic = "function"
	TopicInterfaceBasic    Topic = "interface_basic"
	TopicInterfaceEmbedded Topic = "interface_embedded"
	TopicInterfaceGeneral  Topic = "interface_general"
	TopicInterfaceImpl     Topic = "interface_impl"
	TopicMap               Topic = "map"
	TopicChannel           Topic = "channel"
)

var (
	ErrUnsupportedTopic = errors.New("不支持的类型子主题")
	ErrContentNotFound  = errors.New("未找到该子主题内容")
	ErrQuizUnavailable  = errors.New("当前子主题暂无测验数据")
	ErrKeywordRequired  = errors.New("搜索关键词不能为空")
)

// TypeConcept 描述单个类型主题的元信息。
type TypeConcept struct {
	ID               string   `json:"id"`
	Category         string   `json:"category"`
	Title            string   `json:"title"`
	Summary          string   `json:"summary"`
	GoVersion        string   `json:"goVersion"`
	Rules            []string `json:"rules,omitempty"`
	Keywords         []string `json:"keywords,omitempty"`
	PrintableOutline []string `json:"printableOutline,omitempty"`
}

// TypeRule 表示类型相关的规则或约束。
type TypeRule struct {
	RuleID      string   `json:"ruleId"`
	ConceptID   string   `json:"conceptId"`
	RuleType    string   `json:"ruleType"`
	Description string   `json:"description"`
	References  []string `json:"references,omitempty"`
	Severity    string   `json:"severity,omitempty"`
}

// ExampleCase 用于演示规则的正例或反例。
type ExampleCase struct {
	ID             string   `json:"id"`
	ConceptID      string   `json:"conceptId"`
	Title          string   `json:"title"`
	Code           string   `json:"code"`
	ExpectedOutput string   `json:"expectedOutput"`
	IsValid        bool     `json:"isValid"`
	RuleRef        string   `json:"ruleRef"`
	Notes          []string `json:"notes,omitempty"`
}

// QuizItem 定义单个测验题目。
type QuizItem struct {
	ID          string   `json:"id"`
	ConceptID   string   `json:"conceptId"`
	Stem        string   `json:"stem"`
	Options     []string `json:"options"`
	Answer      string   `json:"answer"`
	Explanation string   `json:"explanation"`
	RuleRef     string   `json:"ruleRef"`
	Difficulty  string   `json:"difficulty,omitempty"`
}

// QuizAnswerFeedback 反馈单题判定。
type QuizAnswerFeedback struct {
	ID          string `json:"id"`
	Correct     bool   `json:"correct"`
	Answer      string `json:"answer"`
	Explanation string `json:"explanation"`
	RuleRef     string `json:"ruleRef"`
}

// QuizResult 汇总测验得分。
type QuizResult struct {
	Score   int                  `json:"score"`
	Total   int                  `json:"total"`
	Details []QuizAnswerFeedback `json:"details"`
}

// ReferenceIndex 用于搜索关键词到内容锚点的映射。
type ReferenceIndex struct {
	Keyword           string            `json:"keyword"`
	ConceptID         string            `json:"conceptId"`
	Summary           string            `json:"summary"`
	PositiveExampleID string            `json:"positiveExampleId,omitempty"`
	NegativeExampleID string            `json:"negativeExampleId,omitempty"`
	Anchors           map[string]string `json:"anchors,omitempty"`
}

// LearningProgress 记录学习进度与测验得分。
type LearningProgress struct {
	UserID            string                `json:"userId"`
	CompletedConcepts []string              `json:"completedConcepts"`
	LastVisited       string                `json:"lastVisited,omitempty"`
	QuizScores        map[string]QuizResult `json:"quizScores,omitempty"`
}

// TopicContent 聚合单个子主题的全部素材。
type TopicContent struct {
	Concept    TypeConcept      `json:"concept"`
	Rules      []TypeRule       `json:"rules"`
	Examples   []ExampleCase    `json:"examples"`
	QuizItems  []QuizItem       `json:"quizItems,omitempty"`
	References []ReferenceIndex `json:"references,omitempty"`
}

var (
	conceptRegistry   = map[Topic]TypeConcept{}
	ruleRegistry      = map[Topic][]TypeRule{}
	exampleRegistry   = map[Topic][]ExampleCase{}
	quizRegistry      = map[Topic][]QuizItem{}
	referenceRegistry = map[string]ReferenceIndex{}
)

// AllTopics 返回当前章节支持的所有子主题，按教学顺序排序。
func AllTopics() []Topic {
	return []Topic{
		TopicBoolean,
		TopicNumeric,
		TopicString,
		TopicArray,
		TopicSlice,
		TopicStruct,
		TopicPointer,
		TopicFunction,
		TopicInterfaceBasic,
		TopicInterfaceEmbedded,
		TopicInterfaceGeneral,
		TopicInterfaceImpl,
		TopicMap,
		TopicChannel,
	}
}

// NormalizeTopic 归一化子主题名称。
func NormalizeTopic(raw string) Topic {
	return Topic(strings.ToLower(strings.TrimSpace(raw)))
}

// IsSupportedTopic 判断是否为支持的子主题。
func IsSupportedTopic(topic Topic) bool {
	for _, t := range AllTopics() {
		if t == topic {
			return true
		}
	}
	return false
}

// RegisterContent 为子主题注册全部素材。
func RegisterContent(topic Topic, content TopicContent) error {
	normalized := NormalizeTopic(string(topic))
	if !IsSupportedTopic(normalized) {
		return ErrUnsupportedTopic
	}
	if content.Concept.ID == "" || content.Concept.Title == "" {
		return fmt.Errorf("子主题缺少必需的概念信息: %s", normalized)
	}
	conceptRegistry[normalized] = content.Concept
	if len(content.Rules) > 0 {
		ruleRegistry[normalized] = content.Rules
	}
	if len(content.Examples) > 0 {
		exampleRegistry[normalized] = content.Examples
	}
	if len(content.QuizItems) > 0 {
		quizRegistry[normalized] = content.QuizItems
	}
	for _, idx := range content.References {
		if idx.Keyword != "" {
			referenceRegistry[strings.ToLower(idx.Keyword)] = idx
		}
	}
	return nil
}

// LoadContent 返回指定子主题的聚合内容。
func LoadContent(topic Topic) (TopicContent, error) {
	normalized := NormalizeTopic(string(topic))
	if !IsSupportedTopic(normalized) {
		return TopicContent{}, ErrUnsupportedTopic
	}
	concept, ok := conceptRegistry[normalized]
	if !ok {
		return TopicContent{}, ErrContentNotFound
	}
	return TopicContent{
		Concept:    concept,
		Rules:      ruleRegistry[normalized],
		Examples:   exampleRegistry[normalized],
		QuizItems:  quizRegistry[normalized],
		References: collectReferencesByConcept(concept.ID),
	}, nil
}

// LoadQuiz 获取指定子主题的测验题目。
func LoadQuiz(topic Topic) ([]QuizItem, error) {
	normalized := NormalizeTopic(string(topic))
	if !IsSupportedTopic(normalized) {
		return nil, ErrUnsupportedTopic
	}
	items, ok := quizRegistry[normalized]
	if !ok || len(items) == 0 {
		return nil, ErrQuizUnavailable
	}
	return items, nil
}

// EvaluateQuiz 根据答案计算得分。
func EvaluateQuiz(topic Topic, answers map[string]string) (QuizResult, error) {
	if len(answers) == 0 {
		return QuizResult{}, fmt.Errorf("未提供答案")
	}
	items, err := LoadQuiz(topic)
	if err != nil {
		return QuizResult{}, err
	}
	score := 0
	var details []QuizAnswerFeedback
	for _, item := range items {
		choice := strings.TrimSpace(answers[item.ID])
		correct := strings.EqualFold(choice, strings.TrimSpace(item.Answer))
		if correct {
			score++
		}
		details = append(details, QuizAnswerFeedback{
			ID:          item.ID,
			Correct:     correct,
			Answer:      item.Answer,
			Explanation: item.Explanation,
			RuleRef:     item.RuleRef,
		})
	}
	return QuizResult{
		Score:   score,
		Total:   len(items),
		Details: details,
	}, nil
}

// SearchReferences 按关键词检索索引。
func SearchReferences(keyword string) ([]ReferenceIndex, error) {
	normalized := strings.ToLower(strings.TrimSpace(keyword))
	if normalized == "" {
		return nil, ErrKeywordRequired
	}
	var results []ReferenceIndex
	for key, idx := range referenceRegistry {
		if strings.Contains(key, normalized) {
			results = append(results, idx)
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Keyword < results[j].Keyword
	})
	return results, nil
}

// NewProgress 创建学习进度记录。
func NewProgress(userID string) *LearningProgress {
	return &LearningProgress{
		UserID:            userID,
		CompletedConcepts: []string{},
		QuizScores:        map[string]QuizResult{},
	}
}

// MarkCompleted 将子主题标记为已完成。
func (p *LearningProgress) MarkCompleted(topic Topic) {
	normalized := string(NormalizeTopic(string(topic)))
	for _, t := range p.CompletedConcepts {
		if t == normalized {
			p.LastVisited = normalized
			return
		}
	}
	p.CompletedConcepts = append(p.CompletedConcepts, normalized)
	p.LastVisited = normalized
}

// RecordQuizScore 记录测验得分。
func (p *LearningProgress) RecordQuizScore(topic Topic, result QuizResult) {
	if p.QuizScores == nil {
		p.QuizScores = map[string]QuizResult{}
	}
	normalized := string(NormalizeTopic(string(topic)))
	p.QuizScores[normalized] = result
	p.LastVisited = normalized
}

func collectReferencesByConcept(conceptID string) []ReferenceIndex {
	var items []ReferenceIndex
	for _, idx := range referenceRegistry {
		if idx.ConceptID == conceptID {
			items = append(items, idx)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Keyword < items[j].Keyword
	})
	return items
}
