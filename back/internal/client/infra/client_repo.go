package infra

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"back/internal/client/domain"
	"back/pkg/repo"
)

// ClientRepo 客户仓储
type ClientRepo struct {
	*repo.Repo[domain.Client]
	db *gorm.DB
}

// NewClientRepo 创建仓储
func NewClientRepo(db *gorm.DB) *ClientRepo {
	return &ClientRepo{
		Repo: repo.NewRepo[domain.Client](db),
		db:   db,
	}
}

// Save 保存客户
func (r *ClientRepo) Save(ctx context.Context, client *domain.Client) error {
	return r.Create(ctx, client)
}

// FindByID 根据 ID 查询
func (r *ClientRepo) FindByID(ctx context.Context, id uint) (*domain.Client, error) {
	client, err := r.GetByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrClientNotFound
	}
	if err != nil {
		return nil, err
	}
	return client, nil
}

// FindAll 查询所有
func (r *ClientRepo) FindAll(ctx context.Context, limit, offset int) ([]*domain.Client, error) {
	clients, err := r.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Client, len(clients))
	for i := range clients {
		result[i] = &clients[i]
	}
	return result, nil
}

// Update 更新客户
func (r *ClientRepo) Update(ctx context.Context, client *domain.Client) error {
	return r.Repo.Update(ctx, client)
}

// Delete 删除客户
func (r *ClientRepo) Delete(ctx context.Context, id uint) error {
	return r.Repo.Delete(ctx, id)
}

// ExistsByID 检查是否存在
func (r *ClientRepo) ExistsByID(ctx context.Context, id uint) (bool, error) {
	return r.Exists(ctx, map[string]interface{}{"id": id})
}

// FindByName 根据名称查询
func (r *ClientRepo) FindByName(ctx context.Context, name string) (*domain.Client, error) {
	client, err := r.First(ctx, map[string]interface{}{"name": name})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrClientNotFound
	}
	if err != nil {
		return nil, err
	}
	return client, nil
}

// ExistsByName 检查名称是否存在
func (r *ClientRepo) ExistsByName(ctx context.Context, name string) (bool, error) {
	return r.Exists(ctx, map[string]interface{}{"name": name})
}

// FindByPhone 根据电话查询
func (r *ClientRepo) FindByPhone(ctx context.Context, phone string) (*domain.Client, error) {
	client, err := r.First(ctx, map[string]interface{}{"phone": phone})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrClientNotFound
	}
	if err != nil {
		return nil, err
	}
	return client, nil
}

// FindByEmail 根据邮箱查询
func (r *ClientRepo) FindByEmail(ctx context.Context, email string) (*domain.Client, error) {
	client, err := r.First(ctx, map[string]interface{}{"email": email})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrClientNotFound
	}
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Count 统计数量
func (r *ClientRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
