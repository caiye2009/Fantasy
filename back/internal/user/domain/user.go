package domain

import (
	"time"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 用户角色常量
const (
	UserRoleAdmin     = "admin"     // 管理员
	UserRoleHR        = "hr"        // 人力资源
	UserRoleSales     = "sales"     // 销售
	UserRoleFollower  = "follower"  // 跟单员
	UserRoleAssistant = "assistant" // 助理
	UserRoleUser      = "user"      // 普通用户
)

// GetAllRoles 获取所有角色（统一维护角色列表）
func GetAllRoles() []string {
	return []string{
		UserRoleAdmin,
		UserRoleHR,
		UserRoleSales,
		UserRoleFollower,
		UserRoleAssistant,
		UserRoleUser,
	}
}

// IsValidRole 验证角色是否有效
func IsValidRole(role string) bool {
	for _, validRole := range GetAllRoles() {
		if validRole == role {
			return true
		}
	}
	return false
}

// 用户状态常量
const (
	UserStatusActive    = "active"    // 激活
	UserStatusInactive  = "inactive"  // 未激活
	UserStatusSuspended = "suspended" // 停用
)

// User 用户聚合根
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	LoginID      string         `gorm:"size:50;not null;uniqueIndex" json:"login_id"`
	Username     string         `gorm:"size:100;not null" json:"username"`
	Department   string         `gorm:"size:100" json:"department"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Email        string         `gorm:"size:100" json:"email"`
	Role         string         `gorm:"size:50;not null;index" json:"role"`
	Status       string         `gorm:"size:20;default:active;index" json:"status"`
	HasInitPass  bool           `gorm:"default:true" json:"has_init_pass"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (User) TableName() string {
	return "users"
}

// Validate 验证用户数据
func (u *User) Validate() error {
	if u.LoginID == "" {
		return ErrLoginIDEmpty
	}
	
	if len(u.LoginID) < 4 || len(u.LoginID) > 20 {
		return ErrLoginIDInvalid
	}
	
	if u.Username == "" {
		return ErrUsernameEmpty
	}
	
	if len(u.Username) < 2 || len(u.Username) > 50 {
		return ErrUsernameInvalid
	}
	
	if u.PasswordHash == "" {
		return ErrPasswordRequired
	}
	
	return nil
}

// SetDefaultPassword 设置默认密码（用于新建用户）
func (u *User) SetDefaultPassword() error {
	defaultPassword := "123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	
	u.PasswordHash = string(hashedPassword)
	u.HasInitPass = true
	return nil
}

// ChangePassword 修改密码
func (u *User) ChangePassword(currentPassword, newPassword string) error {
	// 1. 如果不是初始密码，需要验证当前密码
	if !u.HasInitPass {
		if currentPassword == "" {
			return ErrCurrentPasswordRequired
		}
		
		if err := u.ValidatePassword(currentPassword); err != nil {
			return ErrCurrentPasswordIncorrect
		}
	}
	
	// 2. 验证新密码长度
	if len(newPassword) < 6 {
		return ErrPasswordTooShort
	}
	
	// 3. 生成新密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	
	u.PasswordHash = string(hashedPassword)
	u.HasInitPass = false
	
	return nil
}

// ValidatePassword 验证密码
func (u *User) ValidatePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return ErrPasswordIncorrect
	}
	return nil
}

// UpdateUsername 更新用户名
func (u *User) UpdateUsername(newUsername string) error {
	if newUsername == "" {
		return ErrUsernameEmpty
	}
	if len(newUsername) < 2 || len(newUsername) > 50 {
		return ErrUsernameInvalid
	}

	u.Username = newUsername
	return nil
}

// UpdateDepartment 更新部门
func (u *User) UpdateDepartment(newDepartment string) error {
	if len(newDepartment) > 100 {
		return ErrUsernameInvalid // 可以后续添加专门的 ErrDepartmentInvalid
	}

	u.Department = newDepartment
	return nil
}

// UpdateEmail 更新邮箱
func (u *User) UpdateEmail(newEmail string) error {
	// 简单的邮箱格式验证（实际应该用正则）
	if newEmail != "" && len(newEmail) < 5 {
		return ErrEmailInvalid
	}
	
	u.Email = newEmail
	return nil
}

// UpdateRole 更新角色
func (u *User) UpdateRole(newRole string) error {
	if !IsValidRole(newRole) {
		return ErrInvalidRole
	}
	u.Role = newRole
	return nil
}

// Activate 激活用户
func (u *User) Activate() error {
	if u.Status == UserStatusActive {
		return ErrUserAlreadyActive
	}
	
	u.Status = UserStatusActive
	return nil
}

// Suspend 停用用户
func (u *User) Suspend() error {
	if u.Status == UserStatusSuspended {
		return ErrUserAlreadySuspended
	}
	
	u.Status = UserStatusSuspended
	return nil
}

// IsActive 是否激活
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// CanDelete 是否可以删除
func (u *User) CanDelete() bool {
	// 管理员不能删除（可以根据业务规则调整）
	return u.Role != UserRoleAdmin
}