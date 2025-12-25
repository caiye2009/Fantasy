package endpoint

import (
	"regexp"
	"sync"
)

// Registry 接口注册表
type Registry struct {
	endpoints []*Endpoint
	byName    map[string]*Endpoint // name → endpoint
	byRoute   map[string]*Endpoint // method:path → endpoint
	mu        sync.RWMutex
}

var (
	// GlobalRegistry 全局接口注册表
	GlobalRegistry = NewRegistry()
)

// NewRegistry 创建注册表
func NewRegistry() *Registry {
	return &Registry{
		endpoints: make([]*Endpoint, 0),
		byName:    make(map[string]*Endpoint),
		byRoute:   make(map[string]*Endpoint),
	}
}

// Register 注册接口
func (r *Registry) Register(ep *Endpoint) {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := ep.GetName()
	route := ep.Method + ":" + ep.Path

	r.endpoints = append(r.endpoints, ep)
	r.byName[name] = ep
	r.byRoute[route] = ep
}

// RegisterBatch 批量注册
func (r *Registry) RegisterBatch(endpoints []*Endpoint) {
	for _, ep := range endpoints {
		r.Register(ep)
	}
}

// FindByName 根据名称查找
func (r *Registry) FindByName(name string) *Endpoint {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.byName[name]
}

// FindByRoute 根据 HTTP route 查找
func (r *Registry) FindByRoute(method, path string) *Endpoint {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 1. 先尝试精确匹配
	key := method + ":" + path
	if ep := r.byRoute[key]; ep != nil {
		return ep
	}

	// 2. 标准化路径（/api/v1/user/123 → /api/v1/user/:id）
	normalized := normalizePath(path)
	key = method + ":" + normalized
	return r.byRoute[key]
}

// ListAll 列出所有接口
func (r *Registry) ListAll() []*Endpoint {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*Endpoint, len(r.endpoints))
	copy(result, r.endpoints)
	return result
}

// ListByDomain 按域分组列出
func (r *Registry) ListByDomain() map[string][]*Endpoint {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string][]*Endpoint)
	for _, ep := range r.endpoints {
		result[ep.Domain] = append(result[ep.Domain], ep)
	}
	return result
}

// normalizePath 路径标准化（将数字ID替换为:id）
// 例如：/api/v1/order/123/detail → /api/v1/order/:id/detail
func normalizePath(path string) string {
	re := regexp.MustCompile(`/\d+(/|$)`)
	return re.ReplaceAllString(path, "/:id$1")
}
