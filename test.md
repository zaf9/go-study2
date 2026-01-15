# 测试

1. 进入前端目录，启动前端(npm run dev)
2. 进入后端目录，启动后端(go run main.go -d)
3. 使用chrome devtools mcp连接http://localhost:3000/。注意因为前端的url需要构建，所以等待比较长的时间是正常的，可能需要60-300秒, 所以需要等待一段时间再连接。

账号/密码：admin/GoStudy@789

4. 按照specs/016-progress-quiz-ux下的功能需求和方案设计，逐个在等了后的chrome前端页面点击验证。如果有任何错误或不匹配的地方，请修正。直到完全符合specs/016-progress-quiz-ux的需求

# speckit.specify 

注意：speckit的spec的前缀编号是: 018

为golang的Properties of types and values章节增加学习内容，参考types的所有实现的功能，为Properties of types and values章节实现同样类似的功能。

注意：Properties of types and values章节的每个子章节单独一个.go文件提供学习内容。如果Types章节的子章节再有子子章节，则子子章节仍然是一个.go文件的学习内容

Properties of types and values章节内容如下。如果需要，请联网查询获取足够支撑学golang的Properties of types and values章节的内容。