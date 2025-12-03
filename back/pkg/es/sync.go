package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"

	"back/pkg/log"
)

// ESSync ES 同步服务
type ESSync struct {
	client *elasticsearch.Client
	logger *log.Logger
}

// NewESSync 创建 ES 同步服务
func NewESSync(client *elasticsearch.Client, logger *log.Logger) *ESSync {
	return &ESSync{
		client: client,
		logger: logger.With(log.String("service", "es-sync")),
	}
}

// Index 索引文档 (异步)
func (es *ESSync) Index(entity Indexable) {
	go es.indexSync(entity)
}

// Update 更新文档 (异步)
func (es *ESSync) Update(entity Indexable) {
	go es.updateSync(entity)
}

// Delete 删除文档 (异步)
func (es *ESSync) Delete(indexName, docID string) {
	go es.deleteSync(indexName, docID)
}

// indexSync 同步索引文档
func (es *ESSync) indexSync(entity Indexable) {
	start := time.Now()
	indexName := entity.GetIndexName()
	docID := entity.GetDocumentID()

	// 转换为文档
	doc := entity.ToDocument()
	data, err := json.Marshal(doc)
	if err != nil {
		es.logger.Error("索引失败 - 序列化错误",
			log.String("index", indexName),
			log.String("doc_id", docID),
			log.Error(err),
		)
		return
	}

	// 执行索引
	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: docID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		es.logger.Error("索引失败 - 请求错误",
			log.String("index", indexName),
			log.String("doc_id", docID),
			log.Error(err),
		)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		es.logger.Error("索引失败 - ES 错误",
			log.String("index", indexName),
			log.String("doc_id", docID),
			log.String("status", res.Status()),
			log.String("error", res.String()),
		)
		return
	}

	elapsed := time.Since(start)
	es.logger.Info("索引成功",
		log.String("index", indexName),
		log.String("doc_id", docID),
		log.Duration("elapsed", elapsed),
	)
}

// updateSync 同步更新文档
func (es *ESSync) updateSync(entity Indexable) {
	start := time.Now()
	indexName := entity.GetIndexName()
	docID := entity.GetDocumentID()

	// 转换为文档
	doc := entity.ToDocument()
	data, err := json.Marshal(map[string]interface{}{
		"doc": doc,
	})
	if err != nil {
		es.logger.Error("更新失败 - 序列化错误",
			log.String("index", indexName),
			log.String("doc_id", docID),
			log.Error(err),
		)
		return
	}

	// 执行更新
	req := esapi.UpdateRequest{
		Index:      indexName,
		DocumentID: docID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		es.logger.Error("更新失败 - 请求错误",
			log.String("index", indexName),
			log.String("doc_id", docID),
			log.Error(err),
		)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		es.logger.Error("更新失败 - ES 错误",
			log.String("index", indexName),
			log.String("doc_id", docID),
			log.String("status", res.Status()),
			log.String("error", res.String()),
		)
		return
	}

	elapsed := time.Since(start)
	es.logger.Info("更新成功",
		log.String("index", indexName),
		log.String("doc_id", docID),
		log.Duration("elapsed", elapsed),
	)
}

// deleteSync 同步删除文档
func (es *ESSync) deleteSync(indexName, docID string) {
	start := time.Now()

	// 执行删除
	req := esapi.DeleteRequest{
		Index:      indexName,
		DocumentID: docID,
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		es.logger.Error("删除失败 - 请求错误",
			log.String("index", indexName),
			log.String("doc_id", docID),
			log.Error(err),
		)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		// 404 不算错误 (文档不存在)
		if res.StatusCode != 404 {
			es.logger.Error("删除失败 - ES 错误",
				log.String("index", indexName),
				log.String("doc_id", docID),
				log.String("status", res.Status()),
				log.String("error", res.String()),
			)
			return
		}
	}

	elapsed := time.Since(start)
	es.logger.Info("删除成功",
		log.String("index", indexName),
		log.String("doc_id", docID),
		log.Duration("elapsed", elapsed),
	)
}

// IndexSync 同步索引 (阻塞版本,用于测试)
func (es *ESSync) IndexSync(entity Indexable) error {
	indexName := entity.GetIndexName()
	docID := entity.GetDocumentID()

	doc := entity.ToDocument()
	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("marshal document: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: docID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		return fmt.Errorf("index request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("index error: %s", res.String())
	}

	return nil
}