package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"back/config"
	applog "back/pkg/log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// 支持的所有域索引
var allIndexes = []string{
	"client",
	"material",
	"process",
	"product",
	"supplier",
	"order",
	"plan",
}

func main() {
	log.Println("=== ES Index Clear Tool ===")

	// 解析命令行参数
	args := os.Args[1:]

	var indexesToClear []string
	if len(args) == 0 {
		// 没有参数，删除所有索引
		indexesToClear = allIndexes
		log.Println("No domain specified, will clear ALL indexes")
	} else {
		// 有参数，删除指定的索引
		indexesToClear = args
		log.Printf("Will clear indexes: %v", indexesToClear)
	}

	// 1. 加载配置
	cfg := config.LoadConfig()
	log.Println("✓ Config loaded")

	// 2. 初始化日志
	if err := applog.Init(cfg.LogLevel, cfg.LogFormat); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	log.Println("✓ Logger initialized")

	// 3. 初始化 ES 客户端
	esClient := config.InitElasticsearch(cfg)
	log.Println("✓ Elasticsearch connected")

	// 4. 删除索引
	successCount := 0
	failCount := 0

	for _, index := range indexesToClear {
		log.Printf("Deleting index: %s", index)
		if err := deleteIndex(esClient, index); err != nil {
			log.Printf("✗ Failed to delete index '%s': %v", index, err)
			failCount++
		} else {
			log.Printf("✓ Successfully deleted index: %s", index)
			successCount++
		}
	}

	// 5. 输出统计
	log.Println("=== Clear Summary ===")
	log.Printf("Total: %d", len(indexesToClear))
	log.Printf("Success: %d", successCount)
	log.Printf("Failed: %d", failCount)

	if failCount > 0 {
		log.Fatalf("Clear completed with %d failures", failCount)
	}

	log.Println("=== Clear completed successfully ===")
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

	// 404 表示索引不存在，也算成功
	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("delete index error: %s", res.String())
	}

	return nil
}
