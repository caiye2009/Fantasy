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
	// 1. DTO → Domain Model
	client := &domain.Client{
		CustomNo:     req.CustomNo,
		CustomerCode: req.CustomerCode,
		InputDate:    req.InputDate,
		Sales:        req.Sales,
		CustomName:   req.CustomName,
		StateChNm:    req.StateChNm,
		Country:      req.Country,
		Address:      req.Address,
		AddressEn:    req.AddressEn,
		CustomNameEn: req.CustomNameEn,
		Contactor:    req.Contactor,
		UnitPhone:    req.UnitPhone,
		Mobile:       req.Mobile,
		FaxNum:       req.FaxNum,
		Email:        req.Email,
		PyCustomName: req.PyCustomName,
		CheckRequest: req.CheckRequest,
		CustomStatus: req.CustomStatus,
		DocMan:       req.DocMan,
	}

	// 2. 领域验证
	if err := client.Validate(); err != nil {
		return nil, err
	}

	// 3. 保存到数据库
	if err := s.repo.Save(ctx, client); err != nil {
		return nil, err
	}

	// 4. 异步同步到 ES
	if s.esSync != nil {
		s.esSync.Index(client)
	}

	// 5. Domain Model → DTO
	return toClientResponse(client), nil
}

// Get 获取客户
func (s *ClientService) Get(ctx context.Context, id uint) (*ClientResponse, error) {
	client, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return toClientResponse(client), nil
}

// Update 更新客户
func (s *ClientService) Update(ctx context.Context, id uint, req *UpdateClientRequest) error {
	// 1. 查询客户
	client, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 2. 更新字段
	if req.CustomNo != "" {
		client.CustomNo = req.CustomNo
	}
	if req.CustomerCode != "" {
		client.CustomerCode = req.CustomerCode
	}
	if req.InputDate != nil {
		client.InputDate = req.InputDate
	}
	if req.Sales != "" {
		client.Sales = req.Sales
	}
	if req.CustomName != "" {
		client.CustomName = req.CustomName
	}
	if req.StateChNm != "" {
		client.StateChNm = req.StateChNm
	}
	if req.Country != "" {
		client.Country = req.Country
	}
	if req.Address != "" {
		client.Address = req.Address
	}
	if req.AddressEn != "" {
		client.AddressEn = req.AddressEn
	}
	if req.CustomNameEn != "" {
		client.CustomNameEn = req.CustomNameEn
	}
	if req.Contactor != "" {
		client.Contactor = req.Contactor
	}
	if req.UnitPhone != "" {
		client.UnitPhone = req.UnitPhone
	}
	if req.Mobile != "" {
		client.Mobile = req.Mobile
	}
	if req.FaxNum != "" {
		client.FaxNum = req.FaxNum
	}
	if req.Email != "" {
		client.Email = req.Email
	}
	if req.PyCustomName != "" {
		client.PyCustomName = req.PyCustomName
	}
	if req.CheckRequest != "" {
		client.CheckRequest = req.CheckRequest
	}
	if req.CustomStatus != "" {
		client.CustomStatus = req.CustomStatus
	}
	if req.DocMan != "" {
		client.DocMan = req.DocMan
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

// toClientResponse 将 Domain Model 转换为 Response DTO
func toClientResponse(client *domain.Client) *ClientResponse {
	return &ClientResponse{
		ID:           client.ID,
		CustomNo:     client.CustomNo,
		CustomerCode: client.CustomerCode,
		InputDate:    client.InputDate,
		Sales:        client.Sales,
		CustomName:   client.CustomName,
		StateChNm:    client.StateChNm,
		Country:      client.Country,
		Address:      client.Address,
		AddressEn:    client.AddressEn,
		CustomNameEn: client.CustomNameEn,
		Contactor:    client.Contactor,
		UnitPhone:    client.UnitPhone,
		Mobile:       client.Mobile,
		FaxNum:       client.FaxNum,
		Email:        client.Email,
		PyCustomName: client.PyCustomName,
		CheckRequest: client.CheckRequest,
		CustomStatus: client.CustomStatus,
		DocMan:       client.DocMan,
		CreatedAt:    client.CreatedAt,
		UpdatedAt:    client.UpdatedAt,
	}
}