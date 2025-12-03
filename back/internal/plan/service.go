package plan

import (
	"strconv"
	"back/pkg/es"
)

type PlanService struct {
	planRepo *PlanRepo
	esSync   *es.ESSync
}

func NewPlanService(planRepo *PlanRepo, esSync *es.ESSync) *PlanService {
	return &PlanService{
		planRepo: planRepo,
		esSync:   esSync,
	}
}

func (s *PlanService) Create(req *CreatePlanRequest) (*Plan, error) {
	plan := &Plan{
		PlanNo:      req.PlanNo,
		OrderID:     req.OrderID,
		ProductID:   req.ProductID,
		Quantity:    req.Quantity,
		Status:      "planned",
		ScheduledAt: req.ScheduledAt,
		CreatedBy:   req.CreatedBy,
	}

	if err := s.planRepo.Create(plan); err != nil {
		return nil, err
	}

	// 异步同步到 ES
	s.esSync.Index(plan)

	return plan, nil
}

func (s *PlanService) Get(id uint) (*Plan, error) {
	return s.planRepo.GetByID(id)
}

func (s *PlanService) List() ([]Plan, error) {
	return s.planRepo.List()
}

func (s *PlanService) Update(id uint, req *UpdatePlanRequest) error {
	plan, err := s.planRepo.GetByID(id)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	if req.Status != "" {
		data["status"] = req.Status
		plan.Status = req.Status
	}
	if req.Quantity > 0 {
		data["quantity"] = req.Quantity
		plan.Quantity = req.Quantity
	}
	if req.CompletedAt != nil {
		data["completed_at"] = req.CompletedAt
		plan.CompletedAt = req.CompletedAt
	}

	if err := s.planRepo.Update(id, data); err != nil {
		return err
	}

	// 异步同步到 ES
	s.esSync.Update(plan)

	return nil
}

func (s *PlanService) Delete(id uint) error {
	if err := s.planRepo.Delete(id); err != nil {
		return err
	}

	// 异步删除 ES 文档
	s.esSync.Delete("plans", strconv.Itoa(int(id)))

	return nil
}