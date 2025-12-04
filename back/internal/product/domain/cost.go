package domain

// MaterialCostItem 原料成本明细
type MaterialCostItem struct {
	MaterialID   uint    `json:"material_id"`
	MaterialName string  `json:"material_name"`
	Ratio        float64 `json:"ratio"`
	Price        float64 `json:"price"`
	Cost         float64 `json:"cost"`
}

// ProcessCostItem 工艺成本明细
type ProcessCostItem struct {
	ProcessID   uint    `json:"process_id"`
	ProcessName string  `json:"process_name"`
	Price       float64 `json:"price"`
	Cost        float64 `json:"cost"`
}

// CostBreakdown 成本明细
type CostBreakdown struct {
	Materials []MaterialCostItem `json:"materials"`
	Processes []ProcessCostItem  `json:"processes"`
}

// CostResult 成本计算结果
type CostResult struct {
	ProductID    uint           `json:"product_id"`
	ProductName  string         `json:"product_name"`
	Quantity     float64        `json:"quantity"`
	MaterialCost float64        `json:"material_cost"`
	ProcessCost  float64        `json:"process_cost"`
	UnitCost     float64        `json:"unit_cost"`
	TotalCost    float64        `json:"total_cost"`
	Breakdown    *CostBreakdown `json:"breakdown"`
}