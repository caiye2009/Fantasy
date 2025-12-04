package domain

import (
	"regexp"
	"time"
)

// Client 客户聚合根
type Client struct {
	ID        uint
	Name      string
	Contact   string
	Phone     string
	Email     string
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Validate 验证客户数据
func (c *Client) Validate() error {
	if c.Name == "" {
		return ErrClientNameEmpty
	}
	
	if len(c.Name) < 2 || len(c.Name) > 100 {
		return ErrClientNameInvalid
	}
	
	if c.Contact != "" && len(c.Contact) > 50 {
		return ErrContactTooLong
	}
	
	if c.Phone != "" {
		if len(c.Phone) > 20 {
			return ErrPhoneTooLong
		}
		if !isValidPhone(c.Phone) {
			return ErrPhoneInvalid
		}
	}
	
	if c.Email != "" {
		if len(c.Email) > 100 {
			return ErrEmailTooLong
		}
		if !isValidEmail(c.Email) {
			return ErrEmailInvalid
		}
	}
	
	if c.Address != "" && len(c.Address) > 200 {
		return ErrAddressTooLong
	}
	
	return nil
}

// UpdateName 更新客户名称
func (c *Client) UpdateName(newName string) error {
	if newName == "" {
		return ErrClientNameEmpty
	}
	if len(newName) < 2 || len(newName) > 100 {
		return ErrClientNameInvalid
	}
	
	c.Name = newName
	return nil
}

// UpdateContact 更新联系人
func (c *Client) UpdateContact(newContact string) error {
	if newContact != "" && len(newContact) > 50 {
		return ErrContactTooLong
	}
	
	c.Contact = newContact
	return nil
}

// UpdatePhone 更新电话
func (c *Client) UpdatePhone(newPhone string) error {
	if newPhone != "" {
		if len(newPhone) > 20 {
			return ErrPhoneTooLong
		}
		if !isValidPhone(newPhone) {
			return ErrPhoneInvalid
		}
	}
	
	c.Phone = newPhone
	return nil
}

// UpdateEmail 更新邮箱
func (c *Client) UpdateEmail(newEmail string) error {
	if newEmail != "" {
		if len(newEmail) > 100 {
			return ErrEmailTooLong
		}
		if !isValidEmail(newEmail) {
			return ErrEmailInvalid
		}
	}
	
	c.Email = newEmail
	return nil
}

// UpdateAddress 更新地址
func (c *Client) UpdateAddress(newAddress string) error {
	if newAddress != "" && len(newAddress) > 200 {
		return ErrAddressTooLong
	}
	
	c.Address = newAddress
	return nil
}

// ToDocument 转换为 ES 文档
func (c *Client) ToDocument() map[string]interface{} {
	return map[string]interface{}{
		"id":         c.ID,
		"name":       c.Name,
		"contact":    c.Contact,
		"phone":      c.Phone,
		"email":      c.Email,
		"address":    c.Address,
		"created_at": c.CreatedAt,
		"updated_at": c.UpdatedAt,
	}
}

// GetIndexName ES 索引名称
func (c *Client) GetIndexName() string {
	return "clients"
}

// GetDocumentID ES 文档 ID
func (c *Client) GetDocumentID() string {
	return string(rune(c.ID))
}

// isValidEmail 验证邮箱格式
func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// isValidPhone 验证电话格式
func isValidPhone(phone string) bool {
	// 支持手机号、固话、带区号的固话
	pattern := `^1[3-9]\d{9}$|^0\d{2,3}-?\d{7,8}$|^\d{7,8}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}