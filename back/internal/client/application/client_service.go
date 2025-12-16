package application

import (
	"context"
	"strconv"

	"back/internal/client/domain"
	"back/internal/client/infra"
)

// ESSync ES 同步接口
type ESSync interface {
	Index(doc interface{}) error
	Update(doc interface{}) error
	Delete(indexName, docID string) error
}

// ClientService 客户应用服务
type ClientService struct {
	repo   *infra.ClientRepo
	esSync ESSync
}

// NewClientService 创建客户服务
func NewClientService(repo *infra.ClientRepo, esSync ESSync) *ClientService {
	return &ClientService{
		repo:   repo,
		esSync: esSync,
	}
}

// Create 创建客户
func (s *ClientService) Create(ctx context.Context, req *CreateClientRequest) (*ClientResponse, error) {
	// 1. 检查名称是否重复（可选，根据业务需求）
	exists, err := s.repo.ExistsByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrClientNameExists
	}
	
	// 2. DTO → Domain Model
	client := &domain.Client{
		Name:    req.Name,
		Contact: req.Contact,
		Phone:   req.Phone,
		Email:   req.Email,
		Address: req.Address,
	}

	// 3. 领域验证
	if err := client.Validate(); err != nil {
		return nil, err
	}

	// 4. 保存到数据库
	if err := s.repo.Save(ctx, client); err != nil {
		return nil, err
	}

	// 5. 异步同步到 ES
	if s.esSync != nil {
		s.esSync.Index(client)
	}

	// 6. Domain Model → DTO
	return &ClientResponse{
		ID:        client.ID,
		Name:      client.Name,
		Contact:   client.Contact,
		Phone:     client.Phone,
		Email:     client.Email,
		Address:   client.Address,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}, nil
}

// Get 获取客户
func (s *ClientService) Get(ctx context.Context, id uint) (*ClientResponse, error) {
	client, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &ClientResponse{
		ID:        client.ID,
		Name:      client.Name,
		Contact:   client.Contact,
		Phone:     client.Phone,
		Email:     client.Email,
		Address:   client.Address,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}, nil
}

// Update 更新客户
func (s *ClientService) Update(ctx context.Context, id uint, req *UpdateClientRequest) error {
	// 1. 查询客户
	client, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	
	// 2. 更新字段（通过领域方法）
	if req.Name != "" {
		// 检查新名称是否重复
		if req.Name != client.Name {
			exists, err := s.repo.ExistsByName(ctx, req.Name)
			if err != nil {
				return err
			}
			if exists {
				return domain.ErrClientNameExists
			}
		}
		
		if err := client.UpdateName(req.Name); err != nil {
			return err
		}
	}
	
	if req.Contact != "" {
		if err := client.UpdateContact(req.Contact); err != nil {
			return err
		}
	}
	
	if req.Phone != "" {
		if err := client.UpdatePhone(req.Phone); err != nil {
			return err
		}
	}
	
	if req.Email != "" {
		if err := client.UpdateEmail(req.Email); err != nil {
			return err
		}
	}
	
	if req.Address != "" {
		if err := client.UpdateAddress(req.Address); err != nil {
			return err
		}
	}
	
	// 3. 验证
	if err := client.Validate(); err != nil {
		return err
	}
	
	// 4. 保存
	if err := s.repo.Update(ctx, client); err != nil {
		return err
	}
	
	// 5. 异步同步到 ES
	if s.esSync != nil {
		s.esSync.Update(client)
	}
	
	return nil
}

// Delete 删除客户
func (s *ClientService) Delete(ctx context.Context, id uint) error {
	// 1. 检查是否存在
	exists, err := s.repo.ExistsByID(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrClientNotFound
	}
	
	// 2. 删除
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	
	// 3. 异步删除 ES 文档
	if s.esSync != nil {
		s.esSync.Delete("clients", strconv.Itoa(int(id)))
	}
	
	return nil
}

// Exists 检查客户是否存在（供其他模块调用）
func (s *ClientService) Exists(ctx context.Context, id uint) (bool, error) {
	return s.repo.ExistsByID(ctx, id)
}