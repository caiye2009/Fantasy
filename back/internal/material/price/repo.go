package price

import "gorm.io/gorm"

type MaterialPriceRepo struct {
	db *gorm.DB
}

func NewMaterialPriceRepo(db *gorm.DB) *MaterialPriceRepo {
	return &MaterialPriceRepo{db: db}
}

func (r *MaterialPriceRepo) Create(price *MaterialPrice) error {
	return r.db.Create(price).Error
}

func (r *MaterialPriceRepo) GetMinPrice(materialID uint) (*MaterialPrice, error) {
	var price MaterialPrice
	err := r.db.Where("material_id = ?", materialID).
		Order("price ASC, quoted_at DESC").
		First(&price).Error
	return &price, err
}

func (r *MaterialPriceRepo) GetMaxPrice(materialID uint) (*MaterialPrice, error) {
	var price MaterialPrice
	err := r.db.Where("material_id = ?", materialID).
		Order("price DESC, quoted_at DESC").
		First(&price).Error
	return &price, err
}

func (r *MaterialPriceRepo) GetHistory(materialID uint) ([]MaterialPrice, error) {
	var prices []MaterialPrice
	err := r.db.Where("material_id = ?", materialID).
		Order("quoted_at DESC").
		Find(&prices).Error
	return prices, err
}