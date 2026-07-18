# aidc-platform-service 项目分析

## 一、核心功能与整体架构

### 1.1 项目定位

`aidc-platform-service` 是 AIDC 平台的**管理服务层（BFF）**，作为前端/网关与后端数据之间的桥梁，提供用户管理、权限控制、设备数据查询、计算点管理等核心业务功能。

### 1.2 核心功能

- **用户与认证管理**：集成 Casdoor SSO 实现统一登录，JWT Token 鉴权，AD 域用户同步
- **权限管理（RBAC）**：菜单管理、角色管理、用户-角色-菜单关联
- **位置与设备管理**：DCOM 位置层级查询、BMS 设备/点位查询
- **数据字典**：系统字典数据管理（支持 Redis 缓存）
- **数据同步**：从北向接口同步设备配置数据
- **计算点管理**：DCIM 计算点 CRUD、公式配置、实时值查询
- **指标点服务（MetricPoints）**：计算点数据初始化到 Redis、VictoriaMetrics 历史数据查询
- **北向同步统计**：北向数据同步状态统计
- **测试工具**：测试数据生成器、压力测试服务

### 1.3 架构流程图

```text
                        ┌─────────────────────┐
                        │   前端 / API 网关    │
                        └──────────┬──────────┘
                                   │ gRPC (通过 grpc-gateway)
                                   ▼
                    ┌──────────────────────────────┐
                    │   aidc-platform-service       │
                    │                              │
                    │  ┌─────────────────────────┐ │
                    │  │ AuthInterceptor (Casdoor)│ │  ← 认证拦截器
                    │  └─────────────────────────┘ │
                    │              │                │
                    │  ┌───────────┼──────────────┐│
                    │  │           ▼              ││
                    │  │   gRPC Service 层        ││
                    │  │  ┌──────┐ ┌──────────┐  ││
                    │  │  │ Auth │ │UserService│  ││
                    │  │  └──────┘ └──────────┘  ││
                    │  │  ┌──────┐ ┌──────────┐  ││
                    │  │  │ Dict │ │DcomLoc   │  ││
                    │  │  └──────┘ └──────────┘  ││
                    │  │  ┌──────────────────┐   ││
                    │  │  │ MetricPoints     │   ││
                    │  │  │ DcimMetrics      │   ││
                    │  │  │ NorthSyncStats   │   ││
                    │  │  └──────────────────┘   ││
                    │  └─────────────────────────┘│
                    └──────┬────────┬─────────┬───┘
                           │        │         │
                    ┌──────▼──┐ ┌───▼───┐ ┌───▼──────────┐
                    │PostgreSQL│ │ Redis │ │VictoriaMetrics│
                    └─────────┘ └───────┘ └──────────────┘
```

---

## 二、核心技术栈

| 分类 | 技术 | 用途 |
|------|------|------|
| 语言 | Go 1.22+ | 主语言 |
| RPC 框架 | gRPC + Protobuf | 对外 API |
| 数据库 | PostgreSQL + GORM | 业务数据持久化 |
| 缓存 | Redis (go-redis/v9) | 字典缓存、计算点数据、分布式锁 |
| 时序数据库 | VictoriaMetrics | 历史数据查询（通过 VictoriaQueryBuilder） |
| 认证 | Casdoor SSO + JWT | 统一身份认证 |
| 日志 | klog + xlog | 结构化日志 |
| ID 生成 | Snowflake（灵活版） | 分布式 ID（15位机器ID + 7位序列号） |
| 配置 | envconfig | 环境变量配置 |
| 部署 | Docker + Helm | K8s 部署 |

---

## 三、核心代码文件/目录及作用

```text
aidc-platform-service/
├── cmd/aidc-platform-service/
│   └── main.go                        # 服务入口：初始化组件 + 注册 gRPC 服务
├── pkg/
│   ├── api/                           # gRPC API 实现层（Handler）
│   │   ├── auth_login.go              # 登录认证服务
│   │   ├── auth_menu.go               # 菜单管理服务
│   │   ├── auth_role.go               # 角色管理服务
│   │   ├── user.go                    # 用户管理服务
│   │   ├── ad_user.go                 # AD 域用户服务
│   │   ├── dict.go                    # 数据字典服务
│   │   ├── dcom_location.go           # 位置管理服务
│   │   ├── data_sync.go              # 数据同步服务
│   │   ├── dcim_metrics.go           # DCIM 指标服务
│   │   ├── metric_points.go          # 计算点服务
│   │   ├── north_sync_statistics.go  # 北向同步统计
│   │   ├── test_data_generator.go    # 测试数据生成
│   │   └── stress_test.go            # 压力测试
│   ├── base/
│   │   └── base.go                    # 基础组件初始化（DB/Redis/Casdoor/VMS）
│   ├── config/                        # 配置定义（暂用 polar-common-go）
│   ├── dax/                           # 数据访问层（DAO）
│   │   ├── dcim_metric_points.go     # 计算点 DAO
│   │   ├── bms_device.go             # BMS 设备 DAO
│   │   ├── bms_point.go              # BMS 点位 DAO
│   │   ├── dcom_location.go          # 位置 DAO
│   │   ├── dict_data.go              # 字典 DAO
│   │   └── user_rel_location.go      # 用户-位置关系 DAO
│   ├── middleware/
│   │   └── auth_interceptor.go        # Casdoor 认证拦截器
│   ├── service/
│   │   └── metric_points.go           # 计算点业务逻辑（核心 2000+ 行）
│   └── storage/
│       └── victoria_query_builder.go  # VictoriaMetrics 查询构建器
├── staging/                           # Protobuf 生成代码 / 公共库
├── charts/                            # Helm Chart
├── Makefile                           # 构建脚本
└── Dockerfile                         # 容器镜像
```

---

## 四、核心业务逻辑分析

### 4.1 启动流程

```text
main()
  ├── 加载配置（envconfig，前缀 AIDC_PLATFORM_SERVICE_）
  ├── base.Init()：初始化 PostgreSQL（GORM + 连接池）
  ├── base.InitRedis()：初始化 Redis（支持单机/集群）
  ├── base.InitCasdoor()：初始化 Casdoor SSO 客户端
  ├── base.InitVictoriaMetrics()：初始化 VMS 查询 URL
  ├── 创建 AuthInterceptor（Casdoor + 白名单）
  ├── 创建 VictoriaQueryBuilder
  ├── 注册 12 个 gRPC 服务
  ├── 异步：InitMetricPointsToRedis（计算点数据初始化）
  └── run.DefaultRun()：启动 gRPC 服务器
```

### 4.2 认证流程

```text
gRPC 请求 → AuthInterceptor
  ├── 检查白名单（gRPC method / HTTP path / 后缀匹配）
  │     ├── 命中 → 跳过认证，直接处理
  │     └── 未命中 → 继续认证
  ├── 从 metadata 提取 Authorization Bearer Token
  ├── Casdoor ParseJwtToken 验证
  │     ├── 验证成功 → 提取用户信息，注入 context
  │     └── 验证失败 → 返回 Unauthenticated
  └── 调用 handler 处理业务逻辑
```

### 4.3 计算点初始化到 Redis（核心流程）

```text
InitMetricPointsToRedis()
  ├── 1. 获取分布式锁（Redis SETNX，TTL=5min）
  │     └── 未获取 → 另一个 Pod 正在初始化，跳过
  ├── 2. 启动锁续期 goroutine（每 60s 续期一次）
  ├── 3. 分页查询计算点（每页 1000 条）
  │     └── 对每个计算点：
  │         ├── 写入 formula:{source} Hash（计算点配置）
  │         └── 解析 cal_parameter → 写入 point:{source}:{type} Hash（原始点位）
  ├── 4. Redis Pipeline 批量写入
  └── 5. 释放锁
```

### 4.4 MetricPoints 查询流程

```text
ListMetricPoints(req)
  ├── 1. 从 context 提取用户信息（userId、isAdmin）
  ├── 2. 非管理员：查询用户关联的位置 ID 列表
  ├── 3. 根据 locationId + locationType 展开子位置 ID
  ├── 4. 调用 DAO 分页查询计算点列表
  └── 5. 转换为 Protobuf 响应返回
```

---

## 五、面试常见问题及回答思路

### Q1: 这个服务的认证是怎么做的？

> **回答**：采用 Casdoor SSO + JWT 方案。通过 gRPC UnaryInterceptor 实现统一认证拦截。支持三种白名单匹配方式：gRPC 方法全路径匹配、HTTP 路径匹配、路径后缀匹配。Token 通过 Casdoor SDK 的 `ParseJwtToken` 进行验证和解析，验证通过后将用户信息注入 context 供下游使用。

### Q2: 计算点初始化为什么需要分布式锁？

> **回答**：因为服务部署在 K8s 上，可能有多个 Pod 实例同时启动。如果不加锁，多个 Pod 会重复执行初始化，造成 Redis 重复写入和资源浪费。使用 Redis SETNX 实现分布式锁，同时启动 goroutine 做锁续期（每 60s），防止任务执行时间超过锁 TTL 导致锁失效。

### Q3: 为什么使用 Pipeline 批量写入 Redis？

> **回答**：计算点数量可能达到几万条，如果逐条写入 Redis，每次都是一次网络往返（RTT），性能很差。Pipeline 将多个命令打包成一次网络请求发送，大幅减少网络开销。实测批量写入性能比逐条写入提升 10-50 倍。

### Q4: Snowflake ID 生成器为什么用 15 位机器 ID？

> **回答**：标准 Snowflake 用 10 位机器 ID（最多 1024 台机器）。这里使用 15 位机器 ID + 7 位序列号的灵活配置，支持 32768 台机器（适合大规模 K8s 集群），但每毫秒只能生成 128 个 ID。machineId 优先从环境变量获取，否则从 Pod hostname 哈希生成。

### Q5: 这个服务有哪些设计上的亮点？

> **回答**：
> 1. **分层架构清晰**：api（Handler）→ service（业务逻辑）→ dax（数据访问）→ base（基础设施），职责分离
> 2. **分布式锁 + 锁续期**：计算点初始化使用 SETNX + goroutine 续期，防止多 Pod 重复执行
> 3. **灵活的认证白名单**：支持完整路径、HTTP 路径、后缀三种匹配方式
> 4. **Redis 支持单机/集群双模式**：通过环境变量 `REDIS_SERVER_TYPE` 切换
> 5. **VictoriaMetrics 查询构建器**：封装 PromQL 查询，简化历史数据查询

### Q6: GORM 连接池配置有什么讲究？

> **回答**：
> - `MaxOpenConns=20`：控制最大并发连接数，避免打满数据库
> - `MaxIdleConns=10`（最大连接数的一半）：保持一定数量的空闲连接，减少建连开销
> - `ConnMaxLifetime=1h`：连接定期回收，避免使用过期连接
> - `ConnMaxIdleTime=10min`：空闲连接及时释放，避免占用资源
> 这是生产环境的标准配置，平衡了性能和资源占用。
