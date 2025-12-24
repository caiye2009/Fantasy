# Audit 审计系统 - 使用指南

## 核心概念

每个业务操作都应该有一个**有意义的 action 名字**，例如：
- `orderCreation` - 订单创建
- `departmentAssignment` - 部门分配
- `fabricInputUpdate` - 胚布投入更新
- `userPasswordReset` - 用户密码重置

这些 action 名字在**路由注册时统一定义**，无需在 domain 层定义常量。

## 📝 在路由注册时定义 Action

**推荐方式**：在路由注册时使用 `audit.Mark()` 标记每个操作的 domain 和 action。

### 示例 1：订单模块路由

```go
// internal/order/interfaces/order_handler.go

func RegisterOrderHandlers(rg *gin.RouterGroup, service *application.OrderService) {
    handler := NewOrderHandler(service)

    // 使用 audit.Mark(domain, action) 标记每个路由
    rg.POST("/order",
        audit.Mark("order", "orderCreation"),    // ← 定义 action 名字
        handler.Create)

    rg.POST("/order/:id",
        audit.Mark("order", "orderUpdate"),
        handler.Update)

    rg.DELETE("/order/:id",
        audit.Mark("order", "orderDeletion"),
        handler.Delete)

    rg.POST("/order/:id/assign-department",
        audit.Mark("order", "departmentAssignment"),  // ← 描述性的 action
        handler.AssignDepartment)

    rg.POST("/order/:id/progress/fabric-input",
        audit.Mark("order", "fabricInputUpdate"),     // ← 胚布投入更新
        handler.UpdateFabricInput)

    // GET 请求不需要标记（会自动跳过审计）
    rg.GET("/order/:id", handler.Get)
    rg.GET("/order", handler.List)
}
```

### 示例 2：用户模块路由

```go
func RegisterUserHandlers(rg *gin.RouterGroup, service *application.UserService) {
    handler := NewUserHandler(service)

    rg.POST("/user",
        audit.Mark("user", "userCreation"),
        handler.Create)

    rg.PUT("/user/:id",
        audit.Mark("user", "userUpdate"),
        handler.Update)

    rg.POST("/user/:id/reset-password",
        audit.Mark("user", "passwordReset"),     // ← 密码重置
        handler.ResetPassword)

    rg.POST("/user/:id/activate",
        audit.Mark("user", "userActivation"),    // ← 用户激活
        handler.Activate)

    rg.POST("/user/:id/deactivate",
        audit.Mark("user", "userDeactivation"),  // ← 用户停用
        handler.Deactivate)
}
```

## 🎯 在 Handler 中使用 Recorder

在 handler 中，你**不需要再设置 domain 和 action**（因为已经在路由注册时定义了），只需要：
1. 设置 `resourceID`（可选）
2. 记录 `old` 和 `new` 数据（可选）

### 完整示例：分配部门

```go
func (h *OrderHandler) AssignDepartment(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))

    var req AssignDepartmentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // ===== Audit: 记录旧值 =====
    recorder := audit.Get(c)
    if recorder != nil {
        recorder.SetResourceID(id)  // 设置被操作的资源ID

        // 获取旧值
        oldOrder, err := h.service.Get(c.Request.Context(), uint(id))
        if err == nil {
            recorder.SetOld(map[string]interface{}{
                "order_no":            oldOrder.OrderNo,
                "assigned_department": oldOrder.AssignedDepartment,  // 旧部门
            })
        }
    }
    // ===== Audit End =====

    // 执行业务逻辑
    if err := h.service.AssignDepartment(c.Request.Context(), uint(id), &req, ...); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    // ===== Audit: 记录新值 =====
    if recorder != nil {
        recorder.SetNew(map[string]interface{}{
            "order_no":            oldOrder.OrderNo,
            "assigned_department": req.Department,  // 新部门
        })
    }
    // ===== Audit End =====

    c.JSON(200, gin.H{"message": "分配部门成功"})
    // 中间件会自动调用 recorder.Save() 保存到数据库
}
```

### 简化示例：不需要记录详细数据

如果不需要记录详细的 old/new 数据，甚至可以完全不写 audit 代码：

```go
func (h *OrderHandler) Delete(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))

    // 直接执行业务逻辑，无需任何 audit 代码
    if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"message": "删除成功"})
    // 审计日志会自动记录：
    // - action: "orderDeletion"（来自路由定义）
    // - resource_id: "123"（自动从路径提取）
    // - 其他基本信息（用户、IP、耗时等）
}
```

## 📊 Action 命名规范

建议使用**驼峰命名**，清晰描述操作：

### 基础操作
- `{资源}Creation` - 创建（如 `orderCreation`, `userCreation`）
- `{资源}Update` - 更新（如 `orderUpdate`, `productUpdate`）
- `{资源}Deletion` - 删除（如 `orderDeletion`, `clientDeletion`）

### 业务操作
- `{对象}Assignment` - 分配（如 `departmentAssignment`, `personnelAssignment`）
- `{字段}Update` - 字段更新（如 `fabricInputUpdate`, `statusUpdate`）
- `{对象}Activation` - 激活（如 `userActivation`）
- `{对象}Deactivation` - 停用（如 `userDeactivation`）
- `{操作}Reset` - 重置（如 `passwordReset`）
- `{对象}Addition` - 添加（如 `defectAddition`）

### 中文对照示例

```go
// 订单相关
orderCreation           // 创建订单
orderUpdate             // 更新订单
orderDeletion           // 删除订单
departmentAssignment    // 分配部门
personnelAssignment     // 分配人员
fabricInputUpdate       // 胚布投入更新
productionUpdate        // 生产进度更新
warehouseCheckUpdate    // 验货进度更新
reworkUpdate            // 回修进度更新
defectAddition          // 录入次品

// 用户相关
userCreation            // 创建用户
userUpdate              // 更新用户
userDeletion            // 删除用户
passwordReset           // 密码重置
userActivation          // 激活用户
userDeactivation        // 停用用户
roleAssignment          // 分配角色

// 产品相关
productCreation         // 创建产品
productUpdate           // 更新产品
priceAdjustment         // 价格调整
stockUpdate             // 库存更新
```

## 🔄 迁移指南

### 从 OrderEvent 迁移到 Audit

之前你可能在 `order/domain/event.go` 中定义了这些常量：

```go
// ❌ 旧方式：在 domain 层定义常量
const (
    EventTypeCreateOrder      = "create_order"
    EventTypeAssignDepartment = "assign_department"
    EventTypeUpdateFabricInput = "update_fabric_input"
)
```

**现在不需要这些常量了**，直接在路由注册时定义：

```go
// ✅ 新方式：在路由注册时定义
rg.POST("/order", audit.Mark("order", "orderCreation"), handler.Create)
rg.POST("/order/:id/assign-department", audit.Mark("order", "departmentAssignment"), handler.AssignDepartment)
rg.POST("/order/:id/progress/fabric-input", audit.Mark("order", "fabricInputUpdate"), handler.UpdateFabricInput)
```

### OrderEvent vs AuditLog

两者可以**共存**，各有用途：

| 对比 | AuditLog（系统审计） | OrderEvent（业务审计） |
|------|---------------------|----------------------|
| **适用范围** | 所有业务模块 | 仅订单模块 |
| **粒度** | API 级别 | 业务事件级别 |
| **记录内容** | HTTP 请求/响应信息 | 详细业务变更 |
| **使用场景** | 合规审计、安全追踪 | 业务流程追踪、协作历史 |
| **示例** | "张三调用了 POST /order/123/assign" | "张三分配订单ORD-001到A部门" |

**建议**：保留 OrderEvent 用于详细的业务追踪，同时使用 AuditLog 满足合规审计需求。

## 🎨 完整的模块示例

```go
// internal/product/interfaces/product_handler.go

import "back/pkg/audit"

func RegisterProductHandlers(rg *gin.RouterGroup, service *application.ProductService) {
    handler := NewProductHandler(service)

    // 基础 CRUD
    rg.POST("/product", audit.Mark("product", "productCreation"), handler.Create)
    rg.PUT("/product/:id", audit.Mark("product", "productUpdate"), handler.Update)
    rg.DELETE("/product/:id", audit.Mark("product", "productDeletion"), handler.Delete)

    // 业务操作
    rg.POST("/product/:id/price", audit.Mark("product", "priceAdjustment"), handler.AdjustPrice)
    rg.POST("/product/:id/stock", audit.Mark("product", "stockUpdate"), handler.UpdateStock)
    rg.POST("/product/:id/publish", audit.Mark("product", "productPublication"), handler.Publish)
    rg.POST("/product/:id/unpublish", audit.Mark("product", "productUnpublication"), handler.Unpublish)

    // 查询操作（GET 自动跳过审计）
    rg.GET("/product", handler.List)
    rg.GET("/product/:id", handler.Get)
}
```

## ✅ 优势总结

1. **集中管理**：所有 action 定义在路由注册处，一目了然
2. **无需常量**：不需要在 domain 层定义 `EventType*` 常量
3. **清晰可读**：action 名字直观描述操作，如 `fabricInputUpdate`（胚布投入更新）
4. **易于维护**：新增操作时只需在路由注册时添加一行
5. **零侵入**：handler 代码更简洁，可选地记录详细数据

## 🔍 查询审计日志示例

```sql
-- 查询所有订单创建操作
SELECT * FROM audit_logs
WHERE domain = 'order' AND action = 'orderCreation'
ORDER BY created_at DESC;

-- 查询用户张三今天的所有操作
SELECT * FROM audit_logs
WHERE username = '张三' AND created_at >= CURRENT_DATE
ORDER BY created_at DESC;

-- 查询订单123的所有变更历史
SELECT * FROM audit_logs
WHERE domain = 'order' AND resource_id = '123'
ORDER BY created_at;

-- 统计各类操作的数量
SELECT action, COUNT(*) as count
FROM audit_logs
WHERE domain = 'order'
GROUP BY action
ORDER BY count DESC;
```
