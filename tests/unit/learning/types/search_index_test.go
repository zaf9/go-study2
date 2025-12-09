package types_test

import (
	"testing"

	"go-study2/src/learning/types"
)

// T021/T036: 搜索索引覆盖与边界校验
func TestSearchIndexCoverage(t *testing.T) {
	keywords := []string{"map key", "~int", "interface nil", "slice share", "array length"}
	for _, kw := range keywords {
		results, err := types.SearchReferences(kw)
		if err != nil {
			t.Fatalf("检索失败 %s: %v", kw, err)
		}
		if len(results) == 0 {
			t.Fatalf("关键词 %s 未返回结果", kw)
		}
	}
}
