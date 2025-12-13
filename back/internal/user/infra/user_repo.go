package infra

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"back/internal/user/domain"
	"back/pkg/repo"
)

// UserRepo 用户仓储实现
type UserRepo struct {
	*repo.Repo[domain.User]
	db *gorm.DB
}

// NewUserRepo 创建仓储
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		Repo: repo.NewRepo[domain.User](db),
		db:   db,
	}
}

// Save 保存用户
func (r *UserRepo) Save(ctx context.Context, user *domain.User) error {
	err := r.Create(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

// FindByID 根据 ID 查询
func (r *UserRepo) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	user, err := r.GetByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindAll 查询所有
func (r *UserRepo) FindAll(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	users, err := r.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.User, len(users))
	for i := range users {
		result[i] = &users[i]
	}
	return result, nil
}

// Update 更新用户
func (r *UserRepo) Update(ctx context.Context, user *domain.User) error {
	return r.Repo.Update(ctx, user)
}

// Delete 删除用户
func (r *UserRepo) Delete(ctx context.Context, id uint) error {
	return r.Repo.Delete(ctx, id)
}

// ExistsByID 检查是否存在
func (r *UserRepo) ExistsByID(ctx context.Context, id uint) (bool, error) {
	return r.Exists(ctx, map[string]interface{}{"id": id})
}

// FindByLoginID 根据工号查询
func (r *UserRepo) FindByLoginID(ctx context.Context, loginID string) (*domain.User, error) {
	user, err := r.First(ctx, map[string]interface{}{"login_id": loginID})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ExistsByLoginID 检查工号是否存在
func (r *UserRepo) ExistsByLoginID(ctx context.Context, loginID string) (bool, error) {
	return r.Exists(ctx, map[string]interface{}{"login_id": loginID})
}

// FindByEmail 根据邮箱查询
func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.First(ctx, map[string]interface{}{"email": email})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindByRole 根据角色查询
func (r *UserRepo) FindByRole(ctx context.Context, role string, limit, offset int) ([]*domain.User, error) {
	var result []*domain.User
	err := r.db.WithContext(ctx).
		Where("role = ?", role).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&result).Error
	return result, err
}

// FindByStatus 根据状态查询
func (r *UserRepo) FindByStatus(ctx context.Context, status string, limit, offset int) ([]*domain.User, error) {
	var result []*domain.User
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&result).Error
	return result, err
}

// Count 统计数量
func (r *UserRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}

// GetNextLoginID 获取下一个 login_id（自增，从 1000 开始）
func (r *UserRepo) GetNextLoginID(ctx context.Context) (string, error) {
	var maxLoginID string
	// 只查询纯数字的 login_id（过滤掉 "admin" 等非数字 ID）
	err := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Select("login_id").
		Where("login_id ~ ?", "^[0-9]+$"). // 正则匹配纯数字
		Order("CAST(login_id AS INTEGER) DESC").
		Limit(1).
		Pluck("login_id", &maxLoginID).Error

	if err != nil {
		return "", err
	}

	// 如果没有纯数字的用户，从 1000 开始
	if maxLoginID == "" {
		return "1000", nil
	}

	// 将 login_id 转换为整数并加 1
	var currentID int
	_, err = fmt.Sscanf(maxLoginID, "%d", &currentID)
	if err != nil {
		// 如果转换失败，从 1000 开始
		return "1000", nil
	}

	nextID := currentID + 1
	return fmt.Sprintf("%d", nextID), nil
}

// GetAllDepartments 获取所有唯一的部门列表
func (r *UserRepo) GetAllDepartments(ctx context.Context) ([]string, error) {
	var departments []string
	err := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Distinct("department").
		Where("department != ?", "").
		Order("department ASC").
		Pluck("department", &departments).Error

	return departments, err
}
