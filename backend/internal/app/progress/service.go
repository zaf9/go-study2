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

// TopicProgressSummary 表示主题进度汇总（用于 Dashboard）。
type TopicProgressSummary struct {
	TopicID           string  `json:"topic_id"`
	TopicName         string  `json:"topic_name"`
	DisplayName       string  `json:"display_name"`
	CompletedChapters int     `json:"completed_chapters"`
	TotalChapters     int     `json:"total_chapters"`
	Percentage        float64 `json:"percentage"`
}

// NextChapter 表示推荐继续学习的章节。
type NextChapter struct {
	Topic    string `json:"topic"`
	Chapter  string `json:"chapter"`
	Status   string `json:"status"`
	Progress int    `json:"progress"`
}

// LastLearningRecord 表示用户最后一次学习的记录。
type LastLearningRecord struct {
	TopicID            string `json:"topic_id"`
	TopicName          string `json:"topic_name"`
	TopicDisplayName   string `json:"topic_display_name"`
	ChapterID          string `json:"chapter_id"`
	ChapterName        string `json:"chapter_name"`
	ChapterDisplayName string `json:"chapter_display_name"`
	LastVisitedAt      string `json:"last_visited_at"`
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

	// 根据阅读进度自动更新状态：当滚动进度达到100%或测验通过时，标记为已完成
	est := s.calc.lookupDuration(payload.Topic, payload.Chapter)
	percentFromRead := 0
	if est > 0 {
		percentFromRead = int((payload.ReadDuration * 100) / est)
		if percentFromRead < 0 {
			percentFromRead = 0
		}
		if percentFromRead > 100 {
			percentFromRead = 100
		}
	}
	percent := payload.ScrollProgress
	if percent < percentFromRead {
		percent = percentFromRead
	}

	// 只有在状态不是已完成时才自动更新
	if payload.Status != progressdom.StatusCompleted {
		if percent >= 100 || req.QuizPassed {
			payload.Status = progressdom.StatusCompleted
			if payload.CompletedAt == nil {
				now := time.Now()
				payload.CompletedAt = &now
			}
		} else if payload.Status == progressdom.StatusNotStarted {
			// 如果从未开始变为有进度，更新为学习中
			if percent > 0 || req.ReadDuration > 0 {
				payload.Status = progressdom.StatusInProgress
			}
		}
	}

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
		chapters := TopicChapterOrder[topic]
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

// GetTopicProgressSummary 返回所有主题的进度汇总列表。
func (s *Service) GetTopicProgressSummary(ctx context.Context, userID int64) ([]TopicProgressSummary, error) {
	if userID <= 0 {
		return nil, errors.New("用户信息缺失")
	}

	// 获取用户的整体进度（包含所有主题的进度信息）
	_, topics, err := s.GetOverallProgress(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 获取所有支持的主题 ID（确保即使没有学习记录的主题也显示）
	allTopicIDs := s.sortedTopicIDs()
	topicMap := make(map[string]TopicProgress)
	for _, tp := range topics {
		topicMap[tp.ID] = tp
	}

	// 构建主题进度汇总列表
	summaries := make([]TopicProgressSummary, 0, len(allTopicIDs))
	for _, topicID := range allTopicIDs {
		tp, exists := topicMap[topicID]
		if !exists {
			// 如果主题没有学习记录，创建默认的汇总
			tp = TopicProgress{
				Name:              topicName(topicID),
				ID:                topicID,
				Weight:            s.calc.topicWeight(topicID),
				TotalChapters:     s.calc.topicTotal(topicID),
				Progress:          0,
				CompletedChapters: 0,
			}
		}

		// 计算百分比（保留一位小数）
		percentage := 0.0
		if tp.TotalChapters > 0 {
			percentage = float64(tp.CompletedChapters) / float64(tp.TotalChapters) * 100
			percentage = float64(int(percentage*10+0.5)) / 10 // 四舍五入到一位小数
		}

		summaries = append(summaries, TopicProgressSummary{
			TopicID:           tp.ID,
			TopicName:         tp.Name,
			DisplayName:       tp.Name,
			CompletedChapters: tp.CompletedChapters,
			TotalChapters:     tp.TotalChapters,
			Percentage:        percentage,
		})
	}

	return summaries, nil
}

// GetLastLearningRecord 返回用户最后一次学习的记录。
func (s *Service) GetLastLearningRecord(ctx context.Context, userID int64) (*LastLearningRecord, error) {
	if userID <= 0 {
		return nil, errors.New("用户信息缺失")
	}

	progress, err := s.repo.GetLastLearning(ctx, userID)
	if err != nil {
		return nil, err
	}
	if progress == nil {
		return nil, nil
	}

	// 获取主题显示名称
	topicDisplayName := topicName(progress.Topic)
	topicNameEn := topicDisplayName
	if progress.Topic == "lexical_elements" {
		topicNameEn = "Lexical Elements"
	} else if progress.Topic == "constants" {
		topicNameEn = "Constants"
	} else if progress.Topic == "variables" {
		topicNameEn = "Variables"
	} else if progress.Topic == "types" {
		topicNameEn = "Types"
	}

	// 获取章节显示名称
	chapterDisplayName := chapterDisplayName(progress.Chapter)
	chapterNameEn := chapterDisplayName

	return &LastLearningRecord{
		TopicID:            progress.Topic,
		TopicName:          topicNameEn,
		TopicDisplayName:   topicDisplayName,
		ChapterID:          progress.Chapter,
		ChapterName:        chapterNameEn,
		ChapterDisplayName: chapterDisplayName,
		LastVisitedAt:      progress.LastVisitAt.Format(time.RFC3339),
	}, nil
}

// GetOverview 返回用户的全局学习进度概览 (User Story 4)
// 统一数据源,确保所有页面显示一致的进度统计
func (s *Service) GetOverview(ctx context.Context, userID int64) (*progressdom.ProgressOverview, error) {
	if userID <= 0 {
		return nil, errors.New("用户信息缺失")
	}

	// 调用 repository 的 GetOverview 方法
	overview, err := s.repo.GetOverview(ctx, userID)
	if err != nil {
		return nil, err
	}

	return overview, nil
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
