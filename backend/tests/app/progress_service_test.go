package app

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"testing"
	"time"

	progapp "go-study2/internal/app/progress"
	progressdom "go-study2/internal/domain/progress"
)

func TestProgressService_CreateOrUpdateProgress(t *testing.T) {
	repo := newMemoryProgressRepo()
	calc := progapp.NewCalculator(
		map[string]int{"variables": 30, "constants": 20, "lexical_elements": 25, "types": 25},
		map[string]int{"variables": 2},
	)
	calc.EstimatedDuration["variables/storage"] = 600
	service := progapp.NewService(repo, calc)

	resp, err := service.CreateOrUpdateProgress(ctx(), progapp.UpdateProgressRequest{
		UserID:         1,
		Topic:          "variables",
		Chapter:        "storage",
		ReadDuration:   120,
		ScrollProgress: 60,
		LastPosition:   "120",
	})
	if err != nil {
		t.Fatalf("首次写入失败: %v", err)
	}
	if resp.Status != progressdom.StatusInProgress {
		t.Fatalf("状态应为 in_progress，得到 %s", resp.Status)
	}
	if resp.ReadDuration != 120 || resp.ScrollProgress != 60 {
		t.Fatalf("读写值异常: %+v", resp)
	}

	resp, err = service.CreateOrUpdateProgress(ctx(), progapp.UpdateProgressRequest{
		UserID:         1,
		Topic:          "variables",
		Chapter:        "storage",
		ReadDuration:   480,
		ScrollProgress: 95,
		QuizScore:      90,
		QuizPassed:     true,
		EstimatedSec:   600,
	})
	if err != nil {
		t.Fatalf("二次写入失败: %v", err)
	}
	if resp.Status != progressdom.StatusCompleted {
		t.Fatalf("状态应为 completed，得到 %s", resp.Status)
	}
	if resp.ReadDuration != 600 {
		t.Fatalf("阅读时长应累加为 600，得到 %d", resp.ReadDuration)
	}
	if resp.Overall.Progress <= 0 {
		t.Fatalf("整体进度应大于 0，得到 %d", resp.Overall.Progress)
	}
}

func TestCalculator_StatusAndOverall(t *testing.T) {
	calc := progapp.NewCalculator(nil, map[string]int{"variables": 2})
	progress := progressdom.LearningProgress{
		Topic:          "variables",
		Chapter:        "storage",
		ReadDuration:   500,
		ScrollProgress: 92,
		QuizScore:      90,
		QuizPassed:     true,
		LastPosition:   "900",
		LastVisitAt:    time.Now(),
	}
	status := calc.CalculateChapterStatus(progress, 600)
	if status != progressdom.StatusCompleted {
		t.Fatalf("应判定为 completed，得到 %s", status)
	}
	progress.Status = status

	another := progressdom.LearningProgress{
		Topic:          "variables",
		Chapter:        "pointer",
		ReadDuration:   0,
		ScrollProgress: 0,
		LastVisitAt:    time.Now(),
	}
	overall, topics := calc.CalculateOverallProgress([]progressdom.LearningProgress{progress, another})
	if overall.Progress != 50 {
		t.Fatalf("整体进度应为 50，得到 %d", overall.Progress)
	}
	if len(topics) != 1 || topics[0].Progress != 50 {
		t.Fatalf("主题进度不符合预期: %+v", topics)
	}
}

func TestProgressService_OverallAndNext(t *testing.T) {
	repo := newMemoryProgressRepo()
	calc := progapp.NewCalculator(map[string]int{"variables": 40, "constants": 20}, nil)
	service := progapp.NewService(repo, calc)

	day1 := time.Now().Add(-48 * time.Hour)
	if err := repo.CreateOrUpdate(ctx(), &progressdom.LearningProgress{
		UserID:         1,
		Topic:          "variables",
		Chapter:        "storage",
		Status:         progressdom.StatusCompleted,
		ReadDuration:   600,
		ScrollProgress: 95,
		QuizScore:      90,
		QuizPassed:     true,
		FirstVisitAt:   day1,
		LastVisitAt:    day1,
	}); err != nil {
		t.Fatalf("准备完成记录失败: %v", err)
	}
	if err := repo.CreateOrUpdate(ctx(), &progressdom.LearningProgress{
		UserID:         1,
		Topic:          "variables",
		Chapter:        "static",
		Status:         progressdom.StatusInProgress,
		ReadDuration:   120,
		ScrollProgress: 40,
		LastVisitAt:    time.Now(),
	}); err != nil {
		t.Fatalf("准备进行中记录失败: %v", err)
	}

	overall, topics, err := service.GetOverallProgress(ctx(), 1)
	if err != nil {
		t.Fatalf("获取汇总失败: %v", err)
	}
	if overall.CompletedChapters != 1 {
		t.Fatalf("已完成章节应为 1，得到 %d", overall.CompletedChapters)
	}
	if overall.TotalChapters < 4 {
		t.Fatalf("总章节应包含默认值，得到 %d", overall.TotalChapters)
	}
	if overall.TotalStudyTime != 720 {
		t.Fatalf("学习时长应累加，得到 %d", overall.TotalStudyTime)
	}
	if overall.StudyDays != 2 {
		t.Fatalf("学习天数应去重计算，得到 %d", overall.StudyDays)
	}
	if len(topics) == 0 || topics[0].Weight == 0 || topics[0].Name == "" {
		t.Fatalf("主题摘要应包含权重与名称: %+v", topics)
	}

	next, err := service.GetNextUnfinishedChapter(ctx(), 1)
	if err != nil {
		t.Fatalf("获取下一章节失败: %v", err)
	}
	if next == nil || next.Chapter != "static" {
		t.Fatalf("下一章节应为 static，得到 %+v", next)
	}
}

// ctx 返回带超时的上下文，避免测试泄漏。
func ctx() context.Context {
	return context.Background()
}

// 内存仓储实现 ProgressRepository，便于单元测试。
type memoryProgressRepo struct {
	data map[string]progressdom.LearningProgress
}

func newMemoryProgressRepo() *memoryProgressRepo {
	return &memoryProgressRepo{data: map[string]progressdom.LearningProgress{}}
}

func (m *memoryProgressRepo) key(userID int64, topic, chapter string) string {
	return strings.Join([]string{strconv.FormatInt(userID, 10), topic, chapter}, "|")
}

func (m *memoryProgressRepo) CreateOrUpdate(ctx context.Context, record *progressdom.LearningProgress) error {
	if record == nil {
		return errors.New("record is nil")
	}
	key := m.key(record.UserID, record.Topic, record.Chapter)
	existing, ok := m.data[key]
	if ok {
		record.ReadDuration += existing.ReadDuration
		if record.ScrollProgress < existing.ScrollProgress {
			record.ScrollProgress = existing.ScrollProgress
		}
	}
	if record.Status == progressdom.StatusCompleted && record.QuizPassed && record.CompletedAt == nil {
		now := time.Now()
		record.CompletedAt = &now
	}
	m.data[key] = *record
	return nil
}

func (m *memoryProgressRepo) Get(ctx context.Context, userID int64, topic, chapter string) (*progressdom.LearningProgress, error) {
	key := m.key(userID, topic, chapter)
	if val, ok := m.data[key]; ok {
		cp := val
		return &cp, nil
	}
	return nil, nil
}

func (m *memoryProgressRepo) GetByUser(ctx context.Context, userID int64) ([]progressdom.LearningProgress, error) {
	var list []progressdom.LearningProgress
	for _, v := range m.data {
		if v.UserID == userID {
			list = append(list, v)
		}
	}
	return list, nil
}

func (m *memoryProgressRepo) GetByTopic(ctx context.Context, userID int64, topic string) ([]progressdom.LearningProgress, error) {
	var list []progressdom.LearningProgress
	for _, v := range m.data {
		if v.UserID == userID && v.Topic == topic {
			list = append(list, v)
		}
	}
	return list, nil
}
