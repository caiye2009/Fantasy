# Auth & User 模块重构 - 变更说明

## 更新日期
2025-12-12

## 核心改进

### 1. 统一角色管理 ✅

**问题**：之前角色定义分散在多个地方（domain 枚举、DTO validation、Casbin 策略），添加新角色需要修改多处代码。

**解决方案**：
- 在 `internal/user/domain/user.go` 中统一维护角色列表
- 提供工具函数：
  - `GetAllRoles()`: 获取所有角色枚举
  - `GetAllRoleStrings()`: 获取角色字符串列表
  - `IsValidRole(role string)`: 验证角色是否有效
- 业务层自动调用 `IsValidRole()` 验证

**影响**：
- 添加新角色只需修改 `domain/user.go` 一处
- 所有验证逻辑自动生效

**修改文件**：
- `internal/user/domain/user.go` ✅
- `internal/user/application/user_service.go` ✅
- `internal/user/domain/error.go` ✅

---

### 2. Context 传递优化 ✅

**问题**：Auth 中间件将用户信息放在 gin.Context 中，但业务层难以访问。

**解决方案**：
- 新增 `pkg/auth/context.go` 工具函数
- Auth 中间件同时设置到 `gin.Context` 和 `request.Context`
- 业务层通过 `auth.GetLoginID(ctx)` 和 `auth.GetRole(ctx)` 获取

**工具函数**：
```go
// 获取用户信息（带 bool 返回值）
loginId, ok := auth.GetLoginID(ctx)
role, ok := auth.GetRole(ctx)

// 必须获取（返回空字符串如果不存在）
loginId := auth.MustGetLoginID(ctx)
role := auth.MustGetRole(ctx)

// 设置用户信息（测试用）
ctx = auth.SetLoginID(ctx, "user123")
ctx = auth.SetRole(ctx, "admin")
```

**修改文件**：
- `pkg/auth/context.go` ✅ (新增)
- `pkg/auth/auth.go` ✅

---

### 3. User 模型扩展 ✅

**新增字段**：
- `Department` (string): 用户所属部门

**保留字段**：
- `HasInitPass` (bool): 标记是否为初始密码
- `DeletedAt` (gorm.DeletedAt): 软删除支持

**修改文件**：
- `internal/user/domain/user.go` ✅
- `internal/user/infra/user_po.go` ✅
- `internal/user/application/dto.go` ✅
- `internal/user/application/assembler.go` ✅
- `internal/user/application/user_service.go` ✅

---

### 4. Handler Context 传递检查 ✅

**验证结果**：所有 Handler 都正确使用 `c.Request.Context()` 传递给 Service。

**检查范围**：
- ✅ order_handler.go
- ✅ product_handler.go
- ✅ user_handler.go
- ✅ client_handler.go
- ✅ supplier_handler.go
- ✅ material_handler.go
- ✅ process_handler.go
- ✅ plan_handler.go
- ✅ search_handler.go
- ✅ pricing handlers

---

### 5. 数据库迁移 ✅

**变更**：
- Users 表新增 `department` 字段 (varchar(100))

**迁移方式**：
- 自动迁移：启动应用时 GORM AutoMigrate 自动添加
- 手动迁移：见 `docs/MIGRATION_GUIDE.md`

---

## 文件变更清单

### 新增文件
- `pkg/auth/context.go` - Context 工具函数
- `docs/AUTH_CONTEXT_USAGE.md` - 使用指南
- `docs/MIGRATION_GUIDE.md` - 迁移指南
- `CHANGELOG_AUTH_REFACTOR.md` - 本文档

### 修改文件
- `pkg/auth/auth.go` - 添加 context 设置
- `internal/user/domain/user.go` - 角色管理函数 + Department 字段
- `internal/user/domain/error.go` - ErrInvalidRole
- `internal/user/infra/user_po.go` - Department 字段
- `internal/user/application/dto.go` - Department 字段
- `internal/user/application/assembler.go` - Department 字段
- `internal/user/application/user_service.go` - 角色验证 + Department 更新

---

## API 变更

### CreateUserRequest
```json
{
  "login_id": "emp001",
  "username": "张三",
  "department": "技术部",  // ✅ 新增（可选）
  "email": "zhangsan@example.com",
  "role": "sales"         // ✅ 必填，会验证是否有效
}
```

### UpdateUserRequest
```json
{
  "username": "张三",
  "department": "销售部",  // ✅ 新增（可选）
  "email": "zhangsan@example.com",
  "role": "sales"         // ✅ 可选，会验证是否有效
}
```

### UserResponse
```json
{
  "id": 1,
  "login_id": "emp001",
  "username": "张三",
  "department": "技术部",  // ✅ 新增
  "email": "zhangsan@example.com",
  "role": "sales",
  "status": "active",
  "has_init_pass": false,
  "created_at": "2025-12-12T10:00:00Z",
  "updated_at": "2025-12-12T10:00:00Z"
}
```

---

## 权限体系总结

### 分层防护
```
第 0 层 - JWT
  └─ 验证 token，提取 loginId 和 role

第 1 层 - Casbin
  └─ 检查 role 对资源类型（URL 路径）的操作权限

第 2 层 - ES 查询构建（未来实现）
  └─ 根据 role 动态构建查询条件和过滤器

第 3 层 - 业务层
  └─ 数据归属验证、细粒度权限检查
```

### 当前角色列表
- `admin` - 管理员（全部权限）
- `hr` - 人力资源（用户管理）
- `sales` - 销售（业务操作）
- `follower` - 跟单员（订单跟进）
- `assistant` - 助理（只读）
- `user` - 普通用户

### 添加新角色步骤
1. 修改 `internal/user/domain/user.go`
   ```go
   const (
       UserRoleDeveloper UserRole = "developer"  // 新增
   )

   func GetAllRoles() []UserRole {
       return []UserRole{
           // ...
           UserRoleDeveloper,  // 添加到列表
       }
   }
   ```

2. 修改 `config/casbin.go` 添加策略（如需要）
   ```go
   {"developer", "/api/v1/product/list", "GET"},
   ```

3. 完成！所有验证自动生效

---

## 使用示例

### 在 Service 层获取用户信息
```go
import "back/pkg/auth"

func (s *OrderService) Create(ctx context.Context, req *CreateOrderRequest) (*OrderResponse, error) {
    // 获取当前登录用户
    loginId := auth.MustGetLoginID(ctx)
    role := auth.MustGetRole(ctx)

    // 使用用户信息
    order := &domain.Order{
        OrderNo:   req.OrderNo,
        CreatedBy: loginId,  // 记录创建人
    }

    // 根据角色进行权限控制
    if role != "admin" && role != "sales" {
        return nil, ErrPermissionDenied
    }

    // 保存订单
    return s.repo.Save(ctx, order)
}
```

详细示例见：`docs/AUTH_CONTEXT_USAGE.md`

---

## 测试建议

### 1. 单元测试
```bash
cd back
go test ./internal/user/...
go test ./pkg/auth/...
```

### 2. 集成测试
```bash
# 启动服务
make run

# 测试创建用户（带 department）
curl -X POST http://localhost:8081/api/v1/user/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "login_id": "emp001",
    "username": "测试用户",
    "department": "技术部",
    "role": "sales"
  }'

# 测试无效角色（应该失败）
curl -X POST http://localhost:8081/api/v1/user/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "login_id": "emp002",
    "username": "测试用户2",
    "role": "invalid_role"
  }'
```

### 3. 数据库检查
```sql
-- 检查 department 字段
SELECT id, login_id, username, department, role FROM users LIMIT 5;
```

---

## 性能影响

- **Context 传递**: 无额外性能开销（只是传递指针）
- **角色验证**: O(n) 遍历，n ≤ 10，可忽略
- **数据库**: 新增一个 varchar 字段，影响极小

---

## 向后兼容性

✅ **完全兼容**

- 现有 API 继续工作
- Department 字段可选
- 现有角色保持不变
- 现有数据不受影响（department 为 NULL）

---

## 后续工作（可选）

### 1. ES 权限过滤
在 `internal/search/infra/es_search_repo.go` 中：
```go
func (r *ESSearchRepo) buildESQuery(ctx context.Context, query *domain.SearchQuery) map[string]interface{} {
    loginId := auth.MustGetLoginID(ctx)
    role := auth.MustGetRole(ctx)

    // 根据 role 添加过滤条件
    // ...
}
```

### 2. 操作日志
记录所有创建/更新/删除操作的执行人：
```go
log.Info("Order created",
    "orderId", order.ID,
    "createdBy", auth.MustGetLoginID(ctx),
)
```

### 3. 审计功能
为关键表添加 `created_by` 和 `updated_by` 字段。

---

## 联系与支持

如有问题，请参考：
- 使用指南：`docs/AUTH_CONTEXT_USAGE.md`
- 迁移指南：`docs/MIGRATION_GUIDE.md`
- 项目文档：`CLAUDE.md`
