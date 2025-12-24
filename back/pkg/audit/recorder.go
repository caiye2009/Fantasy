package audit

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	RecorderContextKey = "audit_recorder"
)

// Recorder 审计记录器
type Recorder struct {
	c         *gin.Context
	db        *gorm.DB
	startTime time.Time

	// 基本信息（从 context 自动提取）
	loginID   uint
	username  string

	// 请求信息（自动记录）
	httpMethod  string
	requestPath string
	ipAddress   string
	userAgent   string
	requestID   string

	// 业务信息（可选设置）
	domain     string
	action     string
	resourceID string
	oldData    interface{}
	newData    interface{}
}

// NewRecorder 创建新的审计记录器
func NewRecorder(c *gin.Context, db *gorm.DB) *Recorder {
	// 从 auth context 提取用户信息
	loginID, _ := c.Get("loginId")
	username, exists := c.Get("username")

	// 如果没有 username，尝试构造一个（兼容当前代码）
	if !exists || username == "" {
		if loginIDStr, ok := loginID.(string); ok {
			username = "User" + loginIDStr
		}
	}

	// 转换 loginID
	var loginIDUint uint
	switch v := loginID.(type) {
	case string:
		id, _ := strconv.Atoi(v)
		loginIDUint = uint(id)
	case uint:
		loginIDUint = v
	case int:
		loginIDUint = uint(v)
	}

	// 提取请求ID（如果有的话，可以从 header 中获取）
	requestID := c.GetHeader("X-Request-ID")
	if requestID == "" {
		// 生成简单的请求ID（可以用更复杂的算法）
		requestID = strconv.FormatInt(time.Now().UnixNano(), 36)
	}

	return &Recorder{
		c:           c,
		db:          db,
		startTime:   time.Now(),
		loginID:     loginIDUint,
		username:    username.(string),
		httpMethod:  c.Request.Method,
		requestPath: c.Request.URL.Path,
		ipAddress:   c.ClientIP(),
		userAgent:   c.Request.UserAgent(),
		requestID:   requestID,
	}
}

// SetDomain 设置业务域
func (r *Recorder) SetDomain(domain string) {
	r.domain = domain
}

// SetAction 设置操作动作
func (r *Recorder) SetAction(action string) {
	r.action = action
}

// SetResourceID 设置资源ID
func (r *Recorder) SetResourceID(id interface{}) {
	switch v := id.(type) {
	case string:
		r.resourceID = v
	case uint:
		r.resourceID = strconv.FormatUint(uint64(v), 10)
	case int:
		r.resourceID = strconv.Itoa(v)
	case int64:
		r.resourceID = strconv.FormatInt(v, 10)
	default:
		r.resourceID = ""
	}
}

// SetOld 设置变更前数据
func (r *Recorder) SetOld(data interface{}) {
	r.oldData = data
}

// SetNew 设置变更后数据
func (r *Recorder) SetNew(data interface{}) {
	r.newData = data
}

// Save 保存审计日志到数据库
func (r *Recorder) Save() error {
	// 如果 DB 为 nil，跳过（测试环境可能没有 DB）
	if r.db == nil {
		return nil
	}

	// 获取响应状态码
	statusCode := r.c.Writer.Status()

	// 计算耗时
	duration := time.Since(r.startTime).Milliseconds()

	// 优先使用路由标记的 domain 和 action（通过 audit.Mark 设置）
	if r.domain == "" {
		if domain, exists := r.c.Get("audit_domain"); exists {
			r.domain = domain.(string)
		}
	}
	if r.action == "" {
		if action, exists := r.c.Get("audit_action"); exists {
			r.action = action.(string)
		}
	}

	// 如果还是没有设置 domain 和 action，尝试从路径推断
	if r.domain == "" || r.action == "" {
		r.inferFromPath()
	}

	// 序列化 old 和 new 数据
	var oldDataJSON, newDataJSON string
	if r.oldData != nil {
		if bytes, err := json.Marshal(r.oldData); err == nil {
			oldDataJSON = string(bytes)
		}
	}
	if r.newData != nil {
		if bytes, err := json.Marshal(r.newData); err == nil {
			newDataJSON = string(bytes)
		}
	}

	// 构建审计日志
	auditLog := &AuditLog{
		LoginID:     r.loginID,
		Username:    r.username,
		Domain:      r.domain,
		Action:      r.action,
		ResourceID:  r.resourceID,
		HTTPMethod:  r.httpMethod,
		RequestPath: r.requestPath,
		IPAddress:   r.ipAddress,
		StatusCode:  statusCode,
		DurationMs:  duration,
		UserAgent:   r.userAgent,
		RequestID:   r.requestID,
		OldData:     oldDataJSON,
		NewData:     newDataJSON,
	}

	// 如果请求失败，记录错误信息
	if statusCode >= 400 {
		if err, exists := r.c.Get("error"); exists {
			auditLog.ErrorMessage = err.(string)
		}
	}

	// 保存到数据库
	return r.db.Create(auditLog).Error
}

// inferFromPath 从请求路径推断 domain 和 action
func (r *Recorder) inferFromPath() {
	// 解析路径：/api/v1/{domain}/...
	parts := strings.Split(strings.Trim(r.requestPath, "/"), "/")

	// 至少需要 3 个部分：api/v1/domain
	if len(parts) >= 3 {
		r.domain = parts[2]
	}

	// 根据 HTTP 方法推断 action
	switch r.httpMethod {
	case "POST":
		// POST /api/v1/order -> create
		// POST /api/v1/order/123/assign -> assign
		if len(parts) >= 4 {
			r.action = parts[3] // 可能是 assign-department 之类的
		} else {
			r.action = ActionCreate
		}
	case "PUT", "PATCH":
		r.action = ActionUpdate
	case "DELETE":
		r.action = ActionDelete
	default:
		r.action = strings.ToLower(r.httpMethod)
	}

	// 提取资源ID（如果路径中有数字）
	if r.resourceID == "" && len(parts) >= 4 {
		// /api/v1/order/123 -> resourceID = 123
		if _, err := strconv.Atoi(parts[3]); err == nil {
			r.resourceID = parts[3]
		}
	}
}

// Get 从 gin.Context 获取 Recorder
func Get(c *gin.Context) *Recorder {
	if recorder, exists := c.Get(RecorderContextKey); exists {
		return recorder.(*Recorder)
	}
	return nil
}
