package quiz_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type ansQuestion struct {
	Type   string `yaml:"type"`
	Answer string `yaml:"answer"`
}

type ansFile struct {
	Questions []ansQuestion `yaml:"questions"`
}

func Test_AnswerFormat(t *testing.T) {
	base := findQuizDataPath()
	var entries []string
	_ = filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {
			entries = append(entries, path)
		}
		return nil
	})
	if len(entries) == 0 {
		t.Skip("no yaml files")
	}
	for _, f := range entries {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		var q ansFile
		if err := yaml.Unmarshal(data, &q); err != nil {
			t.Fatalf("unmarshal %s: %v", f, err)
		}
		for i, qq := range q.Questions {
			ans := strings.TrimSpace(strings.ToUpper(qq.Answer))
			if qq.Type == "single" {
				if len(ans) != 1 {
					t.Fatalf("file %s question[%d] single answer length invalid: %s", f, i, ans)
				}
			} else if qq.Type == "multiple" {
				if len(ans) < 2 || len(ans) > 4 {
					t.Fatalf("file %s question[%d] multiple answer length invalid: %s", f, i, ans)
				}
				// 检查是否升序
				parts := strings.Split(ans, "")
				sorted := append([]string(nil), parts...)
				sort.Strings(sorted)
				if strings.Join(sorted, "") != strings.Join(parts, "") {
					t.Fatalf("file %s question[%d] multiple answer not ascending: %s", f, i, ans)
				}
			}
		}
	}
}
