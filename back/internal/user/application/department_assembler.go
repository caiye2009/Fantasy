package application

import (
	"time"
	"back/internal/user/domain"
)

// ToDepartmentDomain 将创建请求转换为领域对象
func ToDepartmentDomain(req *CreateDepartmentRequest) *domain.Department {
	return &domain.Department{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      domain.DepartmentStatusActive,
		ParentID:    req.ParentID,
	}
}

// ToDepartmentResponse 将领域对象转换为响应
func ToDepartmentResponse(dept *domain.Department) *DepartmentResponse {
	var deletedAt *time.Time
	if dept.DeletedAt.Valid {
		deletedAt = &dept.DeletedAt.Time
	}

	return &DepartmentResponse{
		ID:          dept.ID,
		Name:        dept.Name,
		Code:        dept.Code,
		Description: dept.Description,
		Status:      dept.Status,
		ParentID:    dept.ParentID,
		CreatedAt:   dept.CreatedAt,
		UpdatedAt:   dept.UpdatedAt,
		DeletedAt:   deletedAt,
	}
}

// ToDepartmentListResponse 将领域对象列表转换为响应
func ToDepartmentListResponse(departments []*domain.Department, total int64) *DepartmentListResponse {
	responses := make([]*DepartmentResponse, len(departments))
	for i, dept := range departments {
		responses[i] = ToDepartmentResponse(dept)
	}
	return &DepartmentListResponse{
		Total:       total,
		Departments: responses,
	}
}
