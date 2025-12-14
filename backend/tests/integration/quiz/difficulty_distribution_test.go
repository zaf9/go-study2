package quiz

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

type Question struct {
	Difficulty string `yaml:"difficulty"`
}

type QuizFile struct {
	Questions []Question `yaml:"questions"`
}

func Test_DifficultyDistributionPerFile(t *testing.T) {
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
		var q QuizFile
		if err := yaml.Unmarshal(data, &q); err != nil {
			t.Fatalf("unmarshal %s: %v", path, err)
		}
		total := len(q.Questions)
		if total == 0 {
			t.Fatalf("%s: no questions", path)
		}
		easy := 0
		medium := 0
		hard := 0
		for _, qq := range q.Questions {
			switch qq.Difficulty {
			case "easy":
				easy++
			case "medium":
				medium++
			case "hard":
				hard++
			}
		}
		// Targets: easy ~40%, medium ~40%, hard ~20% with ±10%
		if float64(easy) < float64(total)*0.3 || float64(easy) > float64(total)*0.5 {
			t.Fatalf("%s: easy %d/%d outside 40%%±10%%", path, easy, total)
		}
		if float64(medium) < float64(total)*0.3 || float64(medium) > float64(total)*0.5 {
			t.Fatalf("%s: medium %d/%d outside 40%%±10%%", path, medium, total)
		}
		if float64(hard) < float64(total)*0.1 || float64(hard) > float64(total)*0.3 {
			t.Fatalf("%s: hard %d/%d outside 20%%±10%%", path, hard, total)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk quiz_data: %v", err)
	}
}
