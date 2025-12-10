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
		LastVisit:    time.Now(),
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

