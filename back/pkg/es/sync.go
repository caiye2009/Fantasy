// back/pkg/es/sync.go
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

// Index 索引文档（支持自动评分）
func (s *ESSync) Index(doc interface{}) error {
	// 类型断言检查是否实现 Indexable 接口
	indexable, ok := doc.(Indexable)
	if !ok {
		return fmt.Errorf("document does not implement Indexable interface")
	}

	// 获取文档数据
	docData := indexable.ToDocument()

	// 尝试计算 priorityScore（通过接口类型断言）
	if scorer, ok := doc.(PriorityScorer); ok {
		if priorityScore := scorer.CalculatePriorityScore(); priorityScore > 0 {
			docData["priorityScore"] = priorityScore

			if s.logger != nil {
				s.logger.Info("calculated priority score",
					"index", indexable.GetIndexName(),
					"docId", indexable.GetDocumentID(),
					"score", priorityScore,
				)
			}
		}
	}

	data, err := json.Marshal(docData)
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

// Update 更新文档（支持自动评分）
func (s *ESSync) Update(doc interface{}) error {
	// 类型断言检查是否实现 Indexable 接口
	indexable, ok := doc.(Indexable)
	if !ok {
		return fmt.Errorf("document does not implement Indexable interface")
	}

	// 获取文档数据
	docData := indexable.ToDocument()

	// 尝试计算 priorityScore（通过接口类型断言）
	if scorer, ok := doc.(PriorityScorer); ok {
		if priorityScore := scorer.CalculatePriorityScore(); priorityScore > 0 {
			docData["priorityScore"] = priorityScore

			if s.logger != nil {
				s.logger.Info("calculated priority score",
					"index", indexable.GetIndexName(),
					"docId", indexable.GetDocumentID(),
					"score", priorityScore,
				)
			}
		}
	}

	data, err := json.Marshal(map[string]interface{}{
		"doc": docData,
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

// BulkIndex 批量索引文档（更高效，适合大批量操作）
// batchSize 指定每批处理的文档数量（建议 100-500）
// refresh 控制是否在批次完成后刷新索引（reindex 时建议 false，只在最后刷新）
func (s *ESSync) BulkIndex(docs []interface{}, refresh bool) (success, failed int, err error) {
	if len(docs) == 0 {
		return 0, 0, nil
	}

	var buf bytes.Buffer
	successCount := 0
	failedCount := 0

	// 构建 bulk 请求体
	for _, doc := range docs {
		indexable, ok := doc.(Indexable)
		if !ok {
			failedCount++
			continue
		}

		// 获取文档数据
		docData := indexable.ToDocument()

		// 尝试计算 priorityScore（通过接口类型断言）
		if scorer, ok := doc.(PriorityScorer); ok {
			if priorityScore := scorer.CalculatePriorityScore(); priorityScore > 0 {
				docData["priorityScore"] = priorityScore
			}
		}

		// Bulk API 格式：action_and_meta_data\n + optional_source\n
		// Index action
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": indexable.GetIndexName(),
				"_id":    indexable.GetDocumentID(),
			},
		}
		metaData, err := json.Marshal(meta)
		if err != nil {
			failedCount++
			continue
		}
		buf.Write(metaData)
		buf.WriteByte('\n')

		// Document source
		docJSON, err := json.Marshal(docData)
		if err != nil {
			failedCount++
			continue
		}
		buf.Write(docJSON)
		buf.WriteByte('\n')
	}

	// 发送 bulk 请求
	refreshParam := "false"
	if refresh {
		refreshParam = "true"
	}

	res, err := s.client.Bulk(
		bytes.NewReader(buf.Bytes()),
		s.client.Bulk.WithContext(context.Background()),
		s.client.Bulk.WithRefresh(refreshParam),
	)
	if err != nil {
		if s.logger != nil {
			s.logger.Error("bulk request failed", "error", err)
		}
		return 0, len(docs), err
	}
	defer res.Body.Close()

	if res.IsError() {
		if s.logger != nil {
			s.logger.Error("bulk request error", "status", res.Status())
		}
		return 0, len(docs), fmt.Errorf("bulk request error: %s", res.Status())
	}

	// 解析响应
	var bulkRes map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&bulkRes); err != nil {
		if s.logger != nil {
			s.logger.Error("failed to parse bulk response", "error", err)
		}
		return 0, len(docs), err
	}

	// 统计成功和失败
	if items, ok := bulkRes["items"].([]interface{}); ok {
		for _, item := range items {
			if itemMap, ok := item.(map[string]interface{}); ok {
				for _, action := range itemMap {
					if actionMap, ok := action.(map[string]interface{}); ok {
						status := int(actionMap["status"].(float64))
						if status >= 200 && status < 300 {
							successCount++
						} else {
							failedCount++
							if s.logger != nil {
								s.logger.Error("bulk item failed",
									"status", status,
									"error", actionMap["error"],
								)
							}
						}
					}
				}
			}
		}
	}

	return successCount, failedCount, nil
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