package variables

import "fmt"

// ExampleLoadContent_zero 展示零值主题的标题。
func ExampleLoadContent_zero() {
	content, _ := LoadContent(TopicZero)
	fmt.Println(content.Title)
	// Output:
	// 零值与取值规则
}
