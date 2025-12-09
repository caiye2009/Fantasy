// back/pkg/es/scorer.go (新建文件)
package es

import (
	"time"
)

// BaseScorer 基础评分器（其他评分器可以嵌入这个）
type BaseScorer struct{}

// 计算时间衰减分数（越新的分数越高）
func (s *BaseScorer) CalculateTimeBonus(updatedAt time.Time, maxBonus int) int {
	now := time.Now()
	hoursSinceUpdate := now.Sub(updatedAt).Hours()

	// 24小时内：满分
	if hoursSinceUpdate < 24 {
		return maxBonus
	}
	// 24-72小时：递减
	if hoursSinceUpdate < 72 {
		return maxBonus / 2
	}
	// 72小时以上：最低分
	return 0
}

// 计算时间紧急度（越久未处理分数越高）
func (s *BaseScorer) CalculateUrgencyBonus(createdAt time.Time) int {
	now := time.Now()
	hoursSinceCreated := now.Sub(createdAt).Hours()

	// 超过7天：紧急
	if hoursSinceCreated > 168 { // 7天
		return 200
	}
	// 3-7天：较紧急
	if hoursSinceCreated > 72 {
		return 100
	}
	// 1-3天：一般
	if hoursSinceCreated > 24 {
		return 50
	}
	// 24小时内：不紧急
	return 0
}

// DefaultScorer 默认评分器（用于没有特定评分规则的资源）
type DefaultScorer struct {
	BaseScorer
}

func NewDefaultScorer() *DefaultScorer {
	return &DefaultScorer{}
}

func (s *DefaultScorer) CalculateScore(entity interface{}) int {
	// 默认中等优先级
	return 500
}