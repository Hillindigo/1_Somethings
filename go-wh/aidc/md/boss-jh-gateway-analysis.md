# boss-jh-gateway 项目逻辑规则详细分析

> 分析时间：2026-03-02  
> 项目路径：`boss-jh-gateway/`

---

## 一、项目整体架构

`boss-jh-gateway` 是一个 Go 微服务工作区（Monorepo），包含 4 个独立子服务：

| 子服务 | 目录 | 核心职责 |
|--------|------|----------|
| `jh-gateway` | `srvs/jh-gateway/` | 北向网关：与 BMS 系统对接，同步设备拓扑节点 |
| `dcim_metrics` | `srvs/dcim_metrics/` | **核心**：数据点计算引擎，电能/功率计算 |
| `dcim-rack` | `srvs/dcim-rack/` | 机柜资源管理：上下电操作、布线状态管理 |
| `dcim-location` | `srvs/dcim-location/` | 位置管理：IDC 空间位置树管理 |

**公共组件**：
- `database/`：PostgreSQL 连接（pgsql.go）
- `middleware/`：JWT 认证、Casbin 权限
- `pkg/`：Redis、Nacos、定时任务等公共工具
- `orm/`：全局 GORM 实例

---

## 二、dcim_metrics 服务 —— 核心计算引擎

这是整个项目最核心的部分，负责所有数据点的周期计算、历史存储和实时推送。

### 2.1 数据点（DataPoint）模型

```go
// 数据库表：dcim_metric_points
type DcimMetricPoints struct {
    DPID                 string  // 数据点唯一标识
    LocationID           int     // 所属位置 ID
    RackNo               string  // 机柜编号（可为空）
    Name                 string  // 数据点名称
    Lable                string  // 数据点标签（用于精确检索）
    CalIntervalType      string  // 计算间隔类型
    CalFormula           string  // 计算公式（表达式字符串）
    CalMethod            string  // 计算方式
    CalculationParameter string  // 计算参数（JSON）
    Tags                 string  // 标签（管道分隔，用于分类检索）
    Unit                 string  // 单位
}
```

数据点分为 3 类：
- `DpTypeRaw = "00"`：**原始数据点**，来自硬件采集
- `DpTypeCal = "01"`：**计算数据点**，由公式计算产生
- `DpTypeTemp = "03"`：**临时数据点**，计算中间量，不持久化

---

## 三、计算间隔类型规则

系统按固定时间周期触发计算，共 8 种间隔类型：

| 常量 | 值 | 触发时机 | 说明 |
|------|----|---------|------|
| `CalIntervalType10s` | `"00"` | 每 10 秒（`now.Second()%10 == 0`） | 最细粒度，实时监控 |
| `CalIntervalType1min` | `"01"` | 每分钟整（`now.Second() == 0`） | 分钟级统计 |
| `CalIntervalType5min` | `"02"` | 每 5 分钟（`now.Minute()%5 == 0`） | 5分钟聚合 |
| `CalIntervalTypeHour` | `"03"` | 每小时整（`now.Minute() == 0`） | 小时级统计 |
| `CalIntervalTypeDay` | `"04"` | 每天 0 点（`now.Hour() == 0 && now.Minute() == 0`） | 日统计，同时触发计费配置检查 |
| `CalIntervalTypeMon` | `"05"` | 每月 1 日 0 点（`now.Day() == 1 && now.Hour() == 0`） | 月统计 |
| `CalIntervalTypeHD` | `"06"` | 每天 19:00（交接班） | 12小时班次：7:00-19:00 |
| `CalIntervalTypeHDFull` | `"07"` | 每天 07:00（交接班） | 24小时班次：7:00-次日7:00 |

### 继承触发规则（重要！）

**低间隔点自动参与更高间隔的计算**，规则如下：

```
10s  点 → 参与 10s、1min、5min、Hour、Day、Mon、HD、HDFull 所有计算
1min 点 → 参与 1min、5min、Hour、Day、Mon、HD、HDFull
5min 点 → 参与 5min、Hour、Day、Mon、HD、HDFull
Hour 点 → 参与 Hour、Day、Mon、HD、HDFull
HD   点 → 参与 HD、HDFull
HDFull 点 → 只参与 HDFull
Day  点 → 参与 Day、Mon
Mon  点 → 只参与 Mon
```

> **核心逻辑**：间隔越短的点，触发频率越高，也会在更粗粒度的计算周期中被重新计算。

---

## 四、计算方式（CalMethod）规则

### 4.1 实时计算（CalMethodRealTime = "00"）

直接用 `CalFormula` 表达式计算当前值，使用 `expr-lang/expr` 库动态执行：

```go
result, err := calculator.Calculate(Params, CalFormula)
// 实际等价于：
program, _ := expr.Compile(formula, expr.Env(params))
output, _ := expr.Run(program, params)
```

**特殊情况 - 温湿度最大/最小值**：  
当数据点 Tags 含 `HD|MDCTemp` 或 `HD|MDCHum` 时：
- 公式含 `max`：从所有参数中找最大值，同时记录来源位置 Slug（`DpDesc = maxDp.LocationSlug`）
- 公式含 `min`：从所有参数中找最小值，同时记录来源位置 Slug

### 4.2 差值计算（CalMethodDiff = "01"）

用于**电能电表累计量的区间差值**，计算本周期用电量：

```
差值 = 当前电表读数（CalFormula 执行结果） - T0时刻存储值
```

**差值计算具体流程：**

1. 从 Redis 读取上次计算时的基准值 `VT0`（存储在 `TempDpSource` 来源的 `{DPID}_T0` 键）
2. 执行公式得到当前值 `result`
3. 计算差值 `diff = result - VT0`
4. 将当前值 `result` 写回 `{DPID}_T0`（作为下次计算的基准）
5. 将差值 `diff` 作为最终结果输出

**10秒临时点（`_10S` 后缀）特殊逻辑 —— 防电表回滚：**

系统对差值计算的 10 秒点维护一个长度为 5 的滑动窗口 `TempDiffValueCacheMap[DPID][]float64`：

```
规则：当缓存窗口满 5 个且全部 < 0（连续 5 次差值为负），判定为电表回滚/重置
处理：
  1. 将 VVT0（原始 T0 基准）加上窗口第 0 个差值，修正为新基准
  2. 更新 {原始DPID}_T0 的存储值
  3. 才允许输出当次差值
```

> **意义**：防止电表断电重置后，差值出现大幅负数造成数据异常。

### 4.3 平均机柜功耗（特殊公式）

公式为 `total_rack_power/poweron_count` 时，触发特殊逻辑：

```
1. 查询该 Location 下所有开电状态机柜的功耗点位（lable='RACK_P_tol'，power_status IN ('01','02')）
2. 从 Redis 批量读取这些机柜的实时功耗值
3. total_rack_power = 所有机柜功耗之和 / 机柜数量（实际上是平均值）
4. poweron_count 始终设为 1（分母固定为1，结果即为平均功耗）
```

---

## 五、计算依赖排序规则（DAG 拓扑排序）

计算点之间存在依赖关系（A 的结果是 B 的输入参数），系统通过 DAG 拓扑排序保证**先算依赖，再算结果**。

```
CalOrder = 0：无依赖的叶子节点（最先计算）
CalOrder = 1：依赖 0 级节点
CalOrder = N：依赖最深 N-1 级节点
最后 (CalOrder = max+1)：total_rack_power/poweron_count 类型点（需等所有机柜功耗计算完毕）
```

**循环依赖检测**：递归遍历时用 `parentmap` 追踪祖先节点，发现重复则返回错误。

---

## 六、分时段计费类型（电费峰谷平尖）规则

配置表 `power_caltime_config` 存储每小时的计费类型：

| 类型常量 | 值 | 中文 | InfluxDB Tag |
|---------|-----|------|-------------|
| `TimesTypeSharp` | 4 | 尖峰 | `"sharp"` |
| `TimesTypePeak` | 1 | 峰 | `"peak"` |
| `TimesTypeFlat` | 2 | 平 | `"flat"` |
| `TimesTypeValley` | 3 | 谷 | `"valley"` |
| `TimesTypeUnknow` | 0 | 未知 | `"unknow"` |

**计费类型应用场景**：  
当小时级差值计算（`CalIntervalTypeHour + CalMethodDiff`）完成后，写入 InfluxDB 时会打上时间段标签：

```go
tag = dm.PowerCalTimeConf.GetTimeType(time.Now().Add(-10 * time.Minute))
// 取计算时刻前 10 分钟的时间段类型，作为 Tag 写入 InfluxDB
```

**计费配置切换规则**：
- 系统同时只有一条 `enable_status=1` 的启用配置
- 每天 0 点触发 `CheckUpdateCalTimes()`，检查是否有待启用（`enable_status=0`）的配置
- 若待启用配置的 `plan_enable_time` 等于当天，则自动切换（旧配置变 `status=2`，新配置变 `status=1`）

---

## 七、数据质量（Quality）规则

所有数据点携带质量码，计算时传播质量：

| 常量 | 值 | 含义 |
|------|-----|------|
| `QulkOk` | 0 | 正常 |
| `QulkUncertain` | -1 | 未知 |
| `QulkCommDisconnected` | -2 | 通讯中断 |
| `QulkCmdNoResp` | -3 | 无响应 |
| `QulkCmdReadError` | -4 | 通讯错误 |
| `QulkCmdRespError` | -5 | 异常响应 |
| `QulkCmdCrcError` | -6 | 校验错误（参数异常时被传播） |
| `QulkConfigError` | -7 | 配置错误（公式执行失败） |
| `QulkOutOfService` | -8 | 服务外 |
| `QulkOutOfValRange` | -9 | 无效值 |
| `QulkDisable` | -10 | 未启用 |

**质量传播规则**：
- 任意计算参数 `DataPointQuilty != QulkOk`，则结果点质量被标记为 `QulkCmdCrcError`（-6），并在 `DpDesc` 中记录哪个参数异常
- 公式执行失败：结果质量设为 `QulkConfigError`（-7），值保留上一次旧值（`oldDp.DataPointValue`）
- **不中断计算**：即便参数质量异常，仍会继续参与计算（取参数的 `DataPointValue` 转 float64）

---

## 八、数据流向

```
原始数据点（硬件采集）
    ↓  Kafka 消费（topic: D04_34F_Power）
    ↓  解析后写入 Redis（DMP_DATA_POINTS 命名空间）
    ↓
定时触发计算（每 10s/1min/5min/1h/1d/1mon/HD/HDFull）
    ↓  从 Redis 读取参数值
    ↓  按 CalOrder 分级并行计算（最大 100 协程，超时 5s）
    ↓  写回 Redis（计算结果）
    ↓
计算回调（CalculateCallback）
    ├─→ InfluxDB（历史存储，附带时间段 Tag）
    └─→ Kafka 推送（topic: CalDpSource，供下游消费）
```

---

## 九、jh-gateway 服务 —— 北向节点同步规则

与 BMS（楼宇自控系统）对接，同步三类节点：

| NodeType | 中文 | 存储表 |
|----------|------|--------|
| 1（Location） | 空间位置 | `northnodes.Location` |
| 2（Device） | 设备 | `northnodes.BmsDevice` |
| 3（Point） | 数据点位 | `northnodes.BmsPoint` |

**同步触发规则**：
- 服务启动时立即执行一次全量同步
- 每小时 ticker，仅在 `now.Hour() == 0`（每天 0 点）时重新同步
- **版本号比对**：若 BMS 返回的 `version` 与本地存储一致，则跳过同步

**位置层级规则**（通过 Path 解析层级深度）：
```
Path 格式："/root/level1/level2/level3"
level = Path 按 "/" 分割后数组长度 - 1
```

---

## 十、dcim-rack 服务 —— 机柜状态机规则

### 10.1 上下电状态机

机柜开电状态三值：`00=未开电`、`01=测试电`、`02=正式电`

合法操作和状态转换：

| 操作类型 | 触发条件 | 原状态 | 目标状态 |
|---------|---------|--------|---------|
| 测试电 | 客户申请测试 | 00（未开电） | 01（测试电） |
| 测试转正式电 | 测试验收通过 | 01（测试电） | 02（正式电） |
| 正式电 | 直接开正式电 | 00（未开电） | 02（正式电） |
| 测试电下电 | 撤销测试 | 01（测试电） | 00（未开电） |
| 正式电下电 | 合同到期/撤柜 | 02（正式电） | 00（未开电） |

每次状态变更都写入操作记录表 `dcom_rack_operate_records`，记录：操作人、客户服务单号、操作时间、原状态、目标状态。

### 10.2 资源状态机

布线状态（`cabled_status`）：`未布线 → 布线 → 撤销布线`  
废弃状态（`discarded_status`）：`未废弃 → 废弃 → 启用（恢复未废弃）`

**状态约束**：
- 已布线不能重复布线
- 已废弃不能重复废弃
- 未废弃不能执行启用操作
- 未布线不能撤销布线

---

## 十一、数据查询规则

### 11.1 机柜功耗数据点查询

查询按机柜上下电状态过滤，用 `DISTINCT ON` 取最新操作记录判断当前状态：

```sql
WITH filtered_result AS (
  SELECT DISTINCT ON (rack_no) rack_no, target_state as current_status
  FROM dcom_rack_operate_records
  WHERE operate_time < '{EndTime}' AND record_type = '03'
  ORDER BY rack_no, operate_time DESC
)
SELECT rack_no FROM filtered_result WHERE current_status = '02'  -- 正式电
```

### 11.2 平均机柜功耗计算规则

```sql
-- 只统计当前有效（开电）机柜
SELECT A.dpid FROM dcim_metric_points A
  INNER JOIN dcom_location B ON A.location_id = B.id 
    AND B.location_path LIKE '%,{locationID},%'
  INNER JOIN dcom_rack_business C ON C.rack_no = A.rack_no 
    AND C.power_status IN ('01','02')  -- 测试电或正式电
WHERE A.rack_no IS NOT NULL AND A.lable='RACK_P_tol'
```

平均值 = 所有开电机柜功耗之和 ÷ 机柜数量

---

## 十二、关键常量速查

### 数据来源（Redis 命名空间）
```go
CalDpSource  = "DMP_DATA_POINTS"  // 计算结果点
TempDpSource = "Temp_Datapoints"  // 临时中间量（T0基准值等）
```

### 数据点标签（Tags）分类格式
- `RACK|Energy,{更多标签}` —— 机柜能耗相关
- `HD|MDCTemp` —— 模块化数据中心温度
- `HD|MDCHum` —— 模块化数据中心湿度
- `LTG|{子类型}` —— 列头柜相关

### 数据点标识（Lable）关键值
- `RACK_P_tol` —— 机柜总功率
- `RACK_EP_h` —— 机柜小时用电量

---

## 十三、并发控制规则

- 计算工作池最大 **100 协程**，超时 **5 秒**
- 用 `sync.Mutex` 保护 `TempDiffValueCacheMap`（差值缓存写入）
- `PauseSignal` 标志位：重新加载配置时暂停计算，加载完成后恢复
- Kafka 包大小限制：单包超过 **2.56MB** 时递归二分拆包

---

## 十四、总结

`boss-jh-gateway` 的核心是 `dcim_metrics` 服务中的**多周期数据点计算引擎**，其设计亮点：

1. **表达式驱动**：计算公式以字符串存储在数据库，通过 `expr-lang/expr` 动态执行，灵活配置不需改代码
2. **DAG 拓扑保序**：自动处理计算点间的依赖关系，保证顺序正确
3. **差值防回滚**：5窗口滑动检测电表回滚，自动修正基准值
4. **峰谷平尖标记**：小时用电量写 InfluxDB 时自动打时间段标签，支持后续分时段费用分析
5. **质量透传**：参数质量异常不中断计算，但通过质量码和描述信息传播异常状态
