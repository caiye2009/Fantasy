package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"back/internal/user/domain"
	"back/internal/user/infra"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *infra.UserRepo
}

func NewUserService(repo *infra.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, req *UserCreateRequest) (*UserCreateResponse, error) {
	// 生成 LoginID (格式: U + 时间戳后6位)
	loginID := generateLoginID()

	// 默认密码
	defaultPassword := "123"

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户实体
	user := &domain.User{
		LoginID:     loginID,
		Username:    req.Username,
		Password:    string(hashedPassword),
		RealName:    req.Username, // RealName 默认等于 Username
		Email:       "",           // 默认为空
		Role:        req.Role,
		Department:  req.Department,
		IsActive:    true, // 默认激活
		HasInitPass: true, // 需要修改初始密码
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// 返回 LoginID 和用户名给前端
	return &UserCreateResponse{
		LoginID:  loginID,
		Username: req.Username,
	}, nil
}

// generateLoginID 生成登录ID
func generateLoginID() string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("U%d", timestamp%1000000)
}

func (s *UserService) Get(ctx context.Context, id uint) (*UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	return toUserResponse(user), nil
}

func (s *UserService) List(ctx context.Context, limit, offset int) ([]*UserResponse, int64, error) {
	users, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		responses[i] = toUserResponse(&user)
	}

	return responses, count, nil
}

func (s *UserService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrUserNotFound
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

func (s *UserService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrUserNotFound
	}

	return s.repo.Delete(ctx, id)
}

// GetByLoginID retrieves a user by LoginID (required by AuthService)
func (s *UserService) GetByLoginID(ctx context.Context, loginID string) (*domain.User, error) {
	return s.repo.GetByLoginID(ctx, loginID)
}

// ValidatePassword validates user credentials (required by AuthService)
func (s *UserService) ValidatePassword(ctx context.Context, loginID, password string) (*domain.User, error) {
	user, err := s.GetByLoginID(ctx, loginID)
	if err != nil {
		return nil, err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func toUserResponse(user *domain.User) *UserResponse {
	return &UserResponse{
		ID:          user.ID,
		LoginID:     user.LoginID,
		Username:    user.Username,
		Role:        user.Role,
		Department:  user.Department,
		IsActive:    user.IsActive,
		HasInitPass: user.HasInitPass,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
