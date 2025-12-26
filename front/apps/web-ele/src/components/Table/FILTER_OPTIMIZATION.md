# Filter 筛选优化说明

## 问题分析

之前打开页面时会发送多个请求：
1. 每个 filter 都会单独发送聚合请求
2. 所有请求在页面加载时就发送，即使用户不一定会用到

例如 material 管理页面会发送 4 个请求：
- `category` 聚合请求
- `status` 聚合请求
- `unit` 聚合请求
- 数据列表请求

## 优化方案

### 方案 1：懒加载（默认，推荐）

**原理**：只在用户点击下拉框时才加载选项数据

**优点**：
- ✅ 大幅减少初始加载请求数
- ✅ 提升页面加载速度
- ✅ 节省带宽

**使用方式**：
```typescript
const pageConfig: PageConfig = {
  pageType: 'material',
  index: 'material',
  pageSize: 20,
  // 不设置 eagerLoadFilters 或设置为 false（默认）
  eagerLoadFilters: false,
  filters: [
    {
      key: 'category',
      label: '分类',
      type: 'select',
      fetchOptions: async () => {
        // 这个函数只在用户点击下拉框时才会调用
        // ...
      }
    }
  ]
}
```

### 方案 2：一次性加载（可选）

**原理**：页面加载时一次性获取所有 filter 的聚合数据

**优点**：
- ✅ 将多个聚合请求合并为 1 个
- ✅ 用户点击下拉框时无需等待
- ✅ 适合 filter 数量少且用户经常使用的场景

**使用方式**：
```typescript
const pageConfig: PageConfig = {
  pageType: 'material',
  index: 'material',
  pageSize: 20,
  eagerLoadFilters: true, // 开启一次性加载
  filters: [
    {
      key: 'category',
      label: '分类',
      type: 'select',
      // fetchOptions 不再需要，会自动从聚合中获取
    },
    {
      key: 'status',
      label: '状态',
      type: 'select',
    }
  ]
}
```

**注意**：使用此模式时：
- 移除 filter 的 `fetchOptions` 函数
- 确保 filter 的 `key` 对应 ES 中的字段名
- 系统会自动发送一个包含所有字段聚合的请求

### 对比

| 模式 | 初始请求数 | 点击下拉框 | 适用场景 |
|-----|----------|----------|---------|
| 懒加载 | 1（仅数据） | 需加载 | filter 多，用户不一定全用 |
| 一次性加载 | 2（数据+聚合） | 无需加载 | filter 少，用户经常使用 |

## 请求示例

### 懒加载模式
```
页面加载：
1. POST /search - 获取数据列表

用户点击 category 下拉框：
2. POST /search - 获取 category 聚合数据

用户点击 status 下拉框：
3. POST /search - 获取 status 聚合数据
```

### 一次性加载模式
```
页面加载：
1. POST /search - 获取所有聚合数据（包含 category, status, unit 等）
2. POST /search - 获取数据列表

用户点击任何下拉框：
无需额外请求
```

## 后端要求

后端需要支持在一个请求中返回多个聚合结果：

```json
{
  "index": "material",
  "pagination": { "offset": 0, "size": 0 },
  "aggRequests": {
    "category": { "type": "terms", "field": "category", "size": 100 },
    "status": { "type": "terms", "field": "status", "size": 100 },
    "unit": { "type": "terms", "field": "unit", "size": 100 }
  }
}
```

响应：
```json
{
  "items": [],
  "total": 0,
  "aggregations": {
    "category": { "buckets": [...] },
    "status": { "buckets": [...] },
    "unit": { "buckets": [...] }
  }
}
```

## 下一步优化

### 下拉框分页（TODO）

当 filter 选项数量很大时（>100），可以实现分页加载：
- 支持虚拟滚动
- 滚动到底部时加载更多
- 需要后端支持 offset 参数
