package client

import (
	"context"
	"strconv"
	
	"back/pkg/es"
	"back/pkg/repo"
	"gorm.io/gorm"
)

type ClientService struct {
	db     *gorm.DB
	esSync *es.ESSync
}

func NewClientService(db *gorm.DB, esSync *es.ESSync) *ClientService {
	return &ClientService{
		db:     db,
		esSync: esSync,
	}
}

// Create 创建客户
func (s *ClientService) Create(ctx context.Context, req *CreateClientRequest) (*Client, error) {
	client := &Client{
		Name:    req.Name,
		Contact: req.Contact,
		Phone:   req.Phone,
		Email:   req.Email,
		Address: req.Address,
	}

	// 使用通用 repo
	clientRepo := repo.NewRepo[Client](s.db)
	if err := clientRepo.Create(ctx, client); err != nil {
		return nil, err
	}

	// 异步同步到 ES
	s.esSync.Index(client)

	return client, nil
}

// Get 获取客户
func (s *ClientService) Get(ctx context.Context, id uint) (*Client, error) {
	clientRepo := repo.NewRepo[Client](s.db)
	return clientRepo.GetByID(ctx, id)
}

// List 客户列表
func (s *ClientService) List(ctx context.Context, limit, offset int) ([]Client, error) {
	clientRepo := repo.NewRepo[Client](s.db)
	return clientRepo.List(ctx, limit, offset)
}

// Update 更新客户
func (s *ClientService) Update(ctx context.Context, id uint, req *UpdateClientRequest) error {
	clientRepo := repo.NewRepo[Client](s.db)
	
	// 1. 查询客户
	client, err := clientRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 2. 构造更新字段
	fields := make(map[string]interface{})
	if req.Name != "" {
		fields["name"] = req.Name
		client.Name = req.Name
	}
	if req.Contact != "" {
		fields["contact"] = req.Contact
		client.Contact = req.Contact
	}
	if req.Phone != "" {
		fields["phone"] = req.Phone
		client.Phone = req.Phone
	}
	if req.Email != "" {
		fields["email"] = req.Email
		client.Email = req.Email
	}
	if req.Address != "" {
		fields["address"] = req.Address
		client.Address = req.Address
	}

	// 3. 如果没有要更新的字段，直接返回
	if len(fields) == 0 {
		return nil
	}

	// 4. 更新数据库
	if err := clientRepo.UpdateFields(ctx, id, fields); err != nil {
		return err
	}

	// 5. 异步同步到 ES
	s.esSync.Update(client)

	return nil
}

// Delete 删除客户
func (s *ClientService) Delete(ctx context.Context, id uint) error {
	clientRepo := repo.NewRepo[Client](s.db)
	
	// 1. 删除数据库记录
	if err := clientRepo.Delete(ctx, id); err != nil {
		return err
	}

	// 2. 异步删除 ES 文档
	s.esSync.Delete("clients", strconv.Itoa(int(id)))

	return nil
}