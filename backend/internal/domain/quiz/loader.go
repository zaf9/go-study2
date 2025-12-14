package quiz

import (
	"io/fs"
	"os"
	"path/filepath"

	"context"
	"time"

	"gopkg.in/yaml.v3"

	"go-study2/internal/infrastructure/logger"
)

// LoadAllBanks 从指定目录递归加载所有 YAML 题库文件到 repository
func LoadAllBanks(root string, repo *QuizRepository) error {
	start := time.Now()
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logger.LogError(context.Background(), err, "walkDir error", map[string]interface{}{"path": path})
			return err
		}
		if d.IsDir() {
			return nil
		}
		// 仅处理 .yaml/.yml 文件
		ext := filepath.Ext(path)
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			logger.LogWithFields(context.Background(), "ERROR", "quiz.load.open_failed", map[string]interface{}{"path": path, "error": err.Error()})
			return err
		}
		defer f.Close()

		// 解析为 yaml.Node 以便提取每题的行号信息
		var rootNode yaml.Node
		dec := yaml.NewDecoder(f)
		if err := dec.Decode(&rootNode); err != nil {
			logger.LogWithFields(context.Background(), "ERROR", "quiz.load.yaml_decode_failed", map[string]interface{}{"path": path, "error": err.Error()})
			return err
		}

		var bank YAMLBank
		if err := rootNode.Decode(&bank); err != nil {
			logger.LogWithFields(context.Background(), "ERROR", "quiz.load.node_decode_failed", map[string]interface{}{"path": path, "error": err.Error()})
			return err
		}
		// 推断 topic/chapter 来自路径: .../quiz_data/{topic}/{chapter}.yaml
		rel, err := filepath.Rel(root, path)
		if err != nil {
			rel = path
		}
		// 在 Windows filepath.SplitList 使用 ;，改为使用 filepath.ToSlash 并手动分割
		rel2 := filepath.ToSlash(rel)
		segs := splitPath(rel2)
		if len(segs) >= 2 {
			topic := segs[0]
			chapterWithExt := segs[len(segs)-1]
			chapter := chapterWithExt[:len(chapterWithExt)-len(filepath.Ext(chapterWithExt))]
			// 转换为 YAMLQuestion 类型，并收集每题的行号（若存在）
			yamlQs := make([]YAMLQuestion, 0, len(bank.Questions))
			lineMap := map[string]int{}
			// 在 rootNode 中寻找 questions 列表并提取每个 question node 的 id 与行号
			for i := 0; i < len(rootNode.Content); i++ {
				node := rootNode.Content[i]
				if node.Kind == yaml.MappingNode {
					for j := 0; j < len(node.Content); j += 2 {
						keyNode := node.Content[j]
						valNode := node.Content[j+1]
						if keyNode.Value == "questions" && valNode.Kind == yaml.SequenceNode {
							for _, qNode := range valNode.Content {
								// qNode is a mapping of fields for a question
								var qID string
								for k := 0; k < len(qNode.Content); k += 2 {
									kNode := qNode.Content[k]
									vNode := qNode.Content[k+1]
									if kNode.Value == "id" {
										qID = vNode.Value
										lineMap[qID] = vNode.Line
										break
									}
								}
							}
						}
					}
				}
			}

			for _, q := range bank.Questions {
				yamlQs = append(yamlQs, YAMLQuestion{
					ID:          q.ID,
					Type:        q.Type,
					Difficulty:  q.Difficulty,
					Stem:        q.Stem,
					Options:     q.Options,
					Answer:      q.Answer,
					Explanation: q.Explanation,
					Topic:       q.Topic,
					Chapter:     q.Chapter,
				})
			}

			// 在将题库加入 repo 之前进行验证并传入文件与行号映射
			if err := ValidateBankWithSource(bank, path, lineMap); err != nil {
				return err
			}

			repo.AddBank(topic, chapter, yamlQs)
			logger.LogWithFields(context.Background(), "INFO", "quiz.bank.loaded", map[string]interface{}{
				"file":    path,
				"topic":   topic,
				"chapter": chapter,
				"count":   len(yamlQs),
			})
		}
		return nil
	})

	if err == nil {
		logger.LogWithFields(context.Background(), "INFO", "quiz.load.complete", map[string]interface{}{
			"root":     root,
			"duration": time.Since(start).String(),
		})
	}
	return err
}

// splitPath 将 / 分隔路径分割为段
func splitPath(p string) []string {
	if p == "" {
		return []string{}
	}
	var out []string
	cur := ""
	for i := 0; i < len(p); i++ {
		if p[i] == '/' || p[i] == '\\' {
			if cur != "" {
				out = append(out, cur)
				cur = ""
			}
			continue
		}
		cur += string(p[i])
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}
