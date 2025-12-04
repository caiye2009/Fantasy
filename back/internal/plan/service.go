package plan

import (
	"context"
	"strconv"
	
	"back/pkg/es"
	"back/pkg/repo"
	"gorm.io/gorm"
)

type PlanService struct {
	db     *gorm.DB
	esSync *es.ESSync
}

func NewPlanService(db *gorm.DB, esSync *es.ESSync) *PlanService {
	return &PlanService{
		db:     db,
		esSync: esSync,
	}
}

// Create 创建计划
func (s *PlanService) Create(ctx context.Context, req *CreatePlanRequest) (*Plan, error) {
	plan := &Plan{
		PlanNo:      req.PlanNo,
		OrderID:     req.OrderID,
		ProductID:   req.ProductID,
		Quantity:    req.Quantity,
		Status:      "planned",
		ScheduledAt: req.ScheduledAt,
		CreatedBy:   req.CreatedBy,
	}

	planRepo := repo.NewRepo[Plan](s.db)
	if err := planRepo.Create(ctx, plan); err != nil {
		return nil, err
	}

	// 异步同步到 ES
	s.esSync.Index(plan)

	return plan, nil
}

// Get 获取计划
func (s *PlanService) Get(ctx context.Context, id uint) (*Plan, error) {
	planRepo := repo.NewRepo[Plan](s.db)
	return planRepo.GetByID(ctx, id)
}

// List 计划列表
func (s *PlanService) List(ctx context.Context, limit, offset int) ([]Plan, error) {
	planRepo := repo.NewRepo[Plan](s.db)
	return planRepo.List(ctx, limit, offset)
}

// Update 更新计划
func (s *PlanService) Update(ctx context.Context, id uint, req *UpdatePlanRequest) error {
	planRepo := repo.NewRepo[Plan](s.db)
	
	// 1. 查询计划
	plan, err := planRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 2. 构造更新字段
	fields := make(map[string]interface{})
	if req.Status != "" {
		fields["status"] = req.Status
		plan.Status = req.Status
	}
	if req.Quantity > 0 {
		fields["quantity"] = req.Quantity
		plan.Quantity = req.Quantity
	}
	if req.CompletedAt != nil {
		fields["completed_at"] = req.CompletedAt
		plan.CompletedAt = req.CompletedAt
	}

	// 3. 如果没有要更新的字段，直接返回
	if len(fields) == 0 {
		return nil
	}

	// 4. 更新数据库
	if err := planRepo.UpdateFields(ctx, id, fields); err != nil {
		return err
	}

	// 5. 异步同步到 ES
	s.esSync.Update(plan)

	return nil
}

// Delete 删除计划
func (s *PlanService) Delete(ctx context.Context, id uint) error {
	planRepo := repo.NewRepo[Plan](s.db)
	
	if err := planRepo.Delete(ctx, id); err != nil {
		return err
	}

	// 异步删除 ES 文档
	s.esSync.Delete("plans", strconv.Itoa(int(id)))

	return nil
}