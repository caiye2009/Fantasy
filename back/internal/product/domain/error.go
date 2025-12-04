package domain

import "errors"

var (
	// 基础校验
	ErrProductNotFound       = errors.New("product not found")
	ErrProductNameEmpty      = errors.New("product name cannot be empty")
	ErrProductNameInvalid    = errors.New("product name must be between 2 and 100 characters")

	// 材料与工序
	ErrMaterialsRequired     = errors.New("at least one material is required")
	ErrProcessesRequired     = errors.New("at least one process is required")

	// 材料比例
	ErrInvalidMaterialRatio  = errors.New("material ratio must be between 0 and 1")
	ErrMaterialRatioSumNotOne = errors.New("sum of material ratios must equal 1.0")

	// 工艺数量
	ErrInvalidProcessQuantity = errors.New("process quantity must be greater than 0")

	// 状态流转
	ErrCannotSubmit                 = errors.New("can only submit draft products")
	ErrCannotApprove                = errors.New("can only approve submitted products")
	ErrCannotReject                 = errors.New("can only reject submitted products")
	ErrCannotUpdateApprovedProduct  = errors.New("cannot update approved products")
	ErrCannotDeleteApprovedProduct  = errors.New("cannot delete approved products")
)
