# ES Management Tools

用于管理 Elasticsearch 索引的工具集，包括删除索引和重新同步数据。

## esclear - 删除索引工具

删除 Elasticsearch 中的索引。

### 用途

- 删除所有索引以清空数据
- 删除特定域的索引
- 在索引结构变更前清理旧索引

### 使用方法

```bash
cd back

# 删除所有支持的索引
make esclear

# 删除指定域的索引
make esclear domain=client
make esclear domain=material
make esclear domain=process
make esclear domain=product
make esclear domain=supplier
make esclear domain=order
make esclear domain=plan

# 可以指定多个域（空格分隔）
go run ./cmd/tools/esclear.go client material process
```

### 工作流程

1. 连接到 Elasticsearch
2. 删除指定的索引（或全部索引）
3. 输出删除结果统计

### 支持的索引

- `client` - 客户索引
- `material` - 材料索引
- `process` - 工序索引
- `product` - 产品索引
- `supplier` - 供应商索引
- `order` - 订单索引
- `plan` - 计划索引

---

## esreindex - 重建索引工具

从数据库重新同步数据到 Elasticsearch，支持批量重建。

### 用途

当发生以下情况时需要运行此工具：
- ES 索引字段名称发生变化
- 数据直接导入数据库但未同步到 ES
- ES 数据损坏或不一致
- 索引结构变更后需要重建

### 使用方法

```bash
cd back

# 重建所有支持的域（client, material, process）
make esreindex

# 重建指定域
make esreindex domain=client
make esreindex domain=material
make esreindex domain=process

# 可以指定多个域（空格分隔）
go run ./cmd/tools/esreindex.go client material process
```

### 工作流程

1. 连接数据库和 Elasticsearch
2. 删除旧索引（自动清理）
3. 从数据库读取所有记录（包括软删除的记录）
4. 过滤掉已删除的记录
5. 批量索引到 ES
6. 输出索引结果统计

### 当前支持的域

- ✅ `client` - 客户数据
- ✅ `material` - 材料数据
- ✅ `process` - 工序数据
- ⏸️ `product` - 产品数据（待添加）
- ⏸️ `supplier` - 供应商数据（待添加）
- ⏸️ `order` - 订单数据（待添加）
- ⏸️ `plan` - 计划数据（待添加）

### 环境变量

工具会读取以下环境变量（未设置时使用默认值）：

- `DATABASE_DSN`: PostgreSQL 连接字符串
  - 默认: `host=localhost user=postgres password=postgres dbname=fantasy port=5432 sslmode=disable TimeZone=Asia/Shanghai`
- `ES_ADDRESS`: Elasticsearch 地址
  - 默认: `http://localhost:9200`
- `ES_USERNAME`: Elasticsearch 用户名
  - 默认: `elastic`
- `ES_PASSWORD`: Elasticsearch 密码
  - 默认: `123456`
- `LOG_LEVEL`: 日志级别
  - 默认: `info`

### Docker 环境运行

如果使用 docker-compose 启动的服务：

```bash
# 1. 确保服务正在运行
docker-compose up db es

# 2. 在本地运行重建工具（会连接到 Docker 容器中的服务）
cd back
make reindex-clients
```

或者在 Docker 容器中运行：

```bash
docker-compose exec backend go run ./cmd/tools/reindex_clients.go
```

---

## 完整使用示例

### 场景 1: 索引名称变更后重建

```bash
cd back

# 1. 删除所有旧索引
make esclear

# 2. 重建所有支持的域
make esreindex
```

### 场景 2: 只重建特定域

```bash
cd back

# 1. 删除 client 索引
make esclear domain=client

# 2. 重建 client 索引
make esreindex domain=client
```

### 场景 3: 批量重建多个域

```bash
cd back

# 删除并重建 client, material, process
go run ./cmd/tools/esclear.go client material process
go run ./cmd/tools/esreindex.go client material process
```

### 输出示例

**esclear 输出:**
```
=== ES Index Clear Tool ===
Will clear indexes: [client material process]
✓ Config loaded
✓ Logger initialized
✓ Elasticsearch connected
Deleting index: client
✓ Successfully deleted index: client
Deleting index: material
✓ Successfully deleted index: material
Deleting index: process
✓ Successfully deleted index: process
=== Clear Summary ===
Total: 3
Success: 3
Failed: 0
=== Clear completed successfully ===
```

**esreindex 输出:**
```
=== ES Reindex Tool ===
Will reindex domains: [client material process]
✓ Config loaded
✓ Logger initialized
✓ Database connected
✓ Elasticsearch connected

=== Processing domain: client ===
✓ Old index 'client' deleted
→ Fetching clients from database...
  Found 1523 clients
  [100/1523] Progress...
  [200/1523] Progress...
  ...
  [1523/1523] Progress...
✓ Domain 'client' reindexed successfully: 1500 records

=== Processing domain: material ===
✓ Old index 'material' deleted
→ Fetching materials from database...
  Found 856 materials
  [100/856] Progress...
  ...
  [856/856] Progress...
✓ Domain 'material' reindexed successfully: 850 records

=== Processing domain: process ===
✓ Old index 'process' deleted
→ Fetching processes from database...
  Found 42 processes
  [42/42] Progress...
✓ Domain 'process' reindexed successfully: 42 records

=== Overall Summary ===
Domains processed: 3
Total records indexed: 2392
Total failures: 0
=== Reindex completed successfully ===
```

---

## 注意事项

⚠️ **重要提示**：

- **esclear** 会永久删除 ES 索引，请确保可以从数据库恢复
- **esreindex** 会先删除旧索引再重建，索引过程中搜索功能可能短暂不可用
- 建议在业务低峰期运行
- 如果有大量数据（>10万条），索引过程可能需要几分钟
- 软删除的记录不会被索引到 ES

---

## 故障排查

### 连接问题

**问题**: `Failed to connect to database`
- 检查数据库是否启动：`docker-compose ps db`
- 检查 `DATABASE_DSN` 环境变量是否正确

**问题**: `Failed to connect to Elasticsearch`
- 检查 ES 是否启动：`docker-compose ps es`
- 检查 `ES_ADDRESS` 是否正确（默认: `http://localhost:9200`）
- 检查 ES 用户名密码是否正确

### 索引问题

**问题**: `Failed to delete index` (非 404 错误)
- 检查 ES 集群健康状态：`curl http://localhost:9200/_cluster/health`
- 检查索引是否被锁定或正在使用

**问题**: `Failed to index records`
- 查看详细错误日志
- 检查数据是否符合字段验证要求
- 检查 ES 磁盘空间是否充足

**问题**: `Unsupported domain`
- 当前 `esreindex` 仅支持: `client`, `material`, `process`
- 其他域需要先添加到 `esreindex.go` 中

---

## 扩展新域

如需为其他域（如 `product`, `supplier` 等）添加重建索引功能：

### 1. 在 `esreindex.go` 中添加支持

```go
// 1. 添加到支持的域列表
var supportedDomains = []string{"client", "material", "process", "product"}

// 2. 在 reindexDomain 函数中添加 case
func reindexDomain(db *gorm.DB, esSync *es.ESSync, domain string) (success, failed int, err error) {
	switch domain {
	case "client":
		return reindexClients(db, esSync)
	case "material":
		return reindexMaterials(db, esSync)
	case "process":
		return reindexProcesses(db, esSync)
	case "product":  // 新增
		return reindexProducts(db, esSync)
	default:
		return 0, 0, fmt.Errorf("unsupported domain: %s", domain)
	}
}

// 3. 实现具体的重建函数
func reindexProducts(db *gorm.DB, esSync *es.ESSync) (success, failed int, err error) {
	log.Println("→ Fetching products from database...")

	var products []productDomain.Product
	if err := db.Unscoped().Find(&products).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to fetch products: %w", err)
	}

	total := len(products)
	log.Printf("  Found %d products", total)

	if total == 0 {
		return 0, 0, nil
	}

	successCount := 0
	failCount := 0

	for i, product := range products {
		if product.DeletedAt.Valid {
			continue
		}

		if err := esSync.Index(&product); err != nil {
			log.Printf("  [%d/%d] Failed to index product ID=%d: %v", i+1, total, product.ID, err)
			failCount++
		} else {
			successCount++
			if (i+1)%100 == 0 || i+1 == total {
				log.Printf("  [%d/%d] Progress...", i+1, total)
			}
		}
	}

	return successCount, failCount, nil
}
```

### 2. 更新文档

在 README.md 中更新支持的域列表。

---

## 相关文件

- `esclear.go` - 删除索引工具
- `esreindex.go` - 重建索引工具
- `reindex_clients.go` - 遗留的客户重建工具（已废弃，推荐使用 `esreindex`）
- `Makefile` - Make 命令定义
- `README.md` - 本文档
