package application

import "back/internal/process/domain"

// ToProcess DTO → Domain Model
func ToProcess(req *CreateProcessRequest) *domain.Process {
	return &domain.Process{
		Name:        req.Name,
		Description: req.Description,
	}
}

// ToProcessResponse Domain Model → DTO
func ToProcessResponse(p *domain.Process) *ProcessResponse {
	return &ProcessResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// ToProcessListResponse Domain Models → List DTO
func ToProcessListResponse(processes []*domain.Process, total int64) *ProcessListResponse {
	responses := make([]*ProcessResponse, len(processes))
	for i, p := range processes {
		responses[i] = ToProcessResponse(p)
	}
	
	return &ProcessListResponse{
		Total:     total,
		Processes: responses,
	}
}