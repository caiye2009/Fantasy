package infra

import (
	"context"

	"gorm.io/gorm"

	"back/internal/user/domain"
	"back/pkg/repo"
)

type UserRepo struct {
	*repo.Repo[domain.User]
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		Repo: repo.NewRepo[domain.User](db),
		db:   db,
	}
}

func (r *UserRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return r.Repo.Exists(ctx, map[string]interface{}{"id": id})
}

func (r *UserRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}

// GetByLoginID 通过LoginID查找用户
func (r *UserRepo) GetByLoginID(ctx context.Context, loginID string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("login_id = ?", loginID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
