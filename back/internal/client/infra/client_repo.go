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

// FindByCustomNo 根据客户代码查询
func (r *ClientRepo) FindByCustomNo(ctx context.Context, customNo string) (*domain.Client, error) {
	client, err := r.First(ctx, map[string]interface{}{"custom_no": customNo})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrClientNotFound
	}
	if err != nil {
		return nil, err
	}
	return client, nil
}

// ExistsByCustomNo 检查客户代码是否存在
func (r *ClientRepo) ExistsByCustomNo(ctx context.Context, customNo string) (bool, error) {
	return r.Exists(ctx, map[string]interface{}{"custom_no": customNo})
}

// Count 统计数量
func (r *ClientRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
