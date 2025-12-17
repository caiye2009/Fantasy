package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"back/config"
	clientDomain "back/internal/client/domain"
	materialDomain "back/internal/material/domain"
	orderDomain "back/internal/order/domain"
	planDomain "back/internal/plan/domain"
	processDomain "back/internal/process/domain"
	productDomain "back/internal/product/domain"
	supplierDomain "back/internal/supplier/domain"
	applog "back/pkg/log"
	"back/pkg/es"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"gorm.io/gorm"
)

// 支持的域列表
var supportedDomains = []string{"client", "material", "process", "product", "supplier", "order", "plan"}

func main() {
	log.Println("=== ES Reindex Tool ===")

	// 解析命令行参数
	args := os.Args[1:]

	var domainsToReindex []string
	if len(args) == 0 {
		// 没有参数，重建所有支持的域
		domainsToReindex = supportedDomains
		log.Println("No domain specified, will reindex ALL supported domains")
	} else {
		// 有参数，重建指定的域
		domainsToReindex = args
		log.Printf("Will reindex domains: %v", domainsToReindex)
	}

	// 验证域名是否支持
	for _, domain := range domainsToReindex {
		if !isSupported(domain) {
			log.Fatalf("Unsupported domain: %s. Supported domains: %v", domain, supportedDomains)
		}
	}

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

	// 5. 逐个处理域
	totalSuccess := 0
	totalFailed := 0

	for _, domain := range domainsToReindex {
		log.Printf("\n=== Processing domain: %s ===", domain)

		// 删除旧索引
		if err := deleteIndex(esClient, domain); err != nil {
			log.Printf("Warning: Failed to delete old index '%s' (may not exist): %v", domain, err)
		} else {
			log.Printf("✓ Old index '%s' deleted", domain)
			// 等待一下确保索引删除完成
			time.Sleep(2 * time.Second)
		}

		// 重新索引
		success, failed, err := reindexDomain(db, esSync, domain)
		if err != nil {
			log.Printf("✗ Failed to reindex domain '%s': %v", domain, err)
			totalFailed++
			continue
		}

		totalSuccess += success
		log.Printf("✓ Domain '%s' reindexed successfully: %d records", domain, success)

		if failed > 0 {
			totalFailed += failed
			log.Printf("⚠ Domain '%s' had %d failures", domain, failed)
		}
	}

	// 6. 输出总体统计
	log.Println("\n=== Overall Summary ===")
	log.Printf("Domains processed: %d", len(domainsToReindex))
	log.Printf("Total records indexed: %d", totalSuccess)
	log.Printf("Total failures: %d", totalFailed)

	if totalFailed > 0 {
		log.Fatalf("Reindex completed with %d failures", totalFailed)
	}

	log.Println("=== Reindex completed successfully ===")
}

// isSupported 检查域是否支持
func isSupported(domain string) bool {
	for _, d := range supportedDomains {
		if d == domain {
			return true
		}
	}
	return false
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

// reindexDomain 重新索引指定域
func reindexDomain(db *gorm.DB, esSync *es.ESSync, domain string) (success, failed int, err error) {
	switch domain {
	case "client":
		return reindexClients(db, esSync)
	case "material":
		return reindexMaterials(db, esSync)
	case "process":
		return reindexProcesses(db, esSync)
	case "product":
		return reindexProducts(db, esSync)
	case "supplier":
		return reindexSuppliers(db, esSync)
	case "order":
		return reindexOrders(db, esSync)
	case "plan":
		return reindexPlans(db, esSync)
	default:
		return 0, 0, fmt.Errorf("unsupported domain: %s", domain)
	}
}

// reindexClients 重新索引客户数据
func reindexClients(db *gorm.DB, esSync *es.ESSync) (success, failed int, err error) {
	log.Println("→ Fetching clients from database...")

	var clients []clientDomain.Client
	if err := db.Unscoped().Find(&clients).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to fetch clients: %w", err)
	}

	total := len(clients)
	log.Printf("  Found %d clients", total)

	if total == 0 {
		return 0, 0, nil
	}

	// 过滤掉已删除的记录
	var activeClients []interface{}
	for _, client := range clients {
		if !client.DeletedAt.Valid {
			c := client // 复制以避免闭包问题
			activeClients = append(activeClients, &c)
		}
	}

	log.Printf("  Indexing %d active clients...", len(activeClients))

	// 批量索引（每批 200 条）
	return batchIndexDocuments(esSync, activeClients, 200)
}

// reindexMaterials 重新索引材料数据
func reindexMaterials(db *gorm.DB, esSync *es.ESSync) (success, failed int, err error) {
	log.Println("→ Fetching materials from database...")

	var materials []materialDomain.Material
	if err := db.Unscoped().Find(&materials).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to fetch materials: %w", err)
	}

	total := len(materials)
	log.Printf("  Found %d materials", total)

	if total == 0 {
		return 0, 0, nil
	}

	// 过滤掉已删除的记录
	var activeMaterials []interface{}
	for _, material := range materials {
		if !material.DeletedAt.Valid {
			m := material // 复制以避免闭包问题
			activeMaterials = append(activeMaterials, &m)
		}
	}

	log.Printf("  Indexing %d active materials...", len(activeMaterials))

	// 批量索引（每批 200 条）
	return batchIndexDocuments(esSync, activeMaterials, 200)
}

// reindexProcesses 重新索引工序数据
func reindexProcesses(db *gorm.DB, esSync *es.ESSync) (success, failed int, err error) {
	log.Println("→ Fetching processes from database...")

	var processes []processDomain.Process
	if err := db.Unscoped().Find(&processes).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to fetch processes: %w", err)
	}

	total := len(processes)
	log.Printf("  Found %d processes", total)

	if total == 0 {
		return 0, 0, nil
	}

	// 过滤掉已删除的记录
	var activeProcesses []interface{}
	for _, process := range processes {
		if !process.DeletedAt.Valid {
			p := process // 复制以避免闭包问题
			activeProcesses = append(activeProcesses, &p)
		}
	}

	log.Printf("  Indexing %d active processes...", len(activeProcesses))

	// 批量索引（每批 200 条）
	return batchIndexDocuments(esSync, activeProcesses, 200)
}

// reindexProducts 重新索引产品数据
func reindexProducts(db *gorm.DB, esSync *es.ESSync) (success, failed int, err error) {
	log.Println("→ Fetching products from database...")

	var products []productDomain.Product
	if err := db.Unscoped().Find(&products).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to fetch products: %w", err)
	}

	total := len(products)
	log.Printf("  Found %d products", total)

	if total == 0 {
		return 0, 0, nil
	}

	// 过滤掉已删除的记录
	var activeProducts []interface{}
	for _, product := range products {
		if !product.DeletedAt.Valid {
			p := product
			activeProducts = append(activeProducts, &p)
		}
	}

	log.Printf("  Indexing %d active products...", len(activeProducts))

	// 批量索引（每批 200 条）
	return batchIndexDocuments(esSync, activeProducts, 200)
}

// reindexSuppliers 重新索引供应商数据
func reindexSuppliers(db *gorm.DB, esSync *es.ESSync) (success, failed int, err error) {
	log.Println("→ Fetching suppliers from database...")

	var suppliers []supplierDomain.Supplier
	if err := db.Unscoped().Find(&suppliers).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to fetch suppliers: %w", err)
	}

	total := len(suppliers)
	log.Printf("  Found %d suppliers", total)

	if total == 0 {
		return 0, 0, nil
	}

	// 过滤掉已删除的记录
	var activeSuppliers []interface{}
	for _, supplier := range suppliers {
		if !supplier.DeletedAt.Valid {
			s := supplier
			activeSuppliers = append(activeSuppliers, &s)
		}
	}

	log.Printf("  Indexing %d active suppliers...", len(activeSuppliers))

	// 批量索引（每批 200 条）
	return batchIndexDocuments(esSync, activeSuppliers, 200)
}

// reindexOrders 重新索引订单数据
func reindexOrders(db *gorm.DB, esSync *es.ESSync) (success, failed int, err error) {
	log.Println("→ Fetching orders from database...")

	var orders []orderDomain.Order
	if err := db.Unscoped().Find(&orders).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to fetch orders: %w", err)
	}

	total := len(orders)
	log.Printf("  Found %d orders", total)

	if total == 0 {
		return 0, 0, nil
	}

	// 过滤掉已删除的记录
	var activeOrders []interface{}
	for _, order := range orders {
		if !order.DeletedAt.Valid {
			o := order
			activeOrders = append(activeOrders, &o)
		}
	}

	log.Printf("  Indexing %d active orders...", len(activeOrders))

	// 批量索引（每批 200 条）
	return batchIndexDocuments(esSync, activeOrders, 200)
}

// reindexPlans 重新索引计划数据
func reindexPlans(db *gorm.DB, esSync *es.ESSync) (success, failed int, err error) {
	log.Println("→ Fetching plans from database...")

	var plans []planDomain.Plan
	if err := db.Unscoped().Find(&plans).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to fetch plans: %w", err)
	}

	total := len(plans)
	log.Printf("  Found %d plans", total)

	if total == 0 {
		return 0, 0, nil
	}

	// 过滤掉已删除的记录
	var activePlans []interface{}
	for _, plan := range plans {
		if !plan.DeletedAt.Valid {
			p := plan
			activePlans = append(activePlans, &p)
		}
	}

	log.Printf("  Indexing %d active plans...", len(activePlans))

	// 批量索引（每批 200 条）
	return batchIndexDocuments(esSync, activePlans, 200)
}

// batchIndexDocuments 批量索引文档的辅助函数
func batchIndexDocuments(esSync *es.ESSync, docs []interface{}, batchSize int) (totalSuccess, totalFailed int, err error) {
	total := len(docs)

	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		batch := docs[i:end]

		// 最后一批才刷新索引
		refresh := (end == total)

		success, failed, err := esSync.BulkIndex(batch, refresh)
		if err != nil {
			log.Printf("  [%d-%d/%d] Batch failed: %v", i+1, end, total, err)
			return totalSuccess, totalFailed + len(batch), err
		}

		totalSuccess += success
		totalFailed += failed

		log.Printf("  [%d/%d] Batch indexed: %d success, %d failed", end, total, success, failed)
	}

	return totalSuccess, totalFailed, nil
}
