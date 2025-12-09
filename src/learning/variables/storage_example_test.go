package variables

import "fmt"

// ExampleLoadContent_storage 展示存储主题的标题。
func ExampleLoadContent_storage() {
	content, _ := LoadContent(TopicStorage)
	fmt.Println(content.Title)
	// Output:
	// 变量存储与取址
}
