# aidc-process-service 项目分析

## 一、核心功能与整体架构

### 1.1 项目定位

`aidc-process-service` 是 AIDC 平台的**数据处理核心服务**，负责从 Kafka 消费实时设备数据，经过解析处理后双写到 Redis（实时查询）和 VictoriaMetrics（历史存储），同时负责从北向 API 同步设备配置信息。

### 1.2 核心功能

- **Kafka 多 Topic 消费**：同时消费 jinhua、raw、bms-* 等多个 topic 的设备数据流
- **消息解析**：支持多种数据格式（jinhua、raw、bms 等），通过解析器工厂匹配
- **双写存储**：Redis HSET 存最新实时值 + VictoriaMetrics 存时序历史数据
- **限流控制**：为不同 topic 配置独立的 Token Bucket 限流策略
- **北向数据同步**：定时/手动从北向 API 同步位置、设备、点位配置
- **任务管理**：异步同步任务管理，支持提交、状态查询、历史记录
- **gRPC API**：健康检查、同步触发、同步状态查询

### 1.3 架构流程图

```text
  ┌───────────────┐          ┌─────────────────┐
  │ 北向 API      │          │ Kafka Cluster   │
  │ (设备配置)    │          │ (设备实时数据)   │
  └──────┬────────┘          └────────┬────────┘
         │ HTTP/JSON                  │ 多 Topic
         ▼                            ▼
  ┌─────────────────────────────────────────────┐
  │          aidc-process-service                │
  │                                             │
  │  ┌─────────────────┐  ┌──────────────────┐ │
  │  │ SyncService     │  │ Kafka Consumer   │ │
  │  │ • 版本控制      │  │ • 多 Topic 消费  │ │
  │  │ • 分页拉取      │  │ • Token Bucket   │ │
  │  │ • 事务写入 DB   │  │   限流           │ │
  │  └────────┬────────┘  └───────┬──────────┘ │
  │           │                   │             │
  │           ▼                   ▼             │
  │  ┌─────────────┐   ┌──────────────────┐    │
  │  │ Scheduler   │   │ MessageProcessor │    │
  │  │ (Cron 定时) │   │ • 解析器工厂     │    │
  │  └─────────────┘   │ • Worker Pool    │    │
  │                    │ • 并发双写       │    │
  │                    └──────┬───────────┘    │
  │                           │                │
  │              ┌────────────┼────────────┐   │
  │              ▼            ▼            │   │
  │      ┌──────────┐ ┌──────────────┐    │   │
  │      │ Redis    │ │VictoriaMetrics│    │   │
  │      │ Writer   │ │ Writer       │    │   │
  │      └──────────┘ └──────────────┘    │   │
  └───────────────────────────────────────────┘
         │                    │
         ▼                    ▼
  ┌──────────┐      ┌──────────────────┐
  │  Redis   │      │ VictoriaMetrics  │
  │ (实时值) │      │ (历史时序数据)   │
  └──────────┘      └──────────────────┘
```

---

## 二、核心技术栈

| 分类 | 技术 | 用途 |
|------|------|------|
| 语言 | Go 1.22+ | 主语言 |
| RPC 框架 | gRPC + Protobuf | 对外 API（健康检查、同步触发） |
| 数据库 | PostgreSQL + GORM | 设备配置数据、同步任务记录 |
| 消息队列 | Kafka (segmentio/kafka-go) | 消费设备实时数据 |
| 缓存 | Redis (go-redis) | 存储设备最新实时值（Hash 结构） |
| 时序数据库 | VictoriaMetrics | 存储历史时序数据（HTTP import API） |
| 对象存储 | MinIO | 大文件存储 |
| 定时调度 | robfig/cron/v3 | 定时同步北向数据 |
| 限流 | Token Bucket | Kafka 消费限流 |
| 日志 | klog | 结构化日志 |

---

## 三、核心代码文件/目录及作用

```text
aidc-process-service/
├── cmd/aidc-process-service/
│   └── main.go                    # 服务入口：初始化全部组件 + 启动
├── pkg/
│   ├── api/                       # gRPC API 实现
│   │   ├── health.go             # 健康检查
│   │   └── sync.go               # 同步服务（触发/状态/历史）
│   ├── base/
│   │   └── base.go               # 基础组件初始化（DB/Redis/SyncConfig）
│   ├── config/
│   │   ├── config.go             # 主配置结构（70+ 环境变量）
│   │   ├── kafka_config.go       # Kafka 消费者配置 + 限流配置
│   │   ├── redis_config.go       # Redis 配置
│   │   └── victoria_config.go    # VictoriaMetrics 配置
│   ├── kafka/
│   │   ├── consumer.go           # Kafka 多 Topic 消费者（核心）
│   │   └── rate_limiter.go       # Token Bucket 限流器
│   ├── middleware/
│   │   └── auth_interceptor.go   # gRPC 认证拦截器
│   ├── models/
│   │   ├── device.go             # 设备模型
│   │   ├── kafka/                # Kafka 消息模型
│   │   └── northnodes/           # 北向数据模型（Location/Device/Point）
│   ├── parser/
│   │   └── parser_factory.go     # 解析器工厂（按 topic 选择解析器）
│   ├── processor/
│   │   ├── message_processor.go  # 消息处理器（核心：解析 + 双写）
│   │   └── retry.go              # 重试机制
│   ├── scheduler/
│   │   └── scheduler.go          # Cron 定时同步调度器
│   ├── service/
│   │   ├── sync_service.go       # 北向同步核心逻辑（400+ 行）
│   │   └── sync_task_manager.go  # 异步任务管理器
│   ├── storage/
│   │   ├── redis_client.go       # Redis 客户端封装
│   │   ├── redis_writer.go       # Redis Hash 批量写入器
│   │   ├── victoria_client.go    # VictoriaMetrics HTTP 客户端
│   │   ├── victoria_writer.go    # VMS 批量写入器
│   │   └── dedup_cache.go        # 去重缓存（时间窗口去重）
│   └── metrics/                   # 指标计算模块
│       └── scheduler.go          # 指标计算调度器
├── migrations/                    # 数据库迁移脚本
├── Makefile                       # 构建脚本
└── Dockerfile                     # 容器镜像
```

---

## 四、核心业务逻辑分析

### 4.1 启动流程

```text
main()
  ├── 加载配置（envconfig，前缀 AIDC_PROCESS_SERVICE_）
  ├── base.Init()：初始化 PostgreSQL
  ├── base.InitRedis()：初始化 Redis
  ├── base.InitSyncConfig()：加载北向同步配置
  ├── 创建 SyncService + SyncTaskManager
  ├── taskManager.Start()：启动任务队列处理
  ├── initKafkaConsumer()：初始化 Kafka 多 Topic 消费者
  │     ├── 解析 Kafka 配置（Brokers/Topics/GroupID）
  │     ├── 创建 Consumer + RateLimiter
  │     ├── 创建 MessageProcessor（解析器 + 双写器）
  │     ├── 为每个 Topic 注册 Handler
  │     └── consumer.Start()：启动消费
  ├── initMetricsCalculation()：初始化指标计算（可选）
  ├── 配置定时同步调度器（如果 EnableAutoSync）
  ├── 注册 gRPC 服务（Health + Sync）
  ├── 设置信号处理（优雅关闭）
  └── run.DefaultRun()：启动 gRPC 服务器
```

### 4.2 Kafka 消费处理流程（核心）

```text
Kafka Message (topic: jinhua/raw/bms-*)
  │
  ▼
Consumer.consumeTopic()
  ├── FetchMessage()：拉取消息
  ├── RateLimiter.Wait()：Token Bucket 限流等待
  ├── handleMessage()：分发到对应 Handler
  │     │
  │     ▼
  │   MessageProcessor.Handle()
  │     ├── ParserFactory.Parse()：根据 topic 选择解析器
  │     │     ├── jinhua → JinhuaParser
  │     │     ├── raw → RawParser
  │     │     └── bms-* → BmsParser
  │     ├── 提取设备状态 + 数据点值
  │     ├── RedisWriter.Write()
  │     │     └── HSET {source} {point_guid} {json_data}
  │     ├── VictoriaWriter.Write()
  │     │     ├── device_status{device_guid, source} = value
  │     │     └── device_point_value{device_guid, point_guid, source, quality} = value
  │     └── DedupCache 检查去重（60s 窗口）
  │
  └── CommitMessages()：提交 offset（处理成功才提交）
```

### 4.3 北向数据同步流程

```text
SyncAllNorthNodes()
  ├── 1. 获取本地配置版本号
  ├── 2. 请求北向 API 第一页（获取远程版本号）
  ├── 3. 版本对比
  │     ├── 版本相同 & 非强制 → 跳过同步
  │     └── 版本不同或强制 → 继续同步
  ├── 4. 分页拉取（每页 500 条）
  │     └── 对每页数据：
  │         ├── 按 NodeType 分类：
  │         │   ├── NodeType=1 → Location（位置）
  │         │   ├── NodeType=2 → BmsDevice（设备）
  │         │   └── NodeType=3 → BmsPoint（点位）
  │         ├── 事务写入 PostgreSQL（Begin → Create → Commit）
  │         └── 休眠 1s（避免打满北向 API）
  ├── 5. 构建位置层级路径（LocationPath）
  ├── 6. 关联设备与位置（TrimDeviceLocations）
  └── 7. 更新配置版本号
```

### 4.4 限流机制（Token Bucket）

```text
TopicRateLimiter
  ├── 默认限流：10 msg/s
  ├── jinhua：5 msg/s, burst 10
  ├── raw：10 msg/s, burst 20
  └── bms-D04-*：20 msg/s, burst 40

工作原理：
  • 每个 topic 独立的 rate.Limiter
  • 消费消息前调用 Wait() 等待令牌
  • burst 允许突发流量，避免流量尖刺丢消息
```

---

## 五、面试常见问题及回答思路

### Q1: Kafka 消费为什么要用限流？

> **回答**：因为下游存储（Redis、VictoriaMetrics）有写入性能上限。如果 Kafka 消费速度过快，可能导致：
> 1. Redis 连接池被打满，写入超时
> 2. VictoriaMetrics HTTP 写入积压，内存 OOM
> 3. 不同 topic 的数据量差异很大（bms 比 jinhua 多很多），需要差异化限流
>
> 方案是为每个 topic 配置独立的 Token Bucket 限流器，保护下游系统。

### Q2: 为什么选择 Redis Hash 结构存储实时数据？

> **回答**：`HSET {source} {point_guid} {json_data}` 的好处：
> 1. **O(1) 查询**：`HGET source point_guid` 直接获取某个点位的最新值
> 2. **按来源聚合**：同一个 source 的所有点位在一个 Hash 下，方便批量查询
> 3. **内存效率**：Redis Hash 在元素少于 128 个时使用 ziplist 编码，内存占用更小
> 4. **兼容性**：与 boss-jh-gateway 数据格式完全兼容，平滑迁移

### Q3: VictoriaMetrics 写入如何优化性能？

> **回答**：
> 1. **批量写入**：使用缓冲区积累数据，达到 10000 条或 10 秒自动 flush
> 2. **去重缓存（DedupCache）**：基于 60 秒时间窗口，相同设备+点位只写一次，避免 Kafka 重复消费导致的重复写入
> 3. **HTTP import API**：使用 VictoriaMetrics 的批量导入接口，比逐条写入高效 100 倍

### Q4: 北向同步为什么用版本号增量同步？

> **回答**：北向 API 返回的设备配置数据量很大（可能数万条），每次全量同步开销大。通过版本号（version）对比，只有远程版本号变化时才执行同步，大幅减少不必要的同步操作。同时支持 `force=true` 参数强制全量同步。

### Q5: 如何保证消息不丢失？

> **回答**：采用**手动提交 offset** 策略（at-least-once）：
> 1. `FetchMessage()` 拉取消息
> 2. `Handler.Handle()` 处理消息（写 Redis + VMS）
> 3. 处理成功后才 `CommitMessages()` 提交 offset
> 4. 如果处理失败，不提交 offset，下次重新消费
> 5. 通过 DedupCache 去重，防止重复消费导致数据重复

### Q6: 这个服务挂了会怎样？

> **回答**：
> - **Kafka 消费**：使用 Consumer Group，服务重启后自动从上次提交的 offset 继续消费
> - **Redis/VMS 写入失败**：消息不提交 offset，下次重新消费重试
> - **北向同步中断**：支持断点续传（版本号机制），下次启动继续同步
> - **优雅关闭**：监听 SIGTERM 信号，关闭 Kafka Reader、停止调度器、释放数据库连接

### Q7: 消息处理和写入如何保证顺序？

> **回答**：
> - **同一分区内有序**：Kafka 同一 partition 内消息有序
> - **跨分区无序**：不同 partition 间消息无序，但对于 IoT 数据场景，我们只关心**最新值**（HSET 覆盖写），不需要全局有序
> - **VictoriaMetrics**：时序数据按时间戳存储，自动按时间排序
