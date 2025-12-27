package config

import (
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

func InitElasticsearch(cfg *Config) *elasticsearch.Client {
	esCfg := elasticsearch.Config{
		Addresses: cfg.ESAddresses,
		Username:  cfg.ESUsername,
		Password:  cfg.ESPassword,
	}

	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v", err)
	}

	// 测试连接
	res, err := client.Info()
	if err != nil {
		log.Fatalf("Failed to connect to Elasticsearch: %v", err)
	}
	defer res.Body.Close()

	log.Println("✓ Elasticsearch connected")
	return client
}
