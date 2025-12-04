package repo

import (
	"context"
	"gorm.io/gorm"
)

// Repo 通用泛型 Repository
type Repo[T any] struct {
	db *gorm.DB
}

// NewRepo 创建 Repository
func NewRepo[T any](db *gorm.DB) *Repo[T] {
	return &Repo[T]{db: db}
}

// ──────────────────────────────────────
// 基础 CRUD（带 context）
// ──────────────────────────────────────

// Create 创建
func (r *Repo[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// GetByID 根据ID查询
func (r *Repo[T]) GetByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// List 列表查询（带分页）
func (r *Repo[T]) List(ctx context.Context, limit, offset int) ([]T, error) {
	var list []T
	query := r.db.WithContext(ctx)
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	err := query.Find(&list).Error
	return list, err
}

// Update 更新
func (r *Repo[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// UpdateFields 更新部分字段
func (r *Repo[T]) UpdateFields(ctx context.Context, id uint, fields map[string]interface{}) error {
	var entity T
	return r.db.WithContext(ctx).
		Model(&entity).
		Where("id = ?", id).
		Updates(fields).Error
}

// Delete 删除
func (r *Repo[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, id).Error
}

// ──────────────────────────────────────
// 高级查询
// ──────────────────────────────────────

// FindWhere 条件查询
func (r *Repo[T]) FindWhere(ctx context.Context, condition string, args ...interface{}) ([]T, error) {
	var list []T
	err := r.db.WithContext(ctx).Where(condition, args...).Find(&list).Error
	return list, err
}

// First 查询第一条
func (r *Repo[T]) First(ctx context.Context, query map[string]interface{}) (*T, error) {
	var entity T
	db := r.db.WithContext(ctx)
	
	for k, v := range query {
		db = db.Where(k+" = ?", v)
	}
	
	err := db.First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Count 统计数量
func (r *Repo[T]) Count(ctx context.Context, query map[string]interface{}) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(new(T))
	
	for k, v := range query {
		db = db.Where(k+" = ?", v)
	}
	
	err := db.Count(&count).Error
	return count, err
}

// Exists 检查是否存在
func (r *Repo[T]) Exists(ctx context.Context, query map[string]interface{}) (bool, error) {
	count, err := r.Count(ctx, query)
	return count > 0, err
}

// ──────────────────────────────────────
// 批量操作
// ──────────────────────────────────────

// BatchCreate 批量创建
func (r *Repo[T]) BatchCreate(ctx context.Context, entities []T) error {
	return r.db.WithContext(ctx).Create(&entities).Error
}

// BatchDelete 批量删除
func (r *Repo[T]) BatchDelete(ctx context.Context, ids []uint) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, ids).Error
}

// ──────────────────────────────────────
// 事务支持
// ──────────────────────────────────────

// Transaction 执行事务
func (r *Repo[T]) Transaction(ctx context.Context, fn func(*gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

// WithTx 返回带事务的 Repo
func (r *Repo[T]) WithTx(tx *gorm.DB) *Repo[T] {
	return &Repo[T]{db: tx}
}

// ──────────────────────────────────────
// 原生 SQL 支持
// ──────────────────────────────────────

// Raw 原生查询
func (r *Repo[T]) Raw(ctx context.Context, sql string, dest interface{}, args ...interface{}) error {
	return r.db.WithContext(ctx).Raw(sql, args...).Scan(dest).Error
}

// Exec 执行原生 SQL
func (r *Repo[T]) Exec(ctx context.Context, sql string, args ...interface{}) error {
	return r.db.WithContext(ctx).Exec(sql, args...).Error
}

// ──────────────────────────────────────
// 预加载支持
// ──────────────────────────────────────

// GetByIDWithPreload 带预加载的查询
func (r *Repo[T]) GetByIDWithPreload(ctx context.Context, id uint, preloads ...string) (*T, error) {
	var entity T
	db := r.db.WithContext(ctx)
	
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	
	err := db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// ListWithPreload 带预加载的列表查询
func (r *Repo[T]) ListWithPreload(ctx context.Context, limit, offset int, preloads ...string) ([]T, error) {
	var list []T
	db := r.db.WithContext(ctx)
	
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	
	if limit > 0 {
		db = db.Limit(limit)
	}
	if offset > 0 {
		db = db.Offset(offset)
	}
	
	err := db.Find(&list).Error
	return list, err
}