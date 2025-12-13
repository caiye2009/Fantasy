package infra

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"back/internal/process/domain"
	"back/pkg/repo"
)

// ProcessRepo 工序仓储实现
type ProcessRepo struct {
	*repo.Repo[domain.Process]
	db *gorm.DB
}

// NewProcessRepo 创建仓储
func NewProcessRepo(db *gorm.DB) *ProcessRepo {
	return &ProcessRepo{
		Repo: repo.NewRepo[domain.Process](db),
		db:   db,
	}
}

// Save 保存工序
func (r *ProcessRepo) Save(ctx context.Context, process *domain.Process) error {
	return r.Create(ctx, process)
}

// FindByID 根据 ID 查询
func (r *ProcessRepo) FindByID(ctx context.Context, id uint) (*domain.Process, error) {
	process, err := r.GetByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrProcessNotFound
	}
	if err != nil {
		return nil, err
	}
	return process, nil
}

// FindAll 查询所有
func (r *ProcessRepo) FindAll(ctx context.Context, limit, offset int) ([]*domain.Process, error) {
	processes, err := r.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Process, len(processes))
	for i := range processes {
		result[i] = &processes[i]
	}
	return result, nil
}

// Update 更新工序
func (r *ProcessRepo) Update(ctx context.Context, process *domain.Process) error {
	return r.Repo.Update(ctx, process)
}

// Delete 删除工序
func (r *ProcessRepo) Delete(ctx context.Context, id uint) error {
	return r.Repo.Delete(ctx, id)
}

// ExistsByID 检查是否存在
func (r *ProcessRepo) ExistsByID(ctx context.Context, id uint) (bool, error) {
	return r.Exists(ctx, map[string]interface{}{"id": id})
}

// FindByName 根据名称查询
func (r *ProcessRepo) FindByName(ctx context.Context, name string) (*domain.Process, error) {
	process, err := r.First(ctx, map[string]interface{}{"name": name})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrProcessNotFound
	}
	if err != nil {
		return nil, err
	}
	return process, nil
}

// Count 统计数量
func (r *ProcessRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
