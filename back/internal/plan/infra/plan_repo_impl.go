package infra

import (
	"context"
	"errors"
	
	"gorm.io/gorm"
	
	"back/internal/plan/domain"
)

// PlanRepoImpl 计划仓储实现
type PlanRepoImpl struct {
	db *gorm.DB
}

// NewPlanRepoImpl 创建仓储实现
func NewPlanRepoImpl(db *gorm.DB) domain.PlanRepository {
	return &PlanRepoImpl{db: db}
}

// Save 保存计划
func (r *PlanRepoImpl) Save(ctx context.Context, plan *domain.Plan) error {
	po := FromDomain(plan)
	err := r.db.WithContext(ctx).Create(po).Error
	if err != nil {
		return err
	}
	plan.ID = po.ID
	plan.CreatedAt = po.CreatedAt
	plan.UpdatedAt = po.UpdatedAt
	return nil
}

// FindByID 根据 ID 查询
func (r *PlanRepoImpl) FindByID(ctx context.Context, id uint) (*domain.Plan, error) {
	var po PlanPO
	err := r.db.WithContext(ctx).First(&po, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrPlanNotFound
	}
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// FindAll 查询所有
func (r *PlanRepoImpl) FindAll(ctx context.Context, limit, offset int) ([]*domain.Plan, error) {
	var pos []PlanPO
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at DESC").Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	plans := make([]*domain.Plan, len(pos))
	for i, po := range pos {
		plans[i] = po.ToDomain()
	}
	return plans, nil
}

// Update 更新计划
func (r *PlanRepoImpl) Update(ctx context.Context, plan *domain.Plan) error {
	po := FromDomain(plan)
	return r.db.WithContext(ctx).Save(po).Error
}

// Delete 删除计划
func (r *PlanRepoImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&PlanPO{}, id).Error
}

// ExistsByID 检查是否存在
func (r *PlanRepoImpl) ExistsByID(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&PlanPO{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// FindByPlanNo 根据计划编号查询
func (r *PlanRepoImpl) FindByPlanNo(ctx context.Context, planNo string) (*domain.Plan, error) {
	var po PlanPO
	err := r.db.WithContext(ctx).Where("plan_no = ?", planNo).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrPlanNotFound
	}
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// ExistsByPlanNo 检查计划编号是否存在
func (r *PlanRepoImpl) ExistsByPlanNo(ctx context.Context, planNo string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&PlanPO{}).Where("plan_no = ?", planNo).Count(&count).Error
	return count > 0, err
}

// FindByOrderID 根据订单 ID 查询
func (r *PlanRepoImpl) FindByOrderID(ctx context.Context, orderID uint) ([]*domain.Plan, error) {
	var pos []PlanPO
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Order("created_at DESC").Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	plans := make([]*domain.Plan, len(pos))
	for i, po := range pos {
		plans[i] = po.ToDomain()
	}
	return plans, nil
}

// FindByProductID 根据产品 ID 查询
func (r *PlanRepoImpl) FindByProductID(ctx context.Context, productID uint) ([]*domain.Plan, error) {
	var pos []PlanPO
	err := r.db.WithContext(ctx).Where("product_id = ?", productID).Order("created_at DESC").Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	plans := make([]*domain.Plan, len(pos))
	for i, po := range pos {
		plans[i] = po.ToDomain()
	}
	return plans, nil
}

// FindByStatus 根据状态查询
func (r *PlanRepoImpl) FindByStatus(ctx context.Context, status domain.PlanStatus, limit, offset int) ([]*domain.Plan, error) {
	var pos []PlanPO
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	plans := make([]*domain.Plan, len(pos))
	for i, po := range pos {
		plans[i] = po.ToDomain()
	}
	return plans, nil
}

// Count 统计数量
func (r *PlanRepoImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&PlanPO{}).Count(&count).Error
	return count, err
}