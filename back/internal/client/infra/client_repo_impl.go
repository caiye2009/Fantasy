package infra

import (
	"context"
	"errors"
	
	"gorm.io/gorm"
	
	"back/internal/client/domain"
)

// ClientRepoImpl 客户仓储实现
type ClientRepoImpl struct {
	db *gorm.DB
}

// NewClientRepoImpl 创建仓储实现
func NewClientRepoImpl(db *gorm.DB) domain.ClientRepository {
	return &ClientRepoImpl{db: db}
}

// Save 保存客户
func (r *ClientRepoImpl) Save(ctx context.Context, client *domain.Client) error {
	po := FromDomain(client)
	err := r.db.WithContext(ctx).Create(po).Error
	if err != nil {
		return err
	}
	client.ID = po.ID
	client.CreatedAt = po.CreatedAt
	client.UpdatedAt = po.UpdatedAt
	return nil
}

// FindByID 根据 ID 查询
func (r *ClientRepoImpl) FindByID(ctx context.Context, id uint) (*domain.Client, error) {
	var po ClientPO
	err := r.db.WithContext(ctx).First(&po, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrClientNotFound
	}
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// FindAll 查询所有
func (r *ClientRepoImpl) FindAll(ctx context.Context, limit, offset int) ([]*domain.Client, error) {
	var pos []ClientPO
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at DESC").Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	clients := make([]*domain.Client, len(pos))
	for i, po := range pos {
		clients[i] = po.ToDomain()
	}
	return clients, nil
}

// Update 更新客户
func (r *ClientRepoImpl) Update(ctx context.Context, client *domain.Client) error {
	po := FromDomain(client)
	return r.db.WithContext(ctx).Save(po).Error
}

// Delete 删除客户
func (r *ClientRepoImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&ClientPO{}, id).Error
}

// ExistsByID 检查是否存在
func (r *ClientRepoImpl) ExistsByID(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&ClientPO{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// FindByName 根据名称查询
func (r *ClientRepoImpl) FindByName(ctx context.Context, name string) (*domain.Client, error) {
	var po ClientPO
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrClientNotFound
	}
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// ExistsByName 检查名称是否存在
func (r *ClientRepoImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&ClientPO{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

// FindByPhone 根据电话查询
func (r *ClientRepoImpl) FindByPhone(ctx context.Context, phone string) (*domain.Client, error) {
	var po ClientPO
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrClientNotFound
	}
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// FindByEmail 根据邮箱查询
func (r *ClientRepoImpl) FindByEmail(ctx context.Context, email string) (*domain.Client, error) {
	var po ClientPO
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrClientNotFound
	}
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// Count 统计数量
func (r *ClientRepoImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&ClientPO{}).Count(&count).Error
	return count, err
}