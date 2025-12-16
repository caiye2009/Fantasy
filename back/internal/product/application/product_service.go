package application

import (
	"context"
	"strconv"
	
	"back/internal/product/domain"
	"back/internal/product/infra"
)

// ESSync ES 同步接口
type ESSync interface {
	Index(doc interface{}) error
	Update(doc interface{}) error
	Delete(indexName, docID string) error
}

// ProductService 产品应用服务
type ProductService struct {
	repo   *infra.ProductRepo
	esSync ESSync
}

// NewProductService 创建产品服务
func NewProductService(repo *infra.ProductRepo, esSync ESSync) *ProductService {
	return &ProductService{
		repo:   repo,
		esSync: esSync,
	}
}

// Create 创建产品
func (s *ProductService) Create(ctx context.Context, req *CreateProductRequest) (*ProductResponse, error) {
	// 1. DTO → Domain Model
	product := &domain.Product{
		Name:      req.Name,
		Status:    domain.ProductStatusDraft,
		Materials: req.Materials,
		Processes: req.Processes,
	}

	// 2. 领域验证
	if err := product.Validate(); err != nil {
		return nil, err
	}

	// 3. 保存到数据库
	if err := s.repo.Save(ctx, product); err != nil {
		return nil, err
	}

	// 4. 异步同步到 ES
	if s.esSync != nil {
		s.esSync.Index(product)
	}

	// 5. Domain Model → DTO
	return &ProductResponse{
		ID:        product.ID,
		Name:      product.Name,
		Status:    product.Status,
		Materials: product.Materials,
		Processes: product.Processes,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}, nil
}

// Get 获取产品
func (s *ProductService) Get(ctx context.Context, id uint) (*ProductResponse, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &ProductResponse{
		ID:        product.ID,
		Name:      product.Name,
		Status:    product.Status,
		Materials: product.Materials,
		Processes: product.Processes,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}, nil
}

// Update 更新产品
func (s *ProductService) Update(ctx context.Context, id uint, req *UpdateProductRequest) error {
	// 1. 查询产品
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	
	// 2. 更新字段（通过领域方法）
	if req.Name != "" {
		if err := product.UpdateName(req.Name); err != nil {
			return err
		}
	}
	
	if len(req.Materials) > 0 {
		if err := product.UpdateMaterials(req.Materials); err != nil {
			return err
		}
	}
	
	if len(req.Processes) > 0 {
		if err := product.UpdateProcesses(req.Processes); err != nil {
			return err
		}
	}
	
	// 状态更新
	if req.Status != "" {
		switch req.Status {
		case domain.ProductStatusSubmitted:
			if err := product.Submit(); err != nil {
				return err
			}
		case domain.ProductStatusApproved:
			if err := product.Approve(); err != nil {
				return err
			}
		case domain.ProductStatusRejected:
			if err := product.Reject(); err != nil {
				return err
			}
		}
	}
	
	// 3. 验证
	if err := product.Validate(); err != nil {
		return err
	}
	
	// 4. 保存
	if err := s.repo.Update(ctx, product); err != nil {
		return err
	}
	
	// 5. 异步同步到 ES
	if s.esSync != nil {
		s.esSync.Update(product)
	}
	
	return nil
}

// Delete 删除产品
func (s *ProductService) Delete(ctx context.Context, id uint) error {
	// 1. 查询产品
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	
	// 2. 检查是否可以删除（领域规则）
	if !product.CanDelete() {
		return domain.ErrCannotDeleteApprovedProduct
	}
	
	// 3. 删除
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	
	// 4. 异步删除 ES 文档
	if s.esSync != nil {
		s.esSync.Delete("products", strconv.Itoa(int(id)))
	}
	
	return nil
}

// Exists 检查产品是否存在（供其他模块调用）
func (s *ProductService) Exists(ctx context.Context, id uint) (bool, error) {
	return s.repo.ExistsByID(ctx, id)
}