# 完整搜索架构总结

## 📋 核心设计原则

### 1. 双来源定义（Single Source of Truth）
```
┌──────────────────────┐
│ Domain 层            │  ← 字段定义的唯一来源
│ internal/*/domain/   │     - 所有业务字段
│                      │     - CalculatePriorityScore() 可选
└──────────────────────┘
           ↓
┌──────────────────────┐
│ Config 层            │  ← 搜索行为的唯一配置
│ config/search/       │     - query/filter/agg 规则
│                      │     - defaultSort 排序
└──────────────────────┘
```

---

## 🏗️ 三层搜索模型

### Query（全文搜索）
**用途**: 关键词匹配，相关度计算
**特点**: 算分、boost 权重
**触发**: 用户输入搜索框

```yaml
queryFields:
  - field: customName
    boost: 5.0    # 权重越高越重要
  - field: address
    boost: 2.0
```

---

### Filter（结构化过滤）
**用途**: 精确条件过滤
**特点**: 不算分，快速过滤
**触发**: 用户选择下拉筛选

```yaml
filterFields:
  - field: status
    operator: terms    # 多选
  - field: createdAt
    operator: range    # 范围查询
```

---

### Aggregation（动态统计）
**用途**: 基于当前结果集统计，生成筛选项
**特点**: 实时联动，动态去重
**触发**: 每次 query/filter 变化

```yaml
aggregationFields:
  - field: status
    aggType: terms
    size: 50
    supportSearch: true     # 支持下拉框内搜索
    excludeSelf: false      # 是否排除自身条件
```

---

## 🎯 排序机制

### 自动组装的排序逻辑

```
最终排序 = defaultSort（后端自动）+ userSort（用户选择）+ 兜底（id asc）
```

#### 配置示例
```yaml
# config/search/client.yaml
defaultSort:
  - field: priorityScore
    order: desc
    type: computed        # 计算字段（不在 Domain 中）
    missing: _last        # 没有该字段的排最后
```

#### 前端请求
```json
{
  "index": "clients",
  "query": "北京",
  "sort": [
    {"field": "createdAt", "order": "desc"}
  ]
}
```

#### 后端自动组装的 ES 查询
```json
{
  "sort": [
    {"priorityScore": {"order": "desc", "missing": "_last"}},
    {"createdAt": {"order": "desc"}},
    {"id": {"order": "asc"}}
  ]
}
```

---

## 💡 priorityScore 计算机制

### 设计原则
- **可选**：Domain 不强制实现
- **内部化**：逻辑在 Domain 内部
- **默认返回 0**：不需要时自动忽略

### 实现方式

```go
// internal/client/domain/client.go

// CalculatePriorityScore 计算优先级分数（可选）
// 默认返回 0，需要时自己修改
func (c *Client) CalculatePriorityScore() int {
    score := 0

    // 1. 状态评分
    switch c.CustomStatus {
    case "active":    score += 200
    case "potential": score += 100
    case "dormant":   score += 50
    }

    // 2. 时间新鲜度
    if c.InputDate != nil {
        daysSince := int(time.Since(*c.InputDate).Hours() / 24)
        if daysSince < 30 {
            score += 50 - (daysSince / 2)
        }
    }

    // 3. 数据完整度
    if c.Contactor != "" { score += 10 }
    if c.UnitPhone != "" || c.Mobile != "" { score += 10 }
    if c.Email != "" { score += 10 }

    return score
}
```

### ES 同步时自动调用

```go
// pkg/es/sync.go

// 通过反射自动调用 CalculatePriorityScore()
if priorityScore := CalculatePriorityScoreIfExists(doc); priorityScore > 0 {
    docData["priorityScore"] = priorityScore
}
```

**机制**：
- ✅ 有方法且返回 > 0 → 注入 priorityScore 字段
- ✅ 没有方法或返回 0 → 不添加该字段
- ✅ 完全可选，不影响其他 Domain

---

## 🔄 完整数据流

### 1. 前端请求
```json
POST /api/v1/search
{
  "index": "clients",
  "query": "北京",
  "filters": {
    "status": ["active"]
  },
  "sort": [
    {"field": "createdAt", "order": "desc"}
  ],
  "aggRequests": {
    "sales": {
      "search": "",
      "size": 10
    }
  },
  "pagination": {
    "offset": 0,
    "size": 20
  }
}
```

---

### 2. 后端处理流程

```
┌────────────────────────────────────────┐
│ SearchHandler.Search()                 │
│ - 接收前端请求                          │
│ - 验证 JSON 格式                        │
└──────────────────┬─────────────────────┘
                   ↓
┌────────────────────────────────────────┐
│ SearchService.Search()                 │
│ 1. GetConfigByIndex("clients")         │
│    → 获取 queryFields, filterFields 等 │
│                                        │
│ 2. validateRequest()                   │
│    → 验证字段是否在白名单               │
│                                        │
│ 3. buildSortFields()  ← 关键！         │
│    → 自动组装排序：                     │
│      [priorityScore desc,              │
│       createdAt desc,                  │
│       id asc]                          │
│                                        │
│ 4. buildCriteria()                     │
│    → 构建 SearchCriteria               │
└──────────────────┬─────────────────────┘
                   ↓
┌────────────────────────────────────────┐
│ ESSearchRepository.Search()            │
│ 1. buildESQuery()                      │
│    → 构建完整 ES DSL                    │
│                                        │
│ 2. executeSearch()                     │
│    → 调用 ES API                        │
│                                        │
│ 3. parseResponse()                     │
│    → 解析结果和聚合                     │
└──────────────────┬─────────────────────┘
                   ↓
┌────────────────────────────────────────┐
│ 返回结果                                │
│ {                                      │
│   "items": [...],                      │
│   "total": 150,                        │
│   "aggregations": {                    │
│     "sales": {                         │
│       "buckets": [                     │
│         {"key": "张三", "docCount": 50} │
│       ]                                │
│     }                                  │
│   }                                    │
│ }                                      │
└────────────────────────────────────────┘
```

---

### 3. ES 查询执行

```json
{
  "query": {
    "bool": {
      "must": [
        {
          "multi_match": {
            "query": "北京",
            "fields": ["customName^5", "address^2"]
          }
        }
      ],
      "filter": [
        {"terms": {"status": ["active"]}}
      ]
    }
  },
  "sort": [
    {"priorityScore": {"order": "desc", "missing": "_last"}},
    {"createdAt": {"order": "desc"}},
    {"id": {"order": "asc"}}
  ],
  "aggs": {
    "sales": {
      "composite": {
        "size": 10,
        "sources": [
          {"sales": {"terms": {"field": "sales.keyword"}}}
        ]
      }
    }
  },
  "from": 0,
  "size": 20
}
```

---

## 📊 下拉框联动机制

### 工作原理

```
用户输入 query "北京" + 选择 status="active"
   ↓
ES 执行搜索：匹配"北京" + 过滤 status=active
   ↓
在结果集中聚合 sales 字段
   ↓
返回下拉选项：
  - 张三（50 条）
  - 李四（30 条）
  - 王五（20 条）
```

**关键特性**：
1. ✅ 下拉框数据基于**当前筛选条件**
2. ✅ 用户每次修改条件，下拉框自动更新
3. ✅ 支持下拉框内搜索（`search` 参数）
4. ✅ 支持分页加载（`composite aggregation` + `after_key`）
5. ✅ 无限滚动（每页 10 条）

---

### 下拉框搜索示例

```json
// 前端请求
{
  "index": "clients",
  "query": "北京",
  "aggRequests": {
    "sales": {
      "search": "张",    // 下拉框内搜索"张"
      "size": 10,
      "after": null
    }
  }
}

// 后端返回
{
  "aggregations": {
    "sales": {
      "buckets": [
        {"key": "张三", "docCount": 50},
        {"key": "张四", "docCount": 20}
      ],
      "after": {"sales": "张四"},
      "hasMore": true
    }
  }
}

// 前端加载更多（滚动到底部）
{
  "aggRequests": {
    "sales": {
      "search": "张",
      "size": 10,
      "after": {"sales": "张四"}  // 传上次的 after_key
    }
  }
}
```

---

## 🗂️ 分页机制

### 主列表分页（items）

```json
{
  "pagination": {
    "offset": 0,   // 从第几条开始
    "size": 20     // 每页 20 条（最大 100）
  }
}
```

**响应**：
```json
{
  "items": [...],  // 当前页数据
  "total": 1523    // 总条数
}
```

**特点**：
- ✅ 使用 `from` + `size` 分页
- ✅ 最大 100 条/页（后端限制）
- ✅ 返回总条数（前端显示分页器）

---

### 聚合分页（下拉框）

```json
{
  "aggRequests": {
    "sales": {
      "size": 10,                  // 每次加载 10 条
      "after": {"sales": "张三"}    // 上次最后一条的 key
    }
  }
}
```

**响应**：
```json
{
  "aggregations": {
    "sales": {
      "buckets": [...],
      "after": {"sales": "王五"},  // 下次请求用这个
      "hasMore": true             // 是否还有更多
    }
  }
}
```

**特点**：
- ✅ 使用 `composite aggregation`
- ✅ 无限滚动（前端滚到底加载更多）
- ✅ 服务器端分页（不是一次性加载全部）
- ✅ 支持搜索 + 分页同时进行

---

## 📁 文件结构

```
back/
├── config/
│   ├── search_registry.go          # 统一注册 Domain + Config
│   └── search/
│       ├── client.yaml             # Client 搜索配置
│       ├── supplier.yaml
│       └── ...
│
├── internal/
│   ├── client/domain/
│   │   ├── client.go               # Domain 模型（字段定义）
│   │   └── CalculatePriorityScore() # 可选评分函数
│   │
│   └── search/
│       ├── domain/
│       │   ├── search_config.go    # 配置数据结构
│       │   └── search_criteria.go  # 搜索条件
│       ├── application/
│       │   ├── dto.go              # 请求/响应 DTO
│       │   └── search_service.go   # 搜索服务（自动组装排序）
│       ├── infra/
│       │   ├── domain_aware_registry.go  # Domain 感知的注册中心
│       │   ├── es_search_repository.go   # ES 查询执行
│       │   ├── query_builder.go          # Query 构建器
│       │   └── aggregation_builder.go    # Aggregation 构建器
│       └── interfaces/
│           └── search_handler.go         # HTTP Handler
│
└── pkg/es/
    ├── indexable.go                # ES 文档接口
    ├── sync.go                     # ES 同步（自动调用 CalculatePriorityScore）
    └── schema.go                   # Domain schema 提取
```

---

## 🎨 前端使用示例

### 基本搜索
```typescript
const searchClients = async (keyword: string) => {
  const response = await axios.post('/api/v1/search', {
    index: 'clients',
    query: keyword,
    pagination: {
      offset: 0,
      size: 20
    }
  })

  return response.data  // { items, total }
}
```

---

### 带筛选的搜索
```typescript
const searchWithFilters = async (filters: any) => {
  const response = await axios.post('/api/v1/search', {
    index: 'clients',
    query: '',
    filters: {
      status: ['active', 'potential'],
      country: ['CN']
    },
    sort: [
      { field: 'createdAt', order: 'desc' }
    ],
    pagination: {
      offset: 0,
      size: 20
    }
  })

  return response.data
}
```

---

### 加载下拉框选项
```typescript
const loadSalesOptions = async (searchTerm: string, after: any = null) => {
  const response = await axios.post('/api/v1/search', {
    index: 'clients',
    query: '',
    aggRequests: {
      sales: {
        search: searchTerm,
        size: 10,
        after: after
      }
    },
    pagination: { size: 0 }  // 不需要 items
  })

  const aggResult = response.data.aggregations.sales
  return {
    options: aggResult.buckets,
    after: aggResult.after,
    hasMore: aggResult.hasMore
  }
}
```

---

## ⚙️ 配置示例

### 完整的 client.yaml
```yaml
entityType: client
indexName: clients

# 全文搜索字段（按优先级排序）
queryFields:
  - field: customName       # 客户名称（最重要）
    boost: 5.0
  - field: customNameEn     # 英文名称
    boost: 4.0
  - field: customNo         # 客户代码
    boost: 4.0
  - field: contactor        # 联系人
    boost: 3.0
  - field: unitPhone        # 电话
    boost: 2.0
  - field: mobile           # 手机
    boost: 2.0
  - field: email            # 邮箱
    boost: 2.0
  - field: address          # 中文地址
    boost: 1.0

# 过滤字段（type 自动从 Domain 推断）
filterFields:
  - field: id
    operator: term
  - field: customNo
    operator: term
  - field: customName
    operator: match
  - field: sales
    operator: term
  - field: country
    operator: term
  - field: customStatus
    operator: term
  - field: createdAt
    operator: range
  - field: updatedAt
    operator: range

# 聚合字段（用于下拉筛选器）
aggregationFields:
  - field: sales
    aggType: terms
    size: 50
    supportSearch: true
  - field: country
    aggType: terms
    size: 100
    supportSearch: true
  - field: customStatus
    aggType: terms
    size: 10
    supportSearch: false

# 默认排序（后端自动添加，前端不感知）
defaultSort:
  - field: priorityScore
    order: desc
    type: computed        # 计算字段
    missing: _last        # 没有该字段的排最后
```

---

## 🔍 核心优势

### 1. 字段定义单一来源
```
✅ 添加新字段：只需修改 Domain
✅ ES 同步：自动对齐 Domain 字段
✅ 配置验证：启动时自动检查字段是否存在
```

---

### 2. priorityScore 完全可选
```
✅ 需要优先级：在 Domain 中实现 CalculatePriorityScore()
✅ 不需要优先级：不实现或返回 0，自动忽略
✅ 对前端透明：前端完全不知道该字段存在
```

---

### 3. 排序自动组装
```
✅ 后端默认排序（priorityScore）
✅ 用户选择排序（createdAt, updatedAt 等）
✅ 兜底排序（id）
✅ 前端只传用户选择，后端自动组装完整排序
```

---

### 4. 下拉框智能联动
```
✅ 基于当前 query + filters 动态计算
✅ 支持下拉框内搜索
✅ 服务器端分页（每页 10 条）
✅ 无限滚动加载
```

---

### 5. 分页机制完整
```
✅ 主列表分页：from + size（最大 100）
✅ 聚合分页：composite aggregation + after_key
✅ 返回总条数：前端显示分页器
✅ 返回 hasMore：前端判断是否还有更多
```

---

## 🚀 使用流程

### 添加新 Entity 的搜索功能

#### 1. 在 Domain 中定义字段
```go
// internal/supplier/domain/supplier.go
type Supplier struct {
    ID      uint   `json:"id"`
    Name    string `json:"name"`
    Contact string `json:"contact"`
    // ...
}

// 可选：实现优先级评分
func (s *Supplier) CalculatePriorityScore() int {
    return 0  // 默认不需要
}
```

---

#### 2. 注册 Domain
```go
// config/search_registry.go
func registerAllDomains(registry *infra.DomainAwareRegistry) error {
    // ...
    registry.RegisterDomain("supplier", "suppliers", &supplierDomain.Supplier{})
    return nil
}
```

---

#### 3. 创建配置文件
```yaml
# config/search/supplier.yaml
entityType: supplier
indexName: suppliers

queryFields:
  - field: name
    boost: 5.0
  - field: contact
    boost: 3.0

filterFields:
  - field: name
    operator: match

aggregationFields: []

defaultSort: []  # 不需要优先级排序
```

---

#### 4. 前端调用
```typescript
const searchSuppliers = async () => {
  const response = await axios.post('/api/v1/search', {
    index: 'suppliers',  // ← 使用索引名
    query: '关键词',
    pagination: { offset: 0, size: 20 }
  })

  return response.data
}
```

---

## 📖 总结

### 核心设计
1. **双来源定义**：Domain（字段） + Config（行为）
2. **三层搜索**：Query + Filter + Aggregation
3. **自动排序**：defaultSort + userSort + 兜底
4. **可选评分**：CalculatePriorityScore() 默认返回 0
5. **智能联动**：下拉框基于当前条件动态计算
6. **完整分页**：主列表 + 聚合双重分页机制

### 关键特性
- ✅ 字段单一来源，维护简单
- ✅ priorityScore 可选且透明
- ✅ 排序自动组装，前端无感知
- ✅ 下拉框智能联动 + 搜索 + 分页
- ✅ 配置验证，启动时发现错误
- ✅ 类型自动推断，减少配置

### 数据流
```
前端请求 → SearchService → 自动组装排序 → ES 查询 → 返回结果
            ↓
         验证字段（基于 Domain）
            ↓
         构建 DSL（query + filter + agg + sort）
```

---

**完成！🎉 整个搜索架构已经实现并验证通过。**
