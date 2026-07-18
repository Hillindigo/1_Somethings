# AIDC 系统集成分析 - 三个服务间的关联关系

## 系统架构概览

AIDC (工业数据中心) 系统由三个核心微服务组成，形成完整的数据处理和管理平台：

```
┌─────────────────────────────────────────────────────────────────┐
│                        AIDC 系统架构                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  数据源 → APISIX → aidc-pretreat-service → Kafka → aidc-process-service │
│                           ↓                           ↓         │
│                       MinIO 存储              VictoriaMetrics    │
│                                                       ↓         │
│                    aidc-platform-service ← ─ ─ ─ ─ ─ ─          │
│                    (认证、权限、管理)                              │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## 服务职责划分

### aidc-pretreat-service (数据预处理服务)
**角色**: 数据接入层
**职责**: 
- 接收外部数据源的 HTTP 请求
- 根据数据大小进行智能分流处理
- 数据压缩和对象存储
- 生成 Kafka 消息

### aidc-process-service (数据处理服务)
**角色**: 数据处理层
**职责**:
- 消费 Kafka 消息进行数据处理
- 与北向系统进行数据同步
- 时序数据存储和管理
- 数据解析和转换

### aidc-platform-service (平台服务)
**角色**: 平台管理层
**职责**:
- 用户认证和权限管理
- 组织架构和角色管理
- 系统配置和字典管理
- 为其他服务提供基础支撑

## 数据流向分析

### 主数据流

```
1. 数据接入阶段
   外部系统 → APISIX → aidc-pretreat-service
   
2. 数据分流处理
   aidc-pretreat-service → 判断数据大小
   ├── < 10MB: 直接发送到 Kafka (topic: raw-small)
   └── ≥ 10MB: 压缩存储到 MinIO + 发送索引到 Kafka (topic: raw-large-ref)
   
3. 数据处理阶段
   Kafka → aidc-process-service → 数据解析和处理
   ├── 小数据: 直接处理 JSON 内容
   └── 大数据: 从 MinIO 读取并解压处理
   
4. 数据存储阶段
   aidc-process-service → 存储到目标系统
   ├── Redis: 缓存和临时数据
   ├── VictoriaMetrics: 时序数据
   └── PostgreSQL: 关系型数据
```

### 认证授权流

```
1. 用户认证
   客户端 → aidc-platform-service → Casdoor → JWT Token
   
2. 服务间认证
   各服务 → aidc-platform-service → 验证 Token → 用户信息
   
3. 权限控制
   aidc-platform-service → 菜单权限 → API 权限 → 数据权限
```

## 技术集成点

### 1. 消息队列集成 (Kafka)

**aidc-pretreat-service (生产者)**:
- 生产消息到多个 Topic
- 支持消息分区和幂等性
- 批量发送优化

**aidc-process-service (消费者)**:
- 消费多个 Topic 的消息
- 支持限流和重试机制
- 批量处理优化

**Topic 设计**:
```
raw-small: 小于 10MB 的内联数据
raw-large-ref: 大于 10MB 的 MinIO 引用数据
jinhua: 特定数据源 (限流 5 msg/s)
bms-D04-rtdata-shuoyao: BMS 数据 (限流 20 msg/s)
```

### 2. 对象存储集成 (MinIO)

**aidc-pretreat-service**:
- 大文件压缩上传
- 分片上传支持
- 对象路径规划: `{year}/{month}/{day}/{trace_id}.json.zst`

**aidc-process-service**:
- 根据 Kafka 消息中的引用信息读取对象
- 解压缩处理
- 数据解析和转换

### 3. 缓存集成 (Redis)

**aidc-platform-service**:
- 用户会话管理
- 权限信息缓存
- 字典数据缓存

**aidc-process-service**:
- 数据处理结果缓存
- 临时数据存储
- 去重处理

### 4. 数据库集成 (PostgreSQL)

**aidc-platform-service**:
- 用户和权限数据
- 组织架构数据
- 系统配置数据

**aidc-process-service**:
- 同步任务记录
- 设备和节点信息
- 数据处理日志

### 5. 认证集成 (Casdoor)

**统一身份认证**:
- aidc-platform-service 作为认证中心
- 其他服务通过 gRPC 调用验证 Token
- 支持白名单机制

## 服务间通信

### 1. gRPC 通信

**aidc-platform-service** 提供的服务:
```
UserService: 用户管理
AuthService: 认证服务
AuthMenuService: 菜单权限
AuthRoleService: 角色管理
DictService: 字典服务
```

**aidc-process-service** 提供的服务:
```
HealthService: 健康检查
SyncService: 数据同步管理
```

### 2. HTTP 通信

**aidc-pretreat-service** 提供的端点:
```
POST /api/v1/data: 流式分片处理
POST /api/v1/dataMinio: MinIO 流式上传
POST /api/v1/dataOnce: 一次性上传
POST /api/v1/dataCompress: 流式压缩上传
GET /health: 健康检查
```

## 配置管理集成

### 环境变量前缀

```
aidc-pretreat-service: AIDC_PRETREAT_SERVICE_
aidc-process-service: AIDC_PROCESS_SERVICE_
aidc-platform-service: AIDC_PLATFORM_SERVICE_
```

### 共享配置项

**数据库配置**:
- 所有服务都使用 PostgreSQL
- 共享连接池配置模式
- 统一的事务管理策略

**Redis 配置**:
- 共享 Redis 集群
- 统一的缓存策略
- 一致的 Key 命名规范

**监控配置**:
- 统一的 Prometheus 指标
- 共享的日志格式
- 一致的链路追踪

## 部署架构

### Kubernetes 部署

```yaml
# 部署顺序和依赖关系
1. aidc-platform-service (基础服务)
   - 提供认证和权限服务
   - 其他服务的依赖基础

2. aidc-pretreat-service (数据接入)
   - 依赖 Kafka 和 MinIO
   - 独立的数据接入能力

3. aidc-process-service (数据处理)
   - 依赖 aidc-platform-service (认证)
   - 依赖 Kafka (消息消费)
   - 依赖 MinIO (大文件读取)
```

### 服务发现

```
通过 Kubernetes Service 实现服务发现:
- aidc-platform-service.default.svc.cluster.local
- aidc-pretreat-service.default.svc.cluster.local
- aidc-process-service.default.svc.cluster.local
```

## 数据一致性保证

### 1. 事务管理

**aidc-pretreat-service**:
- Kafka 生产者事务
- MinIO 上传原子性

**aidc-process-service**:
- Kafka 消费者事务
- 数据库事务管理

**aidc-platform-service**:
- 用户操作事务
- 权限变更事务

### 2. 幂等性设计

**TraceID 机制**:
- 每个请求生成唯一 TraceID
- 全链路追踪和去重
- 支持重试和恢复

### 3. 错误处理

**重试机制**:
- Kafka 消息重试
- HTTP 请求重试
- 数据库操作重试

**降级策略**:
- 认证服务降级
- 数据处理降级
- 存储服务降级

## 监控和可观测性

### 1. 指标监控

**统一指标**:
```
请求量: aidc_requests_total
响应时间: aidc_request_duration_seconds
错误率: aidc_errors_total
资源使用: aidc_resource_usage
```

### 2. 日志聚合

**日志格式统一**:
```json
{
  "timestamp": "2026-01-28T10:00:00Z",
  "service": "aidc-pretreat-service",
  "trace_id": "uuid",
  "level": "INFO",
  "message": "Request processed",
  "metadata": {...}
}
```

### 3. 链路追踪

**TraceID 传递**:
```
HTTP Request → aidc-pretreat-service (生成 TraceID)
                ↓
Kafka Message → aidc-process-service (继承 TraceID)
                ↓
gRPC Call → aidc-platform-service (传递 TraceID)
```

## 安全集成

### 1. 网络安全

**服务间通信**:
- gRPC TLS 加密
- Kubernetes Network Policy
- 服务网格 (可选)

### 2. 数据安全

**敏感数据处理**:
- 数据传输加密
- 存储加密
- 访问日志记录

### 3. 认证授权

**统一认证**:
- Casdoor 统一身份认证
- JWT Token 验证
- 基于角色的访问控制

## 扩展性设计

### 1. 水平扩展

**无状态设计**:
- 所有服务支持多实例部署
- 负载均衡和故障转移
- 数据库连接池共享

### 2. 垂直扩展

**资源配置**:
- CPU 和内存动态调整
- 存储容量弹性扩展
- 网络带宽优化

### 3. 功能扩展

**插件化架构**:
- 数据处理器插件
- 认证方式扩展
- 存储后端扩展

## 故障处理和恢复

### 1. 故障检测

**健康检查**:
- 服务健康状态监控
- 依赖服务状态检查
- 自动故障发现

### 2. 故障隔离

**熔断机制**:
- 服务间调用熔断
- 数据库连接熔断
- 外部依赖熔断

### 3. 故障恢复

**自动恢复**:
- 服务自动重启
- 数据自动同步
- 状态自动恢复

## 性能优化

### 1. 数据处理优化

**批量处理**:
- Kafka 批量消费
- 数据库批量操作
- 缓存批量更新

### 2. 网络优化

**连接复用**:
- gRPC 连接池
- HTTP 连接复用
- 数据库连接池

### 3. 存储优化

**分层存储**:
- 热数据 Redis 缓存
- 温数据 PostgreSQL
- 冷数据 MinIO 归档

## 总结

AIDC 系统通过三个微服务的协同工作，构建了完整的工业数据处理平台：

**核心优势**:
- 清晰的职责分离和模块化设计
- 高效的数据处理和存储策略
- 完善的认证授权和安全机制
- 良好的扩展性和可维护性
- 统一的监控和运维体系

**技术亮点**:
- 智能数据分流处理 (10MB 阈值)
- 统一身份认证和权限管理
- 高性能消息队列和对象存储
- 完整的链路追踪和监控
- 容器化部署和服务编排

该系统架构为工业物联网数据的接入、处理、存储和管理提供了完整的解决方案，具有良好的性能、可靠性和扩展性。
