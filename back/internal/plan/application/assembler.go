package application

import "back/internal/plan/domain"

// ToPlan DTO → Domain Model
func ToPlan(req *CreatePlanRequest) *domain.Plan {
	return &domain.Plan{
		PlanNo:      req.PlanNo,
		OrderID:     req.OrderID,
		ProductID:   req.ProductID,
		Quantity:    req.Quantity,
		Status:      domain.PlanStatusPlanned,
		ScheduledAt: req.ScheduledAt,
		CreatedBy:   req.CreatedBy,
	}
}

// ToPlanResponse Domain Model → DTO
func ToPlanResponse(p *domain.Plan) *PlanResponse {
	return &PlanResponse{
		ID:          p.ID,
		PlanNo:      p.PlanNo,
		OrderID:     p.OrderID,
		ProductID:   p.ProductID,
		Quantity:    p.Quantity,
		Status:      p.Status,
		ScheduledAt: p.ScheduledAt,
		CompletedAt: p.CompletedAt,
		CreatedBy:   p.CreatedBy,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// ToPlanListResponse Domain Models → List DTO
func ToPlanListResponse(plans []*domain.Plan, total int64) *PlanListResponse {
	responses := make([]*PlanResponse, len(plans))
	for i, p := range plans {
		responses[i] = ToPlanResponse(p)
	}
	
	return &PlanListResponse{
		Total: total,
		Plans: responses,
	}
}