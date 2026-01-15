package properties

// 内容注册入口：逐个子主题在对应文件中注册。
// 拆分文件便于后续扩充示例与规则。

func init() {
	registerAllContent()
}

func registerAllContent() {
	registerRepresentation()
	registerUnderlyingType()
	registerCoreType()
	registerTypeIdentity()
	registerAssignability()
	registerRepresentability()
	registerMethodSet()
	registerSearchIndex()
}
