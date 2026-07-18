# AIDC Pretreat Service 项目分析

## 项目概述

**aidc-pretreat-service** 是一个基于 Go 语言和 Gin 框架开发的数据预处理微服务，主要负责接收 APISIX 转发的 HTTP 请求，根据数据大小进行不同的处理策略。

### 基本信息

- **语言**: Go 1.22.10
- **Web框架**: Gin
- **版本**: 2.0.0
- **主要功能**: HTTP 数据接收、压缩处理、对象存储、消息队列
- **部署方式**: Docker 容器化部署

## 技术栈分析

### 核心依赖

- **Web框架**: `github.com/gin-gonic/gin v1.10.0` - 高性能 HTTP 框架
- **压缩**: `github.com/klauspost/compress v1.17.11` - zstd 压缩算法
- **对象存储**: `github.com/minio/minio-go/v7 v7.0.83` - MinIO 客户端
- **消息队列**: `github.com/IBM/sarama v1.45.0` - Kafka 客户端
- **配置管理**: `github.com/spf13/viper v1.19.0` - 配置解析
- **监控**: `github.com/prometheus/client_golang v1.19.0` - Prometheus 指标
- **日志**: `go.uber.org/zap v1.21.0` - 结构化日志
- **UUID**: `github.com/google/uuid v1.6.0` - 唯一标识生成

## 项目结构分析

### 目录结构

```
aidc-pretreat-service/
├── cmd/aidc-pretreat-service/    # 主程序入口
├── internal/                     # 内部包（不对外暴露）
│   ├── compress/                 # 压缩处理
│   ├── config/                   # 配置管理
│   ├── handler/                  # HTTP 处理器
│   ├── kafka/                    # Kafka 客户端
│   ├── processor/                # 数据处理器
│   ├── storage/                  # 存储层
│   └── worker/                   # 工作池
├── pkg/models/                   # 公共数据模型
├── docs/                         # 文档
├── deploy/                       # 部署配置
└── charts/                       # Helm 部署图表
```

## 核心功能模块

### 1. 数据处理策略

**核心逻辑**: 根据请求体大小采用不同处理策略

**处理流程**:
```
APISIX 转发请求 → Pretreat Service → 判断数据大小
├── < 10MB → 直接发送到 Kafka (topic: raw-small)
└── ≥ 10MB → zstd 压缩 → 上传 MinIO → 发送索引到 Kafka (topic: raw-large-ref)
```

**阈值配置**: 默认 10MB (`DefaultSizeThreshold = 10 * 1024 * 1024`)

### 2. 多种处理器模式

**文件位置**: `internal/processor/`

#### 流式分片处理器 (StreamChunkProcessor)
- **用途**: 处理大数据的流式分片
- **特点**: 内存友好，支持超大文件
- **优化版本**: `stream_chunk_optimized.go` 提供性能优化

#### MinIO 处理器 (MinIOProcessor)
- **用途**: 直接上传到 MinIO 对象存储
- **特点**: 支持多部分上传，断点续传
- **压缩**: 集成 zstd 压缩算法

#### 一次性处理器 (OneTimeProcessor)
- **用途**: 小文件一次性处理
- **特点**: 简单快速，适合小数据量

#### 压缩处理器 (CompressProcessor)
- **用途**: 专门处理需要压缩的数据
- **特点**: 流式压缩上传，节省内存

### 3. HTTP 处理器

**文件位置**: `internal/handler/`

#### 多端点支持
- **`/api/v1/data`**: 流式分片处理
- **`/api/v1/dataMinio`**: MinIO 流式上传
- **`/api/v1/dataOnce`**: 一次性上传处理
- **`/api/v1/dataCompress`**: 流式压缩上传
- **`/health`**: 健康检查
- **`/ready`**: 就绪检查
- **`/live`**: 存活检查

#### 请求处理特性
- 支持大文件流式处理
- 自动生成 TraceID 用于链路追踪
- Prometheus 指标收集
- 错误处理和重试机制

### 4. 消息格式设计

**文件位置**: `pkg/models/messages.go`

#### 统一消息格式 (UnifiedMessage)
```json
{
  "version": "v1",
  "trace_id": "uuid",
  "source": "apisix", 
  "received_at": 1769586077462,
  "data_type": "inline|reference",
  "payload": {...},        // 内联数据（小于10MB）
  "storage": {...},        // 存储引用（大于10MB）
  "meta": {...}           // 元数据信息
}
```

#### 小数据消息 (< 10MB)
- **data_type**: "inline"
- **payload**: 包含完整的原始数据
- **topic**: raw-small

#### 大数据消息 (≥ 10MB)
- **data_type**: "reference"
- **storage**: MinIO 对象信息
- **meta**: 设备数量、点位数量等统计信息
- **topic**: raw-large-ref

### 5. 存储策略

#### MinIO 对象存储
- **Bucket**: aidc-raw-data
- **路径格式**: `{year}/{month}/{day}/{trace_id}.json.zst`
- **压缩算法**: zstd
- **分片上传**: 支持大文件分片上传

#### Kafka 消息队列
- **生产者配置**: 支持幂等性、重试机制
- **批量发送**: 可配置批量大小和超时时间
- **可靠性**: RequiredAcks 确保消息持久化

## 配置管理

### 环境变量配置

服务使用 `AIDC_PRETREAT_SERVICE_` 前缀的环境变量：

#### 服务器配置
- **HOST**: 监听地址
- **PORT**: 监听端口
- **READ_TIMEOUT**: 读取超时
- **WRITE_TIMEOUT**: 写入超时
- **IDLE_TIMEOUT**: 空闲超时

#### MinIO 配置
- **ENDPOINT**: MinIO 服务地址
- **ACCESS_KEY_ID**: 访问密钥ID
- **SECRET_ACCESS_KEY**: 访问密钥
- **BUCKET**: 存储桶名称
- **USE_SSL**: 是否使用 SSL
- **PART_SIZE**: 分片大小

#### Kafka 配置
- **BROKERS**: Broker 地址列表
- **TOPIC**: 主题名称
- **PRODUCER_ID**: 生产者ID
- **MAX_RETRIES**: 最大重试次数
- **ENABLE_IDEMPOTENCE**: 幂等性开关

#### 处理器配置
- **SIZE_THRESHOLD**: 大小阈值（默认10MB）
- **SOURCE**: 数据源标识
- **MAX_BODY_SIZE**: 最大请求体大小

#### 压缩配置
- **LEVEL**: zstd 压缩级别（1-22）

## 性能优化

### 内存管理
- 流式处理避免大文件完全加载到内存
- 分片上传减少内存占用
- 对象池复用减少 GC 压力

### 并发处理
- Worker Pool 模式处理并发请求
- 异步处理避免阻塞 HTTP 响应
- 批量操作提高吞吐量

### 压缩优化
- zstd 算法提供高压缩比和快速压缩
- 流式压缩支持大文件处理
- 可配置压缩级别平衡性能和压缩比

## 监控和可观测性

### Prometheus 指标
- **pretreat_requests_total**: 请求总数统计
- **pretreat_request_duration_seconds**: 请求处理时长
- **pretreat_request_size_bytes**: 请求大小分布
- **pretreat_compression_ratio**: 压缩比统计

### 链路追踪
- 每个请求生成唯一 TraceID
- 全链路日志记录
- 错误和异常追踪

### 健康检查
- **/health**: 详细健康状态
- **/ready**: 就绪状态检查
- **/live**: 存活状态检查

## 部署和运维

### Docker 支持
- 多阶段构建优化镜像大小
- 健康检查配置
- 环境变量配置

### Kubernetes 部署
- Helm Chart 支持
- ConfigMap 和 Secret 管理
- 水平扩展支持

### 配置验证
- 启动时配置验证
- 必需参数检查
- 错误提示和建议

## 架构设计原则

### 异步处理
- HTTP 请求快速响应
- 后台异步处理数据
- 避免长时间阻塞

### 消息不丢失
- 多级缓冲机制
- 重试和错误处理
- 可靠的消息传递

### 高并发支持
- Worker Pool 并发处理
- 无状态设计
- 水平扩展能力

### 可观测性
- 全链路监控
- 详细指标收集
- 结构化日志

### 幂等性
- 支持重复请求处理
- TraceID 去重机制
- 安全的重试策略

## 数据流向

```
数据源 → APISIX → Pretreat Service
                      ↓
                 判断数据大小
                ↙            ↘
        < 10MB              ≥ 10MB
           ↓                   ↓
    直接发送到Kafka      zstd压缩 → MinIO
    (topic: raw-small)         ↓
                        发送索引到Kafka
                     (topic: raw-large-ref)
```

## 总结

aidc-pretreat-service 是一个高性能的数据预处理服务，具有以下特点：

**优势**:
- 智能的数据大小判断和处理策略
- 多种处理器模式适应不同场景
- 高效的压缩和存储机制
- 完善的监控和可观测性
- 良好的并发性能和内存管理
- 灵活的配置管理

**适用场景**:
- 工业物联网数据预处理
- 大数据量的HTTP接入
- 数据压缩和存储优化
- 消息队列数据分发
- 高并发数据处理

该服务在 AIDC 系统中承担数据接入和预处理的关键角色，通过智能的处理策略和高效的技术实现，为后续的数据处理流程提供可靠的数据源。
