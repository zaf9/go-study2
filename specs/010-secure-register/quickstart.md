# Quickstart - 010-secure-register

## 前置
- 后端：Go 1.21+，按仓库 README 运行 `./build.bat`（若存在）或 `go test ./...`，再启动。
- 前端：`npm install`，`npm run build`/`npm run dev`；确保与后端同源或已配置 CORS。

## 本地验证步骤
1) 启动后端，确保数据库可写。  
2) 检查默认管理员：启动完成后请求登录接口，以 admin/gostudy@123 登录，应返回 needPasswordChange=true。  
3) 登录后访问注册接口（携带 JWT）尝试创建新用户，未携带或非管理员应被拒绝。  
4) 强制改密：使用默认口令登录后被引导到改密页，设置符合策略的新密码；改密后旧口令不可再登录。  
5) 重新登录：用新密码登录成功，注册接口可用，权限正常。  
6) 负例：使用长度<8或缺少特殊字符的密码注册/改密，应被前后端共同拒绝。

## 产物位置
- 规范：`specs/010-secure-register/spec.md`
- 计划：`specs/010-secure-register/plan.md`
- 数据模型：`specs/010-secure-register/data-model.md`
- 研究：`specs/010-secure-register/research.md`
- 合同（如有）：`specs/010-secure-register/contracts/`

