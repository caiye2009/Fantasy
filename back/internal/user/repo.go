package user

import (
	"gorm.io/gorm"
	"back/pkg/repo"
)

type UserRepo struct {
	*repo.Repo[User]
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		Repo: repo.NewRepo[User](db),
	}
}

func (r *UserRepo) GetByLoginID(loginID string) (*User, error) {
	var user User
	err := r.DB.Where("login_id = ?", loginID).First(&user).Error
	return &user, err
}

func (r *UserRepo) ListByRole(role string) ([]User, error) {
	var users []User
	err := r.DB.Where("role = ?", role).Find(&users).Error
	return users, err
}

func (r *UserRepo) ListActive() ([]User, error) {
	var users []User
	err := r.DB.Where("status = ?", "active").Find(&users).Error
	return users, err
}