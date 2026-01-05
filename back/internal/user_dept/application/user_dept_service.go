package application

import (
	"context"

	"back/internal/user_dept/domain"
	"back/internal/user_dept/infra"
	"back/pkg/auth"
)

type UserDeptService struct {
	repo *infra.UserDeptRepo
}

func NewUserDeptService(repo *infra.UserDeptRepo) *UserDeptService {
	return &UserDeptService{repo: repo}
}

func (s *UserDeptService) Create(ctx context.Context, req *CreateUserDeptRequest) (*UserDeptResponse, error) {
	createdBy := auth.GetLoginIDFromContext(ctx)

	// 检查部门代码是否已存在
	exists, err := s.repo.ExistsByCode(ctx, req.DeptCode)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrUserDeptCodeExists
	}

	userDept := &domain.UserDept{
		DeptCode:    req.DeptCode,
		DeptName:    req.DeptName,
		Description: req.Description,
		CreatedBy:   createdBy,
	}

	if err := s.repo.Create(ctx, userDept); err != nil {
		return nil, err
	}

	return toUserDeptResponse(userDept), nil
}

func (s *UserDeptService) Get(ctx context.Context, id uint) (*UserDeptResponse, error) {
	userDept, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrUserDeptNotFound
	}
	return toUserDeptResponse(userDept), nil
}

func (s *UserDeptService) List(ctx context.Context, limit, offset int) ([]*UserDeptResponse, int64, error) {
	userDepts, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*UserDeptResponse, len(userDepts))
	for i, userDept := range userDepts {
		responses[i] = toUserDeptResponse(&userDept)
	}

	return responses, count, nil
}

func (s *UserDeptService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrUserDeptNotFound
	}

	delete(updates, "id")
	delete(updates, "createdAt")
	delete(updates, "createdBy")
	delete(updates, "deletedAt")

	if len(updates) == 0 {
		return nil
	}

	return s.repo.UpdateFields(ctx, id, updates)
}

func (s *UserDeptService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrUserDeptNotFound
	}

	return s.repo.Delete(ctx, id)
}

func toUserDeptResponse(userDept *domain.UserDept) *UserDeptResponse {
	return &UserDeptResponse{
		ID:          userDept.ID,
		DeptCode:    userDept.DeptCode,
		DeptName:    userDept.DeptName,
		Description: userDept.Description,
		CreatedBy:   userDept.CreatedBy,
		CreatedAt:   userDept.CreatedAt,
		UpdatedAt:   userDept.UpdatedAt,
	}
}
