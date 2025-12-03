package user

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *UserRepo
}

func NewUserService(userRepo *UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

// Create 创建用户 (默认密码 123, has_init_pass = true)
func (s *UserService) Create(req *CreateUserRequest) (*User, error) {
	// 检查工号是否已存在
	existing, _ := s.userRepo.GetByLoginID(req.LoginID)
	if existing != nil && existing.ID > 0 {
		return nil, errors.New("工号已存在")
	}

	// 默认密码: 123
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
		HasInitPass:  true, // 强制首次登录修改密码
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Get(id uint) (*User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) GetByLoginID(loginID string) (*User, error) {
	return s.userRepo.GetByLoginID(loginID)
}

func (s *UserService) List() ([]User, error) {
	return s.userRepo.List()
}

func (s *UserService) Update(id uint, req *UpdateUserRequest) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	data := make(map[string]interface{})
	if req.Username != "" {
		data["username"] = req.Username
	}
	if req.Email != "" {
		data["email"] = req.Email
	}
	if req.Role != "" {
		data["role"] = req.Role
	}
	if req.Status != "" {
		data["status"] = req.Status
	}

	return s.userRepo.Update(id, data)
}

func (s *UserService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}

func (s *UserService) ChangePassword(id uint, req *ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 如果不是初始密码,需要验证当前密码
	if !user.HasInitPass {
		if req.CurrentPassword == "" {
			return errors.New("请提供当前密码")
		}
		
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
			return errors.New("当前密码错误")
		}
	}

	// 生成新密码 hash
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码并设置 has_init_pass = false
	data := map[string]interface{}{
		"password_hash": string(newPasswordHash),
		"has_init_pass": false,
	}

	return s.userRepo.Update(id, data)
}

func (s *UserService) ValidatePassword(user *User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
}