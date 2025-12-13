package infra

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"back/internal/plan/domain"
	"back/pkg/repo"
)

// PlanRepo 计划仓储实现
type PlanRepo struct {
	*repo.Repo[domain.Plan]
	db *gorm.DB
}

// NewPlanRepo 创建仓储
func NewPlanRepo(db *gorm.DB) *PlanRepo {
	return &PlanRepo{
		Repo: repo.NewRepo[domain.Plan](db),
		db:   db,
	}
}

// Save 保存计划
func (r *PlanRepo) Save(ctx context.Context, plan *domain.Plan) error {
	return r.Create(ctx, plan)
}

// FindByID 根据 ID 查询
func (r *PlanRepo) FindByID(ctx context.Context, id uint) (*domain.Plan, error) {
	plan, err := r.GetByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrPlanNotFound
	}
	if err != nil {
		return nil, err
	}
	return plan, nil
}

// FindAll 查询所有
func (r *PlanRepo) FindAll(ctx context.Context, limit, offset int) ([]*domain.Plan, error) {
	plans, err := r.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Plan, len(plans))
	for i := range plans {
		result[i] = &plans[i]
	}
	return result, nil
}

// Update 更新计划
func (r *PlanRepo) Update(ctx context.Context, plan *domain.Plan) error {
	return r.Repo.Update(ctx, plan)
}

// Delete 删除计划
func (r *PlanRepo) Delete(ctx context.Context, id uint) error {
	return r.Repo.Delete(ctx, id)
}

// ExistsByID 检查是否存在
func (r *PlanRepo) ExistsByID(ctx context.Context, id uint) (bool, error) {
	return r.Exists(ctx, map[string]interface{}{"id": id})
}

// FindByPlanNo 根据计划编号查询
func (r *PlanRepo) FindByPlanNo(ctx context.Context, planNo string) (*domain.Plan, error) {
	plan, err := r.First(ctx, map[string]interface{}{"plan_no": planNo})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrPlanNotFound
	}
	if err != nil {
		return nil, err
	}
	return plan, nil
}

// ExistsByPlanNo 检查计划编号是否存在
func (r *PlanRepo) ExistsByPlanNo(ctx context.Context, planNo string) (bool, error) {
	return r.Exists(ctx, map[string]interface{}{"plan_no": planNo})
}

// FindByOrderID 根据订单 ID 查询
func (r *PlanRepo) FindByOrderID(ctx context.Context, orderID uint) ([]*domain.Plan, error) {
	var result []*domain.Plan
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		Find(&result).Error
	return result, err
}

// FindByProductID 根据产品 ID 查询
func (r *PlanRepo) FindByProductID(ctx context.Context, productID uint) ([]*domain.Plan, error) {
	var result []*domain.Plan
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("created_at DESC").
		Find(&result).Error
	return result, err
}

// FindByStatus 根据状态查询
func (r *PlanRepo) FindByStatus(ctx context.Context, status string, limit, offset int) ([]*domain.Plan, error) {
	var result []*domain.Plan
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&result).Error
	return result, err
}

// Count 统计数量
func (r *PlanRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
