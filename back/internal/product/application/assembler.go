package application

import "back/internal/product/domain"

// ToProduct DTO → Domain Model
func ToProduct(req *CreateProductRequest) *domain.Product {
	return &domain.Product{
		Name:      req.Name,
		Status:    domain.ProductStatusDraft,
		Materials: req.Materials,
		Processes: req.Processes,
	}
}

// ToProductResponse Domain Model → DTO
func ToProductResponse(p *domain.Product) *ProductResponse {
	return &ProductResponse{
		ID:        p.ID,
		Name:      p.Name,
		Status:    string(p.Status),
		Materials: p.Materials,
		Processes: p.Processes,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

// ToProductListResponse Domain Models → List DTO
func ToProductListResponse(products []*domain.Product, total int64) *ProductListResponse {
	responses := make([]*ProductResponse, len(products))
	for i, p := range products {
		responses[i] = ToProductResponse(p)
	}
	
	return &ProductListResponse{
		Total:    total,
		Products: responses,
	}
}