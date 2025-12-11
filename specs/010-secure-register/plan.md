# Implementation Plan: 安全注册与默认管理员（登录后强制改密）

**Branch**: `010-secure-register` | **Date**: 2025-12-11 | **Spec**: [spec.md](./spec.md)  
**Input**: Feature specification from `/specs/010-secure-register/spec.md`

## Summary
- 目标：仅允许已登录且具管理员权限的用户调用注册接口；服务启动时若无 admin 则自动创建默认管理员（admin/gostudy@123），首登强制改密；登录后若标记需改密，必须完成改密后重新登录。
- 口令策略统一：至少8位，且包含大小写字母、数字与特殊字符；需修正现有前后端校验以支持特殊字符。

## Technical Context

**Language/Version**: Go 1.21+（后端现有 GoFrame），TypeScript/React (Next.js) 前端  
**Primary Dependencies**: GoFrame v2（HTTP、JWT 中间件）、golang-jwt/jwt、bcrypt；前端 Ant Design、Axios  
**Storage**: 现有 SQLite（后端已用），复用用户表与刷新令牌表  
**Testing**: Go `testing` + testify；前端 Jest + React Testing Library；覆盖率≥80%  
**Target Platform**: 后端 Linux/Windows；前端现代浏览器  
**Project Type**: Web（backend + frontend）  
**Performance Goals**: 认证接口 p95 < 200ms（不含网络）；前端强制改密流程首屏 < 2s  
**Constraints**: 单实例；前后端同端口静态托管；密码校验需前后端一致  
**Scale/Scope**: 并发 100-1000 用户；默认账号仅一条；注册/改密流量低但需高安全性

## Constitution Check
- Principle I-VII/IX/X: 方案保持单一职责、显式错误处理、最小依赖、口令校验与鉴权完整，HTTPS 与审计要求在规格中已有。
- Principle VIII/XVIII/XIX: 继续使用既有 GoFrame 目录和包 README 约定。
- Principle XI: 完成后更新 README 中的认证与默认管理员说明。
- Principle XVI: 校验与强制改密逻辑采用卫语句，避免深层嵌套。
- Principle XX/XXXVI: 计划 gofmt/vet/lint 与前端 ESLint/Jest，覆盖率≥80%。
- 其余原则无新增风险。

## Project Structure

```text
specs/010-secure-register/
├── plan.md
├── spec.md
├── research.md
├── data-model.md
├── quickstart.md
└── contracts/

backend/
└── internal/
    ├── app/http_server/handler/auth.go        # 增强注册/登录/改密校验与默认管理员逻辑
    ├── domain/user/service.go                 # 口令校验、默认管理员创建、强制改密标记
    ├── infrastructure/repository/..           # 用户/令牌存取
    ├── pkg/password/                          # 口令策略校验
    └── ... (保持现有结构)

frontend/
└── components/auth/                           # LoginForm/RegisterForm 改为8位+特殊字符校验
    └── pages/(auth)/...                       # 登录后强制改密页面与跳转

tests/
├── backend/tests/contract|integration|unit    # 覆盖注册/登录/改密/默认管理员
└── frontend/tests/                            # 覆盖表单校验与强制改密流
```

**Structure Decision**: 继续采用现有前后端分目录与测试分层；不新增子项目。

## Complexity Tracking

无新增复杂度需要豁免；如后续需引入额外依赖，将在实现阶段另行备案。
