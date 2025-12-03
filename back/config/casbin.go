package config

import (
	"log"
	"github.com/casbin/casbin/v2"
)

func InitCasbinPolicies(enforcer *casbin.Enforcer) error {
	policies, err := enforcer.GetPolicy()
	if err != nil {
		return err
	}
	
	if len(policies) > 0 {
		log.Println("✓ Casbin policies already initialized")
		return nil
	}

	log.Println("Initializing Casbin policies...")

	defaultPolicies := [][]string{
		// Admin (全部权限)
		{"admin", "/api/v1/*", "*"},
		
		// HR (用户管理)
		{"hr", "/api/v1/user/create", "POST"},
		{"hr", "/api/v1/user/list", "GET"},
		{"hr", "/api/v1/user/:id", "GET"},
		{"hr", "/api/v1/user/:id", "PUT"},
		
		// Sales (业务员)
		{"sales", "/api/v1/vendor/list", "GET"},
		{"sales", "/api/v1/client/list", "GET"},
		{"sales", "/api/v1/client/create", "POST"},
		{"sales", "/api/v1/material/list", "GET"},
		{"sales", "/api/v1/process/list", "GET"},
		{"sales", "/api/v1/product/list", "GET"},
		{"sales", "/api/v1/product/create", "POST"},
		{"sales", "/api/v1/product/cost", "POST"},
		
		// Follower (跟单员)
		{"follower", "/api/v1/vendor/list", "GET"},
		{"follower", "/api/v1/client/list", "GET"},
		{"follower", "/api/v1/material/list", "GET"},
		{"follower", "/api/v1/process/list", "GET"},
		{"follower", "/api/v1/product/list", "GET"},
		{"follower", "/api/v1/product/cost", "POST"},
		
		// Assistant (助理,只读)
		{"assistant", "/api/v1/vendor/list", "GET"},
		{"assistant", "/api/v1/client/list", "GET"},
		{"assistant", "/api/v1/material/list", "GET"},
		{"assistant", "/api/v1/process/list", "GET"},
		{"assistant", "/api/v1/product/list", "GET"},
	}

	for _, policy := range defaultPolicies {
		_, err := enforcer.AddPolicy(policy)
		if err != nil {
			return err
		}
	}

	log.Printf("✓ Initialized %d Casbin policies", len(defaultPolicies))
	return nil
}