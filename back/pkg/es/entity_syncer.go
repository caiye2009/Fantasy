package es

import (
	"context"
	"log"
)

// EntitySyncer 实体ES同步助手
// 封装通用的ES同步逻辑，避免在每个service中重复编写
type EntitySyncer[T Indexable] struct {
	esSync *ESSync
}

// NewEntitySyncer 创建实体同步助手
func NewEntitySyncer[T Indexable](esSync *ESSync) *EntitySyncer[T] {
	return &EntitySyncer[T]{
		esSync: esSync,
	}
}

// SyncCreate 创建后同步到ES
// entity: 刚创建的实体
func (s *EntitySyncer[T]) SyncCreate(entity T) {
	if s.esSync == nil {
		return
	}

	if err := s.esSync.Index(entity); err != nil {
		log.Printf("Failed to sync entity to ES: %v", err)
		// 不影响主流程，只记录错误
	}
}

// SyncUpdate 更新后同步到ES
// entity: 更新后的完整实体
func (s *EntitySyncer[T]) SyncUpdate(entity T) {
	if s.esSync == nil {
		return
	}

	if err := s.esSync.Update(entity); err != nil {
		log.Printf("Failed to sync entity update to ES: %v", err)
		// 不影响主流程，只记录错误
	}
}

// SyncDelete 删除时从ES删除
// entity: 被删除的实体（需要在删除前获取）
func (s *EntitySyncer[T]) SyncDelete(entity T) {
	if s.esSync == nil {
		return
	}

	if err := s.esSync.Delete(entity.GetIndexName(), entity.GetDocumentID()); err != nil {
		log.Printf("Failed to delete entity from ES: %v", err)
		// 不影响主流程，只记录错误
	}
}

// SyncUpdateWithFetch 更新后同步到ES（带数据获取）
// 当更新操作只修改了部分字段时，需要重新获取完整实体再同步
// ctx: 上下文
// id: 实体ID
// fetcher: 获取完整实体的函数
func (s *EntitySyncer[T]) SyncUpdateWithFetch(ctx context.Context, id uint, fetcher func(context.Context, uint) (T, error)) {
	if s.esSync == nil {
		return
	}

	entity, err := fetcher(ctx, id)
	if err != nil {
		log.Printf("Failed to fetch entity for ES sync: %v", err)
		return
	}

	s.SyncUpdate(entity)
}
