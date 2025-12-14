package quiz

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

type QType struct {
	Type string `yaml:"type"`
}

type QFile struct {
	Questions []QType `yaml:"questions"`
}

func Test_QuestionCountAndTypeDistribution(t *testing.T) {
	base := findQuizDataPath()
	err := filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".yaml" && filepath.Ext(path) != ".yml" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		var q QFile
		if err := yaml.Unmarshal(data, &q); err != nil {
			t.Fatalf("unmarshal %s: %v", path, err)
		}
		if len(q.Questions) < 30 || len(q.Questions) > 50 {
			t.Fatalf("%s: question count %d outside 30-50", path, len(q.Questions))
		}
		single := 0
		multi := 0
		for _, qq := range q.Questions {
			if qq.Type == "single" {
				single++
			} else if qq.Type == "multiple" {
				multi++
			}
		}
		total := single + multi
		if total == 0 {
			t.Fatalf("%s: no single/multiple questions detected", path)
		}
		// target ~50% each, allow ±10%
		min := float64(total) * 0.4
		max := float64(total) * 0.6
		if float64(single) < min || float64(single) > max {
			t.Fatalf("%s: single count %d out of %d (expect ~50%% ±10%%)", path, single, total)
		}
		if float64(multi) < min || float64(multi) > max {
			t.Fatalf("%s: multiple count %d out of %d (expect ~50%% ±10%%)", path, multi, total)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk quiz_data: %v", err)
	}
}
