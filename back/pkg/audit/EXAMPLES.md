# Audit 审计系统 - 其他模块应用示例

## 用户模块 (User)

### 路由注册

```go
// internal/user/interfaces/user_handler.go

import "back/pkg/audit"

func RegisterUserHandlers(rg *gin.RouterGroup, service *application.UserService) {
    handler := NewUserHandler(service)

    // 用户 CRUD
    rg.POST("/user",
        audit.Mark("user", "userCreation"),
        handler.Create)

    rg.PUT("/user/:id",
        audit.Mark("user", "userUpdate"),
        handler.Update)

    rg.DELETE("/user/:id",
        audit.Mark("user", "userDeletion"),
        handler.Delete)

    // 密码相关
    rg.POST("/user/:id/reset-password",
        audit.Mark("user", "passwordReset"),
        handler.ResetPassword)

    rg.POST("/user/:id/change-password",
        audit.Mark("user", "passwordChange"),
        handler.ChangePassword)

    // 用户状态
    rg.POST("/user/:id/activate",
        audit.Mark("user", "userActivation"),
        handler.Activate)

    rg.POST("/user/:id/deactivate",
        audit.Mark("user", "userDeactivation"),
        handler.Deactivate)

    // 角色分配
    rg.POST("/user/:id/assign-role",
        audit.Mark("user", "roleAssignment"),
        handler.AssignRole)

    // 查询操作（GET，自动跳过审计）
    rg.GET("/user", handler.List)
    rg.GET("/user/:id", handler.Get)
}
```

### Handler 示例：密码重置

```go
func (h *UserHandler) ResetPassword(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))

    var req ResetPasswordRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // Audit: 记录操作（不记录密码内容）
    recorder := audit.Get(c)
    if recorder != nil {
        recorder.SetResourceID(id)
        // 注意：不要记录密码到审计日志
        recorder.SetOld(map[string]interface{}{
            "user_id": id,
            "action":  "password_reset_requested",
        })
    }

    if err := h.service.ResetPassword(c.Request.Context(), uint(id)); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    if recorder != nil {
        recorder.SetNew(map[string]interface{}{
            "user_id": id,
            "action":  "password_reset_completed",
        })
    }

    c.JSON(200, gin.H{"message": "密码重置成功"})
}
```

## 产品模块 (Product)

### 路由注册

```go
// internal/product/interfaces/product_handler.go

import "back/pkg/audit"

func RegisterProductHandlers(rg *gin.RouterGroup, service *application.ProductService) {
    handler := NewProductHandler(service)

    // 产品 CRUD
    rg.POST("/product",
        audit.Mark("product", "productCreation"),
        handler.Create)

    rg.PUT("/product/:id",
        audit.Mark("product", "productUpdate"),
        handler.Update)

    rg.DELETE("/product/:id",
        audit.Mark("product", "productDeletion"),
        handler.Delete)

    // 价格管理
    rg.POST("/product/:id/price",
        audit.Mark("product", "priceAdjustment"),
        handler.AdjustPrice)

    // 库存管理
    rg.POST("/product/:id/stock",
        audit.Mark("product", "stockUpdate"),
        handler.UpdateStock)

    // 发布/下架
    rg.POST("/product/:id/publish",
        audit.Mark("product", "productPublication"),
        handler.Publish)

    rg.POST("/product/:id/unpublish",
        audit.Mark("product", "productUnpublication"),
        handler.Unpublish)

    // 查询操作
    rg.GET("/product", handler.List)
    rg.GET("/product/:id", handler.Get)
}
```

### Handler 示例：价格调整

```go
func (h *ProductHandler) AdjustPrice(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))

    var req AdjustPriceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // Audit: 记录价格变更
    recorder := audit.Get(c)
    if recorder != nil {
        recorder.SetResourceID(id)

        product, err := h.service.Get(c.Request.Context(), uint(id))
        if err == nil {
            recorder.SetOld(map[string]interface{}{
                "product_code": product.Code,
                "product_name": product.Name,
                "old_price":    product.Price,
            })
        }
    }

    if err := h.service.AdjustPrice(c.Request.Context(), uint(id), &req); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    if recorder != nil {
        recorder.SetNew(map[string]interface{}{
            "new_price": req.NewPrice,
            "reason":    req.Reason,
        })
    }

    c.JSON(200, gin.H{"message": "价格调整成功"})
}
```

## 客户模块 (Client)

### 路由注册

```go
// internal/client/interfaces/client_handler.go

import "back/pkg/audit"

func RegisterClientHandlers(rg *gin.RouterGroup, service *application.ClientService) {
    handler := NewClientHandler(service)

    // 客户 CRUD
    rg.POST("/client",
        audit.Mark("client", "clientCreation"),
        handler.Create)

    rg.PUT("/client/:id",
        audit.Mark("client", "clientUpdate"),
        handler.Update)

    rg.DELETE("/client/:id",
        audit.Mark("client", "clientDeletion"),
        handler.Delete)

    // 信用额度管理
    rg.POST("/client/:id/credit-limit",
        audit.Mark("client", "creditLimitAdjustment"),
        handler.AdjustCreditLimit)

    // 客户状态
    rg.POST("/client/:id/suspend",
        audit.Mark("client", "clientSuspension"),
        handler.Suspend)

    rg.POST("/client/:id/reactivate",
        audit.Mark("client", "clientReactivation"),
        handler.Reactivate)

    // 查询操作
    rg.GET("/client", handler.List)
    rg.GET("/client/:id", handler.Get)
}
```

## 供应商模块 (Supplier)

### 路由注册

```go
// internal/supplier/interfaces/supplier_handler.go

import "back/pkg/audit"

func RegisterSupplierHandlers(rg *gin.RouterGroup, service *application.SupplierService) {
    handler := NewSupplierHandler(service)

    rg.POST("/supplier",
        audit.Mark("supplier", "supplierCreation"),
        handler.Create)

    rg.PUT("/supplier/:id",
        audit.Mark("supplier", "supplierUpdate"),
        handler.Update)

    rg.DELETE("/supplier/:id",
        audit.Mark("supplier", "supplierDeletion"),
        handler.Delete)

    rg.POST("/supplier/:id/approve",
        audit.Mark("supplier", "supplierApproval"),
        handler.Approve)

    rg.POST("/supplier/:id/reject",
        audit.Mark("supplier", "supplierRejection"),
        handler.Reject)

    rg.GET("/supplier", handler.List)
    rg.GET("/supplier/:id", handler.Get)
}
```

## 计划模块 (Plan)

### 路由注册

```go
// internal/plan/interfaces/plan_handler.go

import "back/pkg/audit"

func RegisterPlanHandlers(rg *gin.RouterGroup, service *application.PlanService) {
    handler := NewPlanHandler(service)

    rg.POST("/plan",
        audit.Mark("plan", "planCreation"),
        handler.Create)

    rg.PUT("/plan/:id",
        audit.Mark("plan", "planUpdate"),
        handler.Update)

    rg.DELETE("/plan/:id",
        audit.Mark("plan", "planDeletion"),
        handler.Delete)

    rg.POST("/plan/:id/approve",
        audit.Mark("plan", "planApproval"),
        handler.Approve)

    rg.POST("/plan/:id/start",
        audit.Mark("plan", "planExecution"),
        handler.Start)

    rg.POST("/plan/:id/complete",
        audit.Mark("plan", "planCompletion"),
        handler.Complete)

    rg.POST("/plan/:id/cancel",
        audit.Mark("plan", "planCancellation"),
        handler.Cancel)

    rg.GET("/plan", handler.List)
    rg.GET("/plan/:id", handler.Get)
}
```

## Action 命名参考

### 命名模式

```
{资源}{操作}
```

### 常用操作词汇

| 英文 | 中文 | 示例 |
|------|------|------|
| Creation | 创建 | userCreation, productCreation |
| Update | 更新 | userUpdate, stockUpdate |
| Deletion | 删除 | orderDeletion, clientDeletion |
| Assignment | 分配 | roleAssignment, departmentAssignment |
| Adjustment | 调整 | priceAdjustment, creditLimitAdjustment |
| Activation | 激活 | userActivation, clientReactivation |
| Deactivation | 停用 | userDeactivation, productUnpublication |
| Approval | 审批 | planApproval, supplierApproval |
| Rejection | 拒绝 | supplierRejection, orderRejection |
| Suspension | 暂停 | clientSuspension, planSuspension |
| Cancellation | 取消 | orderCancellation, planCancellation |
| Completion | 完成 | orderCompletion, planCompletion |
| Execution | 执行 | planExecution, taskExecution |
| Publication | 发布 | productPublication, articlePublication |
| Reset | 重置 | passwordReset, settingsReset |
| Addition | 添加 | defectAddition, itemAddition |
| Removal | 移除 | itemRemoval, memberRemoval |

## 特殊场景处理

### 1. 批量操作

```go
rg.POST("/product/batch-delete",
    audit.Mark("product", "productBatchDeletion"),
    handler.BatchDelete)

func (h *ProductHandler) BatchDelete(c *gin.Context) {
    var req BatchDeleteRequest
    c.ShouldBindJSON(&req)

    recorder := audit.Get(c)
    if recorder != nil {
        recorder.SetResourceID(strings.Join(req.IDs, ","))
        recorder.SetNew(map[string]interface{}{
            "deleted_count": len(req.IDs),
            "product_ids":   req.IDs,
        })
    }

    // ... 执行批量删除
}
```

### 2. 导入导出

```go
rg.POST("/product/import",
    audit.Mark("product", "productDataImport"),
    handler.Import)

rg.POST("/product/export",
    audit.Mark("product", "productDataExport"),
    handler.Export)
```

### 3. 敏感操作（不记录敏感数据）

```go
func (h *UserHandler) ChangePassword(c *gin.Context) {
    // ...

    recorder := audit.Get(c)
    if recorder != nil {
        recorder.SetResourceID(id)
        // ⚠️ 不要记录密码
        recorder.SetNew(map[string]interface{}{
            "action": "password_changed",
            "timestamp": time.Now(),
        })
    }

    // ...
}
```

### 4. 跳过某些路由的审计

```go
// 内部健康检查接口，不需要审计
rg.GET("/health", audit.Skip(), handler.HealthCheck)
```

## 查询审计日志的常用场景

### 1. 查看用户的所有操作

```sql
SELECT
    created_at,
    action,
    domain,
    resource_id,
    status_code
FROM audit_logs
WHERE username = '张三'
ORDER BY created_at DESC
LIMIT 100;
```

### 2. 追踪特定资源的变更历史

```sql
SELECT
    created_at,
    username,
    action,
    old_data,
    new_data
FROM audit_logs
WHERE domain = 'product' AND resource_id = '123'
ORDER BY created_at;
```

### 3. 统计各模块的操作频率

```sql
SELECT
    domain,
    action,
    COUNT(*) as operation_count
FROM audit_logs
WHERE created_at >= NOW() - INTERVAL '30 days'
GROUP BY domain, action
ORDER BY operation_count DESC;
```

### 4. 查找失败的操作

```sql
SELECT
    created_at,
    username,
    domain,
    action,
    request_path,
    status_code,
    error_message
FROM audit_logs
WHERE status_code >= 400
ORDER BY created_at DESC;
```

### 5. 审计敏感操作

```sql
SELECT
    created_at,
    username,
    action,
    resource_id,
    ip_address
FROM audit_logs
WHERE action IN (
    'passwordReset',
    'userDeletion',
    'roleAssignment',
    'creditLimitAdjustment'
)
ORDER BY created_at DESC;
```
