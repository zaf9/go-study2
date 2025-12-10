package variables

import "fmt"

// ExampleLoadContent_static 展示静态类型主题的标题。
func ExampleLoadContent_static() {
	content, _ := LoadContent(TopicStatic)
	fmt.Println(content.Title)
	// Output:
	// 静态类型与可赋值性
}
