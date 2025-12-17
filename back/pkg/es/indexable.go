// back/pkg/es/indexable.go
package es

// Indexable ES 文档接口
type Indexable interface {
	GetIndexName() string
	GetDocumentID() string
	ToDocument() map[string]interface{}
}

// 优先级评分机制说明：
// 如果 Domain 需要业务优先级排序，只需在 Domain 中实现：
//   func (d *DomainModel) CalculatePriorityScore() int
// ES 同步时会自动通过反射调用该方法。
// 如果不实现或返回 0，则不会添加 priorityScore 字段。