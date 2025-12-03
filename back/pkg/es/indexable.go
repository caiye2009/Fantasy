package es

// Indexable ES 可索引接口
type Indexable interface {
	// GetIndexName 获取索引名称
	GetIndexName() string

	// GetDocumentID 获取文档 ID
	GetDocumentID() string

	// ToDocument 转换为 ES 文档
	ToDocument() map[string]interface{}
}