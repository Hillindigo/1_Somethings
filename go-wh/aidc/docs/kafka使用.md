# Kafka 在 AIDC 项目中的使用总结

## 概述

Kafka 在整个 AIDC 系统中承担着**核心消息总线**的角色，贯穿数据从北向采集→预处理→流式消费→计算触发→计算执行→死信补偿的完整链路。下图描述了各服务之间的 Kafka 数据流向：

```
北向设备数据
     │
     ▼
[aidc-pretreat-service]
  Producer → Topic: aidc-raw-data（原始数据预处理后发出）
     │
     ▼
[aidc-process-service]
  Consumer ← Topics: jinhua / raw / bms-D04-rtdata-shuoyao（消费原始数据）
  处理后写入 Redis(R1) + VictoriaMetrics
     │（Redis 数据触发计算）
     ▼
[aidc-metric-service]
  CalcTriggerEvaluator → Producer → Topic: aidc-calc-task（发送计算任务）
  CalcProcessor        ← Consumer ← Topic: aidc-calc-task（消费计算任务）
  DLQHandler           ← Consumer ← Topic: aidc-calc-task-dlq（消费死信）
  DLQHandler           → Producer → Topic: aidc-calc-task（补偿重投）
     │
     ▼
[aidc-platform-service]
  enableKafka 开关（控制 Kafka 推送，当前仅标志位，无独立 Producer/Consumer 实现）

[aidc-north-sync]
  无直接 Kafka Producer/Consumer 代码
  Topic 字段作为北向配置元数据存储于数据库（供其他服务读取）
```

---

## 一、aidc-pretreat-service

### 职责
接收北向原始数据（HTTP 推送），进行**预处理**（大小分片、压缩、MinIO 存储），然后将统一格式的消息**生产**到 Kafka。

### Kafka 角色
**纯 Producer**（只生产，不消费）

### 库依赖
```
github.com/IBM/sarama
```

### Topic 配置

| 配置项 | 环境变量 | 说明 |
|--------|----------|------|
| Brokers | `{PREFIX}_KAFKA_BROKERS` | Broker 地址列表 |
| Topic | `{PREFIX}_KAFKA_TOPIC` | 目标 topic（必填） |
| ProducerID | `{PREFIX}_KAFKA_PRODUCER_ID` | 客户端标识（必填） |

### 生产者实现

服务提供两套生产者，均位于 `internal/kafka/` 下：

#### 1. 普通同步生产者（`producer.go`）
```go
// Producer 封装 sarama.SyncProducer
type Producer struct {
    producer sarama.SyncProducer
    topic    string
}
```

**关键配置：**
- `RequiredAcks`：可配置（0/1/-1），支持精确一次语义
- `EnableIdempotence`：支持幂等生产者，配合 `Net.MaxOpenRequests=1`
- `MaxRetries`：失败重试次数
- `FlushFrequency / FlushMessages / FlushBytes`：批量刷新控制
- `MaxMessageBytes`：单条消息最大字节数

**消息 Key**：使用 `TraceID` 作为分区 Key，保证同一请求的消息路由到同一分区。

**核心方法：**
```go
// Send 发送统一消息（UnifiedMessage）
func (p *Producer) Send(msg *models.UnifiedMessage) error

// SendSmallData（已废弃，使用 Send 替代）
// SendLargeDataRef（已废弃，使用 Send 替代）
```

**消息体结构（`models.UnifiedMessage`）：**
- 小数据（< 阈值）：直接携带完整数据体
- 大数据（≥ 阈值）：数据存 MinIO，消息体仅携带 MinIO 引用路径（大包拆分策略）

#### 2. 批量生产者（`batch_producer.go`）
```go
// BatchProducer 企业级批量 Kafka 生产者
type BatchProducer struct {
    syncProducer  sarama.SyncProducer
    topic         string
    buffer        []*sarama.ProducerMessage
    // 统计：sentCount, errorCount, batchCount, totalLatency
}
```

**特性：**
- 内置消息缓冲区，达到 `BatchSize`（默认 100）自动刷新
- 定时器兜底刷新，默认 `BatchTimeout = 100ms`
- 使用 **Snappy 压缩**（`CompressionSnappy`）减少网络开销
- 批量刷新使用 `sarama.SendMessages` 原子发送多条
- 提供 `StreamBatchSender`：支持并发 Worker + 背压控制的流式高吞吐发送器

**默认参数：**
```go
BatchSize:    100              // 满 100 条触发发送
BatchTimeout: 100ms            // 超时兜底
BufferSize:   10000            // 最大缓冲条数
Concurrency:  4                // 并发发送 Worker 数
Flush.Frequency: 50ms          // sarama 层刷新频率
Flush.Bytes:  1MB              // sarama 层字节阈值
```

### 在主流程中的使用

```go
// main.go
kafkaProducer, err := kafka.NewProducer(&cfg.Kafka)
streamProc := processor.NewOptimizedStreamChunkProcessor(kafkaProducer)
// HTTP 接收 → 预处理 → kafkaProducer.Send(unifiedMessage)
```

---

## 二、aidc-process-service

### 职责
消费多个原始数据 Topic，解析不同格式的设备数据，写入 Redis（R1，供计算触发读取）和 VictoriaMetrics（时序存储）。

### Kafka 角色
**纯 Consumer**（只消费，不生产）

### 库依赖
```
github.com/segmentio/kafka-go
```

### Topic 配置

| 配置项 | 环境变量 | 默认值 | 说明 |
|--------|----------|--------|------|
| Brokers | `KAFKA_BROKERS` | — | Broker 地址（逗号分隔） |
| Topics | `KAFKA_TOPICS` | `jinhua,raw,bms-D04-rtdata-shuoyao` | 订阅 topic 列表 |
| GroupID | `KAFKA_GROUP_ID` | `aidc-process-service` | 消费组 ID |
| MaxRetries | `KAFKA_MAX_RETRIES` | `3` | 最大重试次数 |
| BatchSize | `KAFKA_BATCH_SIZE` | `100` | 批次大小 |
| DefaultRateLimit | `KAFKA_DEFAULT_RATE_LIMIT` | `10` | 默认限速（消息/秒） |

**各 Topic 精细限速配置（硬编码默认值）：**
```go
"jinhua":                 {Rate: 5,  Burst: 10}  // 5条/秒，突发10
"bms-D04-rtdata-shuoyao": {Rate: 20, Burst: 40}  // 20条/秒，突发40
"raw":                    {Rate: 10, Burst: 20}   // 10条/秒，突发20
```

### 消费者实现（`pkg/kafka/consumer.go`）

```go
type Consumer struct {
    config      *config.KafkaConfig
    readers     map[string]*kafka.Reader  // 每个 topic 一个 Reader
    handlers    map[string]MessageHandler // topic → 处理器映射
    rateLimiter *TopicRateLimiter         // 基于 token bucket 的限速器
}
```

**关键特性：**
- **每个 Topic 独立 goroutine**：`go c.consumeTopic(ctx, topic, reader)`
- **手动提交 offset**：`reader.CommitMessages(ctx, msg)`，处理成功后才提交，保证 At-Least-Once
- **限速控制**：消费前通过 `rateLimiter.Wait(ctx, topic)` 阻塞等待令牌，防止下游（Redis/VMS）过载
- **StartOffset**：`kafka.LastOffset`，服务重启从最新 offset 开始（非历史回溯）
- **优雅停止**：`Stop()` 关闭所有 Reader，等待所有 goroutine 退出，超时 30s

**Reader 参数：**
```go
kafka.ReaderConfig{
    Brokers:        brokers,
    Topic:          topic,
    GroupID:        groupID,
    MinBytes:       1,
    MaxBytes:       10 * 1024 * 1024,  // 10MB
    MaxWait:        500ms,
    CommitInterval: 1s,
    StartOffset:    kafka.LastOffset,
}
```

### 消息格式（多 Topic 适配）

| Topic | 消息类型 | 解析类 | 数据来源 |
|-------|----------|--------|----------|
| `jinhua` | 金华格式（通过 APISIX 转发的 HTTP 日志体） | `JinhuaMessage` | 金华北向接口 |
| `raw` | 通用原始格式 | `RawMessage` | 通用北向 |
| `bms-D04-rtdata-shuoyao` | BMS 设备实时数据 | `BmsMessage` | BMS 设备 |

**消息处理器接口：**
```go
type MessageProcessor interface {
    Parse(data []byte) error
    ExtractDevices() ([]models.Device, error)
    GetTopicName() string
    Validate() error
}
```

### 消息处理流程（`pkg/processor/message_processor.go`）

```
消费到消息
    │
    ▼
ParserFactory.GetParser(topic, data)  // 按 topic 选择解析器
    │
    ▼
parser.ExtractDevices()               // 解析出 []Device
    │
    ├──→ RedisWriter.WriteDevices(topic, devices)   // 写 R1 Hash（HSET {topic} {dpid} {json}）
    │
    └──→ VictoriaWriter.WriteDevices(devices)       // 写时序数据库
```

**并发写入**：Redis 和 VictoriaMetrics 写入通过两个 goroutine 并发执行，使用 `sync.WaitGroup` 等待，任一失败则返回错误。

**Worker Pool**：`workerPool chan struct{}`，大小为 `ProcessorConfig.WorkerPoolSize`，防止无限并发。

---

## 三、aidc-metric-service

### 职责
**最复杂的 Kafka 使用方**：同时扮演 Producer 和 Consumer，实现计算任务调度、计算执行和死信补偿。

### Kafka 角色
**Producer + Consumer + DLQ 机制**

### 库依赖
```
github.com/segmentio/kafka-go
```

### Topic 配置

| 配置项 | 环境变量 | 默认值 | 说明 |
|--------|----------|--------|------|
| Brokers | `KAFKA_BROKERS` | — | Broker 地址 |
| CalcTopic | `KAFKA_CALC_TASK_TOPIC` | `aidc-calc-task` | 计算任务 topic |
| DLQTopic | `KAFKA_CALC_DLQ_TOPIC` | `aidc-calc-task-dlq` | 死信 topic |
| GroupID | `KAFKA_CALC_GROUP_ID` | `aidc-metric-service-calc` | 消费组 ID |

### 全局 Writer 初始化（`pkg/base/base.go`）

```go
kafkaWriter = &kafka.Writer{
    Addr:         kafka.TCP(brokerList...),
    Topic:        topic,              // aidc-calc-task
    Balancer:     &kafka.LeastBytes{}, // 最少字节数负载均衡
    MaxAttempts:  3,
    BatchSize:    100,
    BatchTimeout: 10ms,
    Async:        false,             // 同步写入，保证可靠性
}
```

### 3.1 计算任务生产者（CalcTriggerEvaluator）

**文件：** `pkg/calc/trigger.go`

**触发时机：** `process-service` 将原始数据写入 Redis（R1）后，`CalcTriggerEvaluator.Evaluate()` 检查是否有依赖该数据点的计算点，若有则构建并发送 `CalcTaskMessage`。

**消息 Key**：`calc_dpid`（按计算点 ID 分区，保证同一计算点的消息有序）

**两种消息路径：**

| 路径 | 触发条件 | 消息体 |
|------|----------|--------|
| **小依赖路径**（≤50个依赖） | `len(deps) <= LargeDepThreshold` | 消息体内嵌所有依赖值 `Values []DependencyValue` |
| **大依赖路径**（>50个依赖） | `len(deps) > LargeDepThreshold` | 依赖值冻结写入 Redis Snapshot（`calc:snap:{id}`），消息体只含 `SnapshotRef` |

**CoalescingBuffer（去重合并）：** 同一计算点在 100ms 窗口内的多次触发会被合并为一次，避免重复计算。

**CalcTaskMessage 结构：**
```go
type CalcTaskMessage struct {
    CalcDpid       string             // 计算点 ID（分区 Key）
    CalcName       string             // 计算点名称
    CalFormula     string             // 计算公式
    CalMethod      string             // 计算方式（formula/delta/rate）
    Source         string             // 数据源（Redis Hash Key）
    Values         []DependencyValue  // 依赖值（小路径内嵌）
    SnapshotRef    *SnapshotRef       // 大路径快照引用
    TriggerTs      int64              // 触发时间戳（秒）
    ConfigVersion  uint64             // 配置版本（用于幂等 Key）
    IdempotencyKey string             // 幂等 Key = hash(calcDpid+triggerTs+configVersion)
    HopCount       int                // 级联深度（防循环，最大 10 层）
    TriggerType    string             // 触发类型：raw/calc
    LocationSlug   string             // 位置标识
}
```

**级联触发（`TriggerCascade`）：** 计算完成后，自动触发依赖该计算点结果的下游计算点，最大递归深度 10 层（`MaxCascadeDepth`）。

### 3.2 KafkaProducer 封装（`pkg/storage/kafka_producer.go`）

```go
// 支持大包体自动递归对半拆分
const MaxMessageSize = 2560000  // 2.56MB

func (p *KafkaProducer) ProduceData(key string, data []byte) error {
    if len(data) <= MaxMessageSize {
        return p.sendMessage(key, data)
    }
    // 超限则递归对半拆分后分批发送
    return p.splitAndSend(key, data)
}
```

在 `DatapointManager` 中用于将计算结果（数据点快照）发送到 Kafka。

### 3.3 计算任务消费者（CalcProcessor）

**文件：** `pkg/calc/processor.go`

**消费 Topic：** `aidc-calc-task`

**完整处理 Pipeline：**

```
接收 CalcTaskMessage（kafka.Message）
    │
    ├─ 1. 反序列化（失败 → DLQ: deserialize_error）
    │
    ├─ 2. 幂等检查：Redis key "calc:idem:{IdempotencyKey}"（已处理则跳过）
    │
    ├─ 3. 级联深度检查：HopCount >= 10 → DLQ: max_cascade_depth_exceeded
    │
    ├─ 4. 分布式锁：Redis SetNX "calc:lock:{IdempotencyKey}" TTL=30s（并发互斥）
    │
    ├─ 5. 加载依赖值：
    │      ├─ 若有 SnapshotRef → 从 Redis 读取快照
    │      └─ 若轻量 K2 消息（无 Values）→ 双缓存读取：
    │             ├─ R2: MGET "r2:{source}:{dpid}:{ts}"（时间点精确值，TTL=900s）
    │             └─ R1: HMGET "{source}" {dpid}（最新值，超 1h 视为过期）
    │
    ├─ 6. delta/rate 处理：读取 "calc:prev:{calcDpid}" 注入 __prev_value__
    │
    ├─ 7. 公式求值：FormulaEngine.EvaluateAST（失败 → DLQ: formula_error）
    │
    ├─ 8. 写 R1：RedisWriter.WriteDevicesWithTimestampCAS（Lua CAS 防乱序）
    │
    ├─ 9. 写 R2：RedisWriter.WriteR2CalcPoint（SET EX 900，允许覆盖）
    │
    ├─ 10. 写 VMS：VictoriaWriter.WriteDevices（最佳努力，失败不阻塞）
    │
    ├─ 11. 写完整性指标：WriteCompleteness（有效依赖数/总依赖数）
    │
    ├─ 12. 标记幂等完成：SetJSON "calc:idem:{key}" TTL=25h
    │
    └─ 13. 触发级联：TriggerCascade → 下游计算点发送新 CalcTaskMessage
```

**失败时发送到 DLQ（`aidc-calc-task-dlq`）：**
```go
type dlqEnvelope struct {
    OriginalMessage json.RawMessage  // 原始消息
    ErrorType       string           // 错误类型
    ErrorMessage    string           // 错误详情
    FailedAt        int64            // 失败时间戳
    Component       string           // 组件名："calc_processor"
}
```

### 3.4 死信队列消费者（DLQHandler）

**文件：** `pkg/dlq/handler.go`

**消费 Topic：** `aidc-calc-task-dlq`

**按错误类型分类补偿：**

| 错误类型 | 补偿策略 |
|----------|----------|
| `cascade_trigger_failed` | 检查主写幂等 Key → 若已写入则重触发下游级联；否则重投回 `aidc-calc-task` |
| `config_missing` | 等待配置热重载后，更新 ConfigVersion 和 IdempotencyKey，重投 `aidc-calc-task` |
| `snapshot_expired` | 触发 P2 告警（快照 TTL 过期，无法补偿）|
| `formula_error` | 触发 P1 告警（公式错误，需人工介入）|
| `max_cascade_depth_exceeded` | 触发 P1 告警（级联深度超限）|
| `deserialize_error` | 触发 P2 告警 |

**DLQHandler 结构：**
```go
type DLQHandler struct {
    reader      *kafka.Reader   // 消费 aidc-calc-task-dlq
    producer    *kafka.Writer   // 重投回 aidc-calc-task
    redisClient calc.RedisStore
    indexFunc   func() *calc.CalcDependencyIndex
    trigger     CascadeRetrigger
    alerter     Alerter
}
```

---

## 四、aidc-platform-service

### 职责
提供 gRPC 计算点管理服务，从 Redis 读取实时数据对外提供查询，**不直接生产或消费 Kafka 消息**。

### Kafka 角色
**仅有开关控制标志位**

### 相关代码

`pkg/service/dcim_metrics.go` 中的 `enableKafka` 字段：

```go
type DcimMetricsService struct {
    enableKafka bool       // Kafka 推送开关（默认 true）
    mu          sync.RWMutex
}

// HandleEnableKafka 通过 gRPC 接口动态控制 Kafka 推送
func (s *DcimMetricsService) HandleEnableKafka(ctx context.Context, req *apiv1.EnableKafkaRequest) (*apiv1.EnableKafkaResponse, error) {
    switch req.Enable {
    case "enable":  s.enableKafka = true
    case "disable": s.enableKafka = false
    default: // 查询当前状态
    }
}
```

**说明：** `enableKafka` 是为未来扩展预留的控制开关，当前服务本身没有实际的 Kafka Producer/Consumer 实现。该服务通过 **Redis** 直接读取由 `process-service` 和 `metric-service` 写入的数据，不需要通过 Kafka 中转。

---

## 五、aidc-north-sync

### 职责
管理北向接口配置（CRUD），**本身无 Kafka Producer/Consumer 代码**。

### Kafka 角色
**仅作为配置元数据管理方**

### Topic 字段的作用

`NorthConfigure` 数据模型中包含 `Topic` 字段：

```go
type NorthConfigure struct {
    ID            string  `json:"id"`
    MachineRoom   string  `json:"machine_room"`
    NorthAPIURL   string  `json:"north_api_url"`
    NorthToken    string  `json:"north_token"`
    NorthSource   string  `json:"north_source"`     // 数据源标识
    Version       string  `json:"version"`
    Topic         string  `json:"topic"`             // Kafka Topic 名称（配置项）
    ParsingMethod string  `json:"parsing_method"`    // 解析方式
}
```

**说明：** `Topic` 字段存储在 PostgreSQL 中，记录该北向数据源对应的 Kafka Topic 名称，供 `process-service` 等下游服务配置消费时参考，`aidc-north-sync` 自身不直接操作 Kafka。

---

## 六、全链路 Topic 汇总

| Topic 名称 | 生产者 | 消费者 | 用途 |
|-----------|--------|--------|------|
| `aidc-raw-data`（可配置） | aidc-pretreat-service | aidc-process-service（`raw` topic） | 北向原始数据传输 |
| `jinhua` | 外部北向（通过 APISIX） | aidc-process-service | 金华格式原始数据 |
| `raw` | aidc-pretreat-service | aidc-process-service | 通用原始数据 |
| `bms-D04-rtdata-shuoyao` | BMS 设备直推 | aidc-process-service | BMS 设备实时数据 |
| `aidc-calc-task` | aidc-metric-service（CalcTriggerEvaluator + DLQHandler） | aidc-metric-service（CalcProcessor） | 计算任务调度 |
| `aidc-calc-task-dlq` | aidc-metric-service（CalcProcessor） | aidc-metric-service（DLQHandler） | 计算失败死信补偿 |

---

## 七、关键设计模式总结

### 7.1 幂等生产（Idempotent Producer）
`aidc-pretreat-service` 支持 `enable_idempotence=true`，配合 `Net.MaxOpenRequests=1`，保证消息精确一次投递（Exactly Once Semantics）。

### 7.2 大消息处理策略
- **pretreat-service**：超过阈值的数据存 MinIO，Kafka 消息只携带 MinIO 引用路径
- **metric-service KafkaProducer**：超过 2.56MB 的消息自动递归对半拆分发送
- **metric-service CalcTrigger**：依赖超过 50 个时，依赖值冻结到 Redis Snapshot（TTL=24h），消息仅携带 snap_id

### 7.3 消费幂等保障
`CalcProcessor` 通过 Redis 双 Key 实现：
- `calc:idem:{key}`（TTL=25h）：持久幂等，防止同一任务被执行两次
- `calc:lock:{key}`（TTL=30s）：分布式锁，防止多实例并发执行同一任务

### 7.4 死信补偿（DLQ）
`aidc-calc-task-dlq` 消费者（DLQHandler）按错误类型分类处理，支持：
- 重投原始消息到主 topic
- P1/P2 告警上报
- 基于当前时间 + `bypassSkewCheck=true` 的补偿级联触发

### 7.5 限速控制
`aidc-process-service` 使用基于 `golang.org/x/time/rate` 的令牌桶限速器，支持：
- 全局默认限速（`DefaultRateLimit`）
- 每个 Topic 独立精细限速配置

### 7.6 级联触发防护
`CalcTriggerEvaluator` 在级联触发时：
- `HopCount` 递增传递，超过 `MaxCascadeDepth=10` 时丢弃并报 DLQ
- `CascadeRateLimiter` 对每个下游计算点独立限速，防止级联风暴
- 使用上游的 `TriggerTs`（非 `time.Now()`），保证 R2 Key 在级联链路中一致

---

## 八、各服务 Kafka 依赖库对比

| 服务 | 库 | 特性 |
|------|----|------|
| aidc-pretreat-service | `github.com/IBM/sarama` | 功能完整，支持幂等/压缩/批量/事务 |
| aidc-process-service | `github.com/segmentio/kafka-go` | 轻量，原生 Go 实现，Reader/Writer 模型 |
| aidc-metric-service | `github.com/segmentio/kafka-go` | 轻量，原生 Go 实现，Reader/Writer 模型 |
