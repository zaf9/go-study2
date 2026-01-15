package properties

// search.go 提供搜索索引的注册功能。
// 各子主题在注册内容时提供索引项，本文件提供统一的搜索入口。

// registerSearchIndex 注册全局搜索索引。
// 目前索引在各子主题注册时自动添加，本函数保留用于扩展。
func registerSearchIndex() {
	// 搜索索引已在 RegisterContent 中自动注册
	// 本函数保留用于未来的全局索引扩展
}
