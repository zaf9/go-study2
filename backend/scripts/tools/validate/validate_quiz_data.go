package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Question struct {
	ID         string `yaml:"id"`
	Type       string `yaml:"type"`
	Difficulty string `yaml:"difficulty"`
}

type QuizFile struct {
	Questions []Question `yaml:"questions"`
}

func main() {
	base := filepath.Join("backend", "backend", "quiz_data")
	ok := true
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
			fmt.Fprintf(os.Stderr, "read %s: %v\n", path, err)
			ok = false
			return nil
		}
		var q QuizFile
		if err := yaml.Unmarshal(data, &q); err != nil {
			fmt.Fprintf(os.Stderr, "unmarshal %s: %v\n", path, err)
			ok = false
			return nil
		}
		total := len(q.Questions)
		if total < 30 || total > 50 {
			fmt.Printf("FAIL %s: count %d outside 30-50\n", path, total)
			ok = false
		}
		single := 0
		multi := 0
		easy := 0
		medium := 0
		hard := 0
		for _, qq := range q.Questions {
			switch qq.Type {
			case "single":
				single++
			case "multiple":
				multi++
			}
			switch qq.Difficulty {
			case "easy":
				easy++
			case "medium":
				medium++
			case "hard":
				hard++
			}
		}
		if total > 0 {
			if float64(single) < float64(total)*0.4 || float64(single) > float64(total)*0.6 {
				fmt.Printf("FAIL %s: single %d/%d outside 50%%±10%%\n", path, single, total)
				ok = false
			}
			if float64(easy) < float64(total)*0.3 || float64(easy) > float64(total)*0.5 {
				fmt.Printf("FAIL %s: easy %d/%d outside 40%%±10%%\n", path, easy, total)
				ok = false
			}
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "walk error: %v\n", err)
		os.Exit(2)
	}
	if !ok {
		os.Exit(1)
	}
	fmt.Println("All validations passed")
}
