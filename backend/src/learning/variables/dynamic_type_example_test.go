package variables

import "fmt"

// ExampleLoadContent_dynamic 展示动态类型主题的标题。
func ExampleLoadContent_dynamic() {
	content, _ := LoadContent(TopicDynamic)
	fmt.Println(content.Title)
	// Output:
	// 接口动态类型与 nil
}
