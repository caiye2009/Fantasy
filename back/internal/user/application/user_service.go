package application

import (
	"context"
	
	"back/internal/user/domain"
)

// UserService 用户应用服务
type UserService struct {
	repo domain.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Create 创建用户
func (s *UserService) Create(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	// 1. 检查工号是否重复
	exists, err := s.repo.ExistsByLoginID(ctx, req.LoginID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrLoginIDDuplicate
	}
	
	// 2. DTO → Domain Model
	user := ToUser(req)
	
	// 3. 设置默认密码
	if err := user.SetDefaultPassword(); err != nil {
		return nil, err
	}
	
	// 4. 领域验证
	if err := user.Validate(); err != nil {
		return nil, err
	}
	
	// 5. 保存到数据库
	if err := s.repo.Save(ctx, user); err != nil {
		return nil, err
	}
	
	// 6. 返回响应（包含初始密码）
	return &CreateUserResponse{
		Message:         "用户创建成功",
		LoginID:         user.LoginID,
		DefaultPassword: "123",
		User:            ToUserResponse(user),
	}, nil
}

// Get 获取用户
func (s *UserService) Get(ctx context.Context, id uint) (*UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return ToUserResponse(user), nil
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
	
	return ToUserListResponse(users, total), nil
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
	
	if req.Email != "" {
		if err := user.UpdateEmail(req.Email); err != nil {
			return err
		}
	}
	
	if req.Role != "" {
		if err := user.UpdateRole(domain.UserRole(req.Role)); err != nil {
			return err
		}
	}
	
	// 状态更新
	if req.Status != "" {
		targetStatus := domain.UserStatus(req.Status)
		
		switch targetStatus {
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