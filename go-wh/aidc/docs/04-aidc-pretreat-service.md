# aidc-pretreat-service 项目分析

## 一、核心功能与整体架构

### 1.1 项目定位

`aidc-pretreat-service` 是 AIDC 平台的**数据预处理服务**，作为数据采集链路的入口，负责接收外部设备推送的原始数据，根据数据大小选择不同的处理策略（直写/压缩/MinIO 存储），最终将处理后的数据发送到 Kafka 供下游消费。

### 1.2 核心功能

- **多模式数据接收**：支持流式分片、MinIO 上传、一次性上传、压缩上传 4 种模式
- **智能数据分流**：根据数据大小（阈值 10MB）自动选择处理策略
- **Zstd 压缩**：对大数据进行 Zstd 高性能压缩，压缩比可达 3-10 倍
- **异步任务队列**：企业级 Worker Pool + 任务队列架构，支持高并发
- **Kafka 生产者**：将预处理后的数据写入 Kafka topic
- **MinIO 对象存储**：大文件存储到 MinIO，Kafka 只传递引用
- **Prometheus 监控**：内置请求量、延迟、压缩率等指标采集
- **流式分片处理**：支持 APISIX 网关转发的流式分片请求

### 1.3 架构流程图

```text
  ┌───────────────────────────┐
  │ 外部设备 / APISIX 网关     │
  └─────────┬─────────────────┘
            │ HTTP POST
            ▼
  ┌──────────────────────────────────────────────┐
  │          aidc-pretreat-service                │
  │                                              │
  │  ┌────────────────────────────────────────┐  │
  │  │          Gin HTTP Router               │  │
  │  │                                        │  │
  │  │  POST /api/v1/data        → 流式分片  │  │
  │  │  POST /api/v1/dataMinio   → MinIO异步 │  │
  │  │  POST /api/v1/dataOnce    → 一次性异步│  │
  │  │  POST /api/v1/dataCompress→ 压缩异步  │  │
  │  └────────────┬───────────────────────────┘  │
  │               │                              │
  │  ┌────────────▼───────────────────────────┐  │
  │  │       异步任务队列 (TaskQueue)          │  │
  │  │  • 队列容量: 10000                     │  │
  │  │  • Worker 数量: 100                    │  │
  │  │  • 最大待处理: 50000                   │  │
  │  │  • 满时拒绝策略                        │  │
  │  └────────────┬───────────────────────────┘  │
  │               │                              │
  │  ┌────────────▼───────────────────────────┐  │
  │  │       处理器适配器 (ProcessorAdapter)    │  │
  │  │                                        │  │
  │  │  ┌──────────┐  ┌──────────────┐        │  │
  │  │  │小数据     │  │大数据(≥10MB) │        │  │
  │  │  │直接发送   │  │              │        │  │
  │  │  │到 Kafka   │  │ ┌──────────┐│        │  │
  │  │  │          │  │ │Zstd 压缩 ││        │  │
  │  │  │          │  │ └────┬─────┘│        │  │
  │  │  │          │  │      ▼      │        │  │
  │  │  │          │  │ MinIO 存储  │        │  │
  │  │  │          │  │ 或压缩直传  │        │  │
  │  │  └────┬─────┘  └──────┬─────┘        │  │
  │  │       │               │               │  │
  │  │       ▼               ▼               │  │
  │  │  ┌───────────────────────────┐        │  │
  │  │  │    Kafka Producer         │        │  │
  │  │  │    (写入设备数据 topic)   │        │  │
  │  │  └───────────────────────────┘        │  │
  │  └────────────────────────────────────────┘  │
  └──────────────────────────────────────────────┘
            │                    │
            ▼                    ▼
     ┌──────────┐         ┌──────────┐
     │  Kafka   │         │  MinIO   │
     └──────────┘         └──────────┘
```

---

## 二、核心技术栈

| 分类 | 技术 | 用途 |
|------|------|------|
| 语言 | Go 1.22+ | 主语言 |
| HTTP 框架 | Gin | HTTP API 服务 |
| 消息队列 | Kafka (Sarama) | 数据写入下游 |
| 对象存储 | MinIO (minio-go/v7) | 大文件存储 |
| 压缩 | Zstd (klauspost/compress) | 高性能数据压缩 |
| 配置 | Viper | 环境变量配置 |
| 监控 | Prometheus | 请求量/延迟/压缩率指标 |
| 日志 | klog | 结构化日志 |
| ID 生成 | google/uuid | 请求 TraceID |
| 部署 | Docker + Helm | K8s 部署 |

---

## 三、核心代码文件/目录及作用

```text
aidc-pretreat-service/
├── cmd/aidc-pretreat-service/
│   └── main.go                            # 服务入口（233 行，完整启动流程）
├── internal/                               # 内部包（不可被外部引用）
│   ├── async/
│   │   ├── task_queue.go                  # 异步任务队列
│   │   └── processor_adapter.go           # 处理器适配器
│   ├── compress/
│   │   └── zstd.go                        # Zstd 压缩器
│   ├── config/
│   │   └── config.go                      # 配置结构（226 行，7 大类配置）
│   ├── handler/
│   │   ├── handler.go                     # 基础 Handler（Prometheus 指标）
│   │   ├── handler_stream.go             # 流式分片 Handler
│   │   ├── handler_minio.go              # MinIO Handler
│   │   ├── handler_onetime.go            # OneTime Handler
│   │   ├── handler_compress.go           # Compress Handler
│   │   └── handler_async.go              # 异步 Handler
│   ├── kafka/
│   │   └── producer.go                    # Kafka 生产者
│   ├── processor/
│   │   ├── processor.go                   # 核心处理器（大小分流逻辑）
│   │   ├── stream_chunk_processor.go     # 流式分片处理器
│   │   ├── minio_processor.go            # MinIO 处理器
│   │   ├── onetime_processor.go          # 一次性处理器
│   │   └── compress_processor.go         # 压缩处理器
│   └── storage/
│       └── minio_client.go                # MinIO 客户端封装
├── pkg/
│   └── models/
│       ├── message.go                     # Kafka 消息模型
│       ├── process_result.go             # 处理结果模型
│       └── response.go                    # HTTP 响应模型
├── Makefile                               # 构建脚本
└── Dockerfile                             # 容器镜像
```

**注意**：该项目使用 `internal/` 目录（Go 访问控制），而非 `pkg/`，强制内部实现不可被外部引用。

---

## 四、核心业务逻辑分析

### 4.1 启动流程

```text
main()
  ├── godotenv.Load()：加载 .env 文件
  ├── config.LoadDefaultConfig()：Viper 加载环境变量
  │     └── 验证配置（Server/MinIO/Kafka/Processor/Compression/Async）
  ├── compress.NewZstdCompressor()：初始化 Zstd 压缩器
  ├── storage.NewMinIOClient()：初始化 MinIO 客户端
  ├── kafka.NewProducer()：初始化 Kafka 生产者
  │
  ├── 初始化 4 种处理器：
  │   ├── StreamChunkProcessor（流式分片）
  │   ├── MinIOProcessor（MinIO 上传）
  │   ├── OneTimeProcessor（一次性上传）
  │   └── CompressProcessor（压缩上传）
  │
  ├── 初始化异步架构：
  │   ├── TaskQueue（队列容量 10000）
  │   ├── ProcessorAdapter（适配 3 种处理器）
  │   └── taskQueue.Start()（启动 100 个 Worker）
  │
  ├── 注册 Gin 路由（7 个端点 + 3 个健康检查）
  ├── 启动 HTTP Server
  └── 信号处理 → 优雅关闭（30s 超时）
```

### 4.2 数据处理分流逻辑

```text
HTTP POST 请求到达
  │
  ├── /api/v1/data → StreamHandler.HandleAPISIXRequest()
  │     ├── 解析 APISIX 请求格式
  │     ├── 生成 TraceID
  │     └── 流式分片处理 → 写入 Kafka
  │
  ├── /api/v1/dataMinio → AsyncHandler → MinIOProcessor
  │     ├── 小数据（< 10MB）→ 内联数据直发 Kafka
  │     └── 大数据（≥ 10MB）→ Zstd 压缩 → MinIO 上传 → Kafka 发引用
  │
  ├── /api/v1/dataOnce → AsyncHandler → OneTimeProcessor
  │     └── 压缩数据 → MinIO 上传 → Kafka 发引用
  │
  └── /api/v1/dataCompress → AsyncHandler → CompressProcessor
        ├── 小数据 → 直发 Kafka
        └── 大数据 → Zstd 压缩 → 压缩数据直接放入 Kafka 消息
```

### 4.3 异步任务队列架构

```text
AsyncHandler.HandleMinIORequestAsync()
  │
  ├── 1. 读取请求体
  ├── 2. 生成 TraceID
  ├── 3. 构建 Task（含数据、处理类型）
  ├── 4. taskQueue.Submit(task)
  │     ├── 队列未满 → 入队，返回 202 Accepted
  │     └── 队列已满 → 拒绝，返回 503 ServiceUnavailable
  └── 5. 返回 TraceID 给客户端

后台 Worker Pool（100 个 goroutine）：
  Worker → taskQueue.channel → Task
    ├── 根据 Task.Type 选择处理器
    ├── 执行处理（压缩 + 上传 + 发 Kafka）
    └── 记录结果日志
```

### 4.4 大数据压缩处理流程

```text
Processor.processLargeData(data)
  │
  ├── 1. Zstd 压缩（记录压缩时间和吞吐量）
  │     └── 输出：压缩数据 + 压缩比 + 压缩大小
  │
  ├── 2. 构建 StorageInfo
  │     ├── type: "compressed"
  │     ├── original_size / compressed_size
  │     └── compress: 压缩数据字节数组
  │
  ├── 3. 提取元数据（TryExtractMeta）
  │     └── 从原始 JSON 中提取 source、timestamp 等
  │
  ├── 4. 构建引用消息（ReferenceMessage）
  │     └── 包含 TraceID、Source、StorageInfo、Meta
  │
  └── 5. 发送到 Kafka
```

---

## 五、面试常见问题及回答思路

### Q1: 为什么需要预处理服务？直接写 Kafka 不行吗？

> **回答**：
> 1. **数据大小差异大**：设备数据从几 KB 到几百 MB 不等。Kafka 对单条消息有大小限制（默认 1MB，配置后最大 10MB），大数据需要特殊处理
> 2. **格式统一**：外部设备推送的数据格式不统一（APISIX 格式、原始 JSON 等），预处理层统一格式后再写入 Kafka
> 3. **解耦**：将数据接收、压缩、存储的逻辑从下游消费服务中剥离，各服务职责单一
> 4. **流量控制**：通过异步队列 + Worker Pool 控制对下游 Kafka/MinIO 的写入压力

### Q2: 为什么选择 Zstd 压缩算法？

> **回答**：Zstd（Facebook 开源）在压缩比和速度之间取得了最佳平衡：
> - 压缩比：通常 3-10 倍（对 JSON 数据效果尤其好）
> - 压缩速度：比 gzip 快 3-5 倍
> - 解压速度：比 gzip 快 5-10 倍
> - 支持流式压缩和可调压缩级别（1-22）
> 对于 IoT 场景的 JSON 数据，Zstd 是最优选择。

### Q3: 异步任务队列的设计思路？

> **回答**：采用「生产者-消费者」模式：
> - **TaskQueue**：带缓冲的 channel（容量 10000），作为任务缓冲区
> - **Worker Pool**：100 个 goroutine 并发消费任务
> - **背压控制**：队列满时拒绝新请求（返回 503），防止 OOM
> - **优雅关闭**：先停止接收新任务，等待所有 Worker 处理完当前任务
>
> 好处：HTTP 请求快速返回（202），实际处理在后台异步执行，大幅提升吞吐量。

### Q4: 为什么同时保留同步和异步路由？

> **回答**：
> - **异步路由**（推荐）：`/api/v1/dataMinio`、`/api/v1/dataOnce`、`/api/v1/dataCompress` — 高并发场景，请求快速返回
> - **同步路由**（兼容）：`/api/v1/dataMinio/sync`、`/api/v1/dataOnce/sync`、`/api/v1/dataCompress/sync` — 兼容旧版客户端，等待处理完成后返回
>
> 渐进式升级：先部署异步版本，旧客户端继续用同步路由，新客户端切换到异步路由。

### Q5: 这个服务如何监控？

> **回答**：内置 Prometheus 指标：
> - `pretreat_requests_total{status, topic}`：请求总量（按状态和 topic 分）
> - `pretreat_request_duration_seconds{topic}`：请求处理延迟分布
> - `pretreat_request_size_bytes{type}`：请求大小分布（small/large）
> - `pretreat_compression_ratio`：压缩率分布
> - 队列统计 API：`GET /api/v1/queue/stats`（Worker 利用率、队列深度等）

### Q6: 为什么使用 internal 而不是 pkg 目录？

> **回答**：Go 的 `internal` 包有**编译器级别的访问控制**——只有 `internal` 目录的父目录及其子目录可以导入 `internal` 下的包。这意味着预处理服务的内部实现（压缩器、处理器、Kafka 生产者等）不会被其他微服务意外引用，强制保持服务边界清晰。而 `pkg/models` 放在 pkg 下是因为消息模型可能需要被其他服务共享。
