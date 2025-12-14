package quiz_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

type Question struct {
	ID          string   `yaml:"id"`
	Type        string   `yaml:"type"`
	Difficulty  string   `yaml:"difficulty"`
	Stem        string   `yaml:"stem"`
	Options     []string `yaml:"options"`
	Answer      string   `yaml:"answer"`
	Explanation string   `yaml:"explanation"`
	Topic       string   `yaml:"topic"`
	Chapter     string   `yaml:"chapter"`
}

type QuizFile struct {
	Questions []Question `yaml:"questions"`
}

func Test_YAMLFilesHaveRequiredFields(t *testing.T) {
	base := findQuizDataPath()
	entries := []string{}
	err := filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
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
	if err != nil {
		t.Fatalf("walk quiz_data failed: %v", err)
	}
	if len(entries) == 0 {
		t.Fatalf("no yaml files found under %s", base)
	}

	for _, f := range entries {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		var q QuizFile
		if err := yaml.Unmarshal(data, &q); err != nil {
			t.Fatalf("unmarshal %s: %v", f, err)
		}
		if len(q.Questions) < 1 {
			t.Fatalf("file %s has no questions", f)
		}
		for i, qq := range q.Questions {
			if qq.ID == "" {
				t.Fatalf("file %s question[%d] missing id", f, i)
			}
			if qq.Type == "" {
				t.Fatalf("file %s question[%d] missing type", f, i)
			}
			if qq.Difficulty == "" {
				t.Fatalf("file %s question[%d] missing difficulty", f, i)
			}
			if qq.Stem == "" {
				t.Fatalf("file %s question[%d] missing stem", f, i)
			}
			if len(qq.Options) < 2 {
				t.Fatalf("file %s question[%d] has too few options", f, i)
			}
			if qq.Answer == "" {
				t.Fatalf("file %s question[%d] missing answer", f, i)
			}
			if qq.Explanation == "" {
				t.Fatalf("file %s question[%d] missing explanation", f, i)
			}
		}
	}
}

func findQuizDataPath() string {
	candidates := []string{
		filepath.Join("backend", "quiz_data"),
		filepath.Join("backend", "backend", "quiz_data"),
		filepath.Join("..", "backend", "quiz_data"),
		filepath.Join("..", "..", "backend", "quiz_data"),
		filepath.Join("..", "..", "..", "backend", "quiz_data"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	// fallback to original
	return filepath.Join("backend", "backend", "quiz_data")
}
