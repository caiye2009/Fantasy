# 后端接口完整总结（含对象字段）

## 统一响应格式
所有接口返回数据统一包裹在以下格式中:
```json
{
  "code": 0,           // 0=成功, 1=失败
  "msg": "提示信息",    // 错误时返回
  "data": { ... }      // 成功时返回的数据
}
```

---

## 1. auth 包（认证模块）

### 1.1 POST /api/v1/auth/login - 用户登录
**请求体 (LoginRequest)**:
```json
{
  "loginId": "string",      // 必填, 登录ID
  "password": "string"      // 必填, 密码, 最小3位
}
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "accessToken": "string",           // 访问令牌
    "refreshToken": "string",          // 刷新令牌
    "username": "string",              // 用户名
    "role": "string",                  // 角色
    "requirePasswordChange": true      // 是否需要修改密码
  }
}
```

**逻辑**:
- 验证loginID和密码
- 检查用户状态是否为active
- 生成JWT访问令牌（包含loginID和role）和刷新令牌（仅含loginID）
- 返回用户基本信息和令牌

---

### 1.2 POST /api/v1/auth/logout - 退出登录
**成功响应**:
```json
{
  "code": 0,
  "data": {
    "message": "退出成功"
  }
}
```

**逻辑**: 前端清除令牌即可

---

### 1.3 POST /api/v1/auth/refresh - 刷新令牌
**请求体 (RefreshTokenRequest)**:
```json
{
  "refreshToken": "string"  // 必填
}
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "accessToken": "string"  // 新的访问令牌
  }
}
```

**逻辑**: 验证刷新令牌有效性,生成新的访问令牌

---

## 2. user 包（用户管理模块）

### 域模型 User 字段:
```go
ID           uint       // 主键
LoginID      string     // 登录ID, 唯一, 4-20位
Username     string     // 用户名, 2-50位
Department   string     // 部门, 最长100位
PasswordHash string     // 密码哈希
Email        string     // 邮箱, 最长100位
Role         string     // 角色: admin/hr/sales/follower/assistant/user
Status       string     // 状态: active/inactive/suspended
HasInitPass  bool       // 是否为初始密码
CreatedAt    time.Time  // 创建时间
UpdatedAt    time.Time  // 更新时间
```

---

### 2.1 POST /api/v1/user - 创建用户
**请求体 (CreateUserRequest)**:
```json
{
  "username": "string",    // 必填, 2-50位
  "department": "string",  // 必填, 最长100位
  "email": "string",       // 可选, email格式
  "role": "string"         // 必填, 角色枚举值
}
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "login_id": "1000",    // 自动生成的登录ID
    "password": "123"      // 默认密码
  }
}
```

**逻辑**:
- 自动生成loginID（从1000开始自增）
- 设置默认密码"123"
- 验证角色必须是预定义枚举值之一
- 检查loginID唯一性

---

### 2.2 GET /api/v1/user/:id - 查询用户详情
**成功响应**:
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "login_id": "1000",
    "username": "张三",
    "department": "销售部",
    "email": "zhangsan@example.com",
    "role": "sales",
    "status": "active",
    "has_init_pass": false,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### 2.3 GET /api/v1/user - 查询用户列表
**查询参数**:
- `limit`: 每页数量 (默认10)
- `offset`: 偏移量 (默认0)

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "total": 100,
    "users": [
      {
        "id": 1,
        "login_id": "1000",
        "username": "张三",
        "department": "销售部",
        "email": "zhangsan@example.com",
        "role": "sales",
        "status": "active",
        "has_init_pass": false,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

### 2.4 PUT /api/v1/user/:id - 更新用户
**请求体 (UpdateUserRequest)**:
```json
{
  "username": "string",    // 可选, 2-50位
  "department": "string",  // 可选, 最长100位
  "email": "string",       // 可选, email格式
  "role": "string",        // 可选, 角色枚举值
  "status": "string"       // 可选: active/inactive/suspended
}
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "login_id": "1000",
    "username": "张三",
    ...
  }
}
```

**逻辑**: 不能修改admin用户

---

### 2.5 DELETE /api/v1/user/:id - 删除用户
**成功响应**:
```json
{
  "code": 0,
  "data": {
    "message": "删除成功"
  }
}
```

**逻辑**: 禁止删除admin角色的用户

---

### 2.6 POST /api/v1/user/change-password - 修改密码
**请求体 (ChangePasswordRequest)**:
```json
{
  "current_password": "string",  // 非初始密码时必填
  "new_password": "string"       // 必填, 最小6位
}
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "message": "密码修改成功"
  }
}
```

**逻辑**:
- 从context获取当前登录用户的login_id
- 如果是初始密码(has_init_pass=true),无需验证旧密码
- 新密码至少6位
- 修改后has_init_pass设为false

---

### 2.7 GET /api/v1/user/departments - 获取部门列表
**成功响应**:
```json
{
  "code": 0,
  "data": ["销售部", "技术部", "财务部"]
}
```

---

### 2.8 GET /api/v1/user/roles - 获取角色列表
**成功响应**:
```json
{
  "code": 0,
  "data": ["admin", "hr", "sales", "follower", "assistant", "user"]
}
```

---

## 3. department 包（部门管理模块）

### 域模型 Department 字段:
```go
ID          uint        // 主键
Name        string      // 部门名称, 2-100位
Code        string      // 部门编码, 唯一, 2-50位
Description string      // 描述, 最长500位
Status      string      // 状态: active/inactive
ParentID    *uint       // 父部门ID, 可空
CreatedAt   time.Time   // 创建时间
UpdatedAt   time.Time   // 更新时间
```

---

### 3.1 POST /api/v1/department - 创建部门
**请求体 (CreateDepartmentRequest)**:
```json
{
  "name": "string",        // 必填, 2-100位
  "code": "string",        // 可选, 2-50位
  "description": "string", // 可选
  "parent_id": 1           // 可选, 父部门ID
}
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "name": "销售部",
    "code": "SALES",
    "description": "负责产品销售",
    "status": "active",
    "parent_id": null,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**逻辑**: 部门编码必须唯一

---

### 3.2 GET /api/v1/department/:id - 查询部门详情
**成功响应**: 同创建部门响应

---

### 3.3 GET /api/v1/department - 查询部门列表
**查询参数**:
- `status`: active/inactive (可选)
- `page`: 页码 (默认1)
- `page_size`: 每页数量 (默认10)

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "total": 50,
    "departments": [
      {
        "id": 1,
        "name": "销售部",
        "code": "SALES",
        "description": "负责产品销售",
        "status": "active",
        "parent_id": null,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

### 3.4 PUT /api/v1/department/:id - 更新部门
**请求体 (UpdateDepartmentRequest)**:
```json
{
  "name": "string",        // 可选, 2-100位
  "code": "string",        // 可选, 2-50位
  "description": "string", // 可选
  "parent_id": 1           // 可选
}
```

**成功响应**: 返回更新后的部门对象

---

### 3.5 PUT /api/v1/department/:id/deactivate - 停用部门
**成功响应**:
```json
{
  "code": 0,
  "data": {
    "message": "部门已停用"
  }
}
```

**逻辑**: 软删除,将status设为inactive

---

### 3.6 PUT /api/v1/department/:id/activate - 激活部门
**成功响应**:
```json
{
  "code": 0,
  "data": {
    "message": "部门已激活"
  }
}
```

---

## 4. role 包（职位/角色管理模块）

### 域模型 Role 字段:
```go
ID          uint        // 主键
Name        string      // 职位名称, 2-100位
Code        string      // 职位编码, 唯一, 2-50位, 必填
Description string      // 描述, 最长500位
Status      string      // 状态: active/inactive
Level       int         // 职级, 默认0
CreatedAt   time.Time   // 创建时间
UpdatedAt   time.Time   // 更新时间
```

---

### 4.1 POST /api/v1/role - 创建职位
**请求体 (CreateRoleRequest)**:
```json
{
  "name": "string",        // 必填, 2-100位
  "code": "string",        // 必填, 2-50位
  "description": "string", // 可选
  "level": 5               // 可选, 最小0
}
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "name": "高级工程师",
    "code": "SENIOR_ENG",
    "description": "高级技术岗位",
    "status": "active",
    "level": 5,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**逻辑**: 职位编码必须唯一

---

### 4.2 GET /api/v1/role/:id - 查询职位详情
**成功响应**: 同创建职位响应

---

### 4.3 GET /api/v1/role - 查询职位列表
**查询参数**:
- `status`: active/inactive (可选)
- `page`: 页码 (默认1)
- `page_size`: 每页数量 (默认10)

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "total": 20,
    "roles": [
      {
        "id": 1,
        "name": "高级工程师",
        "code": "SENIOR_ENG",
        "description": "高级技术岗位",
        "status": "active",
        "level": 5,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

### 4.4 PUT /api/v1/role/:id - 更新职位
**请求体 (UpdateRoleRequest)**:
```json
{
  "name": "string",        // 可选, 2-100位
  "code": "string",        // 可选, 2-50位
  "description": "string", // 可选
  "level": 6               // 可选, 最小0
}
```

---

### 4.5 PUT /api/v1/role/:id/deactivate - 停用职位
**成功响应**:
```json
{
  "code": 0,
  "data": {
    "message": "职位已停用"
  }
}
```

---

### 4.6 PUT /api/v1/role/:id/activate - 激活职位
**成功响应**:
```json
{
  "code": 0,
  "data": {
    "message": "职位已激活"
  }
}
```

---

## 5. client 包（客户管理模块）

### 域模型 Client 字段:
```go
ID        uint        // 主键
Code      string      // 客户编号, 最长50位
Name      string      // 客户名称, 2-100位, 必填
Contact   string      // 联系人, 最长50位
Phone     string      // 电话, 最长20位, 支持手机和固话格式
Email     string      // 邮箱, 最长100位, email格式
Address   string      // 地址, 最长200位
CreatedAt time.Time   // 创建时间
UpdatedAt time.Time   // 更新时间
```

---

### 5.1 POST /api/v1/client - 创建客户
**请求体 (CreateClientRequest)**:
```json
{
  "name": "string",    // 必填, 2-100位
  "contact": "string", // 可选, 最长50位
  "phone": "string",   // 可选, 最长20位
  "email": "string",   // 可选, email格式
  "address": "string"  // 可选, 最长200位
}
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "name": "ABC公司",
    "contact": "李经理",
    "phone": "13800138000",
    "email": "li@abc.com",
    "address": "北京市朝阳区XXX",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**逻辑**:
- 客户名称必须唯一
- 自动同步到Elasticsearch

---

### 5.2 GET /api/v1/client/:id - 查询客户详情
**成功响应**: 同创建客户响应

---

### 5.3 PUT /api/v1/client/:id - 更新客户
**请求体 (UpdateClientRequest)**:
```json
{
  "name": "string",    // 可选, 2-100位
  "contact": "string", // 可选, 最长50位
  "phone": "string",   // 可选, 最长20位
  "email": "string",   // 可选, email格式
  "address": "string"  // 可选, 最长200位
}
```

**逻辑**: 更新后自动同步ES

---

### 5.4 DELETE /api/v1/client/:id - 删除客户
**成功响应**:
```json
{
  "code": 0,
  "data": {
    "message": "删除成功"
  }
}
```

**逻辑**: 同时删除ES索引

---

## 6. supplier 包（供应商管理模块）

### 域模型 Supplier 字段:
```go
ID        uint        // 主键
Name      string      // 供应商名称, 2-100位, 必填
Contact   string      // 联系人, 最长50位
Phone     string      // 电话, 最长20位
Email     string      // 邮箱, 最长100位
Address   string      // 地址, 最长200位
CreatedAt time.Time   // 创建时间
UpdatedAt time.Time   // 更新时间
```

---

### 6.1 POST /api/v1/supplier - 创建供应商
**请求体 (CreateSupplierRequest)**:
```json
{
  "name": "string",    // 必填, 2-100位
  "contact": "string", // 可选, 最长50位
  "phone": "string",   // 可选, 最长20位
  "email": "string",   // 可选, email格式
  "address": "string"  // 可选, 最长200位
}
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "name": "XYZ材料供应商",
    "contact": "王经理",
    "phone": "13900139000",
    "email": "wang@xyz.com",
    "address": "上海市浦东新区XXX",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**逻辑**: 供应商名称唯一, 自动同步ES

---

### 6.2 GET /api/v1/supplier/:id - 查询供应商详情
**成功响应**: 同创建供应商响应

---

### 6.3 GET /api/v1/supplier - 查询供应商列表
**查询参数**:
- `limit`: 每页数量 (默认10)
- `offset`: 偏移量 (默认0)

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "total": 30,
    "suppliers": [
      {
        "id": 1,
        "name": "XYZ材料供应商",
        "contact": "王经理",
        "phone": "13900139000",
        "email": "wang@xyz.com",
        "address": "上海市浦东新区XXX",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

### 6.4 PUT /api/v1/supplier/:id - 更新供应商
**请求体 (UpdateSupplierRequest)**:
```json
{
  "name": "string",    // 可选, 2-100位
  "contact": "string", // 可选
  "phone": "string",   // 可选
  "email": "string",   // 可选
  "address": "string"  // 可选
}
```

---

### 6.5 DELETE /api/v1/supplier/:id - 删除供应商
**成功响应**:
```json
{
  "code": 0,
  "data": {
    "message": "删除成功"
  }
}
```

---

## 7. material 包（物料管理模块）

### 域模型 Material 字段:
```go
ID          uint        // 主键
Name        string      // 物料名称, 2-100位, 必填
Spec        string      // 规格, 最长200位
Unit        string      // 单位, 最长20位
Description string      // 描述, text类型
CreatedAt   time.Time   // 创建时间
UpdatedAt   time.Time   // 更新时间
```

---

### 7.1 POST /api/v1/material - 创建物料
**请求体 (CreateMaterialRequest)**:
```json
{
  "name": "string",        // 必填, 2-100位
  "spec": "string",        // 可选, 最长200位
  "unit": "string",        // 可选, 最长20位
  "description": "string"  // 可选
}
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "name": "铝型材",
    "spec": "6063-T5",
    "unit": "米",
    "description": "常用铝型材规格",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**逻辑**: 自动同步到ES

---

### 7.2 GET /api/v1/material/:id - 查询物料详情
**成功响应**: 同创建物料响应

---

### 7.3 PUT /api/v1/material/:id - 更新物料
**请求体 (UpdateMaterialRequest)**:
```json
{
  "name": "string",        // 可选, 2-100位
  "spec": "string",        // 可选, 最长200位
  "unit": "string",        // 可选, 最长20位
  "description": "string"  // 可选
}
```

---

### 7.4 DELETE /api/v1/material/:id - 删除物料
**成功响应**:
```json
{
  "code": 0,
  "data": {
    "message": "删除成功"
  }
}
```

**注**: 物料列表查询通过 `/api/v1/search` 接口实现

---

## 8. process 包（工艺管理模块）

### 域模型 Process 字段:
```go
ID          uint        // 主键
Name        string      // 工艺名称, 2-100位, 必填
Description string      // 描述, text类型
CreatedAt   time.Time   // 创建时间
UpdatedAt   time.Time   // 更新时间
```

---

### 8.1 POST /api/v1/process - 创建工艺
**请求体 (CreateProcessRequest)**:
```json
{
  "name": "string",        // 必填, 2-100位
  "description": "string"  // 可选
}
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "name": "阳极氧化",
    "description": "表面处理工艺",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**逻辑**: 自动同步ES

---

### 8.2 GET /api/v1/process/:id - 查询工艺详情
**成功响应**: 同创建工艺响应

---

### 8.3 PUT /api/v1/process/:id - 更新工艺
**请求体 (UpdateProcessRequest)**:
```json
{
  "name": "string",        // 可选, 2-100位
  "description": "string"  // 可选
}
```

---

### 8.4 DELETE /api/v1/process/:id - 删除工艺
**成功响应**:
```json
{
  "code": 0,
  "data": {
    "message": "删除成功"
  }
}
```

**注**: 工艺列表查询通过 `/api/v1/search` 接口

---

## 9. product 包（产品管理模块）

### 域模型 Product 字段:
```go
ID        uint                // 主键
Name      string              // 产品名称, 2-100位, 必填
Status    string              // 状态: draft/submitted/approved/rejected
Materials []MaterialConfig    // 物料配置列表, JSONB存储
Processes []ProcessConfig     // 工艺配置列表, JSONB存储
CreatedAt time.Time           // 创建时间
UpdatedAt time.Time           // 更新时间

// MaterialConfig 结构
MaterialConfig {
  material_id  uint     // 物料ID
  ratio        float64  // 占比(总和必须为1)
}

// ProcessConfig 结构
ProcessConfig {
  process_id   uint     // 工艺ID
  quantity     float64  // 数量(可选)
}
```

---

### 9.1 POST /api/v1/product - 创建产品
**请求体 (CreateProductRequest)**:
```json
{
  "name": "string",           // 必填, 2-100位
  "materials": [              // 必填
    {
      "material_id": 1,
      "ratio": 0.6            // 占比, 所有物料ratio总和必须为1
    },
    {
      "material_id": 2,
      "ratio": 0.4
    }
  ],
  "processes": [              // 必填
    {
      "process_id": 1,
      "quantity": 100
    }
  ]
}
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "name": "铝合金门窗",
    "status": "draft",
    "materials": [
      {
        "material_id": 1,
        "ratio": 0.6
      },
      {
        "material_id": 2,
        "ratio": 0.4
      }
    ],
    "processes": [
      {
        "process_id": 1,
        "quantity": 100
      }
    ],
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**逻辑**:
- 验证所有物料ratio总和必须等于1
- 初始状态为draft
- 自动同步ES

---

### 9.2 GET /api/v1/product/:id - 查询产品详情
**成功响应**: 同创建产品响应

---

### 9.3 PUT /api/v1/product/:id - 更新产品
**请求体 (UpdateProductRequest)**:
```json
{
  "name": "string",                    // 可选, 2-100位
  "status": "string",                  // 可选: draft/submitted/approved/rejected
  "materials": [...],                  // 可选
  "processes": [...]                   // 可选
}
```

**逻辑**: 已审批(approved)的产品不可修改

---

### 9.4 DELETE /api/v1/product/:id - 删除产品
**成功响应**:
```json
{
  "code": 0,
  "data": {
    "message": "删除成功"
  }
}
```

**逻辑**: 已审批的产品不可删除

---

### 9.5 POST /api/v1/product/calculate-cost - 计算产品成本
**请求体 (CalculateCostRequest)**:
```json
{
  "product_id": 1,         // 必填
  "quantity": 100.0,       // 必填, 必须>0
  "use_min_price": true    // 是否使用最低价(true=最低价, false=最高价)
}
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "unit_cost": 150.5,       // 单位成本
    "total_cost": 15050.0,    // 总成本
    "material_costs": [       // 物料成本明细
      {
        "material_id": 1,
        "material_name": "铝型材",
        "ratio": 0.6,
        "price": 80.0,
        "cost": 48.0
      },
      {
        "material_id": 2,
        "material_name": "玻璃",
        "ratio": 0.4,
        "price": 120.0,
        "cost": 48.0
      }
    ],
    "process_costs": [        // 工艺成本明细
      {
        "process_id": 1,
        "process_name": "阳极氧化",
        "price": 54.5,
        "cost": 54.5
      }
    ]
  }
}
```

**逻辑**:
- 物料成本 = Σ(物料价格 × ratio)
- 工艺成本 = Σ(工艺价格)
- 单位成本 = 物料成本 + 工艺成本
- 总成本 = 单位成本 × 数量
- 根据use_min_price选择最低或最高供应商价格

---

## 10. order 包（订单管理模块）

### 域模型 Order 字段:
```go
ID         uint        // 主键
OrderNo    string      // 订单号, 唯一, 最长50位, 必填
ClientID   uint        // 客户ID, 必填
ProductID  uint        // 产品ID, 必填
Quantity   float64     // 数量, decimal(10,2), 必须>0
UnitPrice  float64     // 单价, decimal(10,2), 必须>=0
TotalPrice float64     // 总价, decimal(10,2), 自动计算
Status     string      // 状态: pending/confirmed/production/completed/cancelled
CreatedBy  uint        // 创建人ID, 必填
CreatedAt  time.Time   // 创建时间
UpdatedAt  time.Time   // 更新时间
```

---

### 10.1 POST /api/v1/order - 创建订单
**请求体 (CreateOrderRequest)**:
```json
{
  "order_no": "string",    // 必填, 最长50位
  "client_id": 1,          // 必填
  "product_id": 1,         // 必填
  "quantity": 100.0,       // 必填, 必须>0
  "unit_price": 150.5,     // 必填, 必须>=0
  "created_by": 1          // 必填, 创建人ID
}
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "order_no": "ORD20240101001",
    "client_id": 1,
    "product_id": 1,
    "quantity": 100.0,
    "unit_price": 150.5,
    "total_price": 15050.0,
    "status": "pending",
    "created_by": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**逻辑**:
- 订单号必须唯一
- 总价自动计算 = quantity × unit_price
- 初始状态为pending
- 自动同步ES

---

### 10.2 GET /api/v1/order/:id - 查询订单详情
**成功响应**: 同创建订单响应

---

### 10.3 GET /api/v1/order - 查询订单列表
**查询参数**:
- `limit`: 每页数量 (默认10)
- `offset`: 偏移量 (默认0)

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "total": 200,
    "orders": [
      {
        "id": 1,
        "order_no": "ORD20240101001",
        "client_id": 1,
        "product_id": 1,
        "quantity": 100.0,
        "unit_price": 150.5,
        "total_price": 15050.0,
        "status": "pending",
        "created_by": 1,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

### 10.4 PUT /api/v1/order/:id - 更新订单
**请求体 (UpdateOrderRequest)**:
```json
{
  "status": "confirmed",   // 可选: pending/confirmed/production/completed/cancelled
  "quantity": 120.0,       // 可选, 必须>0
  "unit_price": 160.0      // 可选, 必须>=0
}
```

**逻辑**:
- 状态流转: pending → confirmed → production → completed
- 已完成(completed)或已取消(cancelled)的订单不可修改
- 修改数量或单价会自动重算总价

---

### 10.5 DELETE /api/v1/order/:id - 删除订单
**成功响应**:
```json
{
  "code": 0,
  "data": {
    "message": "删除成功"
  }
}
```

**逻辑**: 已完成的订单不可删除

---

## 11. plan 包（生产计划模块）

### 域模型 Plan 字段:
```go
ID          uint        // 主键
PlanNo      string      // 计划号, 唯一, 最长50位, 必填
OrderID     uint        // 订单ID, 必填
ProductID   uint        // 产品ID, 必填
Quantity    float64     // 数量, decimal(10,2), 必须>0
Status      string      // 状态: planned/in_progress/completed/cancelled
ScheduledAt *time.Time  // 计划时间, 可空
CompletedAt *time.Time  // 完成时间, 可空
CreatedBy   uint        // 创建人ID, 必填
CreatedAt   time.Time   // 创建时间
UpdatedAt   time.Time   // 更新时间
```

---

### 11.1 POST /api/v1/plan - 创建生产计划
**请求体 (CreatePlanRequest)**:
```json
{
  "plan_no": "string",         // 必填, 最长50位
  "order_id": 1,               // 必填
  "product_id": 1,             // 必填
  "quantity": 100.0,           // 必填, 必须>0
  "scheduled_at": "2024-01-15T08:00:00Z",  // 可选
  "created_by": 1              // 必填
}
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "plan_no": "PLAN20240101001",
    "order_id": 1,
    "product_id": 1,
    "quantity": 100.0,
    "status": "planned",
    "scheduled_at": "2024-01-15T08:00:00Z",
    "completed_at": null,
    "created_by": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**逻辑**:
- 计划号必须唯一
- 初始状态为planned
- 自动同步ES

---

### 11.2 GET /api/v1/plan/:id - 查询计划详情
**成功响应**: 同创建计划响应

---

### 11.3 GET /api/v1/plan - 查询计划列表
**查询参数**:
- `limit`: 每页数量 (默认10)
- `offset`: 偏移量 (默认0)

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "total": 150,
    "plans": [
      {
        "id": 1,
        "plan_no": "PLAN20240101001",
        "order_id": 1,
        "product_id": 1,
        "quantity": 100.0,
        "status": "planned",
        "scheduled_at": "2024-01-15T08:00:00Z",
        "completed_at": null,
        "created_by": 1,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

### 11.4 PUT /api/v1/plan/:id - 更新计划
**请求体 (UpdatePlanRequest)**:
```json
{
  "status": "in_progress",     // 可选: planned/in_progress/completed/cancelled
  "quantity": 120.0,           // 可选, 必须>0
  "completed_at": "2024-01-20T18:00:00Z"  // 可选
}
```

**逻辑**:
- 状态流转: planned → in_progress → completed
- 已完成的计划不可修改

---

### 11.5 DELETE /api/v1/plan/:id - 删除计划
**成功响应**:
```json
{
  "code": 0,
  "data": {
    "message": "删除成功"
  }
}
```

**逻辑**: 已完成的计划不可删除

---

## 12. pricing 包（定价管理模块）

### 域模型 SupplierPrice 字段:
```go
ID         uint        // 主键
TargetType string      // 目标类型: material/process
TargetID   uint        // 目标ID(物料ID或工艺ID), 必填
SupplierID uint        // 供应商ID, 必填
Price      float64     // 价格, decimal(10,2), 必须>0
QuotedAt   time.Time   // 报价时间, 必填
CreatedAt  time.Time   // 创建时间
```

### PriceData 结构 (价格缓存):
```go
Price        float64    // 价格
SupplierID   uint       // 供应商ID
SupplierName string     // 供应商名称
QuotedAt     time.Time  // 报价时间
```

---

### 12.1 POST /api/v1/material/price/quote - 物料报价
**请求体 (QuoteRequest)**:
```json
{
  "target_id": 1,      // 必填, 物料ID
  "supplier_id": 1,    // 必填, 供应商ID
  "price": 80.5        // 必填, 必须>0
}
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "message": "报价成功"
  }
}
```

**逻辑**:
- 验证物料ID和供应商ID是否存在
- 创建价格历史记录
- 更新最低/最高价缓存

---

### 12.2 GET /api/v1/material/:id/price - 查询物料价格
**成功响应**:
```json
{
  "code": 0,
  "data": {
    "min": {
      "price": 75.0,
      "supplier_id": 2,
      "supplier_name": "供应商B",
      "quoted_at": "2024-01-10T10:00:00Z"
    },
    "max": {
      "price": 85.0,
      "supplier_id": 1,
      "supplier_name": "供应商A",
      "quoted_at": "2024-01-08T14:00:00Z"
    }
  }
}
```

**逻辑**: 从Redis缓存读取最低和最高价格信息

---

### 12.3 GET /api/v1/material/:id/price/history - 查询物料价格历史
**成功响应**:
```json
{
  "code": 0,
  "data": [
    {
      "price": 80.5,
      "supplier_id": 1,
      "supplier_name": "供应商A",
      "quoted_at": "2024-01-15T10:00:00Z"
    },
    {
      "price": 75.0,
      "supplier_id": 2,
      "supplier_name": "供应商B",
      "quoted_at": "2024-01-10T10:00:00Z"
    }
  ]
}
```

**逻辑**: 返回最近100条报价记录,按报价时间倒序

---

### 12.4 POST /api/v1/process/price/quote - 工艺报价
**请求体 (QuoteRequest)**:
```json
{
  "target_id": 1,      // 必填, 工艺ID
  "supplier_id": 1,    // 必填, 供应商ID
  "price": 54.5        // 必填, 必须>0
}
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "message": "报价成功"
  }
}
```

---

### 12.5 GET /api/v1/process/:id/price - 查询工艺价格
**成功响应**: 同物料价格响应格式

---

### 12.6 GET /api/v1/process/:id/price/history - 查询工艺价格历史
**成功响应**: 同物料价格历史格式

---

## 13. search 包（搜索服务模块）

### 13.1 POST /api/v1/search - Elasticsearch统一搜索
**请求体 (SearchRequest)**:
```json
{
  "query": "铝型材",           // 可选, 搜索关键词
  "indices": [                // 可选, 索引列表
    "materials",
    "processes",
    "products",
    "clients",
    "suppliers",
    "orders",
    "plans"
  ],
  "fields": ["name", "spec"], // 可选, 搜索字段
  "filters": {                // 可选, 过滤条件
    "status": "active"
  },
  "sort": [                   // 可选, 排序
    {
      "field": "created_at",
      "order": "desc"
    }
  ],
  "from": 0,                  // 可选, 偏移量, 默认0
  "size": 20                  // 可选, 返回数量, 默认10
}
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "total": 50,              // 总命中数
    "took": 15,               // 查询耗时(毫秒)
    "max_score": 1.5,         // 最高分数
    "results": [
      {
        "index": "materials",
        "type": "_doc",
        "id": "1",
        "score": 1.5,
        "source": {           // 原始文档数据
          "id": 1,
          "name": "铝型材",
          "spec": "6063-T5",
          "unit": "米",
          "created_at": "2024-01-01T00:00:00Z"
        },
        "highlight": {        // 高亮结果(可选)
          "name": ["<em>铝型材</em>"]
        }
      }
    ]
  }
}
```

**逻辑**:
- 支持跨多个索引搜索
- 支持字段过滤和排序
- 分页查询
- 返回完整文档source

---

### 13.2 GET /api/v1/search/indices - 获取可搜索索引列表
**成功响应**:
```json
{
  "code": 0,
  "data": {
    "indices": [
      {
        "name": "materials",
        "type": "material",
        "fields": ["name", "spec", "unit", "description"]
      },
      {
        "name": "processes",
        "type": "process",
        "fields": ["name", "description"]
      },
      {
        "name": "products",
        "type": "product",
        "fields": ["name", "status"]
      },
      {
        "name": "clients",
        "type": "client",
        "fields": ["name", "contact", "phone", "email"]
      },
      {
        "name": "suppliers",
        "type": "supplier",
        "fields": ["name", "contact", "phone", "email"]
      },
      {
        "name": "orders",
        "type": "order",
        "fields": ["order_no", "status"]
      },
      {
        "name": "plans",
        "type": "plan",
        "fields": ["plan_no", "status"]
      }
    ]
  }
}
```

---

## 14. analytics 包（数据分析模块）

### 14.1 GET /api/v1/return-analysis/customers - 获取客户下拉列表
**成功响应**:
```json
{
  "code": 0,
  "data": [
    {
      "customerNo": "CS0678",
      "customerName": "客户A"
    },
    {
      "customerNo": "CS0679",
      "customerName": "客户B"
    }
  ]
}
```

**逻辑**: 返回所有有已完成订单的客户列表

---

### 14.2 POST /api/v1/return-analysis/analysis - 退货分析
**请求体 (ReturnAnalysisRequest)**:
```json
{
  "customerNo": "CS0678",    // 可选, 空表示所有客户
  "dateRange": {
    "start": "2024-01-01",   // 可选, YYYY-MM-DD格式
    "end": "2024-12-31"      // 可选, start和end必须同时提供或都不提供
  }
}
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "queryConditions": {      // 查询条件回显
      "dateRange": {
        "start": "2024-01-01",
        "end": "2024-12-31"
      },
      "customerNo": "CS0678",
      "customerName": "客户A"
    },
    "meterStats": {           // 米数维度统计
      "totalMeters": 50000.0,
      "returnedMeters": 5000.0,
      "returnRate": "10.00%",
      "orderCount": 80
    },
    "weightStats": {          // 重量维度统计
      "totalWeight": 20000.0,
      "returnedWeight": 2000.0,
      "returnRate": "10.00%",
      "orderCount": 30
    },
    "amountStats": {          // 金额维度统计
      "totalAmountRMB": 120000.0,
      "returnedOrderCount": 15
    },
    "totalOrders": 100        // 总订单数
  }
}
```

**逻辑**:
- 按客户编号筛选（可选）
- 按日期范围筛选（start和end必须同时提供或都不提供）
- 分三个维度统计:
  - 米数: 订单总米数、退货总米数、退货率、涉及订单数
  - 重量: 订单总重量、退货总重量、退货率、涉及订单数
  - 金额: 退款总金额(RMB,USD已按汇率8换算)、有退货的订单数
- 只统计已完成订单

---

## 关键业务规则总结

### 1. 状态流转
- **订单**: pending → confirmed → production → completed
- **计划**: planned → in_progress → completed
- **产品**: draft → submitted → approved/rejected

### 2. 不可变规则
- admin用户不可删除
- 已审批产品不可修改/删除
- 已完成订单/计划不可修改/删除

### 3. 唯一性约束
- User.LoginID
- Department.Code
- Role.Code
- Client.Name
- Supplier.Name
- Order.OrderNo
- Plan.PlanNo

### 4. 计算规则
- 订单总价 = 数量 × 单价
- 产品成本 = Σ(物料价格×比例) + Σ工艺价格
- 物料配比总和必须为1

### 5. Elasticsearch同步
- 所有create/update/delete操作自动同步到ES
- 支持的索引: materials, processes, products, clients, suppliers, orders, plans
