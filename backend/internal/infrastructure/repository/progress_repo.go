package repository

import (
	"context"
	"strconv"
	"strings"
	"time"

	"go-study2/internal/domain/progress"

	"github.com/gogf/gf/v2/database/gdb"
)

// ProgressRepository 使用 GoFrame gdb 持久化学习进度。
type ProgressRepository struct {
	db gdb.DB
}

// NewProgressRepository 创建进度仓储实例。
func NewProgressRepository(db gdb.DB) *ProgressRepository {
	return &ProgressRepository{db: db}
}

// Upsert 写入或更新进度记录。
func (r *ProgressRepository) Upsert(ctx context.Context, record *progress.Progress) error {
	now := time.Now()
	record.LastVisitAt = now
	lastPosition := int64(0)
	if v := strings.TrimSpace(record.LastPosition); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil && parsed >= 0 {
			lastPosition = parsed
		}
	}
	_, err := r.db.Exec(ctx, `
INSERT INTO learning_progress (user_id, topic, chapter, status, last_visit_at, last_position, updated_at)
VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(user_id, topic, chapter)
DO UPDATE SET status=excluded.status,
              last_visit_at=excluded.last_visit_at,
              last_position=excluded.last_position,
              updated_at=CURRENT_TIMESTAMP
`, record.UserID, record.Topic, record.Chapter, record.Status, record.LastVisitAt, lastPosition)
	return err
}

// ListByUser 按用户查询全部进度。
func (r *ProgressRepository) ListByUser(ctx context.Context, userID int64) ([]progress.Progress, error) {
	return r.query(ctx, r.db.Model("learning_progress").Where("user_id", userID))
}

// ListByTopic 查询用户在指定主题下的进度。
func (r *ProgressRepository) ListByTopic(ctx context.Context, userID int64, topic string) ([]progress.Progress, error) {
	return r.query(ctx, r.db.Model("learning_progress").Where("user_id", userID).Where("topic", topic))
}

// CreateOrUpdate 创建或更新学习进度记录
func (r *ProgressRepository) CreateOrUpdate(ctx context.Context, record *progress.LearningProgress) error {
	_, err := r.db.Model("learning_progress").Save(ctx, record)
	return err
}

// Get 获取单个章节的进度
func (r *ProgressRepository) Get(ctx context.Context, userID int64, topic, chapter string) (*progress.LearningProgress, error) {
	record, err := r.db.Model("learning_progress").
		Where("user_id", userID).
		Where("topic", topic).
		Where("chapter", chapter).
		One(ctx)
	if err != nil {
		return nil, err
	}
	var result progress.LearningProgress
	if err := record.Struct(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *ProgressRepository) query(ctx context.Context, model *gdb.Model) ([]progress.Progress, error) {
	records, err := model.OrderDesc("last_visit_at").All(ctx)
	if err != nil {
		return nil, err
	}
	var items []progress.Progress
	if err := records.Structs(&items); err != nil {
		return nil, err
	}
	return items, nil
}

// GetOverview 获取用户的全局学习进度概览
func (r *ProgressRepository) GetOverview(ctx context.Context, userID int64) (*progress.ProgressOverview, error) {
	// 获取所有学习记录
	records, err := r.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 构建章节进度映射
	progressMap := make(map[string]*progress.LearningProgress)
	for i := range records {
		key := records[i].Topic + "/" + records[i].Chapter
		progressMap[key] = &records[i]
	}

	// 统计各主题的进度
	totalChapters := 0
	completedChapters := 0
	inProgressChapters := 0

	topics := make([]progress.TopicProgressSummary, 0, 4)
	for topic := range progress.SupportedTopics {
		chapters, ok := progress.TopicChapterOrder[topic]
		if !ok {
			continue
		}

		topicCompleted := 0
		topicInProgress := 0
		topicTotal := len(chapters)

		for _, chapter := range chapters {
			key := topic + "/" + chapter
			if record, exists := progressMap[key]; exists {
				if record.Status == progress.StatusCompleted {
					topicCompleted++
					completedChapters++
				} else if record.Status == progress.StatusInProgress {
					topicInProgress++
					inProgressChapters++
				}
			}
			totalChapters++
		}

		topics = append(topics, progress.TopicProgressSummary{
			Topic:              topic,
			TotalChapters:      topicTotal,
			CompletedChapters:  topicCompleted,
			InProgressChapters: topicInProgress,
		})
	}

	completionRate := 0.0
	if totalChapters > 0 {
		completionRate = float64(completedChapters) / float64(totalChapters) * 100
	}

	// 查找下一个建议学习的章节
	var nextChapter *progress.NextChapterHint
	for _, topicInfo := range topics {
		if topicInfo.InProgressChapters > 0 {
			// 找到第一个有进行中章节的主题
			chapters, ok := progress.TopicChapterOrder[topicInfo.Topic]
			if !ok {
				continue
			}
			for _, chapter := range chapters {
				key := topicInfo.Topic + "/" + chapter
				if record, exists := progressMap[key]; exists && record.Status == progress.StatusInProgress {
					nextChapter = &progress.NextChapterHint{
						Topic:   topicInfo.Topic,
						Chapter: chapter,
						Title:   chapter, // 简化实现，使用章节名作为标题
					}
					break
				}
			}
			break
		}
	}

	// 如果没有进行中的章节，找第一个未开始的章节
	if nextChapter == nil {
		for _, topicInfo := range topics {
			if topicInfo.CompletedChapters < topicInfo.TotalChapters {
				chapters, ok := progress.TopicChapterOrder[topicInfo.Topic]
				if !ok {
					continue
				}
				for _, chapter := range chapters {
					key := topicInfo.Topic + "/" + chapter
					if record, exists := progressMap[key]; !exists || record.Status == progress.StatusNotStarted {
						nextChapter = &progress.NextChapterHint{
							Topic:   topicInfo.Topic,
							Chapter: chapter,
							Title:   chapter,
						}
						break
					}
				}
				break
			}
		}
	}

	return &progress.ProgressOverview{
		TotalChapters:      totalChapters,
		CompletedChapters:  completedChapters,
		InProgressChapters: inProgressChapters,
		CompletionRate:     completionRate,
		Topics:             topics,
		NextChapter:        nextChapter,
	}, nil
}

func nullableString(v string) interface{} {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	return v
}
