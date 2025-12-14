package quiz_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

type exQuestion struct {
	Explanation string `yaml:"explanation"`
}

type exFile struct {
	Questions []exQuestion `yaml:"questions"`
}

func Test_ExplanationsContainChinese(t *testing.T) {
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
		var q exFile
		if err := yaml.Unmarshal(data, &q); err != nil {
			t.Fatalf("unmarshal %s: %v", f, err)
		}
		for i, qq := range q.Questions {
			if !containsChinese(qq.Explanation) {
				t.Fatalf("file %s question[%d] explanation does not contain Chinese characters", f, i)
			}
		}
	}
}

func containsChinese(s string) bool {
	for _, r := range s {
		if r >= 0x4e00 && r <= 0x9fff {
			return true
		}
	}
	return false
}
