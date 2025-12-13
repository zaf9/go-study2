package progress

// 模块说明：进度计算器负责章节状态判定与整体/主题进度的权重汇总，确保前端展示与规则一致。

import (
	"time"
	"unicode/utf8"

	progressdom "go-study2/internal/domain/progress"

	"go-study2/internal/app/constants"
	"go-study2/internal/app/lexical_elements"
	appconfig "go-study2/internal/config"
)

// Calculator 负责进度与汇总计算。
type Calculator struct {
	TopicWeights      map[string]int
	ChapterTotals     map[string]int
	EstimatedDuration map[string]int64
	// 可配置参数（从 configs/config.yaml 读取）
	charsPerSec        float64
	difficulty         map[string]float64
	minSeconds         int64
	maxSeconds         int64
	completionFraction float64
}

// NewCalculator 创建计算器，未提供的权重与总数将使用默认值。
func NewCalculator(weights map[string]int, totals map[string]int) *Calculator {
	calc := &Calculator{
		TopicWeights:      cloneIntMap(weights, defaultWeights()),
		ChapterTotals:     cloneIntMap(totals, defaultChapterTotals()),
		EstimatedDuration: map[string]int64{},
		// 默认值
		charsPerSec:        5.0,
		difficulty:         map[string]float64{"lexical_elements": 1.0, "constants": 1.1, "variables": 1.0, "types": 1.2},
		minSeconds:         60,
		maxSeconds:         3600,
		completionFraction: 0.5,
	}

	// 尝试从配置加载可配置参数（失败则使用默认值）
	if cfg, err := appconfig.Load(); err == nil && cfg != nil {
		// 如果配置中包含 progress 段则覆盖默认值
		if cfg.Progress.ReadCharsPerSec > 0 {
			calc.charsPerSec = cfg.Progress.ReadCharsPerSec
		}
		if cfg.Progress.Difficulty != nil && len(cfg.Progress.Difficulty) > 0 {
			calc.difficulty = cfg.Progress.Difficulty
		}
		if cfg.Progress.MinSeconds > 0 {
			calc.minSeconds = cfg.Progress.MinSeconds
		}
		if cfg.Progress.MaxSeconds > 0 {
			calc.maxSeconds = cfg.Progress.MaxSeconds
		}
		if cfg.Progress.CompletionFraction > 0 {
			calc.completionFraction = cfg.Progress.CompletionFraction
		}
	}

	// 预计算已知章节预计时长
	calc.initEstimatedDurations()
	return calc
}

// 初始化时预计算已知章节的预计时长，避免运行时缺失导致回退到固定值。
func (c *Calculator) initEstimatedDurations() {
	if c.EstimatedDuration == nil {
		c.EstimatedDuration = map[string]int64{}
	}

	// 词法元素章节
	lexicalMap := map[string]func() string{
		"comments":    lexical_elements.GetCommentsContent,
		"tokens":      lexical_elements.GetTokensContent,
		"semicolons":  lexical_elements.GetSemicolonsContent,
		"identifiers": lexical_elements.GetIdentifiersContent,
		"keywords":    lexical_elements.GetKeywordsContent,
		"operators":   lexical_elements.GetOperatorsContent,
		"integers":    lexical_elements.GetIntegersContent,
		"floats":      lexical_elements.GetFloatsContent,
		"imaginary":   lexical_elements.GetImaginaryContent,
		"runes":       lexical_elements.GetRunesContent,
		"strings":     lexical_elements.GetStringsContent,
	}
	for k, fn := range lexicalMap {
		content := fn()
		c.EstimatedDuration["lexical_elements/"+k] = c.estimateSecondsFromContent(content, "lexical_elements")
	}

	// Constants 章节
	constantsMap := map[string]func() string{
		"boolean":                     constants.GetBooleanContent,
		"rune":                        constants.GetRuneContent,
		"integer":                     constants.GetIntegerContent,
		"floating_point":              constants.GetFloatingPointContent,
		"complex":                     constants.GetComplexContent,
		"string":                      constants.GetStringContent,
		"expressions":                 constants.GetExpressionsContent,
		"typed_untyped":               constants.GetTypedUntypedContent,
		"conversions":                 constants.GetConversionsContent,
		"builtin_functions":           constants.GetBuiltinFunctionsContent,
		"iota":                        constants.GetIotaContent,
		"implementation_restrictions": constants.GetImplementationRestrictionsContent,
	}
	for k, fn := range constantsMap {
		content := fn()
		c.EstimatedDuration["constants/"+k] = c.estimateSecondsFromContent(content, "constants")
	}
}

// CalculateChapterStatus 基于阅读时长、滚动进度与测验结果计算章节状态。
func (c *Calculator) CalculateChapterStatus(p progressdom.LearningProgress, estimatedSeconds int64) string {
	est := estimatedSeconds
	if est <= 0 {
		est = c.lookupDuration(p.Topic, p.Chapter)
	}
	// 要求：即使滚动到 100%，也需要停留至少章节预计时间的 50% 才视为已完成
	needed := int64(float64(est) * 0.5)
	if needed <= 0 {
		needed = est
	}

	// 优先：测验通过且已阅读足够并接近全文时视为已完成（保持原有行为）
	if p.QuizPassed && p.ScrollProgress >= 90 && p.ReadDuration >= needed {
		return progressdom.StatusCompleted
	}
	// 若用户已滚动到100%并且阅读时长达到阈值，也可视为已完成（接受不测验也完成的场景）
	if p.ScrollProgress >= 100 && p.ReadDuration >= needed {
		return progressdom.StatusCompleted
	}
	if p.QuizScore > 0 && !p.QuizPassed {
		return progressdom.StatusTested
	}
	if p.ReadDuration > 0 || p.ScrollProgress > 0 || p.LastPosition != "" {
		return progressdom.StatusInProgress
	}
	return progressdom.StatusNotStarted
}

// CalculateOverallProgress 返回整体与主题维度的进度统计。
func (c *Calculator) CalculateOverallProgress(list []progressdom.LearningProgress) (OverallProgress, []TopicProgress) {
	topicAggregates := map[string]*topicAccumulator{}
	for _, item := range list {
		acc := topicAggregates[item.Topic]
		if acc == nil {
			acc = &topicAccumulator{
				Topic:  item.Topic,
				ChSet:  map[string]struct{}{},
				Weight: c.topicWeight(item.Topic),
				Total:  c.topicTotal(item.Topic),
			}
			topicAggregates[item.Topic] = acc
		}
		acc.add(item)
	}

	var topics []TopicProgress
	var weightedSum float64
	var totalWeight float64
	var totalStudyTime int64
	var completedChapters int
	var totalChapters int
	studyDays := map[string]struct{}{}

	for _, acc := range topicAggregates {
		topicProg := acc.progress()
		topics = append(topics, topicProg)
		weightedSum += float64(topicProg.Progress) * float64(acc.Weight)
		totalWeight += float64(acc.Weight)
		completedChapters += topicProg.CompletedChapters
		totalChapters += topicProg.TotalChapters
		totalStudyTime += acc.ReadDuration
		for day := range acc.StudyDays {
			studyDays[day] = struct{}{}
		}
	}

	overall := OverallProgress{
		Progress:          safeDivide(weightedSum, totalWeight),
		CompletedChapters: completedChapters,
		TotalChapters:     totalChapters,
		StudyDays:         len(studyDays),
		TotalStudyTime:    totalStudyTime,
	}
	return overall, topics
}

func (c *Calculator) lookupDuration(topic, chapter string) int64 {
	if c.EstimatedDuration == nil {
		c.EstimatedDuration = map[string]int64{}
	}
	key := topic + "/" + chapter
	if v, ok := c.EstimatedDuration[key]; ok && v > 0 {
		return v
	}
	// 尝试在首次请求时进行初始化已知章节的预计算
	c.initEstimatedDurations()
	if v, ok := c.EstimatedDuration[key]; ok && v > 0 {
		return v
	}
	// 通用回退：基于章节名长度与主题难度估算一个保守值
	diff := 1.0
	if c.difficulty != nil {
		if d, ok := c.difficulty[topic]; ok {
			diff = d
		}
	}
	// 使用章节标识长度作为保守估算基础
	length := utf8.RuneCountInString(chapter)
	est := int64(float64(length) / c.charsPerSec * diff)
	if est < c.minSeconds {
		est = c.minSeconds
	}
	if est > c.maxSeconds {
		est = c.maxSeconds
	}
	return est
}

// estimateSecondsFromContent 基于内容长度与主题难度估算阅读秒数（使用 Calculator 的配置）
func (c *Calculator) estimateSecondsFromContent(content, topic string) int64 {
	if content == "" {
		return c.minSeconds
	}
	base := float64(utf8.RuneCountInString(content)) / c.charsPerSec
	diff := 1.0
	if c.difficulty != nil {
		if d, ok := c.difficulty[topic]; ok {
			diff = d
		}
	}
	est := int64(base * diff * 1.1) // 加入 10% 额外停顿/示例时间
	if est < c.minSeconds {
		est = c.minSeconds
	}
	if est > c.maxSeconds {
		est = c.maxSeconds
	}
	return est
}

func (c *Calculator) topicWeight(topic string) int {
	if c.TopicWeights == nil {
		return 25
	}
	if w, ok := c.TopicWeights[topic]; ok && w > 0 {
		return w
	}
	return 25
}

func (c *Calculator) topicTotal(topic string) int {
	if c.ChapterTotals == nil {
		return 0
	}
	if v, ok := c.ChapterTotals[topic]; ok && v > 0 {
		return v
	}
	return 0
}

type topicAccumulator struct {
	Topic          string
	Weight         int
	Total          int
	Completed      int
	ChSet          map[string]struct{}
	ReadDuration   int64
	LatestVisit    time.Time
	StudyDays      map[string]struct{}
	CompletedAtAny bool
	ChapStatuses   map[string]string
}

func (a *topicAccumulator) add(p progressdom.LearningProgress) {
	if a.ChSet == nil {
		a.ChSet = map[string]struct{}{}
	}
	if a.StudyDays == nil {
		a.StudyDays = map[string]struct{}{}
	}
	if a.ChapStatuses == nil {
		a.ChapStatuses = map[string]string{}
	}
	key := p.Chapter
	a.ChSet[key] = struct{}{}
	a.ChapStatuses[key] = p.Status
	if p.Status == progressdom.StatusCompleted {
		a.Completed++
		a.CompletedAtAny = true
	}
	a.ReadDuration += p.ReadDuration
	if p.LastVisitAt.After(a.LatestVisit) {
		a.LatestVisit = p.LastVisitAt
	}
	day := p.LastVisitAt.Format("2006-01-02")
	if day != "" {
		a.StudyDays[day] = struct{}{}
	}
}

func (a *topicAccumulator) progress() TopicProgress {
	total := a.Total
	if total <= 0 {
		total = len(a.ChSet)
	}
	if total <= 0 {
		total = 1
	}
	progressValue := int((float64(a.Completed) / float64(total)) * 100)
	return TopicProgress{
		Name:              topicName(a.Topic),
		ID:                a.Topic,
		Weight:            a.Weight,
		Progress:          progressValue,
		CompletedChapters: a.Completed,
		TotalChapters:     total,
		LastVisitAt:       a.LatestVisit,
	}
}

func cloneIntMap(src map[string]int, fallback map[string]int) map[string]int {
	result := map[string]int{}
	for k, v := range fallback {
		result[k] = v
	}
	for k, v := range src {
		if v <= 0 {
			continue
		}
		result[k] = v
	}
	return result
}

func defaultWeights() map[string]int {
	return map[string]int{
		"lexical_elements": 25,
		"constants":        25,
		"variables":        25,
		"types":            25,
	}
}

func safeDivide(sum float64, weight float64) int {
	if weight <= 0 {
		return int(sum)
	}
	return int(sum / weight)
}
