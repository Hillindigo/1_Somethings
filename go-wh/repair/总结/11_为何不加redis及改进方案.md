# 为何 repair-maintenance-service 没有使用 Redis 及改进方案

## 一、根因：主动的设计选择，而非技术限制

从代码看，**不是技术限制，而是主动的设计选择**。错误文件 `pkg/dax/errors.go` 中预定义了 `ErrRedis` 和 `ErrAnalyzeRedisData`，说明**曾规划引入 Redis 但未实施**。

```go
// pkg/dax/errors.go
ErrRedis            = errors.New("redis连接出错，请联系管理员")
ErrAnalyzeRedisData = errors.New("redis数据解析失败")
```

---

## 二、当前数据库查询的实际情况

该服务有**两层数据库**：

- **MySQL（GORM）** — 存储字典表（`dict`）、项目表、用户权限等结构化数据
- **MongoDB** — 存储资产信息（`asset_info`）、仓库资产视图（`warehouse_asset_view`）、出入库记录（`ware_asset_inout`）等大量业务数据

### 高频字典查询（最典型的痛点）

`getChildDictsByParentDictCode` 在 **5 个 service 文件中被调用 32 次**，且每次业务请求都触发多次。以 `GetWarehouseAssets` 为例，单次调用就触发 **4 次独立的字典父级查询**：

```go
// pkg/service/warehouse_asset_view.go
typeCache, err := getChildDictsByParentDictCode(ctx, "EQUIPMENT:TYPE")
hostTypeCache, err := getChildDictsByParentDictCode(ctx, "HOST:TYPE")
statusCache, err := getChildDictsByParentDictCode(ctx, "ASSET:STATUS")
stockStatusCache, err := getChildDictsByParentDictCode(ctx, "WAREHOUSE:STATUS")
```

每次 `getChildDictsByParentDictCode` 内部实际执行 **2 次 MySQL 查询**：

```go
// pkg/api/dict.go
func GetChildDictsByParentDictCode(ctx context.Context, parentDictCode string) (...) {
    dao := dax.NewDict(base.GetDB())
    dict, err := dao.GetDetailByDictCode(ctx, parentDictCode)        // 第1次：按 code 查父节点
    childDicts, err := dao.GetDetailByParentDictId(ctx, dict.DictID) // 第2次：按父ID查子节点
```

还有大量**循环内逐条查询**，例如：

```go
// pkg/service/warehouse_asset_view.go
for typeCode := range typeSet {
    dict, err := dictDao.GetDetailByDictCode(ctx, typeCode)  // N 次查询，N = 资产类型数
```

**结论：单次资产列表请求最多触发 8~16 次 MySQL 字典查询，而字典数据几乎从不变更。**

---

## 三、利弊分析：引入 Redis 是否有价值

### 有利的方面

#### 利一：字典数据是天然的缓存场景（收益最大）

- **字典数据几乎不变**：资产类型、库存状态、操作类型等字典数据由管理员维护，更新频率极低（几周甚至几个月一次）
- **查询频率极高**：每次资产列表接口都要查 4~8 次字典 MySQL
- **数据量小**：单个父节点下的子字典通常 5~20 条，整个字典表估计不超过 500 行
- **预期收益**：字典缓存命中后，每次请求可节省 **8~16 次 MySQL 往返**，响应时间可降低 **30%~60%**

```go
// 理想的缓存结构（示意）
key: "repair-maintenance:dict:children:EQUIPMENT:TYPE"  TTL: 30分钟
value: JSON([{DictCode:"...", DictName:"..."}, ...])
```

#### 利二：仓库列表/项目列表属于低变更高读取数据

`GetWarehouseAssets` 中多次调用：
- `s.warehouseStorage.FindByWarehouseIDs(ctx, warehouseIDs)` — 仓库信息基本不变
- `s.projectStorage.List(ctx, projectFilter)` — 项目信息变更频率低

这类数据也适合短 TTL 缓存（5~10 分钟）。

#### 利三：防止 N+1 查询中的重复 DB 访问

`GetWarehouseAssetInfoByAssetCode` 中资产详情页对同一字典 code 发起多次独立查询（类型、状态、库存状态），Redis 可天然做请求内去重。

#### 利四：预留的错误常量说明架构层面已规划

错误常量已定义，引入不是跨越式变更，风险可控。

---

### 不利的方面

#### 弊一：增加运维复杂度（当前最主要顾虑）

该服务已使用：**MySQL + MongoDB + Minio + gRPC**，再加 Redis 是第 5 个外部依赖。对于中小规模团队，每增加一个中间件就增加：

- 连接管理、健康检查、断路器逻辑
- 运维监控告警
- Redis 故障时的降级处理（必须写）

#### 弊二：缓存一致性维护成本

字典数据虽然变更少，但一旦管理员修改字典，**必须主动失效缓存**，否则出现资产类型显示旧名称、状态名称翻译错误等问题。需要在 `CreateDict`、`UpdateDict`、`DelDictByDictId` 这些写操作后加缓存失效逻辑，改动点分散。

#### 弊三：MongoDB 查询本身已有连接池，延迟可接受

```go
// pkg/base/base.go
opts := options.Client().
    SetMaxPoolSize(100).
    SetMinPoolSize(10).
    SetMaxConnIdleTime(100 * time.Second)
```

MongoDB 连接池配置合理（10~100 连接），对于非超高并发场景，直查 MongoDB 的延迟通常在 **1~5ms**，体验上不明显。

#### 弊四：数据一致性要求高的场景不适合缓存

`BatchUpdateWarehouseAssetStatus`、`UpdateWarehouseAssetStatus` 这类写操作后，如果缓存了资产状态，需要精准失效对应 key，否则前端看到的库存状态是脏数据，这在仓库管理场景中是**不可接受的业务错误**。

#### 弊五：当前没有高并发压力证据

这是一个内部管理系统（仓库管理、资产管理），并发用户数通常有限，**当前性能瓶颈未经过压测验证**，引入 Redis 可能是过度设计。

---

## 四、结论：分场景引入

### 缓存适用性汇总

| 数据类型 | 是否适合缓存 | 建议 TTL | 原因 |
|---|---|---|---|
| 字典子列表（`EQUIPMENT:TYPE` 等） | **强烈推荐** | 30~60 分钟 | 几乎不变，查询超频繁 |
| 仓库基本信息列表 | 推荐 | 10 分钟 | 变更少，多处重复查询 |
| 项目名称映射 | 推荐 | 10 分钟 | 变更少，多处重复查询 |
| 单条资产详情 | 不推荐 | - | 写操作频繁，失效复杂 |
| 出入库记录 | 不推荐 | - | 强一致性要求，数据量大 |
| 用户权限仓库关联 | 谨慎 | 5 分钟 | 变更时必须精准失效 |

**最优先切入点：字典缓存**（改动 1 个函数，全服务所有查询接口自动受益）。

---

## 五、字典缓存改动方案

### 5.1 整体思路

在 `getChildDictsByParentDictCode` 函数加 Redis 二级缓存：

```
请求来了
  ↓
查 Redis
  ↓命中→ 直接返回（跳过 2 次 MySQL）
  ↓未命中
查 MySQL（2次：父节点 + 子列表）
  ↓
写入 Redis（TTL=30min）
  ↓
返回结果

字典写操作（增/改/删）
  ↓
DEL 对应缓存 key
```

**关键设计原则：Redis 不可用时自动降级到直查 MySQL，不影响业务正确性，只是失去缓存加速效果。**

---

### 5.2 第一步：新增 Redis 配置字段

在 `pkg/conf/config.go` 中新增：

```go
RedisSentinelNodes string `default:"localhost:26379" envconfig:"REDIS_SENTINEL_NODES"`
RedisMaster        string `default:"mymaster"        envconfig:"REDIS_MASTER"`
RedisPassword      string `default:""                envconfig:"REDIS_PASSWORD"`
RedisDatabase      int    `default:"0"               envconfig:"REDIS_DATABASE"`
```

---

### 5.3 第二步：新建 `pkg/base/redis.go`

```go
package base

import (
    "context"
    "time"

    "github.com/go-redis/redis/v8"
    "k8s.io/klog/v2"
)

var redisClient *redis.Client

// InitRedis 初始化 Redis 单机客户端
func InitRedis(addr, password string, db int) error {
    rdb := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })
    if err := rdb.Ping(context.Background()).Err(); err != nil {
        return err
    }
    redisClient = rdb
    klog.Info("redis connected")
    return nil
}

// RedisGet 获取字符串值，key 不存在时返回 ("", nil)
func RedisGet(ctx context.Context, key string) (string, error) {
    val, err := redisClient.Get(ctx, key).Result()
    if err == redis.Nil {
        return "", nil
    }
    return val, err
}

// RedisSet 设置字符串值，expire=0 表示永不过期
func RedisSet(ctx context.Context, key, value string, expire time.Duration) error {
    return redisClient.Set(ctx, key, value, expire).Err()
}

// RedisDel 删除一个或多个 key
func RedisDel(ctx context.Context, keys ...string) error {
    return redisClient.Del(ctx, keys...).Err()
}
```

---

### 5.4 第三步：在 `pkg/base/base.go` 的 `Init` 中加入 Redis 初始化

在 `initChannel()` 之前添加（Redis 失败只打警告，不 panic）：

```go
// 初始化 Redis（可选，失败只打警告不 panic，字典查询降级到 MySQL）
redisAddr := strings.Split(c.RedisSentinelNodes, ",")[0]
if err := InitRedis(redisAddr, c.RedisPassword, c.RedisDatabase); err != nil {
    klog.Warningf("init redis failed, dict cache disabled: %v", err)
}
```

---

### 5.5 第四步：新建字典缓存层 `pkg/base/dict_cache.go`（核心）

```go
package base

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    dax "git.sy.com/galaxies/repair-maintenance-service/pkg/dax"
    orm "git.sy.com/galaxies/repair-maintenance-service/pkg/dax/gen/model"
    "k8s.io/klog/v2"
)

const dictCacheTTL = 30 * time.Minute

// dictCacheKey 生成字典子列表的缓存 key
// 格式: repair-maintenance:dict:children:{parentDictCode}
func dictCacheKey(parentDictCode string) string {
    return fmt.Sprintf("repair-maintenance:dict:children:%s", parentDictCode)
}

// dictCodeCacheKey 生成单条字典的缓存 key
// 格式: repair-maintenance:dict:code:{dictCode}
func dictCodeCacheKey(dictCode string) string {
    return fmt.Sprintf("repair-maintenance:dict:code:%s", dictCode)
}

// GetChildDictsCached 带 Redis 缓存的子字典查询
// 优先读 Redis，未命中则查 MySQL 并回写缓存
// Redis 不可用时自动降级到直查 MySQL，不影响业务
func GetChildDictsCached(ctx context.Context, parentDictCode string) ([]*orm.Dict, error) {
    cacheKey := dictCacheKey(parentDictCode)

    // 1. 尝试从 Redis 读取
    if redisClient != nil {
        cached, err := RedisGet(ctx, cacheKey)
        if err == nil && cached != "" {
            var dicts []*orm.Dict
            if jsonErr := json.Unmarshal([]byte(cached), &dicts); jsonErr == nil {
                return dicts, nil
            }
            klog.Warningf("dict cache unmarshal failed for [%s]: %v", parentDictCode, jsonErr)
        }
    }

    // 2. 缓存未命中，查 MySQL（两次查询：先查父节点ID，再查子节点列表）
    dao := dax.NewDict(GetDB())
    parent, err := dao.GetDetailByDictCode(ctx, parentDictCode)
    if err != nil {
        return nil, err
    }
    children, err := dao.GetDetailByParentDictId(ctx, parent.DictID)
    if err != nil {
        return nil, err
    }

    // 3. 回写 Redis 缓存
    if redisClient != nil {
        if data, jsonErr := json.Marshal(children); jsonErr == nil {
            if setErr := RedisSet(ctx, cacheKey, string(data), dictCacheTTL); setErr != nil {
                klog.Warningf("dict cache set failed for [%s]: %v", parentDictCode, setErr)
            }
        }
    }

    return children, nil
}

// GetDictByCodeCached 带 Redis 缓存的单条字典查询（按 dictCode）
func GetDictByCodeCached(ctx context.Context, dictCode string) (*orm.Dict, error) {
    cacheKey := dictCodeCacheKey(dictCode)

    // 1. 尝试从 Redis 读取
    if redisClient != nil {
        cached, err := RedisGet(ctx, cacheKey)
        if err == nil && cached != "" {
            var dict orm.Dict
            if jsonErr := json.Unmarshal([]byte(cached), &dict); jsonErr == nil {
                return &dict, nil
            }
        }
    }

    // 2. 查 MySQL
    dao := dax.NewDict(GetDB())
    dict, err := dao.GetDetailByDictCode(ctx, dictCode)
    if err != nil {
        return nil, err
    }

    // 3. 回写 Redis 缓存
    if redisClient != nil {
        if data, jsonErr := json.Marshal(dict); jsonErr == nil {
            _ = RedisSet(ctx, cacheKey, string(data), dictCacheTTL)
        }
    }

    return dict, nil
}

// InvalidateDictCache 字典数据变更后调用此函数失效相关缓存
// 在 CreateDict、UpdateDict、DelDictByDictId 操作后调用
// dictCode: 变更的字典自身 code
// parentDictCode: 父字典 code（用于失效子列表缓存），不知道时传空字符串
func InvalidateDictCache(ctx context.Context, dictCode, parentDictCode string) {
    if redisClient == nil {
        return
    }
    keys := make([]string, 0, 2)
    if dictCode != "" {
        keys = append(keys, dictCodeCacheKey(dictCode))
    }
    if parentDictCode != "" {
        keys = append(keys, dictCacheKey(parentDictCode))
    }
    if len(keys) > 0 {
        if err := RedisDel(ctx, keys...); err != nil {
            klog.Warningf("dict cache invalidate failed: %v", err)
        }
    }
}
```

---

### 5.6 第五步：service 层替换调用

**原来的调用方式（`warehouse_asset_view.go` 中）：**

```go
// 原代码：每次都查 MySQL，两次 IO
typeCache, err := getChildDictsByParentDictCode(ctx, "EQUIPMENT:TYPE")

// 循环内单条查询：每次都查 MySQL
dict, err := dictDao.GetDetailByDictCode(ctx, dictCode)
```

**改为：**

```go
// 带缓存的子列表查询（命中缓存时跳过2次MySQL）
typeCache, err := base.GetChildDictsCached(ctx, "EQUIPMENT:TYPE")

// 带缓存的单条查询（命中缓存时跳过1次MySQL）
dict, err := base.GetDictByCodeCached(ctx, dictCode)
```

> **注意**：`getChildDictsByParentDictCode` 原返回 proto 类型（`*dictv1.DictTreeDetailResponse`），新的 `GetChildDictsCached` 返回 orm 类型（`*orm.Dict`），需在调用处做类型转换，转换逻辑与 `pkg/api/dict.go` 第 387~408 行相同，可提取为公共工具函数复用。

---

### 5.7 第六步：字典写操作后失效缓存

在 `pkg/api/dict.go` 的三个写操作中加缓存失效。

**`CreateDict`（创建完成后）：**

```go
_, err := s.dictDao.CreateDict(ctx, dict)
if err != nil {
    return nil, status.Errorf(codes.Internal, err.Error())
}
// 失效父字典的子列表缓存
parentDict, _ := s.dictDao.GetDetailByDictId(ctx, req.ParentDictID)
parentCode := ""
if parentDict != nil {
    parentCode = parentDict.DictCode
}
base.InvalidateDictCache(ctx, req.DictCode, parentCode)
```

**`UpdateDict`（更新完成后）：**

```go
_, err := s.dictDao.UpdateDict(ctx, req.DictId, dict)
if err != nil {
    return nil, status.Errorf(codes.Internal, err.Error())
}
// 失效自身和父字典子列表缓存
parentDict, _ := s.dictDao.GetDetailByDictId(ctx, req.ParentDictId)
parentCode := ""
if parentDict != nil {
    parentCode = parentDict.DictCode
}
base.InvalidateDictCache(ctx, dict.DictCode, parentCode)
```

**`DelDictByDictId`（删除完成后）：**

```go
_, err := s.dictDao.DelDictByDictId(ctx, req.DictId, req.Disable)
if err != nil {
    return nil, status.Errorf(codes.Internal, err.Error())
}
// 失效自身缓存
base.InvalidateDictCache(ctx, dict.DictCode, "")
```

---

## 六、改动范围汇总

| 文件 | 改动类型 | 说明 |
|---|---|---|
| `pkg/conf/config.go` | 新增 4 个字段 | Redis 连接配置 |
| `pkg/base/redis.go` | 新建文件 | Redis 客户端基础封装 |
| `pkg/base/dict_cache.go` | 新建文件 | 字典缓存逻辑（核心） |
| `pkg/base/base.go` | 新增 ~5 行 | Init 中加 Redis 初始化 |
| `pkg/api/dict.go` | 新增 ~15 行 | 3 个写操作后加缓存失效 |
| `pkg/service/warehouse_asset_view.go` | 替换调用 ~30 处 | 字典查询改走缓存层 |
| `pkg/service/component_info.go` | 替换调用 ~9 处 | 同上 |
| `pkg/service/host_info.go` | 替换调用 ~5 处 | 同上 |
| `pkg/service/repair_record.go` | 替换调用 ~2 处 | 同上 |

**预期收益：资产列表、出入库查询等核心接口的 MySQL 压力减少 60%~80%，响应时间节省 20~50ms。**
