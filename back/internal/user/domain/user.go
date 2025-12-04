package domain

import (
	"time"	
	"golang.org/x/crypto/bcrypt"
)

// UserRole 用户角色
type UserRole string

const (
	UserRoleAdmin     UserRole = "admin"     // 管理员
	UserRoleHR        UserRole = "hr"        // 人力资源
	UserRoleSales     UserRole = "sales"     // 销售
	UserRoleFollower  UserRole = "follower"  // 跟单员
	UserRoleAssistant UserRole = "assistant" // 助理
	UserRoleUser      UserRole = "user"      // 普通用户
)

// UserStatus 用户状态
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"    // 激活
	UserStatusInactive  UserStatus = "inactive"  // 未激活
	UserStatusSuspended UserStatus = "suspended" // 停用
)

// User 用户聚合根
type User struct {
	ID           uint
	LoginID      string
	Username     string
	PasswordHash string
	Email        string
	Role         UserRole
	Status       UserStatus
	HasInitPass  bool // 是否使用初始密码
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
func (u *User) UpdateRole(newRole UserRole) error {
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