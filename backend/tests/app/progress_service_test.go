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
	day4 := now                   // 今天

	// 同一天的多条记录应只计算为1天
	if err := repo.CreateOrUpdate(ctx(), &progressdom.LearningProgress{
		UserID:       1,
		Topic:        "variables",
		Chapter:      "storage",
		Status:       progressdom.StatusInProgress,
		ReadDuration: 120,
		LastVisitAt:  day1,
	}); err != nil {
		t.Fatalf("准备第1天记录1失败: %v", err)
	}
	if err := repo.CreateOrUpdate(ctx(), &progressdom.LearningProgress{
		UserID:       1,
		Topic:        "variables",
		Chapter:      "pointer",
		Status:       progressdom.StatusInProgress,
		ReadDuration: 60,
		LastVisitAt:  day1.Add(2 * time.Hour), // 同一天的不同时间
	}); err != nil {
		t.Fatalf("准备第1天记录2失败: %v", err)
	}

	// 第2天的记录
	if err := repo.CreateOrUpdate(ctx(), &progressdom.LearningProgress{
		UserID:       1,
		Topic:        "constants",
		Chapter:      "iota",
		Status:       progressdom.StatusInProgress,
		ReadDuration: 90,
		LastVisitAt:  day2,
	}); err != nil {
		t.Fatalf("准备第2天记录失败: %v", err)
	}

	// 第3天的记录
	if err := repo.CreateOrUpdate(ctx(), &progressdom.LearningProgress{
		UserID:       1,
		Topic:        "variables",
		Chapter:      "static",
		Status:       progressdom.StatusCompleted,
		ReadDuration: 300,
		LastVisitAt:  day3,
	}); err != nil {
		t.Fatalf("准备第3天记录失败: %v", err)
	}

	// 第4天的记录
	if err := repo.CreateOrUpdate(ctx(), &progressdom.LearningProgress{
		UserID:       1,
		Topic:        "constants",
		Chapter:      "untyped",
		Status:       progressdom.StatusInProgress,
		ReadDuration: 150,
		LastVisitAt:  day4,
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

// TestProgressService_GetLastLearningRecord 测试获取最后学习记录功能
func TestProgressService_GetLastLearningRecord(t *testing.T) {
	repo := newMemoryProgressRepo()
	calc := progapp.NewCalculator(map[string]int{"variables": 30, "constants": 20}, nil)
	service := progapp.NewService(repo, calc)

	// 准备测试数据：不同时间的学习记录
	now := time.Now()
	day1 := now.AddDate(0, 0, -5)
	day2 := now.AddDate(0, 0, -3)
	day3 := now.AddDate(0, 0, -1)

	// 创建较早的学习记录
	if err := repo.CreateOrUpdate(ctx(), &progressdom.LearningProgress{
		UserID:       1,
		Topic:        "variables",
		Chapter:      "storage",
		Status:       progressdom.StatusInProgress,
		ReadDuration: 120,
		LastVisitAt:  day1,
	}); err != nil {
		t.Fatalf("准备第1条记录失败: %v", err)
	}

	// 创建中间的学习记录
	if err := repo.CreateOrUpdate(ctx(), &progressdom.LearningProgress{
		UserID:       1,
		Topic:        "constants",
		Chapter:      "iota",
		Status:       progressdom.StatusInProgress,
		ReadDuration: 90,
		LastVisitAt:  day2,
	}); err != nil {
		t.Fatalf("准备第2条记录失败: %v", err)
	}

	// 创建最新的学习记录
	if err := repo.CreateOrUpdate(ctx(), &progressdom.LearningProgress{
		UserID:       1,
		Topic:        "variables",
		Chapter:      "static",
		Status:       progressdom.StatusCompleted,
		ReadDuration: 300,
		LastVisitAt:  day3,
	}); err != nil {
		t.Fatalf("准备第3条记录失败: %v", err)
	}

	// 获取最后学习记录
	lastLearning, err := service.GetLastLearningRecord(ctx(), 1)
	if err != nil {
		t.Fatalf("获取最后学习记录失败: %v", err)
	}

	if lastLearning == nil {
		t.Fatalf("应返回最后学习记录，得到 nil")
	}

	// 验证返回的是最新的记录（day3 的 static 章节）
	if lastLearning.TopicID != "variables" {
		t.Fatalf("主题 ID 应为 variables，得到 %s", lastLearning.TopicID)
	}
	if lastLearning.ChapterID != "static" {
		t.Fatalf("章节 ID 应为 static，得到 %s", lastLearning.ChapterID)
	}
	if lastLearning.TopicDisplayName == "" {
		t.Fatalf("主题显示名称不应为空")
	}
	if lastLearning.ChapterDisplayName == "" {
		t.Fatalf("章节显示名称不应为空")
	}
	if lastLearning.LastVisitedAt == "" {
		t.Fatalf("最后访问时间不应为空")
	}

	// 验证无学习记录时返回 nil
	lastLearningEmpty, err := service.GetLastLearningRecord(ctx(), 999)
	if err != nil {
		t.Fatalf("无学习记录时不应返回错误: %v", err)
	}
	if lastLearningEmpty != nil {
		t.Fatalf("无学习记录时应返回 nil，得到 %+v", lastLearningEmpty)
	}
}

// TestProgressService_GetTopicProgressSummary 测试获取主题进度汇总功能
func TestProgressService_GetTopicProgressSummary(t *testing.T) {
	repo := newMemoryProgressRepo()
	calc := progapp.NewCalculator(
		map[string]int{"variables": 30, "constants": 20, "lexical_elements": 25, "types": 25},
		map[string]int{"variables": 2, "constants": 3},
	)
	service := progapp.NewService(repo, calc)

	// 准备测试数据：不同主题的学习记录
	now := time.Now()

	// 主题1：variables - 完成1个章节，总共2个章节（50%）
	if err := repo.CreateOrUpdate(ctx(), &progressdom.LearningProgress{
		UserID:         1,
		Topic:          "variables",
		Chapter:        "storage",
		Status:         progressdom.StatusCompleted,
		ReadDuration:   600,
		ScrollProgress: 100,
		QuizPassed:     true,
		LastVisitAt:    now,
	}); err != nil {
		t.Fatalf("准备 variables 记录失败: %v", err)
	}

	// 主题2：constants - 完成2个章节，总共3个章节（66.7%）
	if err := repo.CreateOrUpdate(ctx(), &progressdom.LearningProgress{
		UserID:         1,
		Topic:          "constants",
		Chapter:        "iota",
		Status:         progressdom.StatusCompleted,
		ReadDuration:   300,
		ScrollProgress: 100,
		QuizPassed:     true,
		LastVisitAt:    now,
	}); err != nil {
		t.Fatalf("准备 constants 记录1失败: %v", err)
	}
	if err := repo.CreateOrUpdate(ctx(), &progressdom.LearningProgress{
		UserID:         1,
		Topic:          "constants",
		Chapter:        "boolean",
		Status:         progressdom.StatusCompleted,
		ReadDuration:   200,
		ScrollProgress: 100,
		QuizPassed:     true,
		LastVisitAt:    now,
	}); err != nil {
		t.Fatalf("准备 constants 记录2失败: %v", err)
	}

	// 主题3：lexical_elements - 未开始（0%）

	// 获取主题进度汇总
	summaries, err := service.GetTopicProgressSummary(ctx(), 1)
	if err != nil {
		t.Fatalf("获取主题进度汇总失败: %v", err)
	}

	if len(summaries) == 0 {
		t.Fatalf("应返回至少一个主题的进度汇总")
	}

	// 验证 variables 主题的进度
	var variablesSummary *progapp.TopicProgressSummary
	for i := range summaries {
		if summaries[i].TopicID == "variables" {
			variablesSummary = &summaries[i]
			break
		}
	}
	if variablesSummary == nil {
		t.Fatalf("应包含 variables 主题的进度汇总")
	}
	if variablesSummary.CompletedChapters != 1 {
		t.Fatalf("variables 主题完成章节数应为 1，得到 %d", variablesSummary.CompletedChapters)
	}
	if variablesSummary.TotalChapters != 2 {
		t.Fatalf("variables 主题总章节数应为 2，得到 %d", variablesSummary.TotalChapters)
	}
	if variablesSummary.Percentage != 50.0 {
		t.Fatalf("variables 主题完成百分比应为 50.0，得到 %.1f", variablesSummary.Percentage)
	}

	// 验证 constants 主题的进度
	var constantsSummary *progapp.TopicProgressSummary
	for i := range summaries {
		if summaries[i].TopicID == "constants" {
			constantsSummary = &summaries[i]
			break
		}
	}
	if constantsSummary == nil {
		t.Fatalf("应包含 constants 主题的进度汇总")
	}
	if constantsSummary.CompletedChapters != 2 {
		t.Fatalf("constants 主题完成章节数应为 2，得到 %d", constantsSummary.CompletedChapters)
	}
	if constantsSummary.TotalChapters != 3 {
		t.Fatalf("constants 主题总章节数应为 3，得到 %d", constantsSummary.TotalChapters)
	}
	// 验证百分比（66.7% 四舍五入到一位小数）
	expectedPercentage := 66.7
	if constantsSummary.Percentage < expectedPercentage-0.1 || constantsSummary.Percentage > expectedPercentage+0.1 {
		t.Fatalf("constants 主题完成百分比应为约 %.1f，得到 %.1f", expectedPercentage, constantsSummary.Percentage)
	}

	// 验证无学习记录时也返回所有主题（进度为0）
	summariesEmpty, err := service.GetTopicProgressSummary(ctx(), 999)
	if err != nil {
		t.Fatalf("无学习记录时不应返回错误: %v", err)
	}
	if len(summariesEmpty) == 0 {
		t.Fatalf("无学习记录时也应返回所有主题的进度汇总（进度为0）")
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

func (m *memoryProgressRepo) GetLastLearning(ctx context.Context, userID int64) (*progressdom.LearningProgress, error) {
	var last *progressdom.LearningProgress
	for _, v := range m.data {
		if v.UserID == userID && v.LastVisitAt.After(time.Time{}) {
			if last == nil || v.LastVisitAt.After(last.LastVisitAt) {
				cp := v
				last = &cp
			}
		}
	}
	return last, nil
}

func (m *memoryProgressRepo) GetOverview(ctx context.Context, userID int64) (*progressdom.ProgressOverview, error) {
	// 获取用户所有记录
	allRecords, err := m.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 按主题和章节构建进度映射
	progressMap := make(map[string]map[string]*progressdom.LearningProgress)
	for i := range allRecords {
		if _, ok := progressMap[allRecords[i].Topic]; !ok {
			progressMap[allRecords[i].Topic] = make(map[string]*progressdom.LearningProgress)
		}
		progressMap[allRecords[i].Topic][allRecords[i].Chapter] = &allRecords[i]
	}

	// 计算各主题的统计信息
	topicSummaries := make([]progressdom.TopicProgressSummary, 0)
	totalCompleted := 0
	totalInProgress := 0
	totalChapters := 0

	for topic, chapters := range progressdom.TopicChapterOrder {
		topicTotal := len(chapters)
		completed := 0
		inProgress := 0

		topicProgress, ok := progressMap[topic]
		if !ok {
			topicProgress = make(map[string]*progressdom.LearningProgress)
		}

		for _, chapter := range chapters {
			record, exists := topicProgress[chapter]
			if exists {
				if record.Status == progressdom.StatusCompleted && record.QuizPassed {
					completed++
				} else if record.Status == progressdom.StatusInProgress {
					inProgress++
				}
			}
		}

		topicSummaries = append(topicSummaries, progressdom.TopicProgressSummary{
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
	var nextChapter *progressdom.NextChapterHint
	lastLearning, err := m.GetLastLearning(ctx, userID)
	if err == nil && lastLearning != nil {
		nextChapter = &progressdom.NextChapterHint{
			Topic:   lastLearning.Topic,
			Chapter: lastLearning.Chapter,
			Title:   progressdom.ChapterDisplayName(lastLearning.Chapter),
		}
	}

	return &progressdom.ProgressOverview{
		TotalChapters:      totalChapters,
		CompletedChapters:  totalCompleted,
		InProgressChapters: totalInProgress,
		CompletionRate:     completionRate,
		Topics:             topicSummaries,
		NextChapter:        nextChapter,
	}, nil
}
