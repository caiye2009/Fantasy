package domain

import "context"

// UserRepository 用户仓储接口
type UserRepository interface {
	// 基础 CRUD
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uint) (*User, error)
	FindAll(ctx context.Context, limit, offset int) ([]*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
	
	// 业务查询
	ExistsByID(ctx context.Context, id uint) (bool, error)
	FindByLoginID(ctx context.Context, loginID string) (*User, error)
	ExistsByLoginID(ctx context.Context, loginID string) (bool, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByRole(ctx context.Context, role UserRole, limit, offset int) ([]*User, error)
	FindByStatus(ctx context.Context, status UserStatus, limit, offset int) ([]*User, error)
	Count(ctx context.Context) (int64, error)
}