package domain

import "errors"

var (
	// 验证错误
	ErrProductIDRequired      = errors.New("产品ID不能为空")
	ErrCategoryRequired       = errors.New("类别不能为空")
	ErrBatchIDRequired        = errors.New("批次ID不能为空")
	ErrInvalidQuantity        = errors.New("数量必须大于0")
	ErrUnitRequired           = errors.New("单位不能为空")
	ErrInvalidUnitCost        = errors.New("单价不能为负数")
	ErrInvalidTotalCost       = errors.New("总成本不能为负数")

	// 业务错误
	ErrInventoryNotFound      = errors.New("库存不存在")
	ErrInsufficientInventory  = errors.New("库存不足")
)
