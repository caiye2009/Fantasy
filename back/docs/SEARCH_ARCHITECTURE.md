# Search Architecture - Domain-First Design

## 架构原则 🎯

### 单一真实来源（Single Source of Truth）

**Domain 是字段定义的唯一来源**
- 所有字段定义在 `internal/*/domain/*.go` 的 Domain 模型中
- 字段名从 `json` tag 自动提取
- 字段类型从 Go 类型自动推断
- **无需在多个地方重复定义字段**

### 配置简化

**只需配置业务逻辑，无需配置字段类型**
- YAML 配置只声明"哪些字段"参与查询/过滤/聚合
- 字段类型自动从 Domain 推断
- 配置验证：启动时自动检查配置的字段是否存在于 Domain

## 架构组件

### 1. Domain 模型（字段定义）

```go
// internal/client/domain/client.go
type Client struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    CustomNo     string    `json:"customNo"`      // ← 字段名从 json tag 提取
    CustomName   string    `json:"customName"`
    Contactor    string    `json:"contactor"`
    UnitPhone    string    `json:"unitPhone"`
    Mobile       string    `json:"mobile"`
    Email        string    `json:"email"`
    CreatedAt    time.Time `json:"createdAt"`
    // ...
}
```

**规则：**
- ✅ 使用 `json` tag 定义字段名（camelCase）
- ✅ 字段名与前端、ES 一致
- ✅ 这是字段定义的**唯一**地方

### 2. Search 配置（业务逻辑）

```yaml
# config/search/client.yaml
entityType: client
indexName: clients

# 全文搜索字段（只配置字段名 + 权重）
queryFields:
  - field: customName     # ← 必须存在于 Domain
    boost: 5.0
  - field: contactor
    boost: 3.0

# 过滤字段（只配置字段名 + 操作符，类型自动推断）
filterFields:
  - field: id
    operator: term        # type 自动推断为 keyword
  - field: customName
    operator: match       # type 自动推断为 text
  - field: createdAt
    operator: range       # type 自动推断为 date

# 聚合字段（只配置字段名 + 聚合类型）
aggregationFields:
  - field: sales
    aggType: terms        # type 自动推断为 keyword
    size: 50
```

**规则：**
- ✅ 只配置业务逻辑（哪些字段用于什么目的）
- ✅ **不需要**配置字段类型（自动推断）
- ✅ 启动时自动验证字段是否存在于 Domain
- ✅ 如果字段不存在，启动失败并报错

### 3. DomainAwareRegistry（自动化注册）

```go
// config/search_registry.go
func InitSearchRegistry() (*infra.DomainAwareRegistry, error) {
    registry := infra.NewDomainAwareRegistry()

    // 1. 注册 Domain（提取字段 schema）
    registry.RegisterDomain("client", "clients", &clientDomain.Client{})

    // 2. 加载配置（自动验证、补全）
    // 自动加载 config/search/*.yaml

    return registry, nil
}
```

**功能：**
- ✅ 从 Domain 提取字段 schema（字段名、类型）
- ✅ 加载 YAML 配置
- ✅ 验证配置的字段是否存在于 Domain
- ✅ 自动推断并补全字段类型
- ✅ 启动时报错，避免运行时问题

## 工作流程

### 添加新字段

#### 旧方式（需要改 3+ 个地方）❌
```
1. 修改 Domain 模型
2. 修改 ES mapping
3. 修改 search config
4. 修改 index_config.go
5. 修改 ToDocument() 方法
```

#### 新方式（只改 1-2 个地方）✅
```
1. 修改 Domain 模型（添加字段 + json tag）
2. 修改 search config（如果需要查询/过滤/聚合该字段）

完成！字段类型自动推断，ES 同步自动对齐。
```

### 示例：添加新字段 `faxNum`

**Step 1: 修改 Domain**
```go
// internal/client/domain/client.go
type Client struct {
    // ... 现有字段 ...
    FaxNum string `json:"faxNum"`  // ← 新增字段
}
```

**Step 2: 修改配置（如需要）**
```yaml
# config/search/client.yaml
queryFields:
  - field: faxNum    # ← 新增查询字段
    boost: 2.0

filterFields:
  - field: faxNum    # ← 新增过滤字段
    operator: term   # type 自动推断为 text
```

**完成！** 🎉
- ES 同步时自动使用新字段（通过 `ToDocument()`）
- 字段类型自动推断
- 启动时自动验证

### 添加新 Entity

**Step 1: 注册 Domain**
```go
// config/search_registry.go
func registerAllDomains(registry *infra.DomainAwareRegistry) error {
    // ... 现有注册 ...

    // 新增
    registry.RegisterDomain("newEntity", "new_entities", &newDomain.NewEntity{})

    return nil
}
```

**Step 2: 创建配置文件**
```yaml
# config/search/new_entity.yaml
entityType: newEntity
indexName: new_entities

queryFields:
  - field: name
    boost: 5.0
# ...
```

**完成！** 🎉

## 字段类型推断规则

| Go 类型 | ES 类型 | Filter 类型 | Agg 类型 |
|---------|---------|-------------|----------|
| `string` | `text` | `text` | `keyword` |
| `int`, `uint`, `int64` | `long` | `numeric` | `numeric` |
| `float64` | `double` | `numeric` | `numeric` |
| `bool` | `boolean` | `keyword` | `keyword` |
| `time.Time` | `date` | `date` | `date` |

## 验证机制

### 启动时验证

```
=== Initializing Search Registry (Domain-Aware) ===
✓ Registered domain schema for 'client' with 24 fields
✓ Registered domain schema for 'vendor' with 18 fields
...
✓ Loaded search config for 'client'
✓ Loaded search config for 'vendor'
...

如果字段不存在：
✗ Failed to load config: filterFields[3]: field 'oldFieldName' not found in domain model client
```

### 配置错误示例

```yaml
# ❌ 错误：字段不存在于 Domain
filterFields:
  - field: nonExistentField  # ← 启动时报错
    operator: term

# ✅ 正确：字段存在于 Domain
filterFields:
  - field: customName
    operator: match
```

## 文件结构

```
back/
├── pkg/es/
│   └── schema.go                        # Domain 字段提取工具
├── internal/
│   ├── client/domain/client.go          # ← 字段定义（唯一来源）
│   └── search/
│       ├── domain/search_config.go      # 配置数据结构
│       └── infra/
│           └── domain_aware_registry.go # ← Domain 感知的注册中心
└── config/
    ├── search_registry.go               # ← 统一注册入口
    └── search/
        ├── client.yaml                  # ← 简化的配置（无 type）
        ├── vendor.yaml
        └── ...
```

## 迁移检查清单

从旧架构迁移到新架构：

- [x] 创建 `pkg/es/schema.go`（字段提取工具）
- [x] 创建 `internal/search/infra/domain_aware_registry.go`
- [x] 创建 `config/search_registry.go`（统一注册）
- [x] 更新 `config/search/client.yaml`（使用新字段名，移除 type）
- [x] 更新 `config/services.go`（使用 DomainAwareRegistry）
- [x] 更新 `internal/search/application/search_service.go`（使用 DomainAwareRegistry）
- [ ] 更新其他 entity 的 YAML 配置
- [ ] 运行 reindex 工具重建 ES 索引

## 最佳实践

### DO ✅

1. **字段定义只在 Domain 中修改**
   ```go
   type Client struct {
       NewField string `json:"newField"`  // ← 只改这里
   }
   ```

2. **配置只声明业务逻辑**
   ```yaml
   queryFields:
     - field: newField    # ← 只配置字段名和权重
       boost: 3.0         # 类型自动推断
   ```

3. **依赖自动验证**
   - 启动时会自动检查配置字段是否存在
   - 无需手动验证

### DON'T ❌

1. **不要在配置中定义字段类型**
   ```yaml
   filterFields:
     - field: customName
       type: text         # ❌ 不需要，会自动推断
       operator: match
   ```

2. **不要在多个地方定义字段**
   - ❌ Domain + ES mapping + config
   - ✅ 只在 Domain 定义

3. **不要跳过启动验证**
   - 配置错误会导致启动失败
   - 这是设计目标（Fail Fast）

## 总结

### 核心优势

1. **单一来源** - Domain 是字段定义的唯一真实来源
2. **自动对齐** - ES 同步、配置验证自动与 Domain 对齐
3. **简化配置** - 只需配置业务逻辑，无需配置字段类型
4. **启动验证** - 配置错误在启动时发现，避免运行时问题
5. **易于维护** - 添加字段只需修改 Domain，配置自动生效

### 从旧字段迁移

```bash
# 1. 更新配置文件（使用新字段名）
vim config/search/client.yaml

# 2. 重启服务（自动验证配置）
make run

# 3. 重建 ES 索引（同步新字段）
make reindex-clients
```

完成！🎉
