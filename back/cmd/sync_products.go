package main

import (
	"log"

	"back/config"
	"back/internal/product/domain"
	"back/pkg/es"
	applog "back/pkg/log"
)

// 临时脚本：将数据库中的产品同步到 ES
func main() {
	log.Println("=== Starting Product Sync to ES ===")

	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库
	db := config.InitDatabase(cfg)

	// 初始化 ES
	esClient := config.InitElasticsearch(cfg)

	// 初始化日志
	if err := applog.Init(cfg.LogLevel, cfg.LogFormat); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// 初始化 ES Sync
	logger := applog.NewLogger()
	loggerAdapter := es.NewLoggerAdapter(logger)
	esSync := es.NewESSync(esClient, loggerAdapter)

	// 从数据库读取所有产品
	var products []domain.Product
	if err := db.Find(&products).Error; err != nil {
		log.Fatalf("Failed to fetch products: %v", err)
	}

	log.Printf("Found %d products in database", len(products))

	if len(products) == 0 {
		log.Println("No products to sync")
		return
	}

	// 同步到 ES
	successCount := 0
	failedCount := 0

	for i := range products {
		product := &products[i]
		log.Printf("Syncing product %d: %s (ID: %d)", i+1, product.Name, product.ID)

		if err := esSync.Index(product); err != nil {
			log.Printf("  ✗ Failed: %v", err)
			failedCount++
		} else {
			log.Printf("  ✓ Success")
			successCount++
		}
	}

	log.Printf("\n=== Sync Complete ===")
	log.Printf("Total: %d", len(products))
	log.Printf("Success: %d", successCount)
	log.Printf("Failed: %d", failedCount)
}
