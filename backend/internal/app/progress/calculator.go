package progress

// 模块说明：进度计算器负责章节状态判定与整体/主题进度的权重汇总，确保前端展示与规则一致。

import (
	"time"

	progressdom "go-study2/internal/domain/progress"
)

// Calculator 负责进度与汇总计算。
type Calculator struct {
	TopicWeights      map[string]int
	ChapterTotals     map[string]int
	EstimatedDuration map[string]int64
}

// NewCalculator 创建计算器，未提供的权重与总数将使用默认值。
func NewCalculator(weights map[string]int, totals map[string]int) *Calculator {
	return &Calculator{
		TopicWeights:      cloneIntMap(weights, defaultWeights()),
		ChapterTotals:     cloneIntMap(totals, defaultChapterTotals()),
		EstimatedDuration: map[string]int64{},
	}
}

// CalculateChapterStatus 基于阅读时长、滚动进度与测验结果计算章节状态。
func (c *Calculator) CalculateChapterStatus(p progressdom.LearningProgress, estimatedSeconds int64) string {
	est := estimatedSeconds
	if est <= 0 {
		est = c.lookupDuration(p.Topic, p.Chapter)
	}
	needed := int64(float64(est) * 0.8)
	if needed <= 0 {
		needed = est
	}

	if p.QuizPassed && p.ScrollProgress >= 90 && p.ReadDuration >= needed {
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
	return 600
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
