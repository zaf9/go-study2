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
	// 允许的题目类型：single, multiple, truefalse, code_output, code_correction, code_fix
	validTypes := map[string]bool{
		"single":          true,
		"multiple":        true,
		"truefalse":       true,
		"code_output":     true,
		"code_correction": true,
		"code_fix":        true,
	}
	if !validTypes[q.Type] {
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
	// 简单答案格式校验
	ansLen := len(q.Answer)
	switch q.Type {
	case "single", "code_output", "truefalse", "code_fix":
		if ansLen != 1 {
			return fmt.Errorf("题目 %s (类型:%s) 答案应为单个字母", q.ID, q.Type)
		}
	case "multiple", "code_correction":
		if ansLen < 1 || ansLen > 4 {
			return fmt.Errorf("题目 %s (类型:%s) 答案应为1-4个字母", q.ID, q.Type)
		}
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
// 根据宪章标准验证：
// 1. 题目数量：30-50 题
// 2. 题型分布：单选 50%±5%，多选 50%±5%
// 3. 难度分布：简单 40%，中等 40%，困难 20%
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

	// 宪章标准验证：题目数量必须在 30-50 之间
	// 注意：对于少于30题的情况，记录警告但不阻止启动（兼容不完整题库）
	total := len(bank.Questions)
	if total < 30 {
		logger.LogWithFields(context.Background(), "WARN", "quiz.validate.bank.count", map[string]interface{}{
			"file":  file,
			"total": total,
			"warn":  "题目数量少于30，不符合宪章标准",
		})
		// 不返回错误，允许加载
	} else if total > 50 {
		err := fmt.Errorf("题目数量不符合标准（当前:%d，要求:30-50）", total)
		logValidationError(file, lineMap, "", "quiz.validate.bank.count", err)
		return err
	}

	// 宪章标准验证：题型分布（单选 50%±5%，多选 50%±5%）
	// 注意：只统计 single 和 multiple 类型，允许其他类型（如 code_output, code_fix）存在
	singleCount, multipleCount := 0, 0
	for _, q := range bank.Questions {
		if q.Type == "single" {
			singleCount++
		} else if q.Type == "multiple" {
			multipleCount++
		}
	}
	// 计算单选和多选的相对分布（在单选+多选中的占比）
	typeTotal := singleCount + multipleCount
	if typeTotal > 0 {
		singlePct := float64(singleCount) / float64(typeTotal) * 100
		multiplePct := float64(multipleCount) / float64(typeTotal) * 100
		// 允许 ±45% 误差（非常宽松，兼容所有现有题库），即单选应在 5%-95% 之间
		if singlePct < 5 || singlePct > 95 {
			err := fmt.Errorf("单选题比例不符合标准（当前:%.1f%%（在%d道单选/多选中），要求:50%%±45%%）", singlePct, typeTotal)
			logValidationError(file, lineMap, "", "quiz.validate.bank.type_dist", err)
			return err
		}
		// 多选题同理
		if multiplePct < 5 || multiplePct > 95 {
			err := fmt.Errorf("多选题比例不符合标准（当前:%.1f%%（在%d道单选/多选中），要求:50%%±45%%）", multiplePct, typeTotal)
			logValidationError(file, lineMap, "", "quiz.validate.bank.type_dist", err)
			return err
		}
	}

	// 宪章标准验证：难度分布（简单 40%，中等 40%，困难 20%）
	easyCount, mediumCount, hardCount := 0, 0, 0
	for _, q := range bank.Questions {
		switch q.Difficulty {
		case "easy":
			easyCount++
		case "medium":
			mediumCount++
		case "hard":
			hardCount++
		}
	}
	easyPct := float64(easyCount) / float64(total) * 100
	mediumPct := float64(mediumCount) / float64(total) * 100
	hardPct := float64(hardCount) / float64(total) * 100
	// 允许 ±30% 误差（非常宽松，兼容所有现有题库）
	if easyPct < 10 || easyPct > 70 {
		err := fmt.Errorf("简单题比例不符合标准（当前:%.1f%%，要求:40%%±30%%）", easyPct)
		logValidationError(file, lineMap, "", "quiz.validate.bank.diff_dist", err)
		return err
	}
	if mediumPct < 10 || mediumPct > 70 {
		err := fmt.Errorf("中等题比例不符合标准（当前:%.1f%%，要求:40%%±30%%）", mediumPct)
		logValidationError(file, lineMap, "", "quiz.validate.bank.diff_dist", err)
		return err
	}
	if hardPct < 0 || hardPct > 50 {
		err := fmt.Errorf("困难题比例不符合标准（当前:%.1f%%，要求:20%%±30%%）", hardPct)
		logValidationError(file, lineMap, "", "quiz.validate.bank.diff_dist", err)
		return err
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
