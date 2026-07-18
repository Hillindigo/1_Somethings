# Redis 使用分析

## 概述

本项目中有两个服务使用了 Redis：

| 服务 | 语言 | Redis 客户端 | 主要用途 |
|---|---|---|---|
| `repair-chassis-service` | Java (Spring Boot) | Spring Data Redis (`RedisTemplate`) | 权限缓存、验证码存储、Token 管理 |
| `repair-workflow-service` | Go | `go-redis/redis/v8` | 分布式锁、跨节点数据传递缓存 |

`repair-maintenance-service` 和 `repair-notify-service` 仅定义了 Redis 相关错误常量，**本身不直接使用 Redis**。

---

## 一、repair-chassis-service（Java）

### 1.1 Redis 配置

**文件：** `repair-chassis-server/src/main/java/sh/bs/admin/server/config/RedisConfig.java`

```java
@Configuration
@EnableCaching
@AutoConfigureBefore(RedisAutoConfiguration.class)
public class RedisConfig extends CachingConfigurerSupport {
    @Bean
    public RedisTemplate<String, Object> redisTemplate(RedisConnectionFactory connectionFactory) {
        RedisTemplate<String, Object> redisTemplate = new RedisTemplate<>();
        redisTemplate.setKeySerializer(new StringRedisSerializer());
        redisTemplate.setHashKeySerializer(new StringRedisSerializer());
        redisTemplate.setValueSerializer(new JdkSerializationRedisSerializer());
        redisTemplate.setHashValueSerializer(new JdkSerializationRedisSerializer());
        redisTemplate.setConnectionFactory(connectionFactory);
        return redisTemplate;
    }
}
```

**要点：**
- Key 使用 `StringRedisSerializer`（字符串序列化）
- Value 使用 `JdkSerializationRedisSerializer`（JDK 序列化）
- 通过 Spring Boot 自动配置读取 `application.yml` 中的 Redis 连接信息

### 1.2 RedisService 工具类

**文件：** `repair-chassis-server/src/main/java/sh/bs/admin/server/utils/RedisService.java`

封装了对 `RedisTemplate` 的所有操作，提供统一的缓存接口：

| 方法 | 说明 |
|---|---|
| `setCacheObject(key, value)` | 缓存对象（无过期时间） |
| `setCacheObject(key, value, timeout, unit)` | 缓存对象（有过期时间） |
| `getCacheObject(key, clazz)` | 获取缓存对象 |
| `deleteObject(key)` | 删除单个缓存 |
| `expire(key, timeout)` | 设置过期时间 |
| `hasKey(key)` | 判断 key 是否存在 |
| `setCacheList / getCacheList` | List 类型缓存 |
| `setCacheSet / getCacheSet` | Set 类型缓存 |
| `setCacheMap / getCacheMap` | Hash 类型缓存 |
| `getKeyExpiration(key)` | 获取 key 剩余 TTL |

**Key 命名规则：** 所有 key 自动加上应用名前缀，格式为 `${spring.application.name}:${key}`，例如：`repair-chassis-service:user123`。

---

### 1.3 具体使用场景

#### 场景一：用户登录 - 权限列表缓存

**触发时机：** 用户调用 `/web/user/login` 登录接口

**文件：** `WebUserController.java`（第 168 行）、`WebCustomerController.java`（第 168 行）

**逻辑：**
1. 用户登录成功后，从数据库查询该用户完整权限列表（`Auth::getPerms`）
2. 以用户名（`user.name`）为 key，将权限列表 JSON 字符串存入 Redis
3. 过期时间与 Token 的 `expires_in` 一致（动态设置）

```java
// 登录时缓存权限
redisService.setCacheObject(
    user.name,
    JSONUtil.toJsonStr(loginUserAuths.stream().map(Auth::getPerms).collect(Collectors.toList())),
    jsonObject.getLong("expires_in"),
    TimeUnit.SECONDS
);
```

**目的：** 避免每次请求都查询数据库获取权限，提升 API 响应速度。

---

#### 场景二：获取用户信息 - 权限缓存刷新

**触发时机：** 用户调用 `/web/user/info` 获取自身信息

**文件：** `WebUserController.java`（第 230 行）

**逻辑：**
- 每次获取用户信息时，主动刷新该用户的权限缓存（TTL 固定为 86400 秒即 24 小时）

```java
redisService.setCacheObject(
    CommonContextHolder.get(UserConstant.USERNAME).toString(),
    JSONUtil.toJsonStr(loginUserAuths.stream().map(Auth::getPerms).collect(Collectors.toList())),
    86400L,
    TimeUnit.SECONDS
);
```

---

#### 场景三：角色/部门变更 - 缓存主动更新

**触发时机：** 管理员修改角色绑定用户、部门绑定用户等操作

**文件：** `RoleServiceImpl.java`（`updateRedisAuth` 方法，第 496 行）、`DepartmentServiceImpl.java`（第 893 行）

**逻辑：**
1. 当角色/部门的用户关系发生变更时，遍历受影响的用户
2. 判断该用户 key 在 Redis 中是否存在（`hasKey`）
3. 若存在，则重新查询并更新其权限缓存（TTL 86400 秒）

```java
public void updateRedisAuth(List<String> userIds) {
    List<User> userList = userService.listByIds(userIds);
    userList.forEach(user -> {
        String username = user.getEmail().replace("@", "").replace(".", "");
        if (redisService.hasKey(username)) {
            List<Auth> auths = authService.getUserAuthsByUserId(user.getUserId());
            redisService.setCacheObject(username, JSONUtil.toJsonStr(
                auths.stream().map(Auth::getPerms).collect(Collectors.toList())
            ), 86400L, TimeUnit.SECONDS);
        }
    });
}
```

**目的：** 保证权限变更后缓存数据的实时性，防止旧权限数据被继续使用。

---

#### 场景四：用户禁用 - 主动删除缓存

**触发时机：** 管理员禁用某个用户账号

**文件：** `WebCustomerController.java`（第 784 行）

**逻辑：**
- 禁用用户时，主动从 Redis 中删除该用户的权限缓存，确保该用户下次请求时无法通过权限校验

```java
if (updStatusParam.getIsForbidden()) {
    redisService.deleteObject(user.name);
}
```

---

#### 场景五：图形验证码存储

**触发时机：** 用户请求获取登录验证码

**文件：** `ICaptchaServiceImpl.java`

**逻辑：**
1. 生成验证码（图片 + 验证答案）
2. 以 `loginId`（UUID）为 key，将验证码对象存入 Redis，TTL 为 **300 秒（5 分钟）**
3. 用户提交验证码时，从 Redis 读取并校验

```java
// 存储验证码
redisService.setCacheObject(loginId, captcha, 300L, TimeUnit.SECONDS);

// 校验验证码
val captcha = redisService.getCacheObject(loginId, PicCaptcha.class);
if (ObjectUtil.isNull(captcha)) {
    throw new ServiceException(PortalError.CAPTCHA_EXPIRED);
}
```

---

#### 场景六：gRPC 接口 - 权限列表查询

**触发时机：** 其他服务通过 gRPC 调用 `getAuthPerms` 接口

**文件：** `GrpcAuthService.java`（第 311 行）

**逻辑：**
- 通过 gRPC 对外暴露权限查询接口，直接从 Redis 中读取当前用户（以 JTI 为 key）的权限列表

```java
String cacheObject = redisService.getCacheObject(
    CommonContextHolder.get(UserConstant.JTI).toString(),
    String.class
);
```

> **注意：** 此处 key 使用 `JTI`（JWT ID），但 `GlobalFilter.java` 中该逻辑已被注释掉，实际登录时以 `username` 为 key 缓存权限。此 gRPC 接口在当前逻辑下存在 key 不一致的风险，可能导致查询结果为空。

---

### 1.4 权限缓存 Key 命名汇总

| Key 格式 | 存储内容 | TTL | 来源 |
|---|---|---|---|
| `{appName}:{username}` | 用户权限列表（JSON 数组） | 与 Token 一致 / 86400s | 登录、获取用户信息 |
| `{appName}:{loginId}` | 图形验证码对象 | 300s | 验证码生成 |

---

## 二、repair-workflow-service（Go）

### 2.1 Redis 配置

**配置文件：** `pkg/conf/config.go`（环境变量）

```go
RedisSentinelNodes string `default:"localhost:23306" envconfig:"REDIS_SENTINEL_NODES"`
RedisMaster        string `default:"mymaster"        envconfig:"REDIS_MASTER"`
RedisPassword      string `default:"xxx"             envconfig:"REDIS_PASSWORD"`
RedisDatabase      int    `default:"5"               envconfig:"REDIS_DATABASE"`
```

**初始化：** `pkg/base/base.go`

```go
err = InitRedis(cf.RedisMaster, cf.RedisPassword, strings.Split(cf.RedisSentinelNodes, ","), cf.RedisDatabase)
```

**客户端实现：** `pkg/base/redis.go`（使用 `go-redis/redis/v8`）

> **注意：** 配置中预留了 Sentinel（哨兵模式）字段 `RedisMaster` 和 `RedisSentinelNodes`，但当前代码实际使用的是单机模式（`redis.NewClient`），Sentinel 连接代码已被注释。

```go
// 当前实际使用单机模式
rdb := redis.NewClient(&redis.Options{
    Addr:     addrs[0],
    Password: password,
    DB:       db,
})
```

---

### 2.2 核心功能：分布式锁

**文件：** `pkg/base/redis.go`

提供三个层次的锁操作：

#### 基础锁操作

```go
// 加锁（SetNX 原子操作，成功返回 true）
func RedisLock(key string, expiration time.Duration) (bool, error)

// 解锁（删除 key）
func RedisUnlock(key string) error
```

#### 自动续期的同步锁

```go
func RedisSyncLock(key string, do func() error) error
```

**实现机制：**
1. 初始锁超时设为 **10 秒**
2. 启动后台 goroutine，每隔 **5 秒** 自动续期（重新加锁 10 秒）
3. 业务逻辑 `do()` 执行完毕后，取消续期并释放锁
4. `defer + recover` 保证 panic 场景下也能释放锁

```go
func RedisSyncLock(key string, do func() error) error {
    ok, err := RedisLock(key, time.Second*10)
    // ...
    go func() {
        for {
            select {
            case <-ctx.Done():
                RedisUnlock(key)
                return
            case <-ticker.C:
                RedisLock(key, time.Second*10) // 自动续期
            }
        }
    }()
    return do()
}
```

**目的：** 防止多个 workflow-service 实例并发处理同一任务，确保同一工作流节点的回调逻辑幂等执行。

---

### 2.3 具体使用场景

#### 场景一：批量更新状态时的仓库参数跨节点传递

**文件：** `pkg/service/task.go`（写入）、`pkg/service/callback_grpc_maintenance.go`（读取）

**背景：** 工作流中的任务节点在触发批量更新设备状态时，需要将"仓库（warehouse）"参数传递给下一个回调处理节点，但两个节点可能运行在不同的服务实例上。

**写入端（task.go）：**
```go
cacheKey := "batch_update_status_warehouse"
if err := base.RedisSet(ctx.Context, cacheKey, warehouseStr, 0); err != nil {
    klog.Errorf("缓存 batch_update_status_warehouse 失败 (TaskID: %s): %v", taskInfo.ID, err)
} else {
    klog.Infof("缓存 batch_update_status_warehouse 成功 (TaskID: %s, Warehouse: %s)", taskInfo.ID, warehouseStr)
}
```

**读取端（callback_grpc_maintenance.go）：**
```go
cacheKey := "batch_update_status_warehouse"
if cachedWarehouse, err := base.RedisGet(ctx.Context, cacheKey); err == nil && cachedWarehouse != "" {
    warehouse = cachedWarehouse
    defer func() {
        if delErr := base.RedisDel(ctx.Context, cacheKey); delErr != nil {
            klog.Errorf("删除缓存 batch_update_status_warehouse 失败: %v", delErr)
        }
    }()
}
```

**用后即删：** 读取完毕后立即通过 `defer RedisDel` 清理 key，避免脏数据残留。

---

### 2.4 Redis 操作封装

**文件：** `pkg/base/redis.go`

| 函数 | 说明 |
|---|---|
| `InitRedis(...)` | 初始化 Redis 连接，Ping 检查可用性 |
| `GetRedisClient()` | 获取全局 Redis 客户端实例 |
| `RedisLock(key, expiration)` | SetNX 加锁 |
| `RedisUnlock(key)` | Del 解锁 |
| `RedisSyncLock(key, do)` | 带自动续期的同步执行锁 |
| `RedisSet(ctx, key, value, expire)` | 设置字符串值 |
| `RedisGet(ctx, key)` | 获取字符串值 |
| `RedisDel(ctx, key)` | 删除 key |

---

## 三、不使用 Redis 的服务

### repair-maintenance-service

仅在 `pkg/dax/errors.go` 中定义了两个 Redis 相关错误常量，**实际业务代码中无 Redis 操作**：

```go
ErrRedis            = errors.New("redis连接出错，请联系管理员")
ErrAnalyzeRedisData = errors.New("redis数据解析失败")
```

这说明该服务在设计阶段曾规划使用 Redis（可能用于设备批量导入场景的数据缓存），但当前版本尚未实现，相关功能应由 `repair-workflow-service` 代理完成。

### repair-notify-service

完全不涉及 Redis，专注于消息通知（RocketMQ 消费 + 钉钉推送）。

---

## 四、整体架构图

```
用户登录请求
    │
    ▼
repair-chassis-service
    ├── 登录成功 ──────────► Redis: SET {username} {权限列表} TTL=token有效期
    ├── 获取用户信息 ────────► Redis: SET {username} {权限列表} TTL=86400s
    ├── 角色/部门变更 ───────► Redis: SET {username} {新权限列表}（仅更新已登录用户）
    ├── 禁用用户 ────────────► Redis: DEL {username}
    ├── 验证码生成 ──────────► Redis: SET {loginId} {验证码对象} TTL=300s
    └── gRPC getAuthPerms ──► Redis: GET {jti}（存在 key 不一致风险）

工作流任务调度
    │
    ▼
repair-workflow-service
    ├── 批量更新设备状态（task.go）
    │       └── Redis: SET batch_update_status_warehouse {仓库参数} TTL=永久
    └── 状态回调处理（callback_grpc_maintenance.go）
            └── Redis: GET batch_update_status_warehouse → DEL（用后即删）
```

---

## 五、问题与优化建议

### 问题一：repair-chassis-service 权限缓存 Key 不一致

- **现象：** 登录时以 `username` 为 key 缓存权限，而 gRPC 接口 `getAuthPerms` 以 `JTI` 为 key 查询权限
- **影响：** `GrpcAuthService.getAuthPerms` 返回空，导致其他服务的权限校验失效
- **建议：** 统一 key 命名策略，优先使用 `username`，或在 gRPC 接口中改为查询 `username` 对应的权限

### 问题二：repair-workflow-service 的跨节点缓存 key 固定

- **现象：** `batch_update_status_warehouse` key 是硬编码的固定字符串，无任务 ID 区分
- **风险：** 若同时存在多个批量任务并发执行，会产生 key 覆盖，导致仓库参数错乱
- **建议：** 将 key 改为 `batch_update_status_warehouse:{taskID}`，按任务隔离缓存

### 问题三：repair-workflow-service 代码注释中的 Sentinel 模式未启用

- **现象：** 配置字段支持 Sentinel（高可用），但代码实际使用单机模式
- **风险：** 单机 Redis 故障时，分布式锁和参数缓存全部失效，影响工作流正常运行
- **建议：** 生产环境启用 Sentinel 或 Cluster 模式，并与配置字段保持一致

### 问题四：repair-chassis-service Value 使用 JDK 序列化

- **现象：** `RedisConfig` 中 Value 使用 `JdkSerializationRedisSerializer`
- **风险：** JDK 序列化的数据不可读，且存在版本兼容性问题；对象字段变更后，旧缓存反序列化可能抛出异常
- **建议：** 改用 `GenericJackson2JsonRedisSerializer` 或 `Jackson2JsonRedisSerializer`，提升可读性和兼容性

### 问题五：repair-workflow-service 缓存设置无过期时间

- **现象：** `task.go` 中 `RedisSet(ctx, cacheKey, warehouseStr, 0)` 第四个参数为 `0`（永不过期）
- **风险：** 若读取端异常退出导致 `defer RedisDel` 未执行，key 会永久残留
- **建议：** 设置合理的超时时间（如 300 秒），作为兜底的自动清理机制
