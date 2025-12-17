package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"back/config"
	"back/internal/client/domain"
	applog "back/pkg/log"
	"back/pkg/es"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"gorm.io/gorm"
)

func main() {
	log.Println("=== Starting Client Reindex Tool ===")

	// 1. 加载配置
	cfg := config.LoadConfig()
	log.Println("✓ Config loaded")

	// 2. 初始化日志
	if err := applog.Init(cfg.LogLevel, cfg.LogFormat); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	logger := applog.NewLogger()
	loggerAdapter := es.NewLoggerAdapter(logger)
	log.Println("✓ Logger initialized")

	// 3. 初始化数据库
	db := config.InitDatabase(cfg)
	log.Println("✓ Database connected")

	// 4. 初始化 ES 客户端
	esClient := config.InitElasticsearch(cfg)
	esSync := es.NewESSync(esClient, loggerAdapter)
	log.Println("✓ Elasticsearch connected")

	// 5. 删除旧索引
	if err := deleteIndex(esClient, "clients"); err != nil {
		log.Printf("Warning: Failed to delete old index (may not exist): %v", err)
	} else {
		log.Println("✓ Old index deleted")
		// 等待一下确保索引删除完成
		time.Sleep(2 * time.Second)
	}

	// 6. 重新索引所有 clients
	if err := reindexClients(db, esSync); err != nil {
		log.Fatalf("Failed to reindex clients: %v", err)
	}

	log.Println("=== Reindex completed successfully ===")
}

// deleteIndex 删除索引
func deleteIndex(esClient *elasticsearch.Client, indexName string) error {
	req := esapi.IndicesDeleteRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("delete index error: %s", res.String())
	}

	return nil
}

// reindexClients 重新索引所有客户数据
func reindexClients(db *gorm.DB, esSync *es.ESSync) error {
	log.Println("=== Starting to reindex clients ===")

	// 1. 从数据库读取所有 clients（包括软删除的）
	var clients []domain.Client
	if err := db.Unscoped().Find(&clients).Error; err != nil {
		return fmt.Errorf("failed to fetch clients: %w", err)
	}

	total := len(clients)
	log.Printf("Found %d clients in database", total)

	if total == 0 {
		log.Println("No clients to index")
		return nil
	}

	// 2. 批量索引到 ES
	successCount := 0
	failCount := 0

	for i, client := range clients {
		// 跳过已删除的记录
		if client.DeletedAt.Valid {
			log.Printf("[%d/%d] Skipping deleted client: ID=%d", i+1, total, client.ID)
			continue
		}

		// 索引到 ES
		if err := esSync.Index(&client); err != nil {
			log.Printf("[%d/%d] Failed to index client ID=%d: %v", i+1, total, client.ID, err)
			failCount++
		} else {
			successCount++
			if (i+1)%100 == 0 || i+1 == total {
				log.Printf("[%d/%d] Indexed client: ID=%d, CustomNo=%s, CustomName=%s",
					i+1, total, client.ID, client.CustomNo, client.CustomName)
			}
		}
	}

	log.Printf("=== Reindex Summary ===")
	log.Printf("Total: %d", total)
	log.Printf("Success: %d", successCount)
	log.Printf("Failed: %d", failCount)
	log.Printf("Skipped (deleted): %d", total-successCount-failCount)

	if failCount > 0 {
		return fmt.Errorf("reindex completed with %d failures", failCount)
	}

	return nil
}
