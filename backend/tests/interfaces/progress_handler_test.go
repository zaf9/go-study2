package interfaces

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	progapp "go-study2/internal/app/progress"
	progressdom "go-study2/internal/domain/progress"
	httpif "go-study2/internal/interfaces/http"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gsession"
)

func TestProgressHandler_Flow(t *testing.T) {
	repo := newMemoryRepo()
	calc := progapp.NewCalculator(map[string]int{"variables": 30}, map[string]int{"variables": 2})
	calc.EstimatedDuration["variables/storage"] = 600
	service := progapp.NewService(repo, calc)
	handler := &httpif.ProgressHandler{Service: service}

	s := ghttp.GetServer(fmt.Sprintf("progress-handler-%d", time.Now().UnixNano()))
	RegisterTestRoutes(s, handler)
	defer s.Shutdown()

	reqBody := map[string]interface{}{
		"topic":             "variables",
		"chapter":           "storage",
		"read_duration":     120,
		"scroll_progress":   70,
		"last_position":     "120",
		"estimated_seconds": 600,
	}
	w := doRequest(t, s, "POST", "/api/v1/progress", reqBody)
	if w.Code != http.StatusOK {
		t.Fatalf("首次上报返回状态码异常: %d", w.Code)
	}
	var resp apiResp
	unmarshalBody(t, w, &resp)
	if resp.Code != 0 || resp.Data["status"] == "" {
		t.Fatalf("首次上报响应异常: %+v", resp)
	}

	reqBody["read_duration"] = 480
	reqBody["scroll_progress"] = 95
	reqBody["quiz_passed"] = true
	reqBody["quiz_score"] = 90
	w = doRequest(t, s, "POST", "/api/v1/progress", reqBody)
	unmarshalBody(t, w, &resp)
	if resp.Data["status"] != progressdom.StatusCompleted {
		t.Fatalf("应返回 completed，得到 %v", resp.Data["status"])
	}

	w = doRequest(t, s, "GET", "/api/v1/progress", nil)
	unmarshalBody(t, w, &resp)
	overall, ok := resp.Data["overall"].(map[string]interface{})
	if !ok || overall["progress"] == nil || overall["totalChapters"] == nil {
		t.Fatalf("整体进度格式不正确: %+v", resp.Data["overall"])
	}
	topicsRaw, ok := resp.Data["topics"].([]interface{})
	if !ok || len(topicsRaw) == 0 {
		t.Fatalf("主题汇总为空: %+v", resp.Data["topics"])
	}
	firstTopic, ok := topicsRaw[0].(map[string]interface{})
	if !ok || firstTopic["weight"] == nil || firstTopic["name"] == nil {
		t.Fatalf("主题字段缺失: %+v", firstTopic)
	}

	w = doRequest(t, s, "GET", "/api/v1/progress/variables", nil)
	unmarshalBody(t, w, &resp)
	topicInfo, _ := resp.Data["topic"].(map[string]interface{})
	var chapters []map[string]interface{}
	if rawChapters, ok := resp.Data["chapters"]; ok {
		raw, _ := json.Marshal(rawChapters)
		_ = json.Unmarshal(raw, &chapters)
	}
	if len(chapters) == 0 {
		t.Fatalf("章节列表为空: %+v", resp)
	}
	if topicInfo["id"] != "variables" {
		t.Fatalf("主题摘要缺失或不匹配: %+v", topicInfo)
	}
}

// TestProgressHandler_DashboardStats 测试 Dashboard 统计数据 API
// 验证 GET /api/v1/progress 返回的数据包含学习天数、完成章节数、总章节数等 Dashboard 所需字段
func TestProgressHandler_DashboardStats(t *testing.T) {
	repo := newMemoryRepo()
	calc := progapp.NewCalculator(
		map[string]int{"variables": 30, "constants": 20, "lexical_elements": 25, "types": 25},
		map[string]int{"variables": 2},
	)
	service := progapp.NewService(repo, calc)
	handler := &httpif.ProgressHandler{Service: service}

	s := ghttp.GetServer(fmt.Sprintf("dashboard-stats-%d", time.Now().UnixNano()))
	RegisterTestRoutes(s, handler)
	defer s.Shutdown()

	// 准备测试数据：不同日期的学习记录
	now := time.Now()
	day1 := now.AddDate(0, 0, -5)
	day2 := now.AddDate(0, 0, -3)
	day3 := now.AddDate(0, 0, -1)

	// 创建不同日期的学习进度记录
	repo.CreateOrUpdate(ctx(), &progressdom.LearningProgress{
		UserID:       1,
		Topic:         "variables",
		Chapter:       "storage",
		Status:        progressdom.StatusCompleted,
		ReadDuration:  600,
		ScrollProgress: 100,
		QuizPassed:    true,
		QuizScore:     90,
		LastVisitAt:   day1,
	})
	repo.CreateOrUpdate(ctx(), &progressdom.LearningProgress{
		UserID:       1,
		Topic:         "variables",
		Chapter:       "pointer",
		Status:        progressdom.StatusInProgress,
		ReadDuration:  120,
		ScrollProgress: 50,
		LastVisitAt:   day2,
	})
	repo.CreateOrUpdate(ctx(), &progressdom.LearningProgress{
		UserID:       1,
		Topic:         "constants",
		Chapter:       "iota",
		Status:        progressdom.StatusCompleted,
		ReadDuration:  300,
		ScrollProgress: 100,
		QuizPassed:    true,
		QuizScore:     85,
		LastVisitAt:   day3,
	})

	// 调用 GET /api/v1/progress API
	w := doRequest(t, s, "GET", "/api/v1/progress", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("API 返回状态码异常: %d, body: %s", w.Code, w.Body.String())
	}

	var resp apiResp
	unmarshalBody(t, w, &resp)
	if resp.Code != 0 {
		t.Fatalf("API 返回错误: code=%d, message=%s", resp.Code, resp.Message)
	}

	// 验证 overall 字段存在且包含 Dashboard 所需字段
	overall, ok := resp.Data["overall"].(map[string]interface{})
	if !ok {
		t.Fatalf("overall 字段缺失或格式不正确: %+v", resp.Data)
	}

	// 验证学习天数（StudyDays）
	studyDays, ok := overall["studyDays"].(float64)
	if !ok {
		t.Fatalf("studyDays 字段缺失或格式不正确: %+v", overall)
	}
	if studyDays != 3 {
		t.Fatalf("学习天数应为 3（3个不同日期），得到 %.0f", studyDays)
	}

	// 验证完成章节数（CompletedChapters）
	completedChapters, ok := overall["completedChapters"].(float64)
	if !ok {
		t.Fatalf("completedChapters 字段缺失或格式不正确: %+v", overall)
	}
	if completedChapters != 2 {
		t.Fatalf("完成章节数应为 2，得到 %.0f", completedChapters)
	}

	// 验证总章节数（TotalChapters）
	totalChapters, ok := overall["totalChapters"].(float64)
	if !ok {
		t.Fatalf("totalChapters 字段缺失或格式不正确: %+v", overall)
	}
	if totalChapters < 2 {
		t.Fatalf("总章节数应至少为 2，得到 %.0f", totalChapters)
	}

	// 验证整体进度（Progress）
	progress, ok := overall["progress"].(float64)
	if !ok {
		t.Fatalf("progress 字段缺失或格式不正确: %+v", overall)
	}
	if progress < 0 || progress > 100 {
		t.Fatalf("整体进度应在 0-100 之间，得到 %.2f", progress)
	}

	// 验证 topics 数组存在
	topicsRaw, ok := resp.Data["topics"].([]interface{})
	if !ok {
		t.Fatalf("topics 字段缺失或格式不正确: %+v", resp.Data)
	}
	if len(topicsRaw) == 0 {
		t.Fatalf("topics 数组不应为空")
	}

	// 验证无学习记录时的响应（使用不同的用户ID）
	emptyRepo := newMemoryRepo()
	emptyService := progapp.NewService(emptyRepo, calc)
	emptyHandler := &httpif.ProgressHandler{Service: emptyService}
	emptyS := ghttp.GetServer(fmt.Sprintf("dashboard-stats-empty-%d", time.Now().UnixNano()))
	RegisterTestRoutes(emptyS, emptyHandler)
	defer emptyS.Shutdown()

	wEmpty := doRequest(t, emptyS, "GET", "/api/v1/progress", nil)
	if wEmpty.Code != http.StatusOK {
		t.Fatalf("无学习记录时 API 应返回成功: code=%d", wEmpty.Code)
	}
	var respEmpty apiResp
	unmarshalBody(t, wEmpty, &respEmpty)
	if respEmpty.Code != 0 {
		t.Fatalf("无学习记录时 API 应返回成功: code=%d", respEmpty.Code)
	}
	overallEmpty, ok := respEmpty.Data["overall"].(map[string]interface{})
	if !ok {
		t.Fatalf("无学习记录时 overall 字段应存在: %+v", respEmpty.Data)
	}
	studyDaysEmpty, _ := overallEmpty["studyDays"].(float64)
	if studyDaysEmpty != 0 {
		t.Fatalf("无学习记录时学习天数应为 0，得到 %.0f", studyDaysEmpty)
	}
}

func RegisterTestRoutes(s *ghttp.Server, handler *httpif.ProgressHandler) {
	sessionStore := gsession.NewStorageMemory()
	s.SetConfig(ghttp.ServerConfig{
		SessionMaxAge: 24 * 3600,
	})
	s.SetSessionStorage(sessionStore)
	httpif.RegisterProgressRoutes(s, handler)
	s.SetPort(0)
	go s.Start()
	time.Sleep(50 * time.Millisecond)
}

type apiResp struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

func doRequest(t *testing.T, s *ghttp.Server, method, path string, body map[string]interface{}) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("编码请求体失败: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "1")
	w := httptest.NewRecorder()
	func() {
		defer func() {
			_ = recover()
		}()
		s.ServeHTTP(w, req)
	}()
	return w
}

func unmarshalBody(t *testing.T, w *httptest.ResponseRecorder, resp *apiResp) {
	t.Helper()
	if err := json.Unmarshal(w.Body.Bytes(), resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
}

// memoryRepo 实现 ProgressRepository，供接口测试使用。
type memoryRepo struct {
	data map[string]progressdom.LearningProgress
}

func newMemoryRepo() *memoryRepo {
	return &memoryRepo{data: map[string]progressdom.LearningProgress{}}
}

func (m *memoryRepo) key(userID int64, topic, chapter string) string {
	return strings.Join([]string{strconv.FormatInt(userID, 10), topic, chapter}, "|")
}

func (m *memoryRepo) CreateOrUpdate(ctx context.Context, record *progressdom.LearningProgress) error {
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
		if existing.QuizScore > 0 && record.QuizScore == 0 {
			record.QuizScore = existing.QuizScore
		}
		if existing.QuizPassed && !record.QuizPassed {
			record.QuizPassed = true
		}
	}
	m.data[key] = *record
	return nil
}

func (m *memoryRepo) Get(ctx context.Context, userID int64, topic, chapter string) (*progressdom.LearningProgress, error) {
	key := m.key(userID, topic, chapter)
	if v, ok := m.data[key]; ok {
		cp := v
		return &cp, nil
	}
	return nil, nil
}

func (m *memoryRepo) GetByUser(ctx context.Context, userID int64) ([]progressdom.LearningProgress, error) {
	var list []progressdom.LearningProgress
	for _, v := range m.data {
		if v.UserID == userID {
			list = append(list, v)
		}
	}
	return list, nil
}

func (m *memoryRepo) GetByTopic(ctx context.Context, userID int64, topic string) ([]progressdom.LearningProgress, error) {
	var list []progressdom.LearningProgress
	for _, v := range m.data {
		if v.UserID == userID && v.Topic == topic {
			list = append(list, v)
		}
	}
	return list, nil
}
