# Auth Context 使用指南

## 概述

本文档说明如何在业务层（Service 层）中使用 Auth 中间件传递的 `loginId` 和 `role` 信息。

## 架构流程

```
请求 → Auth 中间件 → Handler → Service
         ↓
      验证 JWT
      验证 Casbin 权限
      设置 context:
        - loginId
        - role
```

## Context 工具函数

在 `pkg/auth/context.go` 中提供了以下工具函数：

### 获取用户信息

```go
import "back/pkg/auth"

// 获取 loginId（返回值 + bool 表示是否存在）
loginId, ok := auth.GetLoginID(ctx)
if !ok {
    // loginId 不存在（未认证）
}

// 获取 role（返回值 + bool 表示是否存在）
role, ok := auth.GetRole(ctx)
if !ok {
    // role 不存在（未认证）
}

// 必须获取（如果不存在返回空字符串）
loginId := auth.MustGetLoginID(ctx)
role := auth.MustGetRole(ctx)
```

### 设置用户信息（仅在测试或特殊场景使用）

```go
ctx = auth.SetLoginID(ctx, "user123")
ctx = auth.SetRole(ctx, "admin")
```

## 使用示例

### 示例 1: 创建订单时记录创建人

```go
// internal/order/application/order_service.go

func (s *OrderService) Create(ctx context.Context, req *CreateOrderRequest) (*OrderResponse, error) {
    // 从 context 获取当前登录用户
    loginId := auth.MustGetLoginID(ctx)

    // 创建订单，记录创建人
    order := &domain.Order{
        OrderNo:   req.OrderNo,
        ClientID:  req.ClientID,
        CreatedBy: loginId, // 记录创建人
        // ... 其他字段
    }

    if err := s.repo.Save(ctx, order); err != nil {
        return nil, err
    }

    return ToOrderResponse(order), nil
}
```

### 示例 2: 根据角色过滤数据

```go
// internal/order/application/order_service.go

func (s *OrderService) List(ctx context.Context, limit, offset int) (*OrderListResponse, error) {
    // 获取当前用户角色
    role := auth.MustGetRole(ctx)
    loginId := auth.MustGetLoginID(ctx)

    var orders []*domain.Order
    var err error

    // 根据角色决定查询范围
    switch role {
    case "admin", "hr":
        // 管理员可以看到所有订单
        orders, err = s.repo.FindAll(ctx, limit, offset)

    case "sales":
        // 销售只能看到自己创建的订单
        orders, err = s.repo.FindByCreator(ctx, loginId, limit, offset)

    default:
        // 其他角色无权限
        return nil, ErrPermissionDenied
    }

    if err != nil {
        return nil, err
    }

    total, _ := s.repo.Count(ctx)
    return ToOrderListResponse(orders, total), nil
}
```

### 示例 3: 权限校验（业务层）

```go
// internal/order/application/order_service.go

func (s *OrderService) Delete(ctx context.Context, id uint) error {
    // 获取当前用户信息
    loginId := auth.MustGetLoginID(ctx)
    role := auth.MustGetRole(ctx)

    // 查询订单
    order, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return err
    }

    // 业务层权限检查：只有管理员或订单创建人可以删除
    if role != "admin" && order.CreatedBy != loginId {
        return ErrPermissionDenied
    }

    // 删除订单
    return s.repo.Delete(ctx, id)
}
```

### 示例 4: 记录操作日志

```go
// internal/product/application/product_service.go

func (s *ProductService) Update(ctx context.Context, id uint, req *UpdateProductRequest) error {
    loginId := auth.MustGetLoginID(ctx)

    // 查询产品
    product, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return err
    }

    // 更新产品信息
    product.Name = req.Name
    product.UpdatedBy = loginId // 记录更新人

    // 保存
    if err := s.repo.Update(ctx, product); err != nil {
        return err
    }

    // 记录操作日志
    s.logger.Info("Product updated",
        "productId", id,
        "updatedBy", loginId,
    )

    return nil
}
```

## 注意事项

1. **所有 Handler 必须传递 `c.Request.Context()`**
   ```go
   // ✅ 正确
   resp, err := h.service.Create(c.Request.Context(), &req)

   // ❌ 错误
   resp, err := h.service.Create(context.Background(), &req)
   ```

2. **Context 在整个调用链中传递**
   - Handler → Service → Repository
   - 每个函数都应该接收并传递 `ctx context.Context`

3. **分层权限控制**
   - **第 1 层 - JWT**: 验证用户身份
   - **第 2 层 - Casbin**: 检查资源类型级权限（URL 路径）
   - **第 3 层 - 业务层**: 检查数据归属和细粒度权限

4. **未认证请求**
   - 如果请求未通过 Auth 中间件，`loginId` 和 `role` 可能为空
   - 使用 `GetLoginID()` 返回的 `ok` 判断是否存在
   - 或使用 `MustGetLoginID()` 获取空字符串

## 统一角色管理

角色列表在 `internal/user/domain/user.go` 中统一维护：

```go
const (
    UserRoleAdmin     UserRole = "admin"
    UserRoleHR        UserRole = "hr"
    UserRoleSales     UserRole = "sales"
    UserRoleFollower  UserRole = "follower"
    UserRoleAssistant UserRole = "assistant"
    UserRoleUser      UserRole = "user"
)

// 获取所有角色
roles := domain.GetAllRoles()

// 验证角色是否有效
if domain.IsValidRole("sales") {
    // 角色有效
}
```

**添加新角色时，只需修改一处**：
1. 在 `domain/user.go` 中添加新的常量
2. 在 `GetAllRoles()` 函数中添加新角色
3. 所有验证和策略自动生效

## ES 查询权限过滤（未来实现）

未来在 ES 查询中，也可以使用 context 中的用户信息进行权限过滤：

```go
// internal/search/infra/es_search_repo.go

func (r *ESSearchRepo) Search(ctx context.Context, query *domain.SearchQuery) (*domain.SearchResponse, error) {
    // 从 context 获取用户信息
    loginId := auth.MustGetLoginID(ctx)
    role := auth.MustGetRole(ctx)

    // 根据角色添加权限过滤
    esQuery := r.buildESQuery(query, loginId, role)

    // 执行搜索
    return r.executeSearch(ctx, query.GetIndices(), esQuery)
}
```

## 测试

在单元测试中，可以手动设置 context：

```go
func TestCreateOrder(t *testing.T) {
    ctx := context.Background()
    ctx = auth.SetLoginID(ctx, "test_user")
    ctx = auth.SetRole(ctx, "sales")

    service := NewOrderService(mockRepo)
    resp, err := service.Create(ctx, &CreateOrderRequest{
        OrderNo: "ORD001",
    })

    assert.NoError(t, err)
    assert.Equal(t, "test_user", resp.CreatedBy)
}
```
