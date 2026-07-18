# aidc-metric-service 项目分析

## 一、核心功能与整体架构

### 1.1 项目定位

`aidc-metric-service` 是 AIDC 平台的**指标计算服务**，负责对设备原始数据进行公式计算、聚合统计，支持两种计算模式：**定时调度计算**和**流式实时计算**（Kafka calc-task 驱动）。计算结果双写到 Redis 和 VictoriaMetrics。

### 1.2 核心功能

- **定时调度计算**：支持 10s/1min/5min/1h/日/月/交接班 7 种周期
- **流式计算引擎**：消费 Kafka `aidc-calc-task` topic，实时触发公式计算
- **公式引擎**：支持四则运算、聚合函数，公式 AST 预编译
- **拓扑排序**：按数据点依赖关系分层，确保依赖先算、被依赖后算
- **级联触发**：计算完成后自动触发下游依赖的计算点
- **DLQ 补偿机制**：计算失败的任务进入死信队列，自动重试（最多 5 次）
- **热加载**：每 5 分钟从 Redis 重新加载公式配置，无需重启服务
- **双写存储**：计算结果写入 Redis（实时值）+ VictoriaMetrics（历史数据）

### 1.3 架构流程图

```text
  ┌────────────────┐      ┌────────────────────────┐
  │ Redis          │      │ Kafka                  │
  │ formula:*      │      │ aidc-calc-task (主题)   │
  │ (计算点配置)    │      │ aidc-calc-task-dlq     │
  └──────┬─────────┘      └──────────┬─────────────┘
         │                           │
         ▼                           ▼
  ┌──────────────────────────────────────────────────┐
  │            aidc-metric-service                    │
  │                                                  │
  │  ┌───────────────────────────────────────┐       │
  │  │ MetricPointLoader (公式加载器)         │       │
  │  │ • 从 Redis 加载 formula:* 公式        │       │
  │  │ • 预编译 AST                          │       │
  │  │ • 构建依赖索引 CalcDependencyIndex    │       │
  │  │ • 每 5 分钟热加载                     │       │
  │  └───────────────┬───────────────────────┘       │
  │                  │                               │
  │  ┌───────────────▼───────────────┐               │
  │  │  两种计算模式                  │               │
  │  │                               │               │
  │  │  模式 1: 定时调度计算          │               │
  │  │  ┌──────────────────────┐    │               │
  │  │  │ Scheduler            │    │               │
  │  │  │ • 10s/1min/5min/1h  │    │               │
  │  │  │ • 日/月/交接班      │    │               │
  │  │  └──────────┬───────────┘    │               │
  │  │             ▼                │               │
  │  │  ┌──────────────────────┐    │               │
  │  │  │ DatapointManager     │    │               │
  │  │  │ • 拓扑排序按层计算   │    │               │
  │  │  │ • Worker Pool 并发   │    │               │
  │  │  └──────────────────────┘    │               │
  │  │                               │               │
  │  │  模式 2: 流式计算引擎          │               │
  │  │  ┌──────────────────────┐    │               │
  │  │  │ CalcProcessor        │    │               │
  │  │  │ • Kafka 消费 calc-task│    │               │
  │  │  │ • 公式 AST 求值      │    │               │
  │  │  │ • 级联触发下游       │    │               │
  │  │  │ • 失败 → DLQ         │    │               │
  │  │  └──────────────────────┘    │               │
  │  └───────────────┬───────────────┘               │
  │                  │                               │
  │         ┌────────┼────────┐                      │
  │         ▼        ▼        ▼                      │
  │  ┌──────────┐ ┌────────┐ ┌───────────────┐      │
  │  │ Redis    │ │ VMS    │ │ Kafka DLQ     │      │
  │  │ Writer   │ │ Writer │ │ (失败补偿)    │      │
  │  └──────────┘ └────────┘ └───────────────┘      │
  └──────────────────────────────────────────────────┘
```

---

## 二、核心技术栈

| 分类 | 技术 | 用途 |
|------|------|------|
| 语言 | Go 1.22+ | 主语言 |
| RPC 框架 | gRPC + Protobuf | 对外 API（健康检查） |
| 数据库 | PostgreSQL + GORM | 计算点配置数据（DcimMetricPoints 表） |
| 消息队列 | Kafka (segmentio/kafka-go) | calc-task 消费 + DLQ 补偿 + 级联触发 |
| 缓存 | Redis (go-redis) | 实时数据读写、公式配置加载 |
| 时序数据库 | VictoriaMetrics | 计算结果历史存储 |
| 公式引擎 | 自研 formula.FormulaEngine | AST 解析 + 求值 |
| 定时调度 | 自研 Scheduler (time.Ticker) | 7 种周期调度 |
| 日志 | klog | 结构化日志 |
| 配置 | envconfig | 50+ 环境变量配置 |

---

## 三、核心代码文件/目录及作用

```text
aidc-metric-service/
├── cmd/aidc-metric-service/
│   └── main.go                            # 服务入口（680+ 行，含 calcEngine 初始化）
├── pkg/
│   ├── api/
│   │   └── health.go                     # 健康检查 gRPC 服务
│   ├── base/
│   │   └── base.go                       # 基础组件初始化（DB/Redis/VMS/Kafka）
│   ├── calc/                              # 流式计算引擎（核心）
│   │   ├── calc_processor.go             # 计算处理器（消费 calc-task）
│   │   ├── metric_point_loader.go        # 公式加载器（从 Redis 加载）
│   │   ├── calc_dependency_index.go      # 依赖索引（dpid → 上下游关系）
│   │   └── types.go                      # 核心类型定义
│   ├── calc/formula/
│   │   └── formula_engine.go             # 公式引擎（AST 解析 + 求值）
│   ├── config/
│   │   └── config.go                     # 配置结构（50+ 参数）
│   ├── dlq/
│   │   └── dlq_handler.go               # DLQ 死信队列处理器
│   ├── service/
│   │   ├── metric_service.go             # 指标服务（API 层面）
│   │   ├── dcim_metrics_service.go       # DCIM 指标服务
│   │   ├── datapoints/
│   │   │   ├── manager.go               # 数据点管理器（核心 770+ 行）
│   │   │   ├── calculate_task_manager.go # 计算任务管理器（Worker Pool）
│   │   │   ├── datapoint.go             # 运行时数据点
│   │   │   └── sort.go                  # 拓扑排序算法
│   │   └── scheduler/
│   │       └── scheduler.go             # 定时调度器（7 种周期）
│   ├── storage/
│   │   ├── redis_client.go              # Redis 客户端封装
│   │   ├── redis_writer.go              # Redis 写入器
│   │   ├── victoria_writer.go           # VMS 写入器
│   │   ├── victoria_reader.go           # VMS 读取器
│   │   └── kafka_producer.go            # Kafka 生产者
│   └── dax/gen/model/
│       └── dcim_metric_points.go        # 计算点数据模型
├── configs/                              # 配置文件
├── Makefile                              # 构建脚本
└── Dockerfile                            # 容器镜像
```

---

## 四、核心业务逻辑分析

### 4.1 启动流程

```text
main()
  ├── 加载配置（envconfig，前缀 AIDC_METRIC_SERVICE_）
  ├── base.Init()：初始化 PostgreSQL
  ├── base.InitRedis()：初始化 Redis
  ├── base.InitVictoriaMetrics()：初始化 VMS 查询 URL
  ├── base.InitVictoriaMetricsWriter()：初始化 VMS 写入配置
  ├── base.InitKafka()：初始化 Kafka Writer
  ├── 注册 gRPC 服务（HealthService）
  │
  ├── 异步启动：initMetricService()（定时调度计算）
  │     ├── 创建 DatapointManager
  │     │   ├── 加载所有计算点（loadAllDatapoints）
  │     │   ├── 构建索引（buildIndexes）
  │     │   ├── 加载计算参数（loadCalParams）
  │     │   └── 拓扑排序（sortDatapoints）
  │     ├── 创建 Scheduler（7 种周期）
  │     └── scheduler.Start()
  │
  ├── 异步启动：initCalcEngine()（流式计算引擎）
  │     ├── 创建 MetricPointLoader → 从 Redis 加载公式
  │     ├── 创建 FormulaEngine → AST 预编译
  │     ├── 创建 CalcProcessor
  │     ├── 启动 calc-task Kafka consumer
  │     ├── 启动 DLQ consumer + handler
  │     └── 启动热加载 goroutine（每 5 分钟）
  │
  └── run.DefaultRun()：启动 gRPC 服务器
```

### 4.2 定时调度计算流程

```text
Scheduler.schedule10s() / schedule1min() / ...
  │
  ▼
DatapointManager.Calculate(intervalType)
  ├── 获取对应 intervalType 的排序结果 sortedMap
  ├── 按层级 level 顺序迭代（level 0 → level N）
  │     │
  │     ├── loadDatapointValues(dps)：从 Redis 加载依赖点的实时值
  │     │
  │     └── calculateTaskManager.AddCalculateData(dps)
  │           │
  │           ▼
  │         Worker Pool（100 个 Worker）
  │           ├── 对每个数据点执行公式计算
  │           ├── calculateCallback()
  │           │   ├── 读取依赖参数值（Redis HGET）
  │           │   ├── 执行公式（四则运算/聚合函数）
  │           │   ├── 写入 Redis（计算结果）
  │           │   ├── 写入 VictoriaMetrics（历史记录）
  │           │   └── 发送 Kafka calc-task（触发下游）
  │           └── 等待当前层级全部完成
  │
  └── 进入下一层级（确保依赖关系正确）
```

### 4.3 流式计算引擎流程（CalcEngine）

```text
Kafka Message (topic: aidc-calc-task)
  │
  ▼
CalcProcessor.Handle(msg)
  ├── 1. 解析消息：提取 dpid、source、依赖值
  ├── 2. 查询 CalcDependencyIndex → 找到需要计算的计算点
  ├── 3. 从 Redis 读取所有依赖点的最新值
  ├── 4. FormulaEngine.EvaluateAST()：执行公式计算
  │     ├── 成功：
  │     │   ├── RedisWriter.Write()：写入计算结果
  │     │   ├── VictoriaWriter.Write()：写入历史数据
  │     │   └── triggerFunc()：级联触发下游计算点
  │     │         └── 查询 index → 发送新的 calc-task 到 Kafka
  │     └── 失败：
  │         └── 发送到 DLQ topic（aidc-calc-task-dlq）
  └── CommitMessages()：提交 offset
```

### 4.4 拓扑排序算法

```text
sortDatapoints()
  │
  ├── 输入：同一 intervalType 的所有数据点
  ├── 构建依赖图（DAG）
  │     └── 每个数据点的 cal_parameter 中记录了依赖的其他 dpid
  ├── 拓扑排序（Kahn's Algorithm）
  │     ├── Level 0：无依赖的数据点（叶子节点）
  │     ├── Level 1：只依赖 Level 0 的数据点
  │     ├── Level N：依赖 Level N-1 的数据点
  │     └── 检测循环依赖 → 跳过并告警
  └── 输出：map[level]→[]*DataPoint，按层级分组

作用：确保计算顺序正确——先算底层原始指标，再算依赖它们的复合指标
```

### 4.5 DLQ 补偿机制

```text
DLQHandler.Start()
  ├── 消费 aidc-calc-task-dlq topic
  ├── 检查重试次数（retry_count）
  │     ├── < 5 次：
  │     │   ├── 等待 RetryInterval（2 秒）
  │     │   ├── retry_count++
  │     │   └── 重新发送到 aidc-calc-task topic
  │     └── ≥ 5 次：
  │         └── 记录告警日志，放弃重试
  └── CommitMessages()
```

---

## 五、面试常见问题及回答思路

### Q1: 为什么需要拓扑排序？

> **回答**：计算点之间存在依赖关系。例如「楼层总功率」依赖「各机房功率」，「各机房功率」又依赖各设备的原始数据。如果不按顺序计算，可能用到上一轮的过时数据。拓扑排序将数据点按依赖层级分组，逐层计算，确保每个数据点在计算时，它所依赖的数据点已经完成计算。

### Q2: 流式计算引擎和定时调度的区别是什么？

> **回答**：
> - **定时调度**：按固定周期（10s/1min/5min 等）批量计算所有同类型的数据点，适合**周期性汇总统计**（如每小时用电量、每日平均温度）
> - **流式计算**：由 Kafka 消息驱动，当某个原始数据点更新时，立即触发依赖它的计算点重新计算，适合**实时指标**（如实时 PUE、实时功率）
> - 两种模式互补：流式保证实时性，定时保证完整性

### Q3: 公式引擎是怎么实现的？

> **回答**：采用 AST（抽象语法树）方案：
> 1. **解析阶段**：将公式字符串解析为 AST 树，在服务启动时预编译，运行时不再解析
> 2. **求值阶段**：遍历 AST 树，递归求值。支持四则运算（+、-、*、/）和聚合函数（SUM、AVG、MAX、MIN 等）
> 3. **依赖注入**：将依赖点的值作为 `DependencyValue` 数组传入，公式中通过 dpid 引用

### Q4: DLQ 补偿机制解决什么问题？

> **回答**：流式计算中，可能因为依赖数据尚未到达 Redis、公式配置错误、临时网络问题等原因导致计算失败。DLQ 机制将失败消息发送到死信队列，由独立的 DLQHandler 消费并重试（最多 5 次，间隔 2 秒）。这样既不阻塞主消费流程，又能最大限度保证计算完成。

### Q5: 热加载是怎么实现的？

> **回答**：MetricPointLoader 每 5 分钟从 Redis 重新扫描 `formula:*` 键，对比已有索引：
> - 新增的公式：预编译 AST 并加入索引
> - 删除的公式：从索引中移除
> - 修改的公式：重新编译替换
> 整个过程不需要重启服务，实现了运行时动态更新。

### Q6: Worker Pool 是怎么设计的？

> **回答**：CalculateTaskManager 内部维护一个固定大小（默认 100）的 Worker Pool：
> - 使用带缓冲的 channel 作为任务队列
> - Worker goroutine 从 channel 取任务，执行计算回调
> - 支持超时控制（默认 5 秒），防止单个计算任务阻塞整个 Pool
> - 按层级提交任务，等待当前层级全部完成后再提交下一层级

### Q7: 级联触发是怎么实现的？

> **回答**：当某个计算点计算完成并成功写入 Redis 后，通过 `triggerFunc` 查询 `CalcDependencyIndex`，找到所有依赖该计算点的下游计算点，为每个下游点生成新的 `calc-task` 消息发送到 Kafka。下游点的 CalcProcessor 收到消息后再次执行计算，形成级联链条。这样实现了「一个原始数据更新 → 逐层传播计算」的实时效果。
