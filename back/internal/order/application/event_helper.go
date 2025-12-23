package application

import (
	"encoding/json"
	"fmt"
	"back/internal/order/domain"
)

// createEvent 创建事件的辅助函数
func createEvent(orderID uint, eventType string, operatorID uint, operatorName string, operatorRole string, beforeData, afterData interface{}, description string) *domain.OrderEvent {
	var beforeJSON, afterJSON string

	if beforeData != nil {
		if b, err := json.Marshal(beforeData); err == nil {
			beforeJSON = string(b)
		}
	}

	if afterData != nil {
		if b, err := json.Marshal(afterData); err == nil {
			afterJSON = string(b)
		}
	}

	return &domain.OrderEvent{
		OrderID:      orderID,
		EventType:    eventType,
		OperatorID:   operatorID,
		OperatorName: operatorName,
		OperatorRole: operatorRole,
		BeforeData:   beforeJSON,
		AfterData:    afterJSON,
		Description:  description,
	}
}

// getProgressName 获取进度类型的中文名称
func getProgressName(progressType string) string {
	names := map[string]string{
		domain.ProgressTypeFabricInput:    "胚布投入进度",
		domain.ProgressTypeProduction:     "生产进度",
		domain.ProgressTypeWarehouseCheck: "仓库验货进度",
		domain.ProgressTypeRework:         "回修进度",
	}

	if name, ok := names[progressType]; ok {
		return name
	}
	return progressType
}

// getProgressIcon 获取进度类型的图标
func getProgressIcon(progressType string) string {
	icons := map[string]string{
		domain.ProgressTypeFabricInput:    "lucide:package-open",
		domain.ProgressTypeProduction:     "lucide:factory",
		domain.ProgressTypeWarehouseCheck: "lucide:clipboard-check",
		domain.ProgressTypeRework:         "lucide:wrench",
	}

	if icon, ok := icons[progressType]; ok {
		return icon
	}
	return ""
}

// getProgressColor 获取进度类型的颜色
func getProgressColor(progressType string) string {
	colors := map[string]string{
		domain.ProgressTypeFabricInput:    "#409EFF",
		domain.ProgressTypeProduction:     "#67C23A",
		domain.ProgressTypeWarehouseCheck: "#E6A23C",
		domain.ProgressTypeRework:         "#F56C6C",
	}

	if color, ok := colors[progressType]; ok {
		return color
	}
	return "#909399"
}

// calculateOverallProgress 计算订单整体进度
func calculateOverallProgress(progresses []domain.OrderProgress) int {
	existingProgresses := []domain.OrderProgress{}
	for _, p := range progresses {
		if p.Exists {
			existingProgresses = append(existingProgresses, p)
		}
	}

	if len(existingProgresses) == 0 {
		return 0
	}

	totalProgress := 0
	for _, p := range existingProgresses {
		totalProgress += p.Progress
	}

	return totalProgress / len(existingProgresses)
}

// formatQuantityChange 格式化数量变化描述
func formatQuantityChange(progressName string, before, after float64) string {
	return fmt.Sprintf("%s：%.0f → %.0f", progressName, before, after)
}
