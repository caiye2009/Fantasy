// back/pkg/es/indexable.go
package es

// Indexable ES 文档接口
type Indexable interface {
	GetIndexName() string
	GetDocumentID() string
	ToDocument() map[string]interface{}
}

// PriorityScorer 优先级评分器接口
type PriorityScorer interface {
	CalculateScore(entity interface{}) int
}

// ScoredIndexable 带评分的 ES 文档接口
type ScoredIndexable interface {
	Indexable
	GetScorer() PriorityScorer
}