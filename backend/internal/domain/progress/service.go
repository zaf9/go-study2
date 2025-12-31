package progress

import (
	"context"
	"errors"
	"strings"
	"time"
)

// ErrInvalidInput 表示请求参数不合法。
var ErrInvalidInput = errors.New("进度参数不合法")

// Service 封装学习进度的业务逻辑。
type Service struct {
	repo Repository
}

// NewService 创建学习进度服务。
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Save 记录或更新用户进度。
func (s *Service) Save(ctx context.Context, userID int64, topic, chapter, status, position string) (*Progress, error) {
	topic = strings.TrimSpace(topic)
	chapter = strings.TrimSpace(chapter)
	status = strings.TrimSpace(status)
	position = strings.TrimSpace(position)

	if userID <= 0 || !IsSupportedTopic(topic) || chapter == "" || !IsValidStatus(status) {
		return nil, ErrInvalidInput
	}

	record := &Progress{
		UserID:       userID,
		Topic:        topic,
		Chapter:      chapter,
		Status:       status,
		LastVisitAt:  time.Now(),
		LastPosition: position,
	}

	if err := s.repo.Upsert(ctx, record); err != nil {
		return nil, err
	}
	return record, nil
}

// ListAll 返回用户的全部进度。
func (s *Service) ListAll(ctx context.Context, userID int64) ([]Progress, error) {
	if userID <= 0 {
		return nil, ErrInvalidInput
	}
	return s.repo.ListByUser(ctx, userID)
}

// ListByTopic 返回用户在指定主题下的进度。
func (s *Service) ListByTopic(ctx context.Context, userID int64, topic string) ([]Progress, error) {
	topic = strings.TrimSpace(topic)
	if userID <= 0 || !IsSupportedTopic(topic) {
		return nil, ErrInvalidInput
	}
	return s.repo.ListByTopic(ctx, userID, topic)
}

// UpdateChapterStatus 更新章节学习状态，首次访问时标记为"学习中"
func (s *Service) UpdateChapterStatus(ctx context.Context, userID int64, topic, chapter string) error {
	if userID <= 0 || !IsSupportedTopic(topic) || chapter == "" {
		return ErrInvalidInput
	}

	// 查询现有记录
	existing, err := s.Get(ctx, userID, topic, chapter)
	if err != nil {
		return err
	}

	now := time.Now()

	if existing == nil {
		// 首次访问，创建记录，状态为"学习中"
		record := &LearningProgress{
			UserID:       userID,
			Topic:        topic,
			Chapter:      chapter,
			Status:       StatusInProgress,
			FirstVisitAt: now,
			LastVisitAt:  now,
		}
		return s.repo.CreateOrUpdate(ctx, record)
	}

	// 已有记录，仅更新访问时间
	existing.LastVisitAt = now
	// 如果当前状态是"未开始"，更新为"学习中"
	if existing.Status == StatusNotStarted {
		existing.Status = StatusInProgress
		if existing.FirstVisitAt.IsZero() {
			existing.FirstVisitAt = now
		}
	}
	return s.repo.CreateOrUpdate(ctx, existing)
}

// CompleteChapter 完成章节，更新测验成绩和状态
func (s *Service) CompleteChapter(ctx context.Context, userID int64, topic, chapter string, quizScore int, passed bool) error {
	if userID <= 0 || !IsSupportedTopic(topic) || chapter == "" {
		return ErrInvalidInput
	}

	existing, err := s.Get(ctx, userID, topic, chapter)
	if err != nil {
		return err
	}

	if existing == nil {
		return errors.New("章节学习记录不存在，请先学习章节内容")
	}

	// 更新测验成绩
	existing.QuizScore = quizScore
	existing.QuizPassed = passed

	// 如果通过测验，标记为已完成
	if passed {
		existing.Status = StatusCompleted
		now := time.Now()
		existing.CompletedAt = &now
		existing.LastVisitAt = now
	}

	return s.repo.CreateOrUpdate(ctx, existing)
}

// GetTopicProgressWithStatus 获取主题的章节进度详情，包含所有章节的状态
func (s *Service) GetTopicProgressWithStatus(ctx context.Context, userID int64, topic string) (*TopicProgressDetail, error) {
	if userID <= 0 || !IsSupportedTopic(topic) {
		return nil, ErrInvalidInput
	}

	// 获取该主题的学习记录
	records, err := s.repo.ListByTopic(ctx, userID, topic)
	if err != nil {
		return nil, err
	}

	// 构建章节进度映射
	progressMap := make(map[string]*LearningProgress)
	for i := range records {
		progressMap[records[i].Chapter] = &records[i]
	}

	// 获取该主题的所有章节定义
	chapters, ok := TopicChapterOrder[topic]
	if !ok {
		return nil, errors.New("不支持的主题: " + topic)
	}

	// 构建章节状态列表
	chapterInfos := make([]ChapterStatusInfo, 0, len(chapters))
	for _, chapter := range chapters {
		record, exists := progressMap[chapter]
		if exists {
			// 有学习记录
			var lastVisitAt, completedAt *time.Time
			if !record.LastVisitAt.IsZero() {
				lastVisitAt = &record.LastVisitAt
			}
			if record.CompletedAt != nil && !record.CompletedAt.IsZero() {
				completedAt = record.CompletedAt
			}

			chapterInfos = append(chapterInfos, ChapterStatusInfo{
				Chapter:     chapter,
				Status:      record.Status,
				QuizScore:   record.QuizScore,
				QuizPassed:  record.QuizPassed,
				LastVisitAt: lastVisitAt,
				CompletedAt: completedAt,
			})
		} else {
			// 无学习记录，默认为未开始
			chapterInfos = append(chapterInfos, ChapterStatusInfo{
				Chapter:     chapter,
				Status:      StatusNotStarted,
				QuizScore:   0,
				QuizPassed:  false,
				LastVisitAt: nil,
				CompletedAt: nil,
			})
		}
	}

	// 过滤已删除的章节：检查数据库记录中不在当前章节定义中的记录
	// 这些记录可能是之前删除的章节，不应该出现在进度统计中
	// (当前实现已经只显示TopicChapterOrder中定义的章节，自动过滤了删除的章节)

	return &TopicProgressDetail{
		Topic:         topic,
		TotalChapters: len(chapters),
		Chapters:      chapterInfos,
	}, nil
}

// Get 获取单个章节的进度
func (s *Service) Get(ctx context.Context, userID int64, topic, chapter string) (*LearningProgress, error) {
	if userID <= 0 || !IsSupportedTopic(topic) || chapter == "" {
		return nil, ErrInvalidInput
	}

	// 转换为使用 ProgressRepository 接口
	if progRepo, ok := s.repo.(interface {
		Get(context.Context, int64, string, string) (*LearningProgress, error)
	}); ok {
		return progRepo.Get(ctx, userID, topic, chapter)
	}

	// 否则从列表中查找
	list, err := s.repo.ListByTopic(ctx, userID, topic)
	if err != nil {
		return nil, err
	}
	for i := range list {
		if list[i].Chapter == chapter {
			return &list[i], nil
		}
	}
	return nil, nil
}

// GetOverview 获取用户的全局学习进度概览
func (s *Service) GetOverview(ctx context.Context, userID int64) (*ProgressOverview, error) {
	if userID <= 0 {
		return nil, ErrInvalidInput
	}

	// 转换为使用 ProgressRepository 接口
	if progRepo, ok := s.repo.(interface {
		GetOverview(context.Context, int64) (*ProgressOverview, error)
	}); ok {
		return progRepo.GetOverview(ctx, userID)
	}

	// 后备实现：从基本数据构建概览
	return nil, errors.New("GetOverview not supported by repository")
}
