package application

import "back/internal/supplier/domain"

// ToSupplier DTO → Domain Model
func ToSupplier(req *CreateSupplierRequest) *domain.Supplier {
	return &domain.Supplier{
		Name:    req.Name,
		Contact: req.Contact,
		Phone:   req.Phone,
		Email:   req.Email,
		Address: req.Address,
	}
}

// ToSupplierResponse Domain Model → DTO
func ToSupplierResponse(s *domain.Supplier) *SupplierResponse {
	return &SupplierResponse{
		ID:        s.ID,
		Name:      s.Name,
		Contact:   s.Contact,
		Phone:     s.Phone,
		Email:     s.Email,
		Address:   s.Address,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

// ToSupplierListResponse Domain Models → List DTO
func ToSupplierListResponse(suppliers []*domain.Supplier, total int64) *SupplierListResponse {
	responses := make([]*SupplierResponse, len(suppliers))
	for i, s := range suppliers {
		responses[i] = ToSupplierResponse(s)
	}

	return &SupplierListResponse{
		Total:     total,
		Suppliers: responses,
	}
}
