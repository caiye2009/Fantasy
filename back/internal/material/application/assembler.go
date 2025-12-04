package application

import "back/internal/material/domain"

// ToMaterial DTO → Domain Model
func ToMaterial(req *CreateMaterialRequest) *domain.Material {
	return &domain.Material{
		Name:        req.Name,
		Spec:        req.Spec,
		Unit:        req.Unit,
		Description: req.Description,
	}
}

// ToMaterialResponse Domain Model → DTO
func ToMaterialResponse(m *domain.Material) *MaterialResponse {
	return &MaterialResponse{
		ID:          m.ID,
		Name:        m.Name,
		Spec:        m.Spec,
		Unit:        m.Unit,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// ToMaterialListResponse Domain Models → List DTO
func ToMaterialListResponse(materials []*domain.Material, total int64) *MaterialListResponse {
	responses := make([]*MaterialResponse, len(materials))
	for i, m := range materials {
		responses[i] = ToMaterialResponse(m)
	}
	
	return &MaterialListResponse{
		Total:     total,
		Materials: responses,
	}
}