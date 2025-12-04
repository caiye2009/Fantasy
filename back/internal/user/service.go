package user

import (
	"context"
	"errors"
	
	"back/pkg/repo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Create(ctx context.Context, req *CreateUserRequest) (*User, error) {
	userRepo := repo.NewRepo[User](s.db)
	
	existing, _ := userRepo.First(ctx, map[string]interface{}{"login_id": req.LoginID})
	if existing != nil && existing.ID > 0 {
		return nil, errors.New("工号已存在")
	}

	defaultPassword := "123"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		LoginID:      req.LoginID,
		Username:     req.Username,
		PasswordHash: string(passwordHash),
		Email:        req.Email,
		Role:         req.Role,
		Status:       "active",
		HasInitPass:  true,
	}

	if err := userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Get(ctx context.Context, id uint) (*User, error) {
	userRepo := repo.NewRepo[User](s.db)
	return userRepo.GetByID(ctx, id)
}

func (s *UserService) GetByLoginID(ctx context.Context, loginID string) (*User, error) {
	userRepo := repo.NewRepo[User](s.db)
	return userRepo.First(ctx, map[string]interface{}{"login_id": loginID})
}

func (s *UserService) List(ctx context.Context, limit, offset int) ([]User, error) {
	userRepo := repo.NewRepo[User](s.db)
	return userRepo.List(ctx, limit, offset)
}

func (s *UserService) Update(ctx context.Context, id uint, req *UpdateUserRequest) error {
	userRepo := repo.NewRepo[User](s.db)
	
	user, err := userRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("用户不存在")
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	fields := make(map[string]interface{})
	if req.Username != "" {
		fields["username"] = req.Username
	}
	if req.Email != "" {
		fields["email"] = req.Email
	}
	if req.Role != "" {
		fields["role"] = req.Role
	}
	if req.Status != "" {
		fields["status"] = req.Status
	}

	return userRepo.UpdateFields(ctx, id, fields)
}

func (s *UserService) Delete(ctx context.Context, id uint) error {
	userRepo := repo.NewRepo[User](s.db)
	return userRepo.Delete(ctx, id)
}

func (s *UserService) ChangePassword(ctx context.Context, id uint, req *ChangePasswordRequest) error {
	userRepo := repo.NewRepo[User](s.db)
	
	user, err := userRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("用户不存在")
	}

	if !user.HasInitPass {
		if req.CurrentPassword == "" {
			return errors.New("请提供当前密码")
		}
		
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
			return errors.New("当前密码错误")
		}
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	fields := map[string]interface{}{
		"password_hash": string(newPasswordHash),
		"has_init_pass": false,
	}

	return userRepo.UpdateFields(ctx, id, fields)
}

func (s *UserService) ValidatePassword(user *User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
}