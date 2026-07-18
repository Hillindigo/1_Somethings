# MinIO 流式读取 Excel/CSV 深度解析

> **对应场景**：Excel 导入模块中，用户上传文件到 MinIO 后，服务端需要读取并解析文件内容。
>
> **对应源码文件**：
>
> - `pkg/minio/client.go` — MinIO 客户端初始化
> - `pkg/minio/minio.go` — 流式读取核心实现
> - `pkg/service/sync.go` — 消费端（CheckFile goroutine）

---

## 一、为什么需要流式读取

### 问题背景

如果用普通方式读取 Excel：

```go
// 普通方式：全量加载到内存
data, _ := io.ReadAll(minioObject) // 先把整个文件下载到内存
f, _ := excelize.OpenReader(bytes.NewReader(data)) // 再全量解析
rows, _ := f.GetRows("Sheet1") // 全部行全部加载到内存
```

对于一个 10 万行的 Excel 文件（约 50MB），这种方式会：

- 在内存中同时持有：**原始字节（50MB）+ 解析后的结构体（可能 200MB+）**
- 导入过程中内存峰值极高，容易 OOM
- 必须等文件全部下载完才能开始处理第一行

### 流式读取的优势

```text
普通方式：
  下载完整文件(50MB) → 全量解析 → 全量写入DB
  内存峰值: ~250MB，延迟: 下载时间 + 解析时间 + 写入时间（串行）

流式方式：
  下载一块 → 解析一行 → 写入channel → 消费端处理 → 下载下一块
  内存峰值: ~数MB（只有当前处理的行），延迟: 更低（下载与处理并行）
```

---

## 二、MinIO SDK 的 GetObject 返回的是什么

核心在于 `minio.Object` 实现了 `io.Reader` 接口：

```go
// pkg/minio/minio.go:69-75
func (m *minioClient) GetObject(ctx context.Context, fileName string) (*minio.Object, error) {
    object, err := m.MinioCli.MinioCli.GetObject(
        ctx,
        m.MinioCli.BucketName,
        fileName,
        minio.GetObjectOptions{},
    )
    if err != nil {
        return nil, err
    }
    return object, nil
}
```

`minio.Object` 的关键特性：

- **懒加载**：`GetObject` 调用本身**不会立即下载文件**，只是建立了一个 HTTP 连接准备
- **实现 `io.Reader`**：调用 `Read(buf)` 时才真正从 MinIO 服务器按需拉取数据块
- **实现 `io.Seeker`**：支持 Seek 操作（本项目未用到）
- **实现 `io.Closer`**：用完必须 `Close()` 关闭底层连接

正是因为 `minio.Object` 实现了 `io.Reader`，才可以直接把它传给任何接受 `io.Reader` 的库（`excelize`、`csv.Reader`、`bufio.Scanner`）。

---

## 三、三种流式读取实现详解

### 3.1 Excel（xlsx）流式读取：StreamObjectLinesByExcelize

```go
// pkg/minio/minio.go:112-163
func (m *minioClient) StreamObjectLinesByExcelize(
    ctx context.Context,
    fileName string,
) (<-chan []string, <-chan error, error) {

    // 1. 获取 MinIO 对象（懒加载，此时不下载）
    obj, err := m.GetObject(ctx, fileName)
    if err != nil {
        return nil, nil, err
    }

    // 2. 创建带缓冲的 channel
    rows := make(chan []string, 256) // 缓冲 256 行，生产者超前消费者最多 256 行
    errs := make(chan error, 1)      // 容量 1，只需传递一个错误

    // 3. 启动独立 goroutine 异步读取
    go func() {
        defer close(rows)    // goroutine 退出时关闭 channel，通知消费者结束
        defer close(errs)
        defer obj.Close()    // 确保释放 MinIO HTTP 连接

        // 4. excelize.OpenReader 接受 io.Reader，边下载边解析 xlsx 结构
        f, e := excelize.OpenReader(obj)
        if e != nil {
            errs <- e
            return
        }
        defer func() { _ = f.Close() }()

        // 5. 获取第一个 Sheet
        sheets := f.GetSheetList()
        if len(sheets) == 0 {
            errs <- io.ErrUnexpectedEOF
            return
        }

        // 6. 使用 Rows 迭代器（流式，不一次性加载所有行）
        rs, e := f.Rows(sheets[0])
        if e != nil {
            errs <- e
            return
        }
        defer func() { _ = rs.Close() }()

        // 7. 逐行迭代，每行推入 channel
        for rs.Next() {
            record, e := rs.Columns()
            if e != nil {
                errs <- e
                return
            }
            select {
            case rows <- record:  // 推入 channel（若 channel 满则阻塞等消费者）
            case <-ctx.Done():    // 支持 context 取消（用户取消导入时可立即停止）
                errs <- ctx.Err()
                return
            }
        }
    }()

    return rows, errs, nil // 立即返回 channel，调用方异步消费
}
```

**关键点：`excelize.Rows()` 的流式迭代器**

`excelize` 库的 `f.Rows()` 返回的是一个行迭代器，内部使用 XML 流式解析（SAX 风格），不会一次性把所有行加载到内存。每次调用 `rs.Next()` 时，才从底层 `io.Reader`（即 `minio.Object`）读取下一行的 XML 数据并解析。

### 3.2 CSV 流式读取：StreamObjectLinesByCsv

```go
// pkg/minio/minio.go:165-199
func (m *minioClient) StreamObjectLinesByCsv(
    ctx context.Context,
    fileName string,
) (<-chan []string, <-chan error, error) {

    obj, err := m.GetObject(ctx, fileName)
    if err != nil {
        return nil, nil, err
    }

    rows := make(chan []string, 256)
    errs := make(chan error, 1)

    go func() {
        defer close(rows)
        defer close(errs)
        defer obj.Close()

        // csv.NewReader 接受 io.Reader，天然支持流式读取
        r := csv.NewReader(obj)
        r.FieldsPerRecord = -1  // -1 表示允许每行字段数不同（宽松模式）

        for {
            record, e := r.Read()  // 每次只读一行，按需从 MinIO 拉取数据
            if e == io.EOF {
                return
            }
            if e != nil {
                errs <- e
                return
            }
            select {
            case rows <- record:
            case <-ctx.Done():
                errs <- ctx.Err()
                return
            }
        }
    }()

    return rows, errs, nil
}
```

CSV 的流式读取更简单，因为 `encoding/csv` 标准库的 `csv.NewReader` 本身就是逐行读取的，每次 `r.Read()` 只从底层 `io.Reader` 读取一行数据。

### 3.3 纯文本流式读取：StreamObjectLines（辅助方法）

```go
// pkg/minio/minio.go:77-110
func (m *minioClient) StreamObjectLines(ctx context.Context, fileName string) (<-chan string, <-chan error, error) {
    obj, err := m.GetObject(ctx, fileName)

    lines := make(chan string, 256)
    errs := make(chan error, 1)

    go func() {
        defer close(lines)
        defer close(errs)
        defer obj.Close()

        // bufio.Scanner 按行读取，每次只从底层 Reader 读取一小块
        scanner := bufio.NewScanner(obj)
        buf := make([]byte, 512*1024)         // 512KB 初始缓冲
        scanner.Buffer(buf, 100*1024*1024)    // 最大支持 100MB 的单行

        for scanner.Scan() {
            select {
            case lines <- scanner.Text():
            case <-ctx.Done():
                errs <- ctx.Err()
                return
            }
        }
        if scanErr := scanner.Err(); scanErr != nil {
            errs <- scanErr
        }
    }()

    return lines, errs, nil
}
```

---

## 四、消费端：sync.go 的 CheckFile goroutine

流式读取产生的 `linesCh` 被 `CheckFile` goroutine 消费，形成**生产者-消费者**模型：

```go
// pkg/service/sync.go:319-515（核心逻辑简化版）
func CheckFile(ctx context.Context) {
    for {
        select {
        case task := <-taskChannel:  // 有新任务进来
            func() {
                streamCtx, streamCancel := context.WithCancel(ctx)
                defer streamCancel()

                // 1. 根据文件类型选择对应的流式读取方法
                var linesCh <-chan []string
                var errCh  <-chan error
                switch task.FileType {
                case "xlsx":
                    linesCh, errCh, _ = minioService.StreamObjectLinesByExcelize(streamCtx, task.MinioFileName)
                case "csv":
                    linesCh, errCh, _ = minioService.StreamObjectLinesByCsv(streamCtx, task.MinioFileName)
                }

                rowIndex := 0
                var batch []*mongo.AssetInfo  // 积累到 1000 条再批量写入 DB

            streamLoop:
                for linesCh != nil || errCh != nil {
                    select {
                    // 2. 第一行：解析表头
                    case line, ok := <-linesCh:
                        if !ok {
                            linesCh = nil
                            continue
                        }
                        if rowIndex == 0 {
                            headerRows, headerMap, errList = CheckHeaders(streamCtx, task, line)
                            rowIndex++
                            continue
                        }

                        // 3. 数据行：转换 + 校验
                        assetInfo, _ := TransformLine(task, line, headerRows, headerMap, rowIndex, cache)
                        CheckRules(ctx, task, assetInfo, headerMap, cache)

                        batch = append(batch, assetInfo.AssetInfo)

                        // 4. 积累到 BatchSize(1000) 批量写入 import_task_items
                        if len(batch) >= BatchSize {
                            UpdateTaskData(ctx, task, batch)
                            batch = batch[:0]  // 清空，复用切片
                        }
                        rowIndex++

                    // 5. 处理流读取错误
                    case e, ok := <-errCh:
                        if !ok { errCh = nil; continue }
                        if e != nil {
                            streamCancel()  // 取消 MinIO 下载
                            break streamLoop
                        }

                    // 6. 支持 context 取消
                    case <-ctx.Done():
                        return
                    }
                }

                // 7. 处理最后不足 1000 条的剩余数据
                if len(batch) > 0 {
                    UpdateTaskData(ctx, task, batch)
                }
            }()
        }
    }
}
```

---

## 五、整体数据流图

```text
MinIO Server
    │
    │  HTTP/S（按需拉取，分块传输）
    ▼
minio.Object（实现 io.Reader）
    │
    │  每次 Read() 从网络读取一小块数据
    ▼
excelize.OpenReader(obj) / csv.NewReader(obj) / bufio.Scanner(obj)
    │
    │  逐行解析，每解析一行
    ▼
rows channel（缓冲 256 行）
    │                           ← goroutine 边界（生产者/消费者解耦）
    ▼
CheckFile goroutine 消费
    │  TransformLine()  — 字段映射、格式转换
    │  CheckRules()     — 业务校验、查重
    ▼
batch（内存中积累，最多 1000 条）
    │
    ▼
UpdateTaskData()  — 批量写入 import_task_items（MongoDB）
```

---

## 六、关键设计细节

### 6.1 双 channel 模式：数据 + 错误分离

```go
rows := make(chan []string, 256)  // 正常数据
errs := make(chan error, 1)       // 错误信号
```

为什么不用单个 channel 传递带错误的结构体？

- 错误是**稀有事件**，单独 channel 避免每次正常数据都要判断 error 字段
- `errs` 容量为 1：goroutine 只会发送一个致命错误就退出，不需要更大容量
- 消费端用 `select` 同时监听两个 channel，任何一个有数据都能及时响应

消费端的正确关闭检测写法：

```go
// linesCh 和 errCh 都关闭后才退出循环
for linesCh != nil || errCh != nil {
    select {
    case line, ok := <-linesCh:
        if !ok {
            linesCh = nil  // channel 已关闭，置 nil（nil channel 的 select 永远不会被选中）
            continue
        }
        // 处理数据...
    case e, ok := <-errCh:
        if !ok {
            errCh = nil
            continue
        }
        // 处理错误...
    }
}
```

**为什么置 `nil` 而不是 `break`**：`nil` channel 在 `select` 中永远不会被选中，相当于从 `select` 候选列表中移除这个 case，避免已关闭的 channel 被反复选中返回零值。

### 6.2 rows channel 缓冲大小 256 的含义

```go
rows := make(chan []string, 256)
```

- **缓冲 256 行**：生产者（MinIO 下载+解析）最多可以超前消费者（业务处理）256 行
- 如果消费者处理速度 < 生产者解析速度，channel 满后生产者自动阻塞，实现**背压（Backpressure）**
- 如果消费者处理速度 > 生产者（网络慢），消费者在 `select` 中等待，不空转
- 256 是经验值：既能平滑生产/消费速度差异，又不会占用过多内存（256 行 × 每行约 1KB ≈ 256KB）

### 6.3 context 取消的传播链

```go
streamCtx, streamCancel := context.WithCancel(ctx)
defer streamCancel()
```

当任意一个环节出错时，调用 `streamCancel()`：

```text
streamCancel() 被调用
    ↓
streamCtx.Done() channel 关闭
    ↓
生产者 goroutine 的 select 中 <-ctx.Done() 触发
    ↓
生产者 goroutine 退出，defer obj.Close() 执行
    ↓
MinIO HTTP 连接关闭，停止数据下载
```

这确保了即使消费者提前退出（如表头校验失败），也不会有 goroutine 泄露，MinIO 连接也会被及时释放。

### 6.4 excelize 的流式 vs 非流式对比

| 方式 | API | 内存行为 | 适用场景 |
| --- | --- | --- | --- |
| 非流式 | `f.GetRows("Sheet1")` | 一次性加载所有行到 `[][]string` | 文件小、需要随机访问行 |
| 流式 | `f.Rows("Sheet1")` + `rs.Next()` | 每次只在内存中保存当前行 | 文件大、顺序处理、内存敏感 |

项目使用 `f.Rows()` 流式迭代器，是为了避免大文件的内存峰值问题。

---

## 七、完整流程时序图

```text
用户线程          CheckFile goroutine      MinIO goroutine（生产者）     MinIO Server
    │                    │                        │                          │
    │  CreateImportTask  │                        │                          │
    │──────────────────>│                        │                          │
    │                    │  StreamObjectLinesByExcelize()                    │
    │                    │──────────────────────>│                          │
    │                    │                        │  GetObject() ─────────>│
    │                    │                        │  <──── minio.Object ───│
    │                    │  <── (rows ch, errs ch)│                          │
    │                    │                        │  rs.Next() → Read()───>│
    │                    │                        │  <────── 数据块 ────────│
    │                    │  <── rows <- record[0] │                          │
    │  (处理第0行:表头)   │                        │  rs.Next() → Read()───>│
    │                    │                        │  <────── 数据块 ────────│
    │                    │  <── rows <- record[1] │                          │
    │  TransformLine()   │                        │  rs.Next() → Read()───>│
    │  CheckRules()      │                        │  <────── 数据块 ────────│
    │  batch append      │  <── rows <- record[2] │                          │
    │  ...               │                        │  ...                     │
    │  (batch满1000条)   │                        │                          │
    │  UpdateTaskData()  │                        │                          │
    │  (写入MongoDB)     │                        │                          │
    │  batch 清空        │                        │                          │
    │  ...               │                        │  rows channel 关闭       │
    │                    │  <── linesCh = nil     │  errs channel 关闭       │
    │                    │  循环退出               │  obj.Close()             │
    │                    │  处理剩余 batch         │                          │
    │                    │  更新 task.Status=wait  │                          │
```

---

## 八、面试高频追问 Q&A

**Q1：为什么不直接把整个文件下载到本地磁盘再读取？**

> 服务是容器化部署（K8s），没有持久化磁盘；而且增加磁盘 IO 会额外引入延迟。直接通过 `io.Reader` 管道处理，网络 IO 和 CPU 处理并行进行，整体延迟更低。

**Q2：excelize 的 `f.Rows()` 真的是流式的吗？底层是怎么实现的？**

> xlsx 文件本质是一个 zip 压缩包，内部包含 `xl/worksheets/sheet1.xml` 等 XML 文件。`excelize.OpenReader(obj)` 会先将整个 zip 结构读入内存（因为需要解析 zip 目录），但 `f.Rows()` 返回的迭代器使用 `encoding/xml` 的 `xml.Decoder` 对 sheet XML 进行 SAX 风格流式解析，每次 `rs.Next()` 只读取并解析下一个 `<row>` 元素，不会把所有行数据保留在内存。**注意**：zip 目录解析仍需将完整文件读入，所以内存占用不是零，但比 `GetRows` 全量加载行数据要少得多。

**Q3：如果 MinIO 网络抖动，中途断开了怎么办？**

> `minio.Object` 底层的 MinIO Go SDK 在 `Read()` 失败时会返回 error，这个 error 会传播到 `excelize` 或 `csv.Reader`，最终以 error 形式发送到 `errs` channel。消费端收到 error 后调用 `streamCancel()` 取消整个流读取，任务状态标记为 `failed`，用户可以重新上传后再次触发。当前实现没有断点续传/自动重试机制。

**Q4：rows channel 缓冲 256，如果消费者处理很慢，会不会内存溢出？**

> 不会。缓冲满后生产者 goroutine 会阻塞在 `rows <- record` 这一行，停止调用 `rs.Next()`，底层 MinIO 的 HTTP 连接读取也会暂停（TCP 接收窗口会收缩，MinIO 服务端发送速度降低）。这是端到端的背压机制，内存中最多只有 256 行缓冲 + 消费端正在处理的 1000 条 batch，总内存可控。

**Q5：生产者 goroutine 退出时，如果 rows channel 还有未消费的数据怎么办？**

> 生产者调用 `close(rows)` 后，消费端从已关闭的 channel 中仍然可以读取剩余的缓冲数据，直到 channel 为空才会收到零值+`ok=false`。所以生产者 `defer close(rows)` 是安全的，不会丢失已推入 channel 的数据。
