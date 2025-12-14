package quiz_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

type qQuestion struct {
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

type qFile struct {
	Questions []qQuestion `yaml:"questions"`
}

func Test_QuestionHasNineFields(t *testing.T) {
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
		t.Skip("no quiz yaml files found")
	}
	for _, f := range entries {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		var qq qFile
		if err := yaml.Unmarshal(data, &qq); err != nil {
			t.Fatalf("unmarshal %s: %v", f, err)
		}
		for i, q := range qq.Questions {
			if q.ID == "" || q.Type == "" || q.Difficulty == "" || q.Stem == "" || len(q.Options) < 2 || q.Answer == "" || q.Explanation == "" || q.Topic == "" || q.Chapter == "" {
				t.Fatalf("file %s question[%d] missing required fields", f, i)
			}
		}
	}
}
