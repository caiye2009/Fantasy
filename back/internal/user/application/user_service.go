package application

import (
	"context"
	
	"back/internal/user/domain"
	"back/internal/user/infra"
)

// UserService 用户应用服务
type UserService struct {
	repo *infra.UserRepo
}

// NewUserService 创建用户服务
func NewUserService(repo *infra.UserRepo) *UserService {
	return &UserService{repo: repo}
}

// Create 创建用户
func (s *UserService) Create(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	// 1. 验证角色是否有效
	if !domain.IsValidRole(req.Role) {
		return nil, domain.ErrInvalidRole
	}

	// 2. 自动生成 login_id（自增，从 1000 开始）
	loginID, err := s.repo.GetNextLoginID(ctx)
	if err != nil {
		return nil, err
	}

	// 3. DTO → Domain Model
	user := &domain.User{
		LoginID:    loginID,
		Username:   req.Username,
		Department: req.Department,
		Role:       req.Role,
		Status:     domain.UserStatusActive,
	}

	// 4. 设置默认密码（123）
	if err := user.SetDefaultPassword(); err != nil {
		return nil, err
	}

	// 5. 领域验证
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// 6. 保存到数据库
	if err := s.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	// 7. 返回响应（只返回 login_id）
	return &CreateUserResponse{
		LoginID: user.LoginID,
	}, nil
}

// Get 获取用户
func (s *UserService) Get(ctx context.Context, id uint) (*UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:          user.ID,
		LoginID:     user.LoginID,
		Username:    user.Username,
		Department:  user.Department,
		Email:       user.Email,
		Role:        user.Role,
		Status:      user.Status,
		HasInitPass: user.HasInitPass,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

// GetByLoginID 根据工号获取用户
func (s *UserService) GetByLoginID(ctx context.Context, loginID string) (*domain.User, error) {
	return s.repo.FindByLoginID(ctx, loginID)
}

// List 用户列表
func (s *UserService) List(ctx context.Context, limit, offset int) (*UserListResponse, error) {
	users, err := s.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*UserResponse, len(users))
	for i, u := range users {
		responses[i] = &UserResponse{
			ID:          u.ID,
			LoginID:     u.LoginID,
			Username:    u.Username,
			Department:  u.Department,
			Email:       u.Email,
			Role:        u.Role,
			Status:      u.Status,
			HasInitPass: u.HasInitPass,
			CreatedAt:   u.CreatedAt,
			UpdatedAt:   u.UpdatedAt,
		}
	}

	return &UserListResponse{
		Total: total,
		Users: responses,
	}, nil
}

// Update 更新用户
func (s *UserService) Update(ctx context.Context, id uint, req *UpdateUserRequest) error {
	// 1. 查询用户
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	
	// 2. 更新字段（通过领域方法）
	if req.Username != "" {
		if err := user.UpdateUsername(req.Username); err != nil {
			return err
		}
	}

	if req.Department != "" {
		if err := user.UpdateDepartment(req.Department); err != nil {
			return err
		}
	}

	if req.Email != "" {
		if err := user.UpdateEmail(req.Email); err != nil {
			return err
		}
	}

	if req.Role != "" {
		// 验证角色是否有效
		if !domain.IsValidRole(req.Role) {
			return domain.ErrInvalidRole
		}
		if err := user.UpdateRole(req.Role); err != nil {
			return err
		}
	}
	
	// 状态更新
	if req.Status != "" {
		switch req.Status {
		case domain.UserStatusActive:
			if err := user.Activate(); err != nil {
				return err
			}
		case domain.UserStatusSuspended:
			if err := user.Suspend(); err != nil {
				return err
			}
		case domain.UserStatusInactive:
			user.Status = domain.UserStatusInactive
		}
	}
	
	// 3. 验证
	if err := user.Validate(); err != nil {
		return err
	}
	
	// 4. 保存
	return s.repo.Update(ctx, user)
}

// Delete 删除用户
func (s *UserService) Delete(ctx context.Context, id uint) error {
	// 1. 查询用户
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	
	// 2. 检查是否可以删除（领域规则）
	if !user.CanDelete() {
		return domain.ErrCannotDeleteAdmin
	}
	
	// 3. 删除
	return s.repo.Delete(ctx, id)
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(ctx context.Context, id uint, req *ChangePasswordRequest) error {
	// 1. 查询用户
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	
	// 2. 通过领域方法修改密码（包含所有验证逻辑）
	if err := user.ChangePassword(req.CurrentPassword, req.NewPassword); err != nil {
		return err
	}
	
	// 3. 保存
	return s.repo.Update(ctx, user)
}

// ValidatePassword 验证密码（用于登录）
func (s *UserService) ValidatePassword(ctx context.Context, loginID, password string) (*domain.User, error) {
	// 1. 查询用户
	user, err := s.repo.FindByLoginID(ctx, loginID)
	if err != nil {
		return nil, err
	}
	
	// 2. 验证密码
	if err := user.ValidatePassword(password); err != nil {
		return nil, err
	}
	
	// 3. 检查用户状态
	if !user.IsActive() {
		return nil, domain.ErrUserAlreadySuspended
	}
	
	return user, nil
}

// GetAllDepartments 获取所有部门列表
func (s *UserService) GetAllDepartments(ctx context.Context) ([]string, error) {
	return s.repo.GetAllDepartments(ctx)
}

// GetAllRoles 获取所有角色列表
func (s *UserService) GetAllRoles() []string {
	return domain.GetAllRoles()
}