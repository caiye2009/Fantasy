package config

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"

	applog "back/pkg/log"
	"back/pkg/es"
	userDomain "back/internal/user/domain"
)

func Init() error {
	log.Println("=== Initializing Application ===")

	cfg := LoadConfig()
	log.Println("✓ Config loaded")

	// 初始化日志
	if err := applog.Init(cfg.LogLevel, cfg.LogFormat); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	log.Println("✓ Logger initialized")

	db := InitDatabase(cfg)
	rdb := InitRedis(cfg)
	esClient := InitElasticsearch(cfg)

	// 创建 ES 同步服务（使用适配器）
	logger := applog.NewLogger()
	loggerAdapter := es.NewLoggerAdapter(logger)
	esSync := es.NewESSync(esClient, loggerAdapter)
	log.Println("✓ ES Sync initialized")

	// 初始化 Casbin
	enforcer := InitCasbin(cfg, db)
	if err := InitCasbinPolicies(enforcer); err != nil {
		log.Fatalf("Failed to initialize Casbin policies: %v", err)
	}

	jwtWang, authWang := InitAuth(db, rdb, cfg, enforcer)

	log.Println("=== Database Migration ===")
	if err := AutoMigrate(db); err != nil {
		return err
	}

	InitAdminUser(db)

	log.Println("=== Initializing Search Config Registry ===")
	searchRegistry, err := InitSearchRegistry()
	if err != nil {
		log.Fatalf("Failed to initialize search registry: %v", err)
	}
	log.Printf("✓ Search registry initialized with indices: %v", searchRegistry.ListIndices())

	log.Println("=== Initializing Services ===")
	services := InitServices(db, rdb, esClient, jwtWang, esSync, searchRegistry)
	log.Println("✓ Services initialized")

	log.Println("=== Initializing Router ===")
	router := InitRoutes(authWang, services, db)
	log.Println("✓ Router initialized")

	server := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: router,
	}

	go gracefulShutdown(server, logger)

	log.Printf("=== Server started on %s ===\n", cfg.ServerAddr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func InitAdminUser(db *gorm.DB) {
	var count int64
	db.Model(&userDomain.User{}).Where("role = ?", "admin").Count(&count)

	if count > 0 {
		log.Println("✓ Admin user already exists")
		return
	}

	// Admin 密码: admin
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	admin := &userDomain.User{
		LoginID:      "admin",
		Username:     "admin",
		PasswordHash: string(hash),
		Email:        "admin@example.com",
		Role:         userDomain.RoleAdmin,
		Status:       userDomain.UserStatusActive,
		HasInitPass:  false,
	}
	db.Create(admin)
	log.Println("✓ Default admin user created (login_id: 8000, password: admin)")
}

func gracefulShutdown(server *http.Server, logger *applog.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// 同步日志
	if logger != nil {
		logger.Sync()
	}

	log.Println("Server exited")
}