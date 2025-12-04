package infra

import (
	"context"
	"errors"
	
	"gorm.io/gorm"
	
	"back/internal/process/domain"
)

// ProcessRepoImpl 工序仓储实现
type ProcessRepoImpl struct {
	db *gorm.DB
}

// NewProcessRepoImpl 创建仓储实现
func NewProcessRepoImpl(db *gorm.DB) domain.ProcessRepository {
	return &ProcessRepoImpl{db: db}
}

// Save 保存工序
func (r *ProcessRepoImpl) Save(ctx context.Context, process *domain.Process) error {
	po := FromDomain(process)
	err := r.db.WithContext(ctx).Create(po).Error
	if err != nil {
		return err
	}
	process.ID = po.ID
	process.CreatedAt = po.CreatedAt
	process.UpdatedAt = po.UpdatedAt
	return nil
}

// FindByID 根据 ID 查询
func (r *ProcessRepoImpl) FindByID(ctx context.Context, id uint) (*domain.Process, error) {
	var po ProcessPO
	err := r.db.WithContext(ctx).First(&po, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrProcessNotFound
	}
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// FindAll 查询所有
func (r *ProcessRepoImpl) FindAll(ctx context.Context, limit, offset int) ([]*domain.Process, error) {
	var pos []ProcessPO
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	processes := make([]*domain.Process, len(pos))
	for i, po := range pos {
		processes[i] = po.ToDomain()
	}
	return processes, nil
}

// Update 更新工序
func (r *ProcessRepoImpl) Update(ctx context.Context, process *domain.Process) error {
	po := FromDomain(process)
	return r.db.WithContext(ctx).Save(po).Error
}

// Delete 删除工序
func (r *ProcessRepoImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&ProcessPO{}, id).Error
}

// ExistsByID 检查是否存在
func (r *ProcessRepoImpl) ExistsByID(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&ProcessPO{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// FindByName 根据名称查询
func (r *ProcessRepoImpl) FindByName(ctx context.Context, name string) (*domain.Process, error) {
	var po ProcessPO
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrProcessNotFound
	}
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// Count 统计数量
func (r *ProcessRepoImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&ProcessPO{}).Count(&count).Error
	return count, err
}