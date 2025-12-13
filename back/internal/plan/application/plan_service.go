package application

import (
	"context"
	"strconv"
	
	"back/internal/plan/domain"
	"back/internal/plan/infra"
)

// ESSync ES 同步接口
type ESSync interface {
	Index(doc interface{}) error
	Update(doc interface{}) error
	Delete(indexName, docID string) error
}

// PlanService 计划应用服务
type PlanService struct {
	repo   *infra.PlanRepo
	esSync ESSync
}

// NewPlanService 创建计划服务
func NewPlanService(repo *infra.PlanRepo, esSync ESSync) *PlanService {
	return &PlanService{
		repo:   repo,
		esSync: esSync,
	}
}

// Create 创建计划
func (s *PlanService) Create(ctx context.Context, req *CreatePlanRequest) (*PlanResponse, error) {
	// 1. 检查计划编号是否重复
	exists, err := s.repo.ExistsByPlanNo(ctx, req.PlanNo)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrPlanNoDuplicate
	}
	
	// 2. DTO → Domain Model
	plan := ToPlan(req)
	
	// 3. 领域验证
	if err := plan.Validate(); err != nil {
		return nil, err
	}
	
	// 4. 保存到数据库
	if err := s.repo.Save(ctx, plan); err != nil {
		return nil, err
	}
	
	// 5. 异步同步到 ES
	if s.esSync != nil {
		s.esSync.Index(plan)
	}
	
	// 6. Domain Model → DTO
	return ToPlanResponse(plan), nil
}

// Get 获取计划
func (s *PlanService) Get(ctx context.Context, id uint) (*PlanResponse, error) {
	plan, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return ToPlanResponse(plan), nil
}

// List 计划列表
func (s *PlanService) List(ctx context.Context, limit, offset int) (*PlanListResponse, error) {
	plans, err := s.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	
	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}
	
	return ToPlanListResponse(plans, total), nil
}

// Update 更新计划
func (s *PlanService) Update(ctx context.Context, id uint, req *UpdatePlanRequest) error {
	// 1. 查询计划
	plan, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	
	// 2. 更新字段（通过领域方法）
	if req.Status != "" {
		// 根据状态转换调用不同的领域方法
		switch req.Status {
		case domain.PlanStatusInProgress:
			if err := plan.Start(); err != nil {
				return err
			}
		case domain.PlanStatusCompleted:
			if err := plan.Complete(); err != nil {
				return err
			}
		case domain.PlanStatusCancelled:
			if err := plan.Cancel(); err != nil {
				return err
			}
		}
	}
	
	if req.Quantity > 0 {
		if err := plan.UpdateQuantity(req.Quantity); err != nil {
			return err
		}
	}
	
	if req.CompletedAt != nil {
		plan.CompletedAt = req.CompletedAt
	}
	
	// 3. 验证
	if err := plan.Validate(); err != nil {
		return err
	}
	
	// 4. 保存
	if err := s.repo.Update(ctx, plan); err != nil {
		return err
	}
	
	// 5. 异步同步到 ES
	if s.esSync != nil {
		s.esSync.Update(plan)
	}
	
	return nil
}

// Delete 删除计划
func (s *PlanService) Delete(ctx context.Context, id uint) error {
	// 1. 查询计划
	plan, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	
	// 2. 检查是否可以删除（领域规则）
	if !plan.CanDelete() {
		return domain.ErrCannotDeleteCompletedPlan
	}
	
	// 3. 删除
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	
	// 4. 异步删除 ES 文档
	if s.esSync != nil {
		s.esSync.Delete("plans", strconv.Itoa(int(id)))
	}
	
	return nil
}

// GetByOrderID 根据订单 ID 查询计划
func (s *PlanService) GetByOrderID(ctx context.Context, orderID uint) ([]*PlanResponse, error) {
	plans, err := s.repo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	
	responses := make([]*PlanResponse, len(plans))
	for i, p := range plans {
		responses[i] = ToPlanResponse(p)
	}
	
	return responses, nil
}

// GetByProductID 根据产品 ID 查询计划
func (s *PlanService) GetByProductID(ctx context.Context, productID uint) ([]*PlanResponse, error) {
	plans, err := s.repo.FindByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}
	
	responses := make([]*PlanResponse, len(plans))
	for i, p := range plans {
		responses[i] = ToPlanResponse(p)
	}
	
	return responses, nil
}

// GetByStatus 根据状态查询计划
func (s *PlanService) GetByStatus(ctx context.Context, status string, limit, offset int) (*PlanListResponse, error) {
	plans, err := s.repo.FindByStatus(ctx, status, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	return ToPlanListResponse(plans, total), nil
}