package domain

import (
	"fmt"
	"regexp"
	"time"

	"gorm.io/gorm"
)

// Client 客户聚合根
type Client struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CustomNo       string         `gorm:"size:50;uniqueIndex" json:"customNo"`        // 客户代码
	CustomerCode   string         `gorm:"size:50;index" json:"customerCode"`          // 内部编码
	InputDate      *time.Time     `gorm:"type:date" json:"inputDate"`                 // 添加时间
	Sales          string         `gorm:"size:50" json:"sales"`                       // 业务员
	CustomName     string         `gorm:"size:200;not null;index" json:"customName"`  // 客户名称
	StateChNm      string         `gorm:"size:100" json:"stateChNm"`                  // 国家
	Country        string         `gorm:"size:50" json:"country"`                     // 国家代码
	Address        string         `gorm:"size:500" json:"address"`                    // 中文地址
	AddressEn      string         `gorm:"size:500" json:"addressEn"`                  // 英文地址
	CustomNameEn   string         `gorm:"size:200" json:"customNameEn"`               // 英文名称
	Contactor      string         `gorm:"size:100" json:"contactor"`                  // 联系人
	UnitPhone      string         `gorm:"size:50" json:"unitPhone"`                   // 电话
	Mobile         string         `gorm:"size:50" json:"mobile"`                      // 手机
	FaxNum         string         `gorm:"size:50" json:"faxNum"`                      // 传真
	Email          string         `gorm:"size:100" json:"email"`                      // 邮箱
	PyCustomName   string         `gorm:"size:200" json:"pyCustomName"`               // 所属客户
	CheckRequest   string         `gorm:"type:text" json:"checkRequest"`              // 检验要求
	CustomStatus   string         `gorm:"size:50" json:"customStatus"`                // 状态
	DocMan         string         `gorm:"size:50" json:"docMan"`                      // 输入人
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (Client) TableName() string {
	return "clients"
}

// Validate 验证客户数据
func (c *Client) Validate() error {
	if c.CustomName == "" {
		return ErrClientNameEmpty
	}

	if len(c.CustomName) < 2 || len(c.CustomName) > 200 {
		return ErrClientNameInvalid
	}

	if c.CustomNo != "" && len(c.CustomNo) > 50 {
		return fmt.Errorf("custom no too long")
	}

	if c.Contactor != "" && len(c.Contactor) > 100 {
		return ErrContactTooLong
	}

	if c.UnitPhone != "" && len(c.UnitPhone) > 50 {
		return ErrPhoneTooLong
	}

	if c.Mobile != "" && len(c.Mobile) > 50 {
		return ErrPhoneTooLong
	}

	if c.Email != "" {
		if len(c.Email) > 100 {
			return ErrEmailTooLong
		}
		if !isValidEmail(c.Email) {
			return ErrEmailInvalid
		}
	}

	if c.Address != "" && len(c.Address) > 500 {
		return ErrAddressTooLong
	}

	if c.AddressEn != "" && len(c.AddressEn) > 500 {
		return ErrAddressTooLong
	}

	return nil
}


// ToDocument 转换为 ES 文档（小驼峰字段名）
func (c *Client) ToDocument() map[string]interface{} {
	doc := map[string]interface{}{
		"id":           c.ID,
		"customNo":     c.CustomNo,
		"customerCode": c.CustomerCode,
		"sales":        c.Sales,
		"customName":   c.CustomName,
		"stateChNm":    c.StateChNm,
		"country":      c.Country,
		"address":      c.Address,
		"addressEn":    c.AddressEn,
		"customNameEn": c.CustomNameEn,
		"contactor":    c.Contactor,
		"unitPhone":    c.UnitPhone,
		"mobile":       c.Mobile,
		"faxNum":       c.FaxNum,
		"email":        c.Email,
		"pyCustomName": c.PyCustomName,
		"checkRequest": c.CheckRequest,
		"customStatus": c.CustomStatus,
		"docMan":       c.DocMan,
		"createdAt":    c.CreatedAt,
		"updatedAt":    c.UpdatedAt,
	}
	if c.InputDate != nil {
		doc["inputDate"] = c.InputDate
	}
	return doc
}

// GetIndexName ES 索引名称
func (c *Client) GetIndexName() string {
	return "clients"
}

// GetDocumentID ES 文档 ID
func (c *Client) GetDocumentID() string {
	return fmt.Sprintf("%d", c.ID)
}

// isValidEmail 验证邮箱格式
func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// CalculatePriorityScore 计算优先级分数（可选功能，默认返回 0）
// 如果需要业务优先级排序，可在此函数中实现评分逻辑
func (c *Client) CalculatePriorityScore() int {
	score := 0

	// 1. 状态评分（最高 200 分）
	switch c.CustomStatus {
	case "active": // 活跃客户
		score += 200
	case "potential": // 潜在客户
		score += 100
	case "dormant": // 休眠客户
		score += 50
	}

	// 2. 时间新鲜度（最近添加的客户加分，最高 50 分）
	if c.InputDate != nil {
		daysSinceInput := int(time.Since(*c.InputDate).Hours() / 24)
		if daysSinceInput < 30 { // 30 天内添加
			score += 50 - (daysSinceInput / 2)
		}
	}

	// 3. 数据完整度（有联系方式的客户加分）
	if c.Contactor != "" {
		score += 10
	}
	if c.UnitPhone != "" || c.Mobile != "" {
		score += 10
	}
	if c.Email != "" {
		score += 10
	}

	return score
}