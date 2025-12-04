package es

// Indexable ES 文档接口
type Indexable interface {
	GetIndexName() string
	GetDocumentID() string
	ToDocument() map[string]interface{}
}