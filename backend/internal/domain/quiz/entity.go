package quiz

// YAMLQuestion 表示用于从 YAML 解码的题目结构，避免与数据库模型冲突。
type YAMLQuestion struct {
	ID          string   `yaml:"id" json:"id"`
	Type        string   `yaml:"type" json:"type"`             // single/multiple
	Difficulty  string   `yaml:"difficulty" json:"difficulty"` // easy/medium/hard
	Stem        string   `yaml:"stem" json:"stem"`
	Options     []string `yaml:"options" json:"options"`
	Answer      string   `yaml:"answer" json:"answer"`
	Explanation string   `yaml:"explanation" json:"explanation"`
	Topic       string   `yaml:"topic" json:"topic"`
	Chapter     string   `yaml:"chapter" json:"chapter"`
}

// YAMLBank 表示一个章节/文件内的题目集合（用于 YAML 文件）
type YAMLBank struct {
	Questions []YAMLQuestion `yaml:"questions" json:"questions"`
}

// QuizRepository 内存仓储，按 topic->chapter 映射（存放 YAMLQuestion 转换后的业务模型）
type QuizRepository struct {
	banks map[string]map[string][]YAMLQuestion // banks[topic][chapter] = []YAMLQuestion
}

// NewRepository 创建仓储
func NewRepository() *QuizRepository {
	return &QuizRepository{banks: make(map[string]map[string][]YAMLQuestion)}
}

// AddBank 将题目加入仓储
func (r *QuizRepository) AddBank(topic, chapter string, qs []YAMLQuestion) {
	if _, ok := r.banks[topic]; !ok {
		r.banks[topic] = make(map[string][]YAMLQuestion)
	}
	r.banks[topic][chapter] = qs
}

// GetBank 获取题库
func (r *QuizRepository) GetBank(topic, chapter string) ([]YAMLQuestion, bool) {
	if chs, ok := r.banks[topic]; ok {
		qs, ok2 := chs[chapter]
		return qs, ok2
	}
	return nil, false
}
