package entity

import "time"

// QuizQuestion 是题目信息的物理模型（映射数据库表 quiz_questions）。
type QuizQuestion struct {
	Id             int64      `json:"id"             orm:"id"              description:"自增ID"`
	Topic          string     `json:"topic"          orm:"topic"           description:"主题标识"`
	Chapter        string     `json:"chapter"        orm:"chapter"         description:"章节标识"`
	Type           string     `json:"type"           orm:"type"            description:"题型 (single/multiple/truefalse...)"`
	Difficulty     string     `json:"difficulty"     orm:"difficulty"      description:"难度 (easy/medium/hard)"`
	Question       string     `json:"question"       orm:"question"        description:"题干文本"`
	Options        string     `json:"options"        orm:"options"         description:"选项 (JSON 数组文本)"`
	CorrectAnswers string     `json:"correctAnswers" orm:"correct_answers" description:"正确答案 (JSON 数组文本)"`
	Explanation    string     `json:"explanation"    orm:"explanation"     description:"解析内容"`
	CodeSnippet    string     `json:"codeSnippet"    orm:"code_snippet"    description:"代码片段 (可选)"`
	CreatedAt      *time.Time `json:"createdAt"      orm:"created_at"      description:"创建时间"`
	UpdatedAt      *time.Time `json:"updatedAt"      orm:"updated_at"      description:"更新时间"`
}
