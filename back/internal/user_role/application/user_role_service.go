package application

import (
	"context"

	"back/internal/user_role/domain"
	"back/internal/user_role/infra"
)

type UserRoleService struct {
	repo *infra.UserRoleRepo
}

func NewUserRoleService(repo *infra.UserRoleRepo) *UserRoleService {
	return &UserRoleService{repo: repo}
}

func (s *UserRoleService) Create(ctx context.Context, req *CreateUserRoleRequest) (*UserRoleResponse, error) {
	// 检查角色代码是否已存在
	exists, err := s.repo.ExistsByCode(ctx, req.RoleCode)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrUserRoleCodeExists
	}

	userRole := &domain.UserRole{
		RoleCode:    req.RoleCode,
		RoleName:    req.RoleName,
		Description: req.Description,
		CreatedBy:   req.CreatedBy,
	}

	if err := s.repo.Create(ctx, userRole); err != nil {
		return nil, err
	}

	return toUserRoleResponse(userRole), nil
}

func (s *UserRoleService) Get(ctx context.Context, id uint) (*UserRoleResponse, error) {
	userRole, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrUserRoleNotFound
	}
	return toUserRoleResponse(userRole), nil
}

func (s *UserRoleService) List(ctx context.Context, limit, offset int) ([]*UserRoleResponse, int64, error) {
	userRoles, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*UserRoleResponse, len(userRoles))
	for i, userRole := range userRoles {
		responses[i] = toUserRoleResponse(&userRole)
	}

	return responses, count, nil
}

func (s *UserRoleService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrUserRoleNotFound
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

func (s *UserRoleService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrUserRoleNotFound
	}

	return s.repo.Delete(ctx, id)
}

func toUserRoleResponse(userRole *domain.UserRole) *UserRoleResponse {
	return &UserRoleResponse{
		ID:          userRole.ID,
		RoleCode:    userRole.RoleCode,
		RoleName:    userRole.RoleName,
		Description: userRole.Description,
		CreatedBy:   userRole.CreatedBy,
		CreatedAt:   userRole.CreatedAt,
		UpdatedAt:   userRole.UpdatedAt,
	}
}
