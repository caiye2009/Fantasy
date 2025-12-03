package price

import "gorm.io/gorm"

type ProcessPriceRepo struct {
	db *gorm.DB
}

func NewProcessPriceRepo(db *gorm.DB) *ProcessPriceRepo {
	return &ProcessPriceRepo{db: db}
}

func (r *ProcessPriceRepo) Create(price *ProcessPrice) error {
	return r.db.Create(price).Error
}

func (r *ProcessPriceRepo) GetMinPrice(processID uint) (*ProcessPrice, error) {
	var price ProcessPrice
	err := r.db.Where("process_id = ?", processID).
		Order("price ASC, quoted_at DESC").
		First(&price).Error
	return &price, err
}

func (r *ProcessPriceRepo) GetMaxPrice(processID uint) (*ProcessPrice, error) {
	var price ProcessPrice
	err := r.db.Where("process_id = ?", processID).
		Order("price DESC, quoted_at DESC").
		First(&price).Error
	return &price, err
}

func (r *ProcessPriceRepo) GetHistory(processID uint) ([]ProcessPrice, error) {
	var prices []ProcessPrice
	err := r.db.Where("process_id = ?", processID).
		Order("quoted_at DESC").
		Find(&prices).Error
	return prices, err
}