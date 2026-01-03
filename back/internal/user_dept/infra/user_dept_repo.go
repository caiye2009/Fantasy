package infra

import (
	"context"
	"gorm.io/gorm"

	"back/internal/user_dept/domain"
)

type UserDeptRepo struct {
	db *gorm.DB
}

func NewUserDeptRepo(db *gorm.DB) *UserDeptRepo {
	return &UserDeptRepo{db: db}
}

func (r *UserDeptRepo) Create(ctx context.Context, userDept *domain.UserDept) error {
	return r.db.WithContext(ctx).Create(userDept).Error
}

func (r *UserDeptRepo) GetByID(ctx context.Context, id uint) (*domain.UserDept, error) {
	var userDept domain.UserDept
	if err := r.db.WithContext(ctx).First(&userDept, id).Error; err != nil {
		return nil, err
	}
	return &userDept, nil
}

func (r *UserDeptRepo) List(ctx context.Context, limit, offset int) ([]domain.UserDept, error) {
	var userDepts []domain.UserDept
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&userDepts).Error; err != nil {
		return nil, err
	}
	return userDepts, nil
}

func (r *UserDeptRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.UserDept{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UserDeptRepo) UpdateFields(ctx context.Context, id uint, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&domain.UserDept{}).Where("id = ?", id).Updates(updates).Error
}

func (r *UserDeptRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.UserDept{}, id).Error
}

func (r *UserDeptRepo) Exists(ctx context.Context, id uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.UserDept{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserDeptRepo) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.UserDept{}).Where("dept_code = ?", code).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
