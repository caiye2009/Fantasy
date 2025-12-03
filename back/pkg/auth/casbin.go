package auth

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type CasbinWang struct {
	enforcer *casbin.Enforcer
}

func NewCasbinWang(db *gorm.DB, rdb *redis.Client, modelPath string) (*CasbinWang, error) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	m, err := model.NewModelFromFile(modelPath)
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	return &CasbinWang{enforcer: enforcer}, nil
}

func (c *CasbinWang) CheckPermission(subject, object, action string) (bool, error) {
	return c.enforcer.Enforce(subject, object, action)
}

// GetEnforcer 返回 enforcer 实例
func (c *CasbinWang) GetEnforcer() *casbin.Enforcer {
	return c.enforcer
}