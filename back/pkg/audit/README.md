# Audit 审计系统

## 概述

Audit 审计系统是一个统一的审计日志中间件，用于记录所有业务操作的审计信息。

## 特性

- ✅ **自动化审计**：中间件自动记录所有非 GET 请求
- ✅ **统一调用点**：在 handler 执行后统一保存审计日志
- ✅ **灵活记录**：业务层可选择性地记录详细的 old/new 数据
- ✅ **零侵入**：即使不设置 old/new，也会记录基本信息（路径、方法、IP 等）
- ✅ **智能推断**：自动从路径推断 domain 和 action

## 架构

```
请求流程：
┌─────────┐     ┌──────────┐     ┌─────────┐     ┌──────────┐
│ Request │ --> │   Auth   │ --> │ Handler │ --> │  Audit   │
│         │     │Middleware│     │         │     │Middleware│
└─────────┘     └──────────┘     └─────────┘     └──────────┘
                    ▲                                  │
                    │                                  │
                 认证鉴权                           统一保存
                    │                              审计日志
                    │                                  │
                 失败则 Abort                          ▼
                 不会到达 Audit                   Database
```

## 数据库表结构

```sql
CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    login_id INTEGER NOT NULL,           -- 操作人ID
    username VARCHAR(100) NOT NULL,       -- 操作人姓名
    domain VARCHAR(50) NOT NULL,          -- 业务域（order, user, product等）
    action VARCHAR(100) NOT NULL,         -- 操作动作（create, update, delete等）
    resource_id VARCHAR(100),             -- 被操作的资源ID
    http_method VARCHAR(10) NOT NULL,     -- HTTP方法
    request_path VARCHAR(500) NOT NULL,   -- 请求路径
    ip_address VARCHAR(50),               -- 客户端IP
    status_code INTEGER NOT NULL,         -- 响应状态码
    error_message TEXT,                   -- 错误信息（如果失败）
    duration_ms BIGINT NOT NULL,          -- 操作耗时（毫秒）
    user_agent VARCHAR(500),              -- 客户端信息
    request_id VARCHAR(100),              -- 请求追踪ID
    old_data JSONB,                       -- 变更前数据（JSON）
    new_data JSONB,                       -- 变更后数据（JSON）
    created_at TIMESTAMP NOT NULL,
    INDEX idx_login_id (login_id),
    INDEX idx_domain (domain),
    INDEX idx_action (action),
    INDEX idx_resource_id (resource_id),
    INDEX idx_request_id (request_id),
    INDEX idx_created_at (created_at)
);
```

## 使用方法

### 1. 在路由注册时定义 Action（推荐 ⭐）

每个业务操作都应该有一个**有意义的 action 名字**，在路由注册时使用 `audit.Mark()` 统一定义：

```go
import "back/pkg/audit"

func RegisterOrderHandlers(rg *gin.RouterGroup, service *application.OrderService) {
    handler := NewOrderHandler(service)

    // 使用 audit.Mark(domain, action) 定义每个操作的名字
    rg.POST("/order",
        audit.Mark("order", "orderCreation"),    // ← 创建订单
        handler.Create)

    rg.POST("/order/:id/assign-department",
        audit.Mark("order", "departmentAssignment"),  // ← 分配部门
        handler.AssignDepartment)

    rg.POST("/order/:id/progress/fabric-input",
        audit.Mark("order", "fabricInputUpdate"),     // ← 胚布投入更新
        handler.UpdateFabricInput)

    // GET 请求不需要标记（会自动跳过审计）
    rg.GET("/order/:id", handler.Get)
}
```

**Action 命名建议**：
- `orderCreation` - 订单创建
- `departmentAssignment` - 部门分配
- `fabricInputUpdate` - 胚布投入更新
- `passwordReset` - 密码重置
- `userActivation` - 用户激活

查看 [USAGE.md](./USAGE.md) 获取更多命名示例和完整指南。

### 2. 基础使用（零配置）

如果不使用 `audit.Mark()`，中间件会自动从路径和 HTTP 方法推断 domain 和 action：

**无需在业务代码中做任何修改**，即可记录以下信息：
- 操作人信息（从 auth context 自动提取）
- 请求信息（路径、方法、IP、User-Agent 等）
- 响应信息（状态码、耗时等）
- 自动推断的 domain 和 action（如 `order` + `create`）

### 3. 高级使用（记录详细的 old/new 数据）

如果需要记录详细的业务数据变更，可以在 handler 中使用 `audit.Recorder`（**不需要再设置 domain 和 action**，只需设置 old/new）：

```go
import "back/pkg/audit"

func (h *OrderHandler) Update(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))

    var req UpdateOrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // ===== 获取 audit recorder =====
    recorder := audit.Get(c)
    if recorder != nil {
        recorder.SetResourceID(id)  // 设置资源ID

        // 获取并记录旧值
        oldOrder, err := h.service.Get(c.Request.Context(), uint(id))
        if err == nil {
            recorder.SetOld(map[string]interface{}{
                "order_no":   oldOrder.OrderNo,
                "quantity":   oldOrder.Quantity,
                "status":     oldOrder.Status,
            })
        }
    }

    // 执行业务逻辑
    newOrder, err := h.service.Update(c.Request.Context(), uint(id), &req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    // ===== 记录新值 =====
    if recorder != nil {
        recorder.SetNew(map[string]interface{}{
            "order_no":   newOrder.OrderNo,
            "quantity":   newOrder.Quantity,
            "status":     newOrder.Status,
        })
    }

    c.JSON(200, newOrder)
    // 中间件会自动调用 recorder.Save() 保存审计日志
}
```

### 4. Recorder API（可选使用）

```go
// 获取 recorder
recorder := audit.Get(c)

// 设置资源ID（推荐）
recorder.SetResourceID(123)
recorder.SetResourceID("order-123")

// 设置变更前数据（可选）
recorder.SetOld(oldData)

// 设置变更后数据（可选）
recorder.SetNew(newData)

// 设置业务域和操作动作（通常不需要，因为在路由注册时已定义）
// recorder.SetDomain(audit.DomainOrder)
// recorder.SetAction("customAction")
```

### 5. 预定义常量（不常用）

```go
// 业务域
audit.DomainOrder      // "order"
audit.DomainUser       // "user"
audit.DomainProduct    // "product"
audit.DomainClient     // "client"
audit.DomainSupplier   // "supplier"
audit.DomainMaterial   // "material"
audit.DomainProcess    // "process"
audit.DomainPricing    // "pricing"
audit.DomainPlan       // "plan"

// 操作动作
audit.ActionCreate     // "create"
audit.ActionUpdate     // "update"
audit.ActionDelete     // "delete"
audit.ActionAssign     // "assign"
```

## 自动推断规则

### Domain 推断

从请求路径的第三部分提取：
- `/api/v1/order/123` → domain: `"order"`
- `/api/v1/user/456` → domain: `"user"`

### Action 推断

根据 HTTP 方法：
- `POST /api/v1/order` → action: `"create"`
- `PUT /api/v1/order/123` → action: `"update"`
- `DELETE /api/v1/order/123` → action: `"delete"`
- `POST /api/v1/order/123/assign` → action: `"assign"`

### Resource ID 推断

从路径中提取数字：
- `/api/v1/order/123` → resource_id: `"123"`

## 示例

### 示例 1：创建订单（零配置）

**路由注册**：
```go
rg.POST("/order", audit.Mark("order", "orderCreation"), handler.Create)
```

**Handler 代码**（不需要任何 audit 代码）：
```go
func (h *OrderHandler) Create(c *gin.Context) {
    var req CreateOrderRequest
    c.ShouldBindJSON(&req)

    order, _ := h.service.Create(c.Request.Context(), &req)
    c.JSON(200, order)
}
```

**自动记录的审计日志**：
```json
{
  "login_id": 1001,
  "username": "张三",
  "domain": "order",
  "action": "orderCreation",
  "http_method": "POST",
  "request_path": "/api/v1/order",
  "ip_address": "192.168.1.100",
  "status_code": 200,
  "duration_ms": 150
}
```

### 示例 2：分配部门（记录详细变更）

**路由注册**：
```go
rg.POST("/order/:id/assign-department",
    audit.Mark("order", "departmentAssignment"),
    handler.AssignDepartment)
```

**Handler 代码**（参考上面的"高级使用"部分）

**记录的审计日志**：
```json
{
  "login_id": 1001,
  "username": "张三",
  "domain": "order",
  "action": "departmentAssignment",
  "resource_id": "123",
  "http_method": "POST",
  "request_path": "/api/v1/order/123/assign-department",
  "ip_address": "192.168.1.100",
  "status_code": 200,
  "duration_ms": 200,
  "old_data": {
    "order_no": "ORD-001",
    "assigned_department": ""
  },
  "new_data": {
    "order_no": "ORD-001",
    "assigned_department": "生产一部"
  }
}
```

## 注意事项

1. **GET 请求自动跳过**：只读操作不记录审计日志
2. **认证失败不记录**：auth 中间件失败时会 Abort，不会到达 audit 中间件
3. **不影响业务流程**：审计日志保存失败不会影响业务操作，只会打印错误日志
4. **性能影响**：审计记录是同步的，但通常耗时很短（< 10ms）

## 查询审计日志示例

```sql
-- 查询某个用户的所有操作
SELECT * FROM audit_logs WHERE login_id = 1001 ORDER BY created_at DESC;

-- 查询某个订单的所有变更历史
SELECT * FROM audit_logs WHERE domain = 'order' AND resource_id = '123' ORDER BY created_at;

-- 查询失败的操作
SELECT * FROM audit_logs WHERE status_code >= 400 ORDER BY created_at DESC;

-- 查询今天的所有删除操作
SELECT * FROM audit_logs
WHERE action = 'delete'
  AND created_at >= CURRENT_DATE
ORDER BY created_at DESC;

-- 分析操作耗时
SELECT domain, action, AVG(duration_ms) as avg_duration
FROM audit_logs
GROUP BY domain, action
ORDER BY avg_duration DESC;
```

## 与 OrderEvent 的区别

| 对比 | AuditLog（系统审计） | OrderEvent（业务审计） |
|------|---------------------|----------------------|
| **粒度** | 粗粒度 API 调用 | 细粒度业务事件 |
| **范围** | 所有业务模块 | 仅订单模块 |
| **记录内容** | 通用请求响应信息 | 详细业务变更 |
| **使用场景** | 合规审计、安全追踪 | 业务流程追踪 |
| **示例** | "张三调用了 POST /order/123/assign" | "张三分配订单到A部门" |

两者可以并存，各司其职。
