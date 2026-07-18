# 项目命名 MoJian（墨笺） — 改动总结

## 变更概览

项目正式命名为 **MoJian（墨笺）**，同步更新后端 module 路径、前端品牌标识及配置文件。

## 路由变更对照

无接口变更。

## 后端改动

### 1. Module 路径重命名 — `go.mod`
- `module github.com/goagent/blog` → `module github.com/goagent/mojian`

### 2. 全量 import 路径更新 — 21 个 Go 文件
- 所有 `github.com/goagent/blog/...` import 替换为 `github.com/goagent/mojian/...`
- 涉及文件：
  - `cmd/server/main.go`
  - `cmd/createadmin/main.go`
  - `internal/handler/` 下 6 个文件（article、captcha、category、comment、tag、user）
  - `internal/service/` 下 5 个文件（article、category、comment、tag、user）
  - `internal/repository/` 下 5 个文件（article、category、comment、tag、user）
  - `internal/middleware/auth.go`
  - `internal/database/database.go`
  - `internal/router/router.go`

### 3. 日志文件名更新 — `config/config.yaml`
- `logs/blog.log` → `logs/mojian.log`

## 前端改动

### 1. 品牌标识更新 — `frontend/src/components/Navbar.jsx`
- `Editorial Ink` → `墨笺MoJian`

### 2. 页面标题更新 — `frontend/index.html`
- `<title>Editorial Ink</title>` → `<title>墨笺 MoJian</title>`

### 3. 项目名更新 — `frontend/package.json`
- `"name": "blog-frontend"` → `"name": "mojian-frontend"`

### 4. 构建产物重新生成 — `frontend/dist/`
- 执行 `npm run build` 重新构建，dist 中品牌名已同步更新

## 数据库变更

无新增表、无新增字段、无数据迁移。

## 破坏性变更

- **Go module 路径变更**：`github.com/goagent/blog` → `github.com/goagent/mojian`，所有下游依赖需同步更新 import
- **日志文件路径变更**：`logs/blog.log` → `logs/mojian.log`，日志采集配置需同步调整

## 部署注意事项

- 部署前确认 `config/config.yaml` 中日志路径已更新
- 如有日志采集/监控服务引用了 `blog.log`，需同步修改采集路径
- 无其他特殊操作，可直接部署

## 遗留问题

无
