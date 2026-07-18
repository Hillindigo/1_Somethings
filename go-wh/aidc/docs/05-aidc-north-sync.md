# aidc-north-sync 项目分析

## 一、核心功能与整体架构

### 1.1 项目定位

`aidc-north-sync` 是从 `aidc-process-service` 中独立出来的**北向数据同步服务**，专注于从北向 API 拉取位置、设备、点位配置数据并同步到本地 PostgreSQL。相比 process-service 中内嵌的同步模块，north-sync 支持**多配置管理**、**动态调度**、**北向接口 CRUD** 等高级特性。

### 1.2 核心功能

- **多配置同步管理**：支持创建多个同步配置（`north_sync_config`），每个配置对应一个北向接口数据源
- **北向接口配置管理**：CRUD 管理北向接口信息（`north_configure`），包括 API URL、Token、数据源标识
- **动态定时调度**：每个配置可独立配置 Cron 表达式，运行时动态启停定时任务
- **手动触发同步**：通过 HTTP API 手动触发指定配置的同步任务
- **同步任务管理**：异步任务队列、任务状态查询、历史记录分页查询
- **版本控制增量同步**：基于版本号对比，只有远程版本变化时才执行全量同步
- **数据分类存储**：按节点类型（位置/设备/点位）分类存储到不同表

### 1.3 架构流程图

```text
  ┌──────────────────┐
  │ 前端 / API 调用  │
  └────────┬─────────┘
           │ HTTP REST API
           ▼
  ┌──────────────────────────────────────────────┐
  │            aidc-north-sync                    │
  │                                              │
  │  ┌────────────────────────────────────────┐  │
  │  │          Gin HTTP Router               │  │
  │  │                                        │  │
  │  │  配置管理:                              │  │
  │  │    GET/POST/PUT/DELETE /api/v1/*       │  │
  │  │    POST /api/v1/toggle-auto-sync      │  │
  │  │                                        │  │
  │  │  同步任务:                              │  │
  │  │    POST /api/v1/sync-now              │  │
  │  │    GET  /api/v1/status                │  │
  │  │    GET  /api/v1/history               │  │
  │  │                                        │  │
  │  │  北向配置:                              │  │
  │  │    /api/v1/north-configure/*          │  │
  │  └────────────┬───────────────────────────┘  │
  │               │                              │
  │  ┌────────────▼───────────────────────────┐  │
  │  │    Service 层                           │  │
  │  │  ┌───────────────┐ ┌────────────────┐  │  │
  │  │  │ConfigService  │ │NorthConfigure  │  │  │
  │  │  │• CRUD 配置    │ │Service         │  │  │
  │  │  │• 调度器联动   │ │• CRUD 北向接口 │  │  │
  │  │  └───────┬───────┘ └────────────────┘  │  │
  │  │          │                              │  │
  │  │  ┌───────▼────────────────────────────┐│  │
  │  │  │      Scheduler (动态多配置)        ││  │
  │  │  │  map[configName] → Cron Instance   ││  │
  │  │  │  • StartConfig / StopConfig        ││  │
  │  │  │  • InitFromDB (启动加载)           ││  │
  │  │  └───────┬────────────────────────────┘│  │
  │  │          │                              │  │
  │  │  ┌───────▼────────────────────────────┐│  │
  │  │  │    SyncTaskManager                 ││  │
  │  │  │  • 异步任务队列 (chan)             ││  │
  │  │  │  • 动态创建 SyncService            ││  │
  │  │  │  • 任务状态持久化                  ││  │
  │  │  └───────┬────────────────────────────┘│  │
  │  │          │                              │  │
  │  │  ┌───────▼────────────────────────────┐│  │
  │  │  │    SyncService                     ││  │
  │  │  │  • HTTP 请求北向 API               ││  │
  │  │  │  • 分页拉取数据                    ││  │
  │  │  │  • 分类存储（位置/设备/点位）      ││  │
  │  │  │  • 版本号对比                      ││  │
  │  │  └────────────────────────────────────┘│  │
  │  └────────────────────────────────────────┘  │
  └──────────────────────────────────────────────┘
           │
           ▼
  ┌──────────────┐     ┌─────────────────┐
  │ PostgreSQL   │     │ 北向 API 接口   │
  │ (配置+数据)  │     │ (设备配置源)    │
  └──────────────┘     └─────────────────┘
```

---

## 二、核心技术栈

| 分类 | 技术 | 用途 |
|------|------|------|
| 语言 | Go 1.22+ | 主语言 |
| HTTP 框架 | Gin | REST API 服务 |
| 数据库 | PostgreSQL + GORM | 配置数据 + 同步数据持久化 |
| 定时调度 | robfig/cron/v3 | 多配置动态定时任务 |
| 配置 | kelseyhightower/envconfig | 环境变量配置 |
| 日志 | klog | 结构化日志 |
| HTTP 客户端 | 自封装 utils.HTTPClient | 请求北向 API（超时+重试） |
| 部署 | Docker + Helm | K8s 部署 |

---

## 三、核心代码文件/目录及作用

```text
aidc-north-sync/
├── cmd/aidc-north-sync/
│   └── main.go                        # 服务入口（127 行，简洁清晰）
├── pkg/
│   ├── api/
│   │   ├── router.go                 # 路由注册（统一管理所有端点）
│   │   ├── handler.go                # 同步任务 Handler（SyncNow/Status/History）
│   │   ├── config_handler.go         # 配置管理 Handler（CRUD + ToggleAutoSync）
│   │   └── north_configure_handler.go # 北向配置 Handler（CRUD）
│   ├── base/
│   │   └── base.go                   # 基础组件初始化（PostgreSQL）
│   ├── config/
│   │   └── config.go                 # 配置结构（DB/Server/北向接口）
│   ├── dao/                           # 数据访问层（DAO）
│   │   ├── config_dao.go            # north_sync_config 表操作
│   │   ├── task_dao.go              # sync_task 表操作
│   │   ├── north_configure_dao.go   # north_configure 表操作
│   │   └── north_nodes_dao.go       # 位置/设备/点位表操作
│   ├── models/
│   │   ├── north_sync_config.go     # 同步配置模型
│   │   ├── north_configure.go       # 北向接口配置模型
│   │   ├── sync_task.go             # 同步任务模型
│   │   └── northnodes/              # 北向数据模型（Location/Device/Point）
│   ├── scheduler/
│   │   └── scheduler.go             # 动态多配置调度器（核心）
│   ├── service/
│   │   ├── config_service.go        # 配置管理服务（含调度器联动）
│   │   ├── north_configure_service.go # 北向配置管理服务
│   │   ├── sync_service.go          # 同步核心逻辑（490 行）
│   │   └── sync_task_manager.go     # 异步任务管理器（490 行）
│   └── utils/
│       └── httpclient.go            # HTTP 客户端（超时+重试）
├── migrations/                        # 数据库迁移脚本
├── Makefile                           # 构建脚本
└── Dockerfile                         # 容器镜像
```

---

## 四、核心业务逻辑分析

### 4.1 启动流程

```text
main()
  ├── envconfig.Process("AIDC_NORTH_SYNC")：加载配置
  ├── base.Init()：初始化 PostgreSQL
  ├── AutoMigrate()：自动建表
  │     └── NorthConfigure / NorthSyncConfig / SyncTask
  │
  ├── 创建 DAO 层实例
  │     ├── ConfigDAO（north_sync_config 表）
  │     ├── TaskDAO（sync_task 表）
  │     ├── NorthConfigureDAO（north_configure 表）
  │     └── NorthNodesDAO（位置/设备/点位表）
  │
  ├── 创建 SyncTaskManager → Start()
  │
  ├── 创建 Scheduler → InitFromDB()
  │     └── 遍历所有 enable_auto_sync=true 的配置
  │         └── 为每个配置启动独立的 Cron 定时任务
  │
  ├── 创建 Service 层 + Handler 层
  ├── 注册 Gin 路由
  ├── 启动 HTTP Server
  └── 信号处理 → 优雅关闭
```

### 4.2 动态多配置调度器（核心亮点）

```text
Scheduler
  │
  ├── entries: map[configName] → {cronInstance, configName, cronExpr}
  │
  ├── InitFromDB(ctx)
  │     └── 从数据库加载所有 enable_auto_sync=true 的配置
  │         └── 对每个配置调用 StartConfig()
  │
  ├── StartConfig(ctx, configName)
  │     ├── 如果已存在 → 先停止旧的
  │     ├── 从数据库读取配置（每次执行时重新读取最新）
  │     ├── 验证 enable_auto_sync + sync_cron
  │     ├── 创建独立的 cron.Cron 实例
  │     ├── AddFunc(cronExpr, 同步函数)
  │     │     └── 每次执行时从 DB 读取最新配置 → SubmitTask
  │     └── 启动 Cron
  │
  ├── StopConfig(configName)
  │     └── 停止并删除对应的 Cron 实例
  │
  └── StopAll()
        └── 停止所有 Cron 实例

关键设计：
  • 每个配置一个独立的 Cron 实例（互不影响）
  • 执行时从 DB 重新读取配置（热更新）
  • 通过 ConfigService 联动（创建/更新/删除配置时自动启停调度器）
```

### 4.3 同步任务管理器

```text
SyncTaskManager
  │
  ├── SubmitTask(config, force, triggeredBy)  [异步]
  │     ├── 查询 north_configure 获取北向接口配置
  │     ├── 创建 TaskRequest（含 API URL/Token/Source）
  │     ├── 创建 SyncTask 记录（status=pending）
  │     ├── 入队 taskQueue channel
  │     └── 异步消费 resultChan
  │
  ├── SubmitTaskAndWait(config, force, triggeredBy)  [同步等待]
  │     └── 同上，但等待 resultChan 返回结果
  │
  ├── processTaskQueue()  [后台 goroutine]
  │     └── 循环从 taskQueue 取任务 → executeTask()
  │
  └── executeTask(taskReq)
        ├── 更新状态 → running
        ├── createSyncServiceFromRequest()  [动态创建]
        │     └── 根据 taskReq 中的配置创建 SyncService 实例
        ├── syncService.SyncAllNorthNodesWithStats()
        │     └── 分页拉取 + 分类存储 + 版本对比
        ├── 更新状态 → completed / failed
        ├── 记录统计：locationCount/deviceCount/pointCount
        └── 通知 resultChan
```

### 4.4 配置管理与调度器联动

```text
ConfigService.Create(params)
  ├── 校验：enable_auto_sync=true 时 sync_cron 必填
  ├── 校验：north_configure_id 存在且 URL/Token 非空
  ├── 入库：configDAO.Create()
  └── 联动：reloadScheduler(configName)
        ├── enable_auto_sync=true → scheduler.StartConfig()
        └── enable_auto_sync=false → scheduler.StopConfig()

ConfigService.Update(params)
  ├── 更新数据库
  └── 联动：reloadScheduler(configName)

ConfigService.Delete(name)
  ├── 删除数据库记录
  └── 联动：scheduler.StopConfig(name)

ConfigService.ToggleAutoSync(params)
  ├── 更新 enable_auto_sync 字段
  └── 联动：
        ├── enable=true → scheduler.StartConfig()
        └── enable=false → scheduler.StopConfig()
```

### 4.5 数据库表结构

```text
north_configure（北向接口配置表）
  ├── id (PK)
  ├── machine_room      # 机房标识
  ├── north_api_url     # 北向 API 地址
  ├── north_token       # 认证 Token
  ├── north_source      # 数据源标识
  ├── version           # 当前版本号
  ├── topic             # 关联 Kafka topic
  └── parsing_method    # 数据解析方式

north_sync_config（同步配置表）
  ├── id (PK)
  ├── name (UNIQUE)     # 配置名称
  ├── sync_cron         # Cron 表达式
  ├── enable_auto_sync  # 是否自动同步
  ├── north_configure_id # 关联北向配置
  ├── sync_idc_info     # 是否同步 IDC 信息
  ├── http_timeout      # HTTP 超时
  └── http_retry        # HTTP 重试次数

sync_task（同步任务表）
  ├── task_id (UUID)
  ├── config_name       # 关联配置名
  ├── status            # pending/running/completed/failed
  ├── triggered_by      # auto/manual/api
  ├── location_count / device_count / point_count
  ├── started_at / completed_at / duration_seconds
  └── error_message
```

---

## 五、面试常见问题及回答思路

### Q1: 为什么要把同步服务从 process-service 独立出来？

> **回答**：
> 1. **职责分离**：process-service 的核心职责是 Kafka 消费和数据处理，同步逻辑是附加功能，独立后各服务更聚焦
> 2. **多配置支持**：独立后支持多个同步配置（多机房、多数据源），process-service 只支持单一配置
> 3. **独立扩缩容**：同步任务是 IO 密集型（HTTP 请求），与消费处理的资源需求不同
> 4. **运维友好**：同步服务有独立的 API（配置 CRUD、手动触发），方便运维人员管理

### Q2: 动态调度器是怎么设计的？

> **回答**：核心设计是 `map[configName] → cronInstance`，每个同步配置对应一个独立的 Cron 实例：
> - **启动时**：从 DB 加载所有 `enable_auto_sync=true` 的配置，为每个配置创建 Cron
> - **运行时**：ConfigService 的 CRUD 操作自动联动调度器（创建配置 → 启动 Cron，删除配置 → 停止 Cron）
> - **热更新**：Cron 执行时从 DB 重新读取最新配置，不需要重启服务
> - **互不影响**：每个配置独立的 Cron 实例，一个配置的调度不会影响其他配置

### Q3: 为什么 SyncService 是动态创建的？

> **回答**：因为支持多配置，每个配置对应不同的北向 API URL、Token、Source。SyncTaskManager 在执行任务时，根据任务请求中的配置参数动态创建 SyncService 实例（`createSyncServiceFromRequest`），而不是共享一个全局实例。这样不同配置的同步互不干扰。

### Q4: 数据库表设计有什么讲究？

> **回答**：采用**配置与执行分离**的三表设计：
> - `north_configure`：存储北向接口的**连接信息**（URL/Token/Source），是基础配置
> - `north_sync_config`：存储**同步策略**（Cron 表达式、是否自动同步），关联 north_configure
> - `sync_task`：存储**执行记录**（状态、统计、耗时），关联 config_name
>
> 好处：修改北向接口信息不影响同步策略，修改同步策略不影响历史记录。

### Q5: 如何保证同步任务不重复执行？

> **回答**：SyncTaskManager 使用 channel 作为任务队列，内部维护 `currentTask` 字段标记当前正在执行的任务。`SubmitTask` 时检查 `HasRunningTask()`，如果有任务正在执行则拒绝新任务。同时 channel 容量为 10，超过则返回「task queue is full」错误。

### Q6: 与 process-service 中同步模块的区别？

> **回答**：

| 特性 | process-service 同步 | north-sync |
|------|---------------------|------------|
| 配置方式 | 环境变量（单一配置） | 数据库管理（多配置） |
| 调度方式 | 单一 Cron 实例 | 动态多 Cron（按配置） |
| 管理 API | gRPC（TriggerSync） | REST API（完整 CRUD） |
| 北向配置 | 固定（环境变量） | 动态管理（north_configure 表） |
| 运行时修改 | 需要重启 | 热更新（DB 读取 + 调度联动） |
| 适用场景 | 单机房、单数据源 | 多机房、多数据源 |

### Q7: 这个项目的分层架构是怎样的？

> **回答**：严格的四层架构：
> 1. **API 层**（`pkg/api/`）：HTTP Handler，请求参数校验和响应封装
> 2. **Service 层**（`pkg/service/`）：业务逻辑，含调度器联动、配置校验
> 3. **DAO 层**（`pkg/dao/`）：数据库操作，封装 GORM 查询
> 4. **Model 层**（`pkg/models/`）：数据模型定义
>
> Service 层通过接口（`SchedulerReloader`）依赖调度器，避免 service 和 scheduler 包的循环依赖。
