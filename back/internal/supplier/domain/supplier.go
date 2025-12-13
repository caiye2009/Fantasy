package domain

import (
	"regexp"
	"time"

	"gorm.io/gorm"
)

// Supplier 供应商聚合根
type Supplier struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null;index" json:"name"`
	Contact   string         `gorm:"size:50" json:"contact"`
	Phone     string         `gorm:"size:20;index" json:"phone"`
	Email     string         `gorm:"size:100;index" json:"email"`
	Address   string         `gorm:"size:200" json:"address"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (Supplier) TableName() string {
	return "suppliers"
}

// Validate 验证供应商数据
func (s *Supplier) Validate() error {
	if s.Name == "" {
		return ErrSupplierNameEmpty
	}

	if len(s.Name) < 2 || len(s.Name) > 100 {
		return ErrSupplierNameInvalid
	}

	if s.Contact != "" && len(s.Contact) > 50 {
		return ErrContactTooLong
	}

	if s.Phone != "" {
		if len(s.Phone) > 20 {
			return ErrPhoneTooLong
		}
		// 简单的手机号验证（可以根据业务调整）
		if !isValidPhone(s.Phone) {
			return ErrPhoneInvalid
		}
	}

	if s.Email != "" {
		if len(s.Email) > 100 {
			return ErrEmailTooLong
		}
		if !isValidEmail(s.Email) {
			return ErrEmailInvalid
		}
	}

	if s.Address != "" && len(s.Address) > 200 {
		return ErrAddressTooLong
	}

	return nil
}

// UpdateName 更新供应商名称
func (s *Supplier) UpdateName(newName string) error {
	if newName == "" {
		return ErrSupplierNameEmpty
	}
	if len(newName) < 2 || len(newName) > 100 {
		return ErrSupplierNameInvalid
	}

	s.Name = newName
	return nil
}

// UpdateContact 更新联系人
func (s *Supplier) UpdateContact(newContact string) error {
	if newContact != "" && len(newContact) > 50 {
		return ErrContactTooLong
	}

	s.Contact = newContact
	return nil
}

// UpdatePhone 更新电话
func (s *Supplier) UpdatePhone(newPhone string) error {
	if newPhone != "" {
		if len(newPhone) > 20 {
			return ErrPhoneTooLong
		}
		if !isValidPhone(newPhone) {
			return ErrPhoneInvalid
		}
	}

	s.Phone = newPhone
	return nil
}

// UpdateEmail 更新邮箱
func (s *Supplier) UpdateEmail(newEmail string) error {
	if newEmail != "" {
		if len(newEmail) > 100 {
			return ErrEmailTooLong
		}
		if !isValidEmail(newEmail) {
			return ErrEmailInvalid
		}
	}

	s.Email = newEmail
	return nil
}

// UpdateAddress 更新地址
func (s *Supplier) UpdateAddress(newAddress string) error {
	if newAddress != "" && len(newAddress) > 200 {
		return ErrAddressTooLong
	}

	s.Address = newAddress
	return nil
}

// ToDocument 转换为 ES 文档
func (s *Supplier) ToDocument() map[string]interface{} {
	return map[string]interface{}{
		"id":         s.ID,
		"name":       s.Name,
		"contact":    s.Contact,
		"phone":      s.Phone,
		"email":      s.Email,
		"address":    s.Address,
		"created_at": s.CreatedAt,
		"updated_at": s.UpdatedAt,
	}
}

// GetIndexName ES 索引名称
func (s *Supplier) GetIndexName() string {
	return "suppliers"
}

// GetDocumentID ES 文档 ID
func (s *Supplier) GetDocumentID() string {
	return string(rune(s.ID))
}

// isValidEmail 验证邮箱格式
func isValidEmail(email string) bool {
	// 简单的邮箱正则验证
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// isValidPhone 验证电话格式
func isValidPhone(phone string) bool {
	// 简单的电话号码验证（支持手机和固话）
	// 可以根据业务需求调整正则
	pattern := `^1[3-9]\d{9}$|^0\d{2,3}-?\d{7,8}$|^\d{7,8}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}
