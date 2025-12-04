package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ESSync struct {
	client *elasticsearch.Client
	logger Logger
}

type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}

func NewESSync(client *elasticsearch.Client, logger Logger) *ESSync {
	return &ESSync{
		client: client,
		logger: logger,
	}
}

// Index 索引文档（参数改为 interface{}）
func (s *ESSync) Index(doc interface{}) error {
	// 类型断言检查是否实现 Indexable 接口
	indexable, ok := doc.(Indexable)
	if !ok {
		return fmt.Errorf("document does not implement Indexable interface")
	}
	
	data, err := json.Marshal(indexable.ToDocument())
	if err != nil {
		if s.logger != nil {
			s.logger.Error("marshal document failed", "error", err)
		}
		return err
	}

	req := esapi.IndexRequest{
		Index:      indexable.GetIndexName(),
		DocumentID: indexable.GetDocumentID(),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), s.client)
	if err != nil {
		if s.logger != nil {
			s.logger.Error("index document failed", "error", err)
		}
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		if s.logger != nil {
			s.logger.Error("index error", "status", res.Status())
		}
		return fmt.Errorf("index error: %s", res.String())
	}

	return nil
}

// Update 更新文档（参数改为 interface{}）
func (s *ESSync) Update(doc interface{}) error {
	// 类型断言检查是否实现 Indexable 接口
	indexable, ok := doc.(Indexable)
	if !ok {
		return fmt.Errorf("document does not implement Indexable interface")
	}
	
	data, err := json.Marshal(map[string]interface{}{
		"doc": indexable.ToDocument(),
	})
	if err != nil {
		if s.logger != nil {
			s.logger.Error("marshal document failed", "error", err)
		}
		return err
	}

	req := esapi.UpdateRequest{
		Index:      indexable.GetIndexName(),
		DocumentID: indexable.GetDocumentID(),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), s.client)
	if err != nil {
		if s.logger != nil {
			s.logger.Error("update document failed", "error", err)
		}
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		if s.logger != nil {
			s.logger.Error("update error", "status", res.Status())
		}
		return fmt.Errorf("update error: %s", res.String())
	}

	return nil
}

// Delete 删除文档
func (s *ESSync) Delete(indexName, docID string) error {
	req := esapi.DeleteRequest{
		Index:      indexName,
		DocumentID: docID,
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), s.client)
	if err != nil {
		if s.logger != nil {
			s.logger.Error("delete document failed", "error", err)
		}
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		// 404 not found 不算错误
		if res.StatusCode == 404 {
			log.Printf("Document not found: %s/%s", indexName, docID)
			return nil
		}
		if s.logger != nil {
			s.logger.Error("delete error", "status", res.Status())
		}
		return fmt.Errorf("delete error: %s", res.String())
	}

	return nil
}