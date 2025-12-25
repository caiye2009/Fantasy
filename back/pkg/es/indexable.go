// back/pkg/es/indexable.go
package es

// Indexable ES 文档接口
type Indexable interface {
	GetIndexName() string
	GetDocumentID() string
	ToDocument() map[string]interface{}
}

// PriorityScorer 优先级评分接口
// 如果 Domain 需要业务优先级排序，只需实现此接口
// ES 同步时会自动调用 CalculatePriorityScore() 方法
// 如果不实现或返回 0，则不会添加 priorityScore 字段
type PriorityScorer interface {
	CalculatePriorityScore() int
}