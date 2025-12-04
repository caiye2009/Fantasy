package domain

import "context"

// ClientRepository 客户仓储接口
type ClientRepository interface {
	// 基础 CRUD
	Save(ctx context.Context, client *Client) error
	FindByID(ctx context.Context, id uint) (*Client, error)
	FindAll(ctx context.Context, limit, offset int) ([]*Client, error)
	Update(ctx context.Context, client *Client) error
	Delete(ctx context.Context, id uint) error
	
	// 业务查询
	ExistsByID(ctx context.Context, id uint) (bool, error)
	FindByName(ctx context.Context, name string) (*Client, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	FindByPhone(ctx context.Context, phone string) (*Client, error)
	FindByEmail(ctx context.Context, email string) (*Client, error)
	Count(ctx context.Context) (int64, error)
}