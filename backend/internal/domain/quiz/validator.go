package quiz

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go-study2/internal/infrastructure/logger"
)

// ValidateQuestion 做基础验证，返回错误（仅针对 YAMLQuestion）
func ValidateQuestion(q YAMLQuestion) error {
	if strings.TrimSpace(q.ID) == "" {
		return errors.New("题目缺少 id")
	}
	if q.Type != "single" && q.Type != "multiple" {
		return fmt.Errorf("题目 %s 类型非法: %s", q.ID, q.Type)
	}
	if q.Difficulty != "easy" && q.Difficulty != "medium" && q.Difficulty != "hard" {
		return fmt.Errorf("题目 %s 难度非法: %s", q.ID, q.Difficulty)
	}
	if strings.TrimSpace(q.Stem) == "" {
		return fmt.Errorf("题目 %s 题干为空", q.ID)
	}
	if len(q.Options) < 2 {
		return fmt.Errorf("题目 %s 选项不足（至少2个）", q.ID)
	}
	if strings.TrimSpace(q.Answer) == "" {
		return fmt.Errorf("题目 %s 答案为空", q.ID)
	}
	// 简单答案格式校验：single -> 1 字母，multiple -> 2-4 字母
	ansLen := len(q.Answer)
	if q.Type == "single" && ansLen != 1 {
		return fmt.Errorf("题目 %s 单选答案应为单个字母", q.ID)
	}
	if q.Type == "multiple" && (ansLen < 2 || ansLen > 4) {
		return fmt.Errorf("题目 %s 多选答案应为2-4个字母", q.ID)
	}
	return nil
}

// ValidateBank 验证整个题库文件（YAML），不包含文件/行号信息。
// 向后兼容：内部调用会转发到 ValidateBankWithSource
func ValidateBank(bank YAMLBank) error {
	return ValidateBankWithSource(bank, "", nil)
}

// ValidateBankWithSource 验证整个题库文件，若传入 file 与 lineMap, 则在日志/错误中包含文件名与题目行号。
// lineMap 可选：映射 questionID -> yaml line number
func ValidateBankWithSource(bank YAMLBank, file string, lineMap map[string]int) error {
	seen := map[string]struct{}{}
	for _, q := range bank.Questions {
		if err := ValidateQuestion(q); err != nil {
			logValidationError(file, lineMap, q.ID, "quiz.validate.question", err)
			return fmt.Errorf("验证题目失败: %s", err.Error())
		}
		if _, ok := seen[q.ID]; ok {
			err := fmt.Errorf("题目 id 重复: %s", q.ID)
			logValidationError(file, lineMap, q.ID, "quiz.validate.bank", err)
			return err
		}
		seen[q.ID] = struct{}{}
	}
	return nil
}

// logValidationError 将验证错误以结构化方式写入日志（包含可选的文件与行号）
func logValidationError(file string, lineMap map[string]int, qid string, event string, err error) {
	fields := map[string]interface{}{"id": qid, "error": err.Error()}
	if file != "" {
		fields["file"] = file
	}
	if lineMap != nil {
		if ln, ok := lineMap[qid]; ok {
			fields["line"] = ln
		}
	}
	logger.LogWithFields(context.Background(), "ERROR", event, fields)
}
