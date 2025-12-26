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

// TestProgressService_StudyDaysCalculation 测试学习天数计算逻辑
// 验证不同日期的学习记录能正确去重计算学习天数
func TestProgressService_StudyDaysCalculation(t *testing.T) {
	repo := newMemoryProgressRepo()
	calc := progapp.NewCalculator(map[string]int{"variables": 30, "constants": 20}, nil)
	service := progapp.NewService(repo, calc)

	// 准备不同日期的学习记录
	now := time.Now()
	day1 := now.AddDate(0, 0, -5) // 5天前
	day2 := now.AddDate(0, 0, -3) // 3天前
	day3 := now.AddDate(0, 0, -1) // 1天前
	day4 := now                    // 今天

	// 同一天的多条记录应只计算为1天
	if err := repo.CreateOrUpdate(ctx(), &progressdom.LearningProgress{
		UserID:       1,
		Topic:         "variables",
		Chapter:       "storage",
		Status:        progressdom.StatusInProgress,
		ReadDuration:  120,
		LastVisitAt:   day1,
	}); err != nil {
		t.Fatalf("准备第1天记录1失败: %v", err)
	}
	if err := repo.CreateOrUpdate(ctx(), &progressdom.LearningProgress{
		UserID:       1,
		Topic:         "variables",
		Chapter:       "pointer",
		Status:        progressdom.StatusInProgress,
		ReadDuration:  60,
		LastVisitAt:   day1.Add(2 * time.Hour), // 同一天的不同时间
	}); err != nil {
		t.Fatalf("准备第1天记录2失败: %v", err)
	}

	// 第2天的记录
	if err := repo.CreateOrUpdate(ctx(), &progressdom.LearningProgress{
		UserID:       1,
		Topic:         "constants",
		Chapter:       "iota",
		Status:        progressdom.StatusInProgress,
		ReadDuration:  90,
		LastVisitAt:   day2,
	}); err != nil {
		t.Fatalf("准备第2天记录失败: %v", err)
	}

	// 第3天的记录
	if err := repo.CreateOrUpdate(ctx(), &progressdom.LearningProgress{
		UserID:       1,
		Topic:         "variables",
		Chapter:       "static",
		Status:        progressdom.StatusCompleted,
		ReadDuration:  300,
		LastVisitAt:   day3,
	}); err != nil {
		t.Fatalf("准备第3天记录失败: %v", err)
	}

	// 第4天的记录
	if err := repo.CreateOrUpdate(ctx(), &progressdom.LearningProgress{
		UserID:       1,
		Topic:         "constants",
		Chapter:       "untyped",
		Status:        progressdom.StatusInProgress,
		ReadDuration:  150,
		LastVisitAt:   day4,
	}); err != nil {
		t.Fatalf("准备第4天记录失败: %v", err)
	}

	// 获取整体进度并验证学习天数
	overall, _, err := service.GetOverallProgress(ctx(), 1)
	if err != nil {
		t.Fatalf("获取汇总失败: %v", err)
	}

	// 应该有4个不同的学习天数（day1算1天，day2、day3、day4各1天）
	expectedDays := 4
	if overall.StudyDays != expectedDays {
		t.Fatalf("学习天数应为 %d，得到 %d。记录日期: day1=%v, day2=%v, day3=%v, day4=%v",
			expectedDays, overall.StudyDays, day1.Format("2006-01-02"), day2.Format("2006-01-02"),
			day3.Format("2006-01-02"), day4.Format("2006-01-02"))
	}

	// 验证无学习记录时学习天数为0
	overallEmpty, _, err := service.GetOverallProgress(ctx(), 999)
	if err != nil {
		t.Fatalf("获取空用户汇总失败: %v", err)
	}
	if overallEmpty.StudyDays != 0 {
		t.Fatalf("无学习记录时学习天数应为 0，得到 %d", overallEmpty.StudyDays)
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
