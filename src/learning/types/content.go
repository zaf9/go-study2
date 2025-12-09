package types

// 内容注册入口：逐个子主题在对应文件中注册。
// 拆分文件便于后续扩充示例与规则。

func init() {
	registerAllContent()
}

func registerAllContent() {
	registerBoolean()
	registerNumeric()
	registerString()
	registerArray()
	registerSlice()
	registerStruct()
	registerPointer()
	registerFunction()
	registerInterfaceBasic()
	registerInterfaceEmbedded()
	registerInterfaceGeneral()
	registerInterfaceImpl()
	registerMapType()
	registerChannel()
	registerSearchIndex()
}
