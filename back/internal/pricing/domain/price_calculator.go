package domain

// PriceCalculator 价格计算器（领域服务）
type PriceCalculator struct{}

// NewPriceCalculator 创建价格计算器
func NewPriceCalculator() *PriceCalculator {
	return &PriceCalculator{}
}

// CalculateAverage 计算平均价格
func (c *PriceCalculator) CalculateAverage(prices []*SupplierPrice) float64 {
	if len(prices) == 0 {
		return 0
	}
	
	total := 0.0
	for _, p := range prices {
		total += p.Price
	}
	
	return total / float64(len(prices))
}

// FindLowest 找到最低价格
func (c *PriceCalculator) FindLowest(prices []*SupplierPrice) *SupplierPrice {
	if len(prices) == 0 {
		return nil
	}
	
	lowest := prices[0]
	for _, p := range prices[1:] {
		if p.IsLowerThan(lowest.Price) {
			lowest = p
		}
	}
	
	return lowest
}

// FindHighest 找到最高价格
func (c *PriceCalculator) FindHighest(prices []*SupplierPrice) *SupplierPrice {
	if len(prices) == 0 {
		return nil
	}
	
	highest := prices[0]
	for _, p := range prices[1:] {
		if p.IsHigherThan(highest.Price) {
			highest = p
		}
	}
	
	return highest
}