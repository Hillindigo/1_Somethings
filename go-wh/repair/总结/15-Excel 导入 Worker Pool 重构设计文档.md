# Excel 导入 Worker Pool 重构设计文档

> **目标**：针对当前 Excel 导入功能存在的数据不一致、并发安全、错误丢失等问题，用 Worker Pool 思路重新设计一套健壮的导入方案。
>
> **阅读对象**：Go 初学者，用大白话 + 图解说明。

---

## 一、先搞清楚：现在的系统到底有什么问题？

在讲新方案之前，你需要知道老代码"烂"在哪里。我把核心问题按严重程度分成三类：

### 🔴 会丢数据 / 数据错乱（必须修）

| # | 问题 | 用大白话说 |
|---|------|-----------|
| 1 | **任务无重入防护** | 用户快速点两次"确认导入"，会启动两个 goroutine 同时处理同一个任务，导致资产被**重复插入**。就像两个人同时往同一个箱子里放苹果，最后苹果数是对的两倍。 |
| 2 | **主数据写成功、副作用写失败，没有补偿** | 资产 `BulkCreate` 成功了，但后续写变更历史、写 SN 绑定、写仓库记录如果失败了，只打了一行日志就跳过了。结果就是：资产存在，但历史记录缺失、SN 绑定断裂，查不到来龙去脉。 |
| 3 | **handler 失败后 item 变"僵尸"** | 一批数据被标记成 `processing` 后，如果处理失败了，这些数据永远不会再被拾取（因为下次只查 `pending`），也不会被标记 `failed`。这些数据就"消失"了。 |

### 🟠 数据不准确（应该修）

| # | 问题 | 用大白话说 |
|---|------|-----------|
| 4 | **导入计数不准** | `imported` 和 `failed` 两个计数器没有加锁，异步 goroutine 和主流程同时在改，就像两个人同时在同一张纸上记数，最终数字可能是错的。 |
| 5 | **confirmeds 列表有空值** | `make([]string, len(items))` 先创建了 N 个空字符串，再 `append` 新的，结果前 N 个都是空的。前端拿到一堆空字符串。 |
| 6 | **静默跳过不记录** | 有些行因为项目不存在被直接 `continue` 跳过了，但 `failedCount` 没加 1，用户也看不到错误信息。 |

### 🟡 系统健壮性差（建议修）

| # | 问题 | 用大白话说 |
|---|------|-----------|
| 7 | **用了 `context.Background()`** | 服务要关机了，但导入任务还在跑，数据写到一半被强制中断，造成半成品数据。 |
| 8 | **cancel() 没用 defer** | 如果中间 panic 了，`cancel()` 不会被调用，会导致内存泄漏。 |
| 9 | **taskChannel 阻塞** | channel 满了时，API 请求线程会被卡住，用户的 HTTP 请求一直等着。 |

---

## 二、新方案的总体设计思路

### 2.1 一句话总结

> 把导入任务看成一条**流水线**，每个阶段职责单一、错误可追溯、失败可恢复。

### 2.2 流水线全景图

先看整体，后面再逐个拆解：

```text
┌─────────────────────────────────────────────────────────────────────┐
│                        用户操作                                      │
│  ① 上传 Excel → ② 预览/确认 → ③ 查询进度 → ④ 查看结果              │
└──────┬──────────────┬──────────────┬──────────────┬─────────────────┘
       │              │              │              │
       ▼              ▼              │              │
  ┌─────────┐   ┌──────────┐        │              │
  │ 文件校验  │   │ 确认导入  │        │              │
  │ 阶段     │   │ 阶段     │        │              │
  └────┬─────┘   └────┬─────┘        │              │
       │              │              │              │
       ▼              ▼              ▼              ▼
  ┌───────────────────────────────────────────────────────────────────┐
  │                     MongoDB（数据库层）                             │
  │                                                                   │
  │  import_tasks 表          import_task_items 表     asset_info 表  │
  │  ┌──────────────┐        ┌───────────────────┐   ┌─────────────┐ │
  │  │ id           │        │ id                │   │ asset_code  │ │
  │  │ status       │◄───────│ task_id           │   │ factory_sn  │ │
  │  │ total_count  │        │ status            │   │ ...         │ │
  │  │ import_count │        │ type (new/dup)    │   └─────────────┘ │
  │  │ failed_count │        │ asset_info (嵌入) │                   │
  │  │ error_msgs   │        │ error_msg         │                   │
  │  │ ...          │        │ retry_count       │ ← 新增字段        │
  │  └──────────────┘        └───────────────────┘                   │
  └───────────────────────────────────────────────────────────────────┘
       │              │
       ▼              ▼
  ┌─────────┐   ┌───────────────────────────────────────────────┐
  │ 解析     │   │            导入引擎（核心重构区域）              │
  │ Worker   │   │                                               │
  │          │   │  ┌─────────────┐   ┌────────────────────┐    │
  │ 流式读取  │   │  │ 主写入循环   │   │  副作用 Worker Pool  │    │
  │ Excel    │   │  │ (单goroutine│──►│  (8个并发槽位)       │    │
  │ 逐行校验  │   │  │  串行处理)  │   │  写历史/写绑定/算完整度│   │
  │ 分类入库  │   │  └──────┬──────┘   └─────────┬──────────┘    │
  │          │   │         │                     │               │
  │          │   │         ▼                     ▼               │
  │          │   │  ┌─────────────┐   ┌────────────────────┐    │
  │          │   │  │ 补偿管理器   │   │  结果聚合器          │    │
  │          │   │  │ (失败重试)  │   │  (精确计数)          │    │
  │          │   │  └─────────────┘   └────────────────────┘    │
  └─────────┘   └───────────────────────────────────────────────┘
```

### 2.3 设计原则（记住这四条就够了）

| 原则 | 含义 | 举例 |
|------|------|------|
| **幂等性** | 同一个操作执行多次，结果和执行一次一样 | 同一个资产不管导入几次，数据库里只有一条 |
| **最终一致性** | 即使中间有失败，通过重试/补偿最终数据是完整的 | 变更历史写失败了，补偿任务会重新写 |
| **可观测性** | 每一步的结果都能被追踪和查询 | 用户能看到"第 35 行导入失败：原厂 SN 已存在" |
| **优雅降级** | 副作用失败不影响主数据写入 | SN 绑定写失败了，但资产本身已经成功入库 |

---

## 三、各阶段详细设计

### 阶段一：文件校验 & 解析（CheckFile）

#### 这个阶段做什么？

用户上传 Excel 后，系统需要：
1. 校验文件格式
2. 流式读取每一行
3. 校验每行数据是否合法
4. 把"新资产"和"重复资产"分开存到中间表

#### 数据流图

```text
用户上传 Excel
       │
       ▼
  ┌──────────────────────┐
  │ 1. 基础校验            │
  │    - 文件格式 xlsx/csv │
  │    - MinIO 文件存在    │
  │    - 用户权限          │
  └──────────┬───────────┘
             │
             ▼
  ┌──────────────────────┐
  │ 2. 创建 ImportTask    │
  │    status = pending   │
  │    存入 MongoDB       │
  └──────────┬───────────┘
             │
             ▼
  ┌──────────────────────────────────────────────────────┐
  │ 3. 异步解析（通过 taskChannel 投递）                    │
  │                                                      │
  │    MinIO ──流式读取──► 逐行解析 ──► 校验规则            │
  │                                      │               │
  │                              ┌───────┴────────┐      │
  │                              │                │      │
  │                              ▼                ▼      │
  │                        校验通过           校验失败     │
  │                         │                  │         │
  │                    ┌────┴────┐         记录错误       │
  │                    │         │         中断解析       │
  │                    ▼         ▼                       │
  │              查数据库判断                              │
  │              FactorySN 是否存在                       │
  │                    │                                 │
  │              ┌─────┴─────┐                           │
  │              │           │                           │
  │              ▼           ▼                           │
  │         不存在 →      已存在 →                        │
  │      type = "new"   type = "duplicate"              │
  │              │           │                           │
  │              └─────┬─────┘                           │
  │                    ▼                                 │
  │           批量写入 import_task_items                  │
  │           (每 1000 条一批)                            │
  │                    │                                 │
  │                    ▼                                 │
  │           更新 ImportTask                            │
  │           status = wait（等待用户确认）                │
  └──────────────────────────────────────────────────────┘
```

#### 关键改进点

**老问题**：`CreateImportTask` 用死循环等 channel 不满，会阻塞 API。

**新方案**：改用带超时的投递，超时直接返回错误：

```go
// 新版：带超时的任务投递
func (q *TaskQueue) Enqueue(ctx context.Context, task *mongo.ImportTask) error {
    select {
    case q.ch <- task:
        return nil
    case <-ctx.Done():
        return fmt.Errorf("投递任务超时，系统繁忙请稍后重试")
    case <-time.After(5 * time.Second):
        return fmt.Errorf("任务队列已满，请稍后重试")
    }
}
```

---

### 阶段二：确认导入（StartImport）—— 重构核心

这是改动最大的部分，我拆成多个子模块来讲。

#### 3.2.1 任务锁：防止重复启动

**老问题**：用户点两次"确认"，两个 goroutine 同时跑。

**新方案**：用 MongoDB 的 `findOneAndUpdate` 做**乐观锁**。

```text
什么是乐观锁？

就像你去食堂打菜，看到最后一个鸡腿，你伸手去拿：
- 如果你拿到了 → 成功，鸡腿是你的
- 如果被别人抢先拿了 → 失败，你知道没拿到

在代码里：
- "看到鸡腿" = 查到 status 是 wait
- "拿鸡腿"   = 把 status 改成 running
- 这两步合成一个原子操作，数据库保证只有一个人能成功
```

```go
// StartImport 新版：用原子操作抢锁
func (i *importWorkerService) StartImport(ctx context.Context, taskID string, isCover bool) error {
    if taskID == "" {
        return errors.New("任务ID不能为空")
    }

    // 原子操作：只有 status=wait 的任务才能被改成 running
    // 如果两个请求同时来，只有一个能成功修改，另一个会发现 status 已经不是 wait 了
    updated, err := i.taskService.FindOneAndUpdate(ctx,
        bson.M{
            "_id":    taskID,
            "status": mongo.TaskStatusWait, // 前置条件：必须是 wait 状态
        },
        bson.M{
            "$set": bson.M{
                "status":      mongo.TaskStatusRunning,
                "update_time": time.Now(),
            },
        },
    )
    if err != nil {
        return fmt.Errorf("启动导入失败: %w", err)
    }
    if updated == nil {
        return errors.New("任务不在可导入状态（可能已在运行中或已完成）")
    }

    // 拿到锁了，安全地启动 goroutine
    userInfo := common.GetAuthCtxUserInfoFromCtx(ctx)
    go func() {
        runner := newStartImportRunner(i, updated, userInfo)
        runner.Run()
    }()

    return nil
}
```

**效果**：无论用户点多少次，只有第一次能成功把 `wait → running`，后续请求都会被拦住。

#### 3.2.2 主写入循环（runPhase）：修复僵尸数据问题

**老问题**：handler 失败后，item 卡在 `processing` 状态永远不被处理。

**新方案**：失败时**把 item 状态回退到 `pending`** + 增加重试计数器。

```text
状态机对比：

老版本:
  pending → processing → completed
                      ↘ (失败了？没人管，卡死在 processing)

新版本:
  pending → processing → completed
      ↑         │
      │         ▼
      └──── failed_once (retry_count < 3，自动回退到 pending 重试)
                │
                ▼ (retry_count >= 3)
             failed (永久失败，记录错误原因)
```

```go
func (r *startImportRunner) runPhase(itemType string, handler batchHandler) error {
    for {
        // 1. 原子操作：查询并标记 processing（合二为一！）
        //    用 MongoDB 的 UpdateMany + 返回修改过的文档
        items, err := r.svc.taskService.ClaimPendingItems(
            r.ctx,             // 用可取消的 context，不再是 Background()
            r.task.ID.Hex(),
            itemType,
            r.limit,           // 每批 1000
        )
        if err != nil {
            return fmt.Errorf("领取待处理数据失败: %w", err)
        }
        if len(items) == 0 {
            return nil // 全部处理完毕
        }

        // 2. 执行业务处理
        result := handler(items)

        // 3. 根据结果更新状态
        if result.err != nil {
            // ★ 关键改进：失败时回退状态，而不是丢弃
            r.svc.taskService.ResetItemsToPending(r.ctx, result.failedItemIDs)
            r.svc.taskService.IncrementRetryCount(r.ctx, result.failedItemIDs)

            // 超过最大重试次数的标记为永久失败
            r.svc.taskService.MarkExceededAsFailed(r.ctx, result.failedItemIDs, 3)
        }

        if len(result.succeededItemIDs) > 0 {
            r.svc.taskService.MarkItemsCompleted(r.ctx, result.succeededItemIDs)
        }

        // 4. 精确更新进度
        r.updateProgress()
    }
}
```

#### 3.2.3 `ClaimPendingItems`：原子"领取"操作

```text
为什么要"原子领取"？

老代码分两步：
  第 1 步：查出 1000 条 pending 的数据
  第 2 步：把这 1000 条改成 processing

如果第 1 步和第 2 步之间出了什么事（比如服务崩了），
这 1000 条数据已经被查出来了但没改状态，
下次查询可能会再查到它们，导致重复处理。

新方案把两步合成一步：
  "找到 pending 的数据，同时改成 processing，返回给我"
  这是一个原子操作，要么全成功，要么全不动。
```

```go
// ClaimPendingItems 原子地"领取"一批待处理的 item
// 使用 MongoDB 的 aggregate + $merge 或 findAndModify 实现
func (s *importTaskStorage) ClaimPendingItems(
    ctx context.Context, taskID string, itemType string, limit int64,
) ([]*ImportTaskItem, error) {
    now := time.Now()
    filter := bson.M{
        "task_id":     taskID,
        "type":        itemType,
        "status":      TaskItemStatusPending,
        "retry_count": bson.M{"$lt": 3}, // 重试次数小于 3 次
    }
    update := bson.M{
        "$set": bson.M{
            "status":      TaskItemStatusProcessing,
            "claimed_at":  now,   // 记录被领取的时间（用于超时检测）
            "update_time": now,
        },
    }

    // 分批领取：循环 findOneAndUpdate，直到拿满 limit 条或没有更多
    var items []*ImportTaskItem
    for i := int64(0); i < limit; i++ {
        var item ImportTaskItem
        err := collection.FindOneAndUpdate(ctx, filter, update,
            options.FindOneAndUpdate().SetReturnDocument(options.After),
        ).Decode(&item)
        if err != nil {
            if err == mongodriver.ErrNoDocuments {
                break // 没有更多 pending 的了
            }
            return nil, err
        }
        items = append(items, &item)
    }

    return items, nil
}
```

> **小贴士**：`FindOneAndUpdate` 是 MongoDB 的原子操作，等价于"查询 + 修改"一步完成。每次只处理一条，循环 1000 次。虽然不如 `UpdateMany` 快，但**绝对安全**。如果你需要更快，后面的"性能优化"章节会介绍批量版本。

#### 3.2.4 processNewBatch 改进：确保副作用可追溯

**老问题**：副作用（写历史、写绑定）失败了，只打日志，无法恢复。

**新方案**：引入**副作用结果收集器**，失败的副作用记入数据库，由定时补偿任务重跑。

```text
什么是"副作用"？

主操作：把资产数据写入 asset_info 表（这是核心的，不能丢）
副作用：写变更历史、写 SN 绑定、算完整度（这些是附带的，但也很重要）

以前：副作用失败了，日志里记一行就完了，没人管
现在：副作用失败了，记到一个"待补偿"表里，有定时任务会重试
```

```go
func (r *startImportRunner) processNewBatch(items []*ImportTaskItem) *batchResult {
    result := &batchResult{
        succeededItemIDs: make([]primitive.ObjectID, 0),
        failedItemIDs:    make([]primitive.ObjectID, 0),
    }

    // 1. 批量生成资产编号
    codes, err := utils.GenerateBatchAssetCodes(r.ctx, r.prefix, len(items), r.assetStorage)
    if err != nil {
        // 生成编号失败，整批回退
        for _, it := range items {
            result.failedItemIDs = append(result.failedItemIDs, it.ID)
        }
        result.err = err
        return result
    }

    // 2. 组装资产对象
    toInsert := make([]*mongo.AssetInfo, 0, len(items))
    for idx, it := range items {
        asset := it.AssetInfo
        asset.AssetCode = codes[idx]
        asset.AssetType = r.task.AssetType
        asset.CreateTime = time.Now()
        asset.UpdateTime = time.Now()
        toInsert = append(toInsert, asset)
    }

    // 3. 批量写入（主操作）
    if err := r.assetStorage.BulkCreate(r.ctx, toInsert); err != nil {
        for _, it := range items {
            result.failedItemIDs = append(result.failedItemIDs, it.ID)
        }
        result.err = err
        return result
    }

    // 主操作成功，记录成功的 item
    for _, it := range items {
        result.succeededItemIDs = append(result.succeededItemIDs, it.ID)
    }
    // ★ 用原子操作更新计数，避免数据竞争
    atomic.AddInt64(&r.imported, int64(len(toInsert)))

    // 4. 副作用提交到 Worker Pool（带结果收集）
    inserted := toInsert
    r.dispatcher.Submit(func() {
        for _, asset := range inserted {
            if err := r.writeSideEffects(asset); err != nil {
                // ★ 失败不是打日志完事，而是记到补偿表
                r.compensator.Record(CompensationTask{
                    Type:      "side_effect",
                    AssetCode: asset.AssetCode,
                    TaskID:    r.task.ID.Hex(),
                    Error:     err.Error(),
                    CreatedAt: time.Now(),
                })
            }
        }
    })

    return result
}
```

#### 3.2.5 副作用 Worker Pool 改进：panic 安全 + 超时控制

```go
// 新版 dispatcher：增加 panic 恢复和超时
type sideEffectDispatcher struct {
    sem    chan struct{}
    wg     sync.WaitGroup
    ctx    context.Context    // ★ 新增：绑定可取消的 context
    cancel context.CancelFunc
}

func newSideEffectDispatcher(ctx context.Context, maxConcurrent int) *sideEffectDispatcher {
    if maxConcurrent <= 0 {
        maxConcurrent = 8
    }
    dCtx, cancel := context.WithCancel(ctx)
    return &sideEffectDispatcher{
        sem:    make(chan struct{}, maxConcurrent),
        ctx:    dCtx,
        cancel: cancel,
    }
}

func (d *sideEffectDispatcher) Submit(fn func()) {
    d.wg.Add(1)

    select {
    case d.sem <- struct{}{}:
        // 拿到槽位
    case <-d.ctx.Done():
        // 服务要关了，不再接受新任务
        d.wg.Done()
        return
    }

    go func() {
        defer func() {
            // ★ 捕获 panic，防止一个任务崩溃导致整个系统挂掉
            if r := recover(); r != nil {
                klog.Errorf("[dispatcher] panic recovered: %v", r)
            }
            <-d.sem
            d.wg.Done()
        }()
        fn()
    }()
}

// WaitWithTimeout 带超时的等待
func (d *sideEffectDispatcher) WaitWithTimeout(timeout time.Duration) error {
    done := make(chan struct{})
    go func() {
        d.wg.Wait()
        close(done)
    }()

    select {
    case <-done:
        return nil
    case <-time.After(timeout):
        d.cancel() // 超时了，取消所有还在跑的副作用
        return fmt.Errorf("等待副作用完成超时(%v)", timeout)
    }
}
```

#### 3.2.6 精确计数器：用 atomic 替代裸读写

```go
type progressCounter struct {
    imported int64
    failed   int64
    skipped  int64
}

// 线程安全的增减方法
func (p *progressCounter) AddImported(n int64) { atomic.AddInt64(&p.imported, n) }
func (p *progressCounter) AddFailed(n int64)   { atomic.AddInt64(&p.failed, n) }
func (p *progressCounter) AddSkipped(n int64)  { atomic.AddInt64(&p.skipped, n) }

func (p *progressCounter) Snapshot() (imported, failed, skipped int64) {
    return atomic.LoadInt64(&p.imported),
           atomic.LoadInt64(&p.failed),
           atomic.LoadInt64(&p.skipped)
}
```

```text
为什么要用 atomic？

想象一个计数器，变量值是 100。
两个 goroutine 同时执行 counter++：
  goroutine A 读到 100，准备写入 101
  goroutine B 也读到 100，准备写入 101
  结果：counter = 101（少加了一次！）

用 atomic.AddInt64：
  这是 CPU 级别的原子操作，"读 + 改 + 写"在一个时钟周期内完成
  不可能被打断，所以两个 goroutine 加完后一定是 102
```

---

### 阶段三：补偿机制（解决"最终一致性"）

#### 为什么需要补偿？

```text
真实场景：
  1. BulkCreate 写入了 1000 条资产 ✅
  2. 写变更历史时，gRPC 调用超时了 ❌
  3. 写 SN 绑定时成功了 ✅
  4. 算完整度时 MongoDB 连接断了 ❌

结果：资产数据完整，但变更历史和完整度数据缺失。

补偿机制的作用：
  → 把第 2 步和第 4 步的失败记录下来
  → 后台定时任务每隔 30 秒检查一次
  → 自动重试，直到成功或达到最大重试次数
```

#### 补偿表设计

```text
compensation_tasks 表结构：
┌────────────────┬──────────┬───────────────────────────────┐
│ 字段            │ 类型     │ 说明                          │
├────────────────┼──────────┼───────────────────────────────┤
│ _id            │ ObjectID │ 主键                          │
│ task_id        │ string   │ 关联的导入任务 ID              │
│ type           │ string   │ 补偿类型：                     │
│                │          │   change_history（变更历史）    │
│                │          │   sn_binding（SN绑定）         │
│                │          │   completeness（完整度）        │
│                │          │   warehouse（仓库记录）         │
│ asset_code     │ string   │ 关联的资产编号                  │
│ payload        │ bson.M   │ 重试所需的参数                  │
│ status         │ string   │ pending / completed / failed   │
│ retry_count    │ int      │ 已重试次数                      │
│ max_retry      │ int      │ 最大重试次数（默认 5）           │
│ error          │ string   │ 最后一次失败的错误信息           │
│ created_at     │ time     │ 创建时间                       │
│ next_retry_at  │ time     │ 下次重试时间（指数退避）         │
│ completed_at   │ time     │ 完成时间                       │
└────────────────┴──────────┴───────────────────────────────┘
```

#### 补偿执行器

```go
// CompensationRunner 定时扫描并执行补偿任务
func (c *CompensationRunner) Run(ctx context.Context) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            c.processPendingCompensations(ctx)
        }
    }
}

func (c *CompensationRunner) processPendingCompensations(ctx context.Context) {
    // 查询：status=pending 且 next_retry_at <= 当前时间 且 retry_count < max_retry
    tasks, _ := c.storage.FindPendingCompensations(ctx, time.Now(), 100)

    for _, task := range tasks {
        var err error
        switch task.Type {
        case "change_history":
            err = c.retryChangeHistory(ctx, task)
        case "sn_binding":
            err = c.retrySNBinding(ctx, task)
        case "completeness":
            err = c.retryCompleteness(ctx, task)
        }

        if err != nil {
            // 失败：增加重试次数，用指数退避算下次重试时间
            nextRetry := time.Now().Add(time.Duration(math.Pow(2, float64(task.RetryCount))) * time.Minute)
            c.storage.UpdateRetry(ctx, task.ID, task.RetryCount+1, nextRetry, err.Error())
        } else {
            // 成功：标记完成
            c.storage.MarkCompleted(ctx, task.ID)
        }
    }
}
```

```text
什么是"指数退避"？

第 1 次失败 → 2^1 = 2 分钟后重试
第 2 次失败 → 2^2 = 4 分钟后重试
第 3 次失败 → 2^3 = 8 分钟后重试
第 4 次失败 → 2^4 = 16 分钟后重试
第 5 次失败 → 达到最大重试次数，标记为永久失败，人工介入

好处：避免在下游服务故障时疯狂重试，给对方恢复的时间。
```

---

### 阶段四：优雅停机（Graceful Shutdown）

**老问题**：用 `context.Background()`，服务关停时导入任务写到一半。

**新方案**：

```text
服务收到关停信号（SIGTERM）
       │
       ▼
  ┌──────────────────────────────┐
  │ 1. 停止接受新的导入请求        │
  │    (taskChannel 关闭)         │
  │                              │
  │ 2. 等待正在运行的 runner 完成  │
  │    (最多等 2 分钟)            │
  │                              │
  │ 3. 如果 2 分钟内没完成：       │
  │    - 取消 runner 的 context    │
  │    - 正在 processing 的 item  │
  │      不用管，启动恢复机制处理   │
  │                              │
  │ 4. 等待 dispatcher 完成       │
  │    (最多再等 30 秒)           │
  │                              │
  │ 5. 安全退出                   │
  └──────────────────────────────┘
```

```go
// runner 使用可取消的 context
func newStartImportRunner(svc *importWorkerService, task *mongo.ImportTask,
    userInfo *common.AuthCtxUserInfo, shutdownCtx context.Context,
) *startImportRunner {
    runnerCtx, cancel := context.WithCancel(shutdownCtx)
    return &startImportRunner{
        ctx:        runnerCtx,     // ★ 不再是 context.Background()
        cancel:     cancel,
        dispatcher: newSideEffectDispatcher(runnerCtx, 8),
        // ...
    }
}
```

---

### 阶段五：启动恢复（Startup Recovery）

**问题场景**：服务被强制杀死（kill -9），有些 item 卡在 `processing` 状态。

```go
// RecoverStuckItems 在服务启动时调用
func RecoverStuckItems(ctx context.Context, taskService ImportTaskStorage) {
    // 查找：status=processing 且 claimed_at 在 10 分钟之前的 item
    // 说明这些 item 被"领走"后处理方已经不在了
    stuckItems, err := taskService.FindStuckProcessingItems(ctx, 10*time.Minute)
    if err != nil {
        klog.Errorf("[Recovery] find stuck items failed: %v", err)
        return
    }

    if len(stuckItems) == 0 {
        return
    }

    klog.Infof("[Recovery] found %d stuck processing items, resetting to pending", len(stuckItems))

    // 重置为 pending，增加 retry_count
    ids := make([]primitive.ObjectID, len(stuckItems))
    for i, item := range stuckItems {
        ids[i] = item.ID
    }

    taskService.ResetItemsToPending(ctx, ids)
    taskService.IncrementRetryCount(ctx, ids)

    // 同时检查对应的 task，如果 task 状态是 running 但实际没有 runner 在跑
    // 标记为 failed，让用户重新触发
    taskService.MarkOrphanedTasksAsFailed(ctx)
}
```

---

## 四、新版任务状态机（完整版）

```text
                          ┌─────────┐
                          │ pending │  用户刚上传，等待解析
                          └────┬────┘
                               │ taskChannel 消费
                               ▼
                          ┌──────────┐
                     ┌────│ parsing  │  正在解析 Excel
                     │    └────┬─────┘
                     │         │ 解析完成
               解析失败        ▼
                     │    ┌──────────┐
                     │    │  wait    │  等待用户确认
                     │    └────┬─────┘
                     │         │ 用户点击确认（原子锁）
                     │         ▼
                     │    ┌──────────┐
                     │    │ running  │  正在导入
                     │    └────┬─────┘
                     │         │
                     │    ┌────┴────────────┐
                     │    │                 │
                     ▼    ▼                 ▼
                ┌──────────┐          ┌───────────┐
                │  failed  │          │ completed │
                └──────────┘          └───────────┘

import_task_item 的状态机：

                ┌─────────┐
                │ pending │
                └────┬────┘
                     │ ClaimPendingItems（原子领取）
                     ▼
                ┌──────────────┐
                │ processing   │
                └────┬────┬────┘
                     │    │
                成功  │    │ 失败
                     │    │
                     ▼    ▼
              ┌──────────┐  ┌──────────┐
              │completed │  │ pending  │ ← 回退重试（retry_count++）
              └──────────┘  └────┬─────┘
                                 │ retry_count >= 3
                                 ▼
                           ┌──────────┐
                           │  failed  │  永久失败
                           └──────────┘
```

---

## 五、新老方案对比总结

| 维度 | 老方案 | 新方案 |
|------|--------|--------|
| **防重复启动** | 无防护，可重复触发 | `findOneAndUpdate` 原子锁 |
| **领取 item** | 查询 + 修改分两步（非原子） | `FindOneAndUpdate` 原子领取 |
| **失败处理** | item 卡死在 processing | 回退到 pending + 重试计数 |
| **副作用失败** | 打日志，数据丢失 | 记入补偿表，定时重试 |
| **计数器** | 裸读写，有数据竞争 | `atomic` 原子操作 |
| **停机** | `context.Background()`，无法中断 | 可取消 context + 优雅停机 |
| **恢复** | 无 | 启动时扫描卡死 item，自动恢复 |
| **panic** | dispatcher 里 panic 会 crash | `recover()` 捕获，不影响其他任务 |
| **超时** | `cancel()` 可能泄漏 | 统一 `defer cancel()` |
| **进度** | 不精确 | atomic 计数 + 每批更新 |

---

## 六、文件结构规划

```text
pkg/service/
├── import_worker.go                         # API 入口层（几乎不变）
├── import_worker_runner.go                  # ★ 新文件：runner 主逻辑
│   ├── startImportRunner 结构体
│   ├── Run()
│   ├── runPhase()
│   ├── processNewBatch()
│   └── processDuplicateBatch()
├── import_worker_dispatcher.go              # ★ 新文件：Worker Pool
│   ├── sideEffectDispatcher
│   ├── Submit()
│   └── WaitWithTimeout()
├── import_worker_compensator.go             # ★ 新文件：补偿机制
│   ├── CompensationRunner
│   ├── Record()
│   └── processPendingCompensations()
├── import_worker_recovery.go                # ★ 新文件：启动恢复
│   └── RecoverStuckItems()
├── import_worker_counter.go                 # ★ 新文件：线程安全计数器
│   └── progressCounter
├── import_worker_start_import_patterns.go   # 策略模式（保留，小改）
└── sync.go                                  # 文件解析（保留，小改 TaskQueue）
```

---

## 七、面试怎么讲这个方案

### 30 秒快速版

> 我用 Worker Pool + 批量插入重构了 Excel 导入功能。主写入用单 goroutine 串行 BulkCreate 保证数据安全，副作用（写历史、写绑定）交给 8 并发的 Worker Pool 异步处理。通过 MongoDB findOneAndUpdate 实现任务原子锁防重复启动，失败的 item 自动回退重试，副作用失败记入补偿表由定时任务重跑，确保最终一致性。

### 如果面试官追问

**Q：为什么不用消息队列（如 Kafka）来做？**

> 当前是单 Pod 处理，任务量不大（万级），用内存 channel + goroutine 足够了。引入 Kafka 会增加运维复杂度。如果未来需要多 Pod 水平扩展，可以将 taskChannel 替换为 Kafka topic，核心的 runner 逻辑不需要改。

**Q：补偿机制和事务有什么区别？**

> 事务要求"要么全成功要么全失败"（强一致性），但跨服务调用（gRPC 写历史、写绑定）无法用单个 MongoDB 事务覆盖。补偿机制追求的是"最终一致性"——允许中间状态短暂不一致，但通过重试机制最终达到一致。这是分布式系统的常用做法。

**Q：`findOneAndUpdate` 原子领取性能够用吗？**

> 每次领取一条，1000 条需要 1000 次 `findOneAndUpdate`。在 MongoDB 本地网络下单次约 1ms，1000 条约 1 秒。如果嫌慢，可以改用 `bulkWrite` + 版本号乐观锁的方式批量领取。但在当前场景（万级数据、非实时要求），1 秒完全可以接受。
