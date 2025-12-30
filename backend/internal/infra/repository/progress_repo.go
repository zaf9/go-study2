package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go-study2/internal/domain/progress"

	"github.com/gogf/gf/v2/database/gdb"
)

// ProgressRepository 使用 GoFrame gdb 实现进度表的持久化。
type ProgressRepository struct {
	db gdb.DB
}

// NewProgressRepository 创建仓储实例。
func NewProgressRepository(db gdb.DB) *ProgressRepository {
	return &ProgressRepository{db: db}
}

// CreateOrUpdate 写入或更新学习进度，累加时长并避免状态回退。
func (r *ProgressRepository) CreateOrUpdate(ctx context.Context, record *progress.LearningProgress) error {
	if record == nil {
		return errors.New("record is nil")
	}
	lastPosition := int64(0)
	if record.LastPosition != "" {
		if parsed, err := strconv.ParseInt(record.LastPosition, 10, 64); err == nil {
			lastPosition = parsed
		}
	}

	_, err := r.db.Exec(ctx, `
INSERT INTO learning_progress (user_id, topic, chapter, status, read_duration, scroll_progress, last_position, quiz_score, quiz_passed, first_visit_at, last_visit_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT(user_id, topic, chapter)
DO UPDATE SET
    read_duration = learning_progress.read_duration + excluded.read_duration,
    scroll_progress = excluded.scroll_progress,
    last_position = excluded.last_position,
    quiz_score = COALESCE(excluded.quiz_score, learning_progress.quiz_score),
    quiz_passed = CASE WHEN excluded.quiz_passed = 1 THEN 1 ELSE learning_progress.quiz_passed END,
    status = CASE
                 WHEN learning_progress.status = 'completed' AND excluded.status != 'completed' THEN learning_progress.status
                 WHEN learning_progress.status = 'tested' AND excluded.status = 'in_progress' THEN learning_progress.status
                 ELSE excluded.status
             END,
    last_visit_at = CURRENT_TIMESTAMP,
    completed_at = CASE
                       WHEN excluded.status = 'completed' AND excluded.quiz_passed = 1 THEN COALESCE(learning_progress.completed_at, CURRENT_TIMESTAMP)
                       ELSE learning_progress.completed_at
                   END,
    updated_at = CURRENT_TIMESTAMP
`, record.UserID, record.Topic, record.Chapter, record.Status, record.ReadDuration, record.ScrollProgress, lastPosition, record.QuizScore, record.QuizPassed)
	return err
}

// Get 返回指定章节的进度。
func (r *ProgressRepository) Get(ctx context.Context, userID int64, topic, chapter string) (*progress.LearningProgress, error) {
	one, err := r.db.Model("learning_progress").
		Where("user_id", userID).
		Where("topic", topic).
		Where("chapter", chapter).
		One(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if one.IsEmpty() {
		return nil, nil
	}
	var item progress.LearningProgress
	if err := one.Struct(&item); err != nil {
		return nil, err
	}
	return &item, nil
}

// GetByUser 返回用户的全部进度。
func (r *ProgressRepository) GetByUser(ctx context.Context, userID int64) ([]progress.LearningProgress, error) {
	records, err := r.db.Model("learning_progress").
		Where("user_id", userID).
		OrderDesc("last_visit_at").
		All(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []progress.LearningProgress{}, nil
		}
		return nil, err
	}
	var items []progress.LearningProgress
	if err := records.Structs(&items); err != nil {
		return nil, err
	}
	return items, nil
}

// GetByTopic 返回指定主题下的进度。
func (r *ProgressRepository) GetByTopic(ctx context.Context, userID int64, topic string) ([]progress.LearningProgress, error) {
	records, err := r.db.Model("learning_progress").
		Where("user_id", userID).
		Where("topic", topic).
		OrderDesc("last_visit_at").
		All(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []progress.LearningProgress{}, nil
		}
		return nil, err
	}
	var items []progress.LearningProgress
	if err := records.Structs(&items); err != nil {
		return nil, err
	}
	return items, nil
}

// GetLastLearning 返回用户最后一次学习的记录（按 last_visit_at 降序排序的第一条）。
func (r *ProgressRepository) GetLastLearning(ctx context.Context, userID int64) (*progress.LearningProgress, error) {
	one, err := r.db.Model("learning_progress").
		Where("user_id", userID).
		Where("last_visit_at IS NOT NULL").
		OrderDesc("last_visit_at").
		One(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if one.IsEmpty() {
		return nil, nil
	}
	var item progress.LearningProgress
	if err := one.Struct(&item); err != nil {
		return nil, err
	}
	return &item, nil
}

// GetOverview 计算用户的全局学习进度概览
func (r *ProgressRepository) GetOverview(ctx context.Context, userID int64) (*progress.ProgressOverview, error) {
	// 获取用户所有学习记录
	allRecords, err := r.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 按主题和章节构建进度映射
	progressMap := make(map[string]map[string]*progress.LearningProgress) // topic -> chapter -> record
	for _, record := range allRecords {
		if _, ok := progressMap[record.Topic]; !ok {
			progressMap[record.Topic] = make(map[string]*progress.LearningProgress)
		}
		progressMap[record.Topic][record.Chapter] = &record
	}

	// 计算各主题的统计信息
	topicSummaries := make([]progress.TopicProgressSummary, 0)
	totalCompleted := 0
	totalInProgress := 0
	totalChapters := 0

	for topic, chapters := range progress.TopicChapterOrder {
		topicTotal := len(chapters)
		completed := 0
		inProgress := 0

		topicProgress, ok := progressMap[topic]
		if !ok {
			topicProgress = make(map[string]*progress.LearningProgress)
		}

		for _, chapter := range chapters {
			record, exists := topicProgress[chapter]
			if exists {
				if record.Status == progress.StatusCompleted && record.QuizPassed {
					completed++
				} else if record.Status == progress.StatusInProgress {
					inProgress++
				}
			}
		}

		topicSummaries = append(topicSummaries, progress.TopicProgressSummary{
			Topic:              topic,
			TotalChapters:      topicTotal,
			CompletedChapters:  completed,
			InProgressChapters: inProgress,
		})

		totalChapters += topicTotal
		totalCompleted += completed
		totalInProgress += inProgress
	}

	// 计算完成率
	completionRate := 0.0
	if totalChapters > 0 {
		completionRate = float64(totalCompleted) / float64(totalChapters) * 100
	}

	// 查找下一个建议学习的章节
	var nextChapter *progress.NextChapterHint
	lastLearning, err := r.GetLastLearning(ctx, userID)
	if err == nil && lastLearning != nil {
		// 简单逻辑：返回最后学习的章节
		nextChapter = &progress.NextChapterHint{
			Topic:   lastLearning.Topic,
			Chapter: lastLearning.Chapter,
			Title:   progress.ChapterDisplayName(lastLearning.Chapter),
		}
	}

	return &progress.ProgressOverview{
		TotalChapters:      totalChapters,
		CompletedChapters:  totalCompleted,
		InProgressChapters: totalInProgress,
		CompletionRate:     completionRate,
		Topics:             topicSummaries,
		NextChapter:        nextChapter,
	}, nil
}

// GetByUserAndTopic 获取用户在指定主题的所有章节进度，包含未开始的章节
func (r *ProgressRepository) GetByUserAndTopic(ctx context.Context, userID int64, topic string) (*progress.TopicProgressDetail, error) {
	// 获取该主题的所有学习记录
	records, err := r.GetByTopic(ctx, userID, topic)
	if err != nil {
		return nil, err
	}

	// 构建章节进度映射
	progressMap := make(map[string]*progress.LearningProgress)
	for i := range records {
		progressMap[records[i].Chapter] = &records[i]
	}

	// 获取该主题的所有章节定义
	chapters, ok := progress.TopicChapterOrder[topic]
	if !ok {
		return nil, fmt.Errorf("不支持的主题: %s", topic)
	}

	// 构建章节状态列表
	chapterInfos := make([]progress.ChapterStatusInfo, 0, len(chapters))
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

			chapterInfos = append(chapterInfos, progress.ChapterStatusInfo{
				Chapter:     chapter,
				Status:      record.Status,
				QuizScore:   record.QuizScore,
				QuizPassed:  record.QuizPassed,
				LastVisitAt: lastVisitAt,
				CompletedAt: completedAt,
			})
		} else {
			// 无学习记录，默认为未开始
			chapterInfos = append(chapterInfos, progress.ChapterStatusInfo{
				Chapter:     chapter,
				Status:      progress.StatusNotStarted,
				QuizScore:   0,
				QuizPassed:  false,
				LastVisitAt: nil,
				CompletedAt: nil,
			})
		}
	}

	return &progress.TopicProgressDetail{
		Topic:         topic,
		TotalChapters: len(chapters),
		Chapters:      chapterInfos,
	}, nil
}
