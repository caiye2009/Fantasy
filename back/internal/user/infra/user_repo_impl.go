package infra

import (
	"context"
	"errors"
	
	"gorm.io/gorm"
	
	"back/internal/user/domain"
)

// UserRepoImpl 用户仓储实现
type UserRepoImpl struct {
	db *gorm.DB
}

// NewUserRepoImpl 创建仓储实现
func NewUserRepoImpl(db *gorm.DB) domain.UserRepository {
	return &UserRepoImpl{db: db}
}

// Save 保存用户
func (r *UserRepoImpl) Save(ctx context.Context, user *domain.User) error {
	po := FromDomain(user)
	err := r.db.WithContext(ctx).Create(po).Error
	if err != nil {
		return err
	}
	user.ID = po.ID
	user.CreatedAt = po.CreatedAt
	user.UpdatedAt = po.UpdatedAt
	return nil
}

// FindByID 根据 ID 查询
func (r *UserRepoImpl) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	var po UserPO
	err := r.db.WithContext(ctx).First(&po, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// FindAll 查询所有
func (r *UserRepoImpl) FindAll(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	var pos []UserPO
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at DESC").Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	users := make([]*domain.User, len(pos))
	for i, po := range pos {
		users[i] = po.ToDomain()
	}
	return users, nil
}

// Update 更新用户
func (r *UserRepoImpl) Update(ctx context.Context, user *domain.User) error {
	po := FromDomain(user)
	return r.db.WithContext(ctx).Save(po).Error
}

// Delete 删除用户
func (r *UserRepoImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&UserPO{}, id).Error
}

// ExistsByID 检查是否存在
func (r *UserRepoImpl) ExistsByID(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&UserPO{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// FindByLoginID 根据工号查询
func (r *UserRepoImpl) FindByLoginID(ctx context.Context, loginID string) (*domain.User, error) {
	var po UserPO
	err := r.db.WithContext(ctx).Where("login_id = ?", loginID).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// ExistsByLoginID 检查工号是否存在
func (r *UserRepoImpl) ExistsByLoginID(ctx context.Context, loginID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&UserPO{}).Where("login_id = ?", loginID).Count(&count).Error
	return count > 0, err
}

// FindByEmail 根据邮箱查询
func (r *UserRepoImpl) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var po UserPO
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// FindByRole 根据角色查询
func (r *UserRepoImpl) FindByRole(ctx context.Context, role domain.UserRole, limit, offset int) ([]*domain.User, error) {
	var pos []UserPO
	err := r.db.WithContext(ctx).
		Where("role = ?", role).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	users := make([]*domain.User, len(pos))
	for i, po := range pos {
		users[i] = po.ToDomain()
	}
	return users, nil
}

// FindByStatus 根据状态查询
func (r *UserRepoImpl) FindByStatus(ctx context.Context, status domain.UserStatus, limit, offset int) ([]*domain.User, error) {
	var pos []UserPO
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	users := make([]*domain.User, len(pos))
	for i, po := range pos {
		users[i] = po.ToDomain()
	}
	return users, nil
}

// Count 统计数量
func (r *UserRepoImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&UserPO{}).Count(&count).Error
	return count, err
}