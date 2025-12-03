package repo

import "gorm.io/gorm"

type Repo[T any] struct {
	DB *gorm.DB
}

func NewRepo[T any](db *gorm.DB) *Repo[T] {
	return &Repo[T]{DB: db}
}

func (r *Repo[T]) Create(entity *T) error {
	return r.DB.Create(entity).Error
}

func (r *Repo[T]) GetByID(id uint) (*T, error) {
	var entity T
	err := r.DB.First(&entity, id).Error
	return &entity, err
}

func (r *Repo[T]) List() ([]T, error) {
	var list []T
	err := r.DB.Find(&list).Error
	return list, err
}

func (r *Repo[T]) Update(id uint, data map[string]interface{}) error {
	var entity T
	return r.DB.Model(&entity).Where("id = ?", id).Updates(data).Error
}

func (r *Repo[T]) Delete(id uint) error {
	var entity T
	return r.DB.Delete(&entity, id).Error
}

func (r *Repo[T]) FindWhere(condition string, args ...interface{}) ([]T, error) {
	var list []T
	err := r.DB.Where(condition, args...).Find(&list).Error
	return list, err
}