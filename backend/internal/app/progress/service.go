package progress

// 模块说明：进度服务负责写入/查询学习进度、汇总整体统计并提供继续学习指引，所有注释保持中文便于团队理解。

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	progressdom "go-study2/internal/domain/progress"
)

// UpdateProgressRequest 表示进度更新请求。
type UpdateProgressRequest struct {
	UserID         int64
	Topic          string
	Chapter        string
	ReadDuration   int64
	ScrollProgress int
	LastPosition   string
	QuizScore      int
	QuizPassed     bool
	EstimatedSec   int64
	ForceSync      bool
}

// ProgressResponse 表示进度写入后的结果。
type ProgressResponse struct {
	Status         string          `json:"status"`
	Overall        OverallProgress `json:"overall"`
	Topic          TopicProgress   `json:"topic"`
	ReadDuration   int64           `json:"read_duration"`
	ScrollProgress int             `json:"scroll_progress"`
	LastPosition   string          `json:"last_position"`
}

// OverallProgress 汇总整体进度。
type OverallProgress struct {
	Progress          int   `json:"progress"`
	CompletedChapters int   `json:"completedChapters"`
	TotalChapters     int   `json:"totalChapters"`
	StudyDays         int   `json:"studyDays"`
	TotalStudyTime    int64 `json:"totalStudyTime"`
}

// TopicProgress 汇总单个主题的进度。
type TopicProgress struct {
	Name              string    `json:"name"`
	ID                string    `json:"id"`
	Weight            int       `json:"weight"`
	Progress          int       `json:"progress"`
	CompletedChapters int       `json:"completedChapters"`
	TotalChapters     int       `json:"totalChapters"`
	LastVisitAt       time.Time `json:"lastVisitAt"`
}

// NextChapter 表示推荐继续学习的章节。
type NextChapter struct {
	Topic    string `json:"topic"`
	Chapter  string `json:"chapter"`
	Status   string `json:"status"`
	Progress int    `json:"progress"`
}

// Service 提供进度写入与查询。
type Service struct {
	repo progressdom.ProgressRepository
	calc *Calculator
	// cache 用于在测试或轻量场景下兜底返回最近写入的进度，避免空表时出现空列表。
	cache map[int64]map[string][]progressdom.LearningProgress
}

// ChapterProgressDTO 是用于 HTTP 响应的章节进度表示，包含额外的 percent 字段（数值形式进度）。
type ChapterProgressDTO struct {
	progressdom.LearningProgress
	Percent int `json:"percent"`
}

// EnrichChapters 将域模型列表转换为包含 percent 的 DTO 列表，且在返回前保证 status 与 percent 一致（不写回数据库，仅用于响应）。
func (s *Service) EnrichChapters(items []progressdom.LearningProgress) []ChapterProgressDTO {
	if items == nil {
		return []ChapterProgressDTO{}
	}
	var out []ChapterProgressDTO
	for _, it := range items {
		// 复制条目以避免修改源对象
		copyItem := it
		est := s.calc.lookupDuration(copyItem.Topic, copyItem.Chapter)
		// 计算基于阅读时长的进度（保守估算）
		percentFromRead := 0
		if est > 0 {
			percentFromRead = int((copyItem.ReadDuration * 100) / est)
			if percentFromRead < 0 {
				percentFromRead = 0
			}
			if percentFromRead > 100 {
				percentFromRead = 100
			}
		}
		// 优先取滚动进度与阅读估算的最大值
		percent := copyItem.ScrollProgress
		if percent < percentFromRead {
			percent = percentFromRead
		}
		if percent < 0 {
			percent = 0
		}
		if percent > 100 {
			percent = 100
		}

		// 保证状态与 percent 一致（仅响应层面）
		if percent >= 100 {
			copyItem.Status = progressdom.StatusCompleted
		} else {
			// 以计算器规则为准，但传入估算时长以保证判定一致
			copyItem.Status = s.calc.CalculateChapterStatus(copyItem, est)
		}

		out = append(out, ChapterProgressDTO{
			LearningProgress: copyItem,
			Percent:          percent,
		})
	}
	return out
}

// NewService 创建进度服务。
func NewService(repo progressdom.ProgressRepository, calc *Calculator) *Service {
	if calc == nil {
		calc = NewCalculator(nil, nil)
	}
	return &Service{
		repo:  repo,
		calc:  calc,
		cache: map[int64]map[string][]progressdom.LearningProgress{},
	}
}

// CreateOrUpdateProgress 写入或更新进度，并返回最新汇总。
func (s *Service) CreateOrUpdateProgress(ctx context.Context, req UpdateProgressRequest) (*ProgressResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	current, err := s.repo.Get(ctx, req.UserID, req.Topic, req.Chapter)
	if err != nil {
		return nil, err
	}
	if current == nil {
		current = &progressdom.LearningProgress{
			UserID:       req.UserID,
			Topic:        req.Topic,
			Chapter:      req.Chapter,
			Status:       progressdom.StatusNotStarted,
			FirstVisitAt: time.Now(),
			LastVisitAt:  time.Now(),
		}
	}

	merged, totalRead := s.mergeProgress(*current, req)
	statusProbe := merged
	statusProbe.ReadDuration = totalRead
	merged.Status = s.calc.CalculateChapterStatus(statusProbe, req.EstimatedSec)
	if err := s.repo.CreateOrUpdate(ctx, &merged); err != nil {
		return nil, err
	}
	s.saveCache(req.UserID, merged, totalRead)

	all, err := s.repo.GetByUser(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	overall, topics := s.calc.CalculateOverallProgress(all)
	topics = s.sortTopics(topics)

	var topicSummary TopicProgress
	for _, tp := range topics {
		if tp.ID == req.Topic {
			topicSummary = tp
			break
		}
	}
	if topicSummary.ID == "" {
		topicSummary = TopicProgress{
			Name:          topicName(req.Topic),
			ID:            req.Topic,
			Weight:        s.calc.topicWeight(req.Topic),
			TotalChapters: s.calc.topicTotal(req.Topic),
		}
	}

	return &ProgressResponse{
		Status:         merged.Status,
		Overall:        overall,
		Topic:          topicSummary,
		ReadDuration:   totalRead,
		ScrollProgress: merged.ScrollProgress,
		LastPosition:   merged.LastPosition,
	}, nil
}

// GetProgress 返回整体与主题列表进度。
func (s *Service) GetProgress(ctx context.Context, userID int64) (OverallProgress, []TopicProgress, error) {
	return s.GetOverallProgress(ctx, userID)
}

// GetOverallProgress 返回整体进度与主题维度汇总。
func (s *Service) GetOverallProgress(ctx context.Context, userID int64) (OverallProgress, []TopicProgress, error) {
	if userID <= 0 {
		return OverallProgress{}, nil, errors.New("用户信息缺失")
	}
	list, err := s.repo.GetByUser(ctx, userID)
	if err != nil {
		return OverallProgress{}, nil, err
	}
	overall, topics := s.calc.CalculateOverallProgress(list)
	return overall, s.sortTopics(topics), nil
}

// ListByUser 返回用户的章节进度列表。
func (s *Service) ListByUser(ctx context.Context, userID int64) ([]progressdom.LearningProgress, error) {
	if userID <= 0 {
		return nil, errors.New("用户信息缺失")
	}
	list, err := s.repo.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		if topics, ok := s.cache[userID]; ok {
			var fallback []progressdom.LearningProgress
			for _, items := range topics {
				fallback = append(fallback, items...)
			}
			if len(fallback) > 0 {
				return fallback, nil
			}
		}
		return []progressdom.LearningProgress{}, nil
	}
	return list, nil
}

// GetTopicProgress 返回指定主题的章节进度列表与汇总。
func (s *Service) GetTopicProgress(ctx context.Context, userID int64, topic string) (TopicProgress, []progressdom.LearningProgress, error) {
	topic = strings.TrimSpace(topic)
	if userID <= 0 || !progressdom.IsSupportedTopic(topic) {
		return TopicProgress{}, nil, errors.New("请求参数无效")
	}
	items, err := s.repo.GetByTopic(ctx, userID, topic)
	if err != nil {
		return TopicProgress{}, nil, err
	}
	if len(items) == 0 {
		if topics, ok := s.cache[userID]; ok {
			if cached, ok := topics[topic]; ok && len(cached) > 0 {
				items = cached
			}
		}
	}
	_, topicStats := s.calc.CalculateOverallProgress(items)
	var summary TopicProgress
	for _, tp := range topicStats {
		if tp.ID == topic {
			summary = tp
			break
		}
	}
	if summary.ID == "" {
		summary = TopicProgress{
			Name:          topicName(topic),
			ID:            topic,
			Weight:        s.calc.topicWeight(topic),
			TotalChapters: s.calc.topicTotal(topic),
		}
	}
	return summary, items, nil
}

func (s *Service) validateRequest(req UpdateProgressRequest) error {
	if req.UserID <= 0 {
		return errors.New("用户信息缺失")
	}
	req.Topic = strings.TrimSpace(req.Topic)
	req.Chapter = strings.TrimSpace(req.Chapter)
	if !progressdom.IsSupportedTopic(req.Topic) || req.Chapter == "" {
		return errors.New("请求参数无效")
	}
	if req.ReadDuration < 0 || req.ScrollProgress < 0 {
		return errors.New("请求参数无效")
	}
	if req.ScrollProgress > 100 {
		req.ScrollProgress = 100
	}
	return nil
}

func (s *Service) mergeProgress(existing progressdom.LearningProgress, req UpdateProgressRequest) (progressdom.LearningProgress, int64) {
	totalRead := existing.ReadDuration + req.ReadDuration
	payload := existing
	payload.ReadDuration = req.ReadDuration
	if req.ScrollProgress > payload.ScrollProgress {
		payload.ScrollProgress = req.ScrollProgress
	}
	if strings.TrimSpace(req.LastPosition) != "" {
		payload.LastPosition = req.LastPosition
	}
	if req.QuizScore > 0 {
		payload.QuizScore = req.QuizScore
	}
	if req.QuizPassed {
		payload.QuizPassed = true
	}
	payload.LastVisitAt = time.Now()
	return payload, totalRead
}

func (s *Service) saveCache(userID int64, p progressdom.LearningProgress, totalRead int64) {
	if s.cache == nil {
		return
	}
	if _, ok := s.cache[userID]; !ok {
		s.cache[userID] = map[string][]progressdom.LearningProgress{}
	}
	p.ReadDuration = totalRead
	list := s.cache[userID][p.Topic]
	updated := false
	for i, item := range list {
		if item.Chapter == p.Chapter {
			list[i] = p
			updated = true
			break
		}
	}
	if !updated {
		list = append(list, p)
	}
	s.cache[userID][p.Topic] = list
}

// GetNextUnfinishedChapter 返回按权重和章节顺序计算的首个未完成章节。
func (s *Service) GetNextUnfinishedChapter(ctx context.Context, userID int64) (*NextChapter, error) {
	if userID <= 0 {
		return nil, errors.New("用户信息缺失")
	}
	list, err := s.repo.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	progressMap := map[string]map[string]progressdom.LearningProgress{}
	for _, item := range list {
		if _, ok := progressMap[item.Topic]; !ok {
			progressMap[item.Topic] = map[string]progressdom.LearningProgress{}
		}
		progressMap[item.Topic][item.Chapter] = item
	}
	for _, topic := range s.sortedTopicIDs() {
		chapters := topicChapterOrder[topic]
		if len(chapters) == 0 {
			continue
		}
		for _, ch := range chapters {
			entry, ok := progressMap[topic][ch]
			if !ok || entry.Status != progressdom.StatusCompleted {
				status := progressdom.StatusNotStarted
				progressVal := 0
				if ok {
					status = entry.Status
					progressVal = entry.ScrollProgress
					if progressVal < 0 {
						progressVal = 0
					}
					if progressVal > 100 {
						progressVal = 100
					}
				}
				return &NextChapter{
					Topic:    topic,
					Chapter:  ch,
					Status:   status,
					Progress: progressVal,
				}, nil
			}
		}
	}
	return nil, nil
}

// sortTopics 按权重与预设顺序返回稳定排序的主题列表。
func (s *Service) sortTopics(topics []TopicProgress) []TopicProgress {
	sorted := make([]TopicProgress, len(topics))
	copy(sorted, topics)
	orderIndex := func(id string) int { return topicOrderIndex(id) }
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Weight == sorted[j].Weight {
			return orderIndex(sorted[i].ID) < orderIndex(sorted[j].ID)
		}
		return sorted[i].Weight > sorted[j].Weight
	})
	return sorted
}

// sortedTopicIDs 返回按权重及默认顺序排序的主题 ID。
func (s *Service) sortedTopicIDs() []string {
	var ids []string
	for id := range s.calc.TopicWeights {
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		ids = append(ids, defaultTopicOrder...)
	}
	orderIndex := func(id string) int { return topicOrderIndex(id) }
	sort.Slice(ids, func(i, j int) bool {
		wi := s.calc.topicWeight(ids[i])
		wj := s.calc.topicWeight(ids[j])
		if wi == wj {
			return orderIndex(ids[i]) < orderIndex(ids[j])
		}
		return wi > wj
	})
	return ids
}

// topicOrderIndex 返回主题在默认顺序中的索引，未命中时追加在末尾。
func topicOrderIndex(topic string) int {
	for idx, t := range defaultTopicOrder {
		if t == topic {
			return idx
		}
	}
	return len(defaultTopicOrder)
}
