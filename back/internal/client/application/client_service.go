package application

import (
	"context"

	"back/internal/client/domain"
	"back/internal/client/infra"
)

// ClientService 客户应用服务
type ClientService struct {
	repo *infra.ClientRepo
}

// NewClientService 创建服务
func NewClientService(repo *infra.ClientRepo) *ClientService {
	return &ClientService{repo: repo}
}

// Create 创建客户
func (s *ClientService) Create(ctx context.Context, req *CreateClientRequest) (*ClientResponse, error) {
	client := &domain.Client{
		ClientCode:    req.ClientCode,
		ClientName:    req.ClientName,
		ClientCountry: req.ClientCountry,
		CreatedBy:     req.CreatedBy,
	}

	if err := s.repo.Create(ctx, client); err != nil {
		return nil, err
	}

	return toClientResponse(client), nil
}

// Get 获取客户
func (s *ClientService) Get(ctx context.Context, id uint) (*ClientResponse, error) {
	client, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrClientNotFound
	}
	return toClientResponse(client), nil
}

// List 列表查询
func (s *ClientService) List(ctx context.Context, limit, offset int) ([]*ClientResponse, int64, error) {
	clients, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*ClientResponse, len(clients))
	for i, client := range clients {
		responses[i] = toClientResponse(&client)
	}

	return responses, count, nil
}

// Update 部分更新
func (s *ClientService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 检查是否存在
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrClientNotFound
	}

	// 过滤不允许更新的字段
	delete(updates, "id")
	delete(updates, "createdAt")
	delete(updates, "createdBy")
	delete(updates, "deletedAt")

	// 如果没有要更新的字段，直接返回
	if len(updates) == 0 {
		return nil
	}

	// 部分更新
	return s.repo.UpdateFields(ctx, id, updates)
}

// Delete 删除客户
func (s *ClientService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrClientNotFound
	}

	return s.repo.Delete(ctx, id)
}

// DTO 转换
func toClientResponse(client *domain.Client) *ClientResponse {
	return &ClientResponse{
		ID:            client.ID,
		ClientCode:    client.ClientCode,
		ClientName:    client.ClientName,
		ClientCountry: client.ClientCountry,
		CreatedBy:     client.CreatedBy,
		CreatedAt:     client.CreatedAt,
		UpdatedAt:     client.UpdatedAt,
	}
}
