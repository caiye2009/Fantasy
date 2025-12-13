# 数据库迁移指南

## 更新内容

### 1. User 表新增字段

- **department** (varchar(100)): 用户所属部门

### 2. 角色管理优化

- 角色列表统一在 `internal/user/domain/user.go` 中维护
- 新增角色验证函数 `IsValidRole()`
- 创建/更新用户时会自动验证角色是否有效

## 迁移步骤

### 方式 1: 自动迁移（推荐）

GORM 会自动检测结构体变化并更新表结构：

```bash
cd back
make run
```

或使用 Docker：

```bash
docker-compose up backend
```

启动时会自动执行 `AutoMigrate`，添加 `department` 字段。

### 方式 2: 手动 SQL（可选）

如果需要手动执行 SQL：

```sql
-- 添加 department 字段
ALTER TABLE users ADD COLUMN department VARCHAR(100);
```

## 验证迁移

### 1. 检查表结构

```sql
-- PostgreSQL
\d users

-- 或
SELECT column_name, data_type, character_maximum_length
FROM information_schema.columns
WHERE table_name = 'users';
```

应该看到：
- `id` (integer)
- `login_id` (varchar)
- `username` (varchar)
- `department` (varchar) ✅ **新增**
- `password_hash` (varchar)
- `email` (varchar)
- `role` (varchar)
- `status` (varchar)
- `has_init_pass` (boolean)
- `created_at` (timestamp)
- `updated_at` (timestamp)
- `deleted_at` (timestamp) - 软删除

### 2. 测试 API

创建用户时传入 department：

```bash
curl -X POST http://localhost:8081/api/v1/user/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "login_id": "emp001",
    "username": "张三",
    "department": "技术部",
    "email": "zhangsan@example.com",
    "role": "sales"
  }'
```

## 回滚（如需要）

如果需要回滚：

```sql
-- 删除 department 字段
ALTER TABLE users DROP COLUMN department;
```

## 注意事项

1. **现有数据**：
   - 现有用户的 `department` 字段会是 `NULL` 或空字符串
   - 可以通过 API 更新用户时补充部门信息

2. **非破坏性变更**：
   - 本次迁移只是**新增字段**，不会删除或修改现有数据
   - 所有现有功能保持兼容

3. **Casbin 策略**：
   - Casbin 权限策略保持不变
   - 角色列表与 Casbin 策略保持一致

## 生产环境部署建议

1. **备份数据库**：
   ```bash
   pg_dump -U postgres -d fantasy > backup_$(date +%Y%m%d_%H%M%S).sql
   ```

2. **先在测试环境验证**：
   ```bash
   # 测试环境
   docker-compose -f docker-compose.test.yml up backend
   ```

3. **分步部署**：
   - 先部署数据库迁移（添加 department 字段）
   - 验证无误后再部署新代码
   - 最后逐步更新用户的 department 信息

4. **监控日志**：
   ```bash
   docker-compose logs -f backend | grep -i "migration"
   ```

## 常见问题

### Q1: 迁移失败怎么办？

检查日志：
```bash
docker-compose logs backend | tail -50
```

常见原因：
- 数据库连接失败
- 权限不足
- 字段已存在

### Q2: 如何批量更新现有用户的部门？

```sql
-- 示例：根据 role 批量设置部门
UPDATE users SET department = '销售部' WHERE role = 'sales';
UPDATE users SET department = '人力资源部' WHERE role = 'hr';
UPDATE users SET department = '技术部' WHERE role = 'follower';
```

### Q3: 角色验证失败怎么办？

确保角色值在以下列表中：
- `admin`
- `hr`
- `sales`
- `follower`
- `assistant`
- `user`

如果需要新增角色，修改 `internal/user/domain/user.go` 中的角色定义。
