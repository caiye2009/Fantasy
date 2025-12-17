# ES Reindex Tools

工具集用于重建 Elasticsearch 索引并重新同步数据。

## reindex_clients

重建 clients 索引，将数据库中的所有客户数据重新同步到 Elasticsearch。

### 用途

当发生以下情况时需要运行此工具：
- ES 索引字段名称发生变化（如从旧字段 name, contact, phone 改为新字段 customName, contactor, unitPhone）
- 数据直接导入数据库但未同步到 ES
- ES 数据损坏或不一致

### 使用方法

#### 方式 1: 使用 Makefile（推荐）

```bash
cd back
make reindex-clients
```

#### 方式 2: 直接运行

```bash
cd back
go run ./cmd/tools/reindex_clients.go
```

### 工作流程

1. 连接数据库和 Elasticsearch
2. 删除旧的 `clients` 索引（清除旧字段名数据）
3. 从数据库读取所有客户记录
4. 使用新字段名（customName, customNo, contactor, unitPhone 等）批量索引到 ES
5. 输出索引结果统计

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

### 输出示例

```
=== Starting Client Reindex Tool ===
✓ Config loaded
✓ Logger initialized
✓ Database connected (PostgreSQL)
✓ Elasticsearch connected
✓ Old index deleted
=== Starting to reindex clients ===
Found 1523 clients in database
[100/1523] Indexed client: ID=100, CustomNo=C100, CustomName=示例客户
[200/1523] Indexed client: ID=200, CustomNo=C200, CustomName=测试公司
...
[1523/1523] Indexed client: ID=1523, CustomNo=C1523, CustomName=最后一个客户
=== Reindex Summary ===
Total: 1523
Success: 1500
Failed: 0
Skipped (deleted): 23
=== Reindex completed successfully ===
```

### 注意事项

⚠️ **重要提示**：
- 此工具会**删除**现有的 `clients` 索引，请确保可以从数据库恢复
- 索引过程中 ES 搜索功能可能短暂不可用
- 建议在业务低峰期运行
- 如果有大量数据（>10万条），索引过程可能需要几分钟

### 故障排查

**问题**: `Failed to connect to database`
- 检查数据库是否启动：`docker-compose ps db`
- 检查 `DATABASE_DSN` 环境变量是否正确

**问题**: `Failed to connect to Elasticsearch`
- 检查 ES 是否启动：`docker-compose ps es`
- 检查 `ES_ADDRESS` 是否正确
- 检查 ES 用户名密码是否正确

**问题**: `Failed to index client`
- 查看详细错误日志
- 检查 ES 集群健康状态：`curl http://localhost:9200/_cluster/health`
- 检查数据是否符合字段验证要求

## 扩展

如需为其他实体（如 vendors, materials 等）创建重建索引工具，可参考 `reindex_clients.go` 的实现。
