package infra

import (
	"context"
	"gorm.io/gorm"

	"back/internal/user_role/domain"
)

type UserRoleRepo struct {
	db *gorm.DB
}

func NewUserRoleRepo(db *gorm.DB) *UserRoleRepo {
	return &UserRoleRepo{db: db}
}

func (r *UserRoleRepo) Create(ctx context.Context, userRole *domain.UserRole) error {
	return r.db.WithContext(ctx).Create(userRole).Error
}

func (r *UserRoleRepo) GetByID(ctx context.Context, id uint) (*domain.UserRole, error) {
	var userRole domain.UserRole
	if err := r.db.WithContext(ctx).First(&userRole, id).Error; err != nil {
		return nil, err
	}
	return &userRole, nil
}

func (r *UserRoleRepo) List(ctx context.Context, limit, offset int) ([]domain.UserRole, error) {
	var userRoles []domain.UserRole
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&userRoles).Error; err != nil {
		return nil, err
	}
	return userRoles, nil
}

func (r *UserRoleRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.UserRole{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UserRoleRepo) UpdateFields(ctx context.Context, id uint, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&domain.UserRole{}).Where("id = ?", id).Updates(updates).Error
}

func (r *UserRoleRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.UserRole{}, id).Error
}

func (r *UserRoleRepo) Exists(ctx context.Context, id uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.UserRole{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRoleRepo) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.UserRole{}).Where("role_code = ?", code).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
