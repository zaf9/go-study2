package progress

import (
	"testing"
	"time"

	appconfig "go-study2/internal/config"
	progressdom "go-study2/internal/domain/progress"
)

// 测试：估算函数与完成判定边界
func TestEstimateAndBoundary(t *testing.T) {
	calc := NewCalculator(nil, nil)

	// 确认已为已知章节预计算预计时长
	est := calc.lookupDuration("lexical_elements", "comments")
	if est <= 0 {
		t.Fatalf("预计时长应为正值，got=%d", est)
	}
	if est < 60 || est > 3600 {
		t.Fatalf("预计时长应在合理范围内（60-3600），got=%d", est)
	}

	// 计算所需阈值（50%）
	needed := int64(float64(est) * 0.5)
	if needed <= 0 {
		t.Fatalf("计算得到的 needed 不应小于等于 0")
	}

	// 边界测试：scroll=100 但停留时间不足 -> 不应为 completed
	p := progressdom.LearningProgress{
		Topic:          "lexical_elements",
		Chapter:        "comments",
		ScrollProgress: 100,
		ReadDuration:   needed - 1,
		LastVisitAt:    time.Now(),
	}
	status := calc.CalculateChapterStatus(p, 0)
	if status == progressdom.StatusCompleted {
		t.Errorf("停留不足但滚动为100，状态不应为 completed")
	}

	// 现在满足停留阈值 -> 应为 completed
	p.ReadDuration = needed
	status = calc.CalculateChapterStatus(p, 0)
	if status != progressdom.StatusCompleted {
		t.Errorf("满足停留阈值后应为 completed，got=%s", status)
	}

	// 测验覆盖情形：QuizPassed + Scroll >=90 + ReadDuration >= needed -> completed
	p2 := progressdom.LearningProgress{
		Topic:          "lexical_elements",
		Chapter:        "comments",
		ScrollProgress: 90,
		ReadDuration:   needed,
		QuizPassed:     true,
		LastVisitAt:    time.Now(),
	}
	status = calc.CalculateChapterStatus(p2, 0)
	if status != progressdom.StatusCompleted {
		t.Errorf("测验通过且阅读足够应为 completed，got=%s", status)
	}

	// 测验已提交但未通过 -> tested
	p3 := progressdom.LearningProgress{
		Topic:          "lexical_elements",
		Chapter:        "comments",
		ScrollProgress: 50,
		ReadDuration:   0,
		QuizScore:      60,
		QuizPassed:     false,
		LastVisitAt:    time.Now(),
	}
	status = calc.CalculateChapterStatus(p3, 0)
	if status != progressdom.StatusTested {
		t.Errorf("测验已提交但未通过应为 tested，got=%s", status)
	}
}

// TestConfigBinding 验证 NewCalculator 会从 configs/config.yaml 的 progress 段读取配置并覆盖默认值
func TestConfigBinding(t *testing.T) {
	cfg, err := appconfig.Load()
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}
	calc := NewCalculator(nil, nil)

	// 仅当配置中存在 progress 段时才进行严格断言
	if cfg.Progress.ReadCharsPerSec > 0 {
		if calc.charsPerSec != cfg.Progress.ReadCharsPerSec {
			t.Fatalf("charsPerSec 未从配置绑定: want=%v got=%v", cfg.Progress.ReadCharsPerSec, calc.charsPerSec)
		}
	}
	if cfg.Progress.MinSeconds > 0 {
		if calc.minSeconds != cfg.Progress.MinSeconds {
			t.Fatalf("minSeconds 未从配置绑定: want=%v got=%v", cfg.Progress.MinSeconds, calc.minSeconds)
		}
	}
	if cfg.Progress.MaxSeconds > 0 {
		if calc.maxSeconds != cfg.Progress.MaxSeconds {
			t.Fatalf("maxSeconds 未从配置绑定: want=%v got=%v", cfg.Progress.MaxSeconds, calc.maxSeconds)
		}
	}
	if cfg.Progress.CompletionFraction > 0 {
		if calc.completionFraction != cfg.Progress.CompletionFraction {
			t.Fatalf("completionFraction 未从配置绑定: want=%v got=%v", cfg.Progress.CompletionFraction, calc.completionFraction)
		}
	}
}
