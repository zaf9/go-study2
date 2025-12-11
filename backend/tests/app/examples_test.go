package app

import (
	"fmt"
)

// Example 展示进度服务输出示例，便于文档化与 go test 示例收集。
func Example() {
	fmt.Println("Progress service 示例：传入用户ID与章节，返回状态与汇总。")
	// Output:
	// Progress service 示例：传入用户ID与章节，返回状态与汇总。
}

// Example_quiz 展示测验服务流程示例。
func Example_quiz() {
	fmt.Println("Quiz service 示例：先获取题目 session，再提交答案获得得分与详情。")
	// Output:
	// Quiz service 示例：先获取题目 session，再提交答案获得得分与详情。
}
