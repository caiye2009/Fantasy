package auth

import (
	"context"
	"fmt"
	"time"
	"github.com/redis/go-redis/v9"
)

const (
	// JWT 白名单的 Redis key 前缀
	whitelistKeyPrefix = "jwt:whitelist"
	// JWT 过期时间（和 accessToken 一致）
	jwtExpiration = 1 * time.Hour
)

// WhitelistManager JWT 白名单管理器
type WhitelistManager struct {
	rdb *redis.Client
}

// NewWhitelistManager 创建白名单管理器
func NewWhitelistManager(rdb *redis.Client) *WhitelistManager {
	return &WhitelistManager{rdb: rdb}
}

// AddToWhitelist 添加 JWT 到白名单
// 一个 loginId + source 组合只能有一个有效的 JTI（新的会覆盖旧的）
func (w *WhitelistManager) AddToWhitelist(loginID string, source string, jti string) error {
	ctx := context.Background()
	key := w.buildKey(loginID, source)

	// 设置白名单，过期时间和 JWT 一致
	return w.rdb.Set(ctx, key, jti, jwtExpiration).Err()
}

// IsInWhitelist 检查 JWT 是否在白名单中
func (w *WhitelistManager) IsInWhitelist(loginID string, source string, jti string) (bool, error) {
	ctx := context.Background()
	key := w.buildKey(loginID, source)

	// 从 Redis 获取当前白名单中的 JTI
	currentJTI, err := w.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		// key 不存在，说明不在白名单
		return false, nil
	}
	if err != nil {
		return false, err
	}

	// 比较 JTI 是否匹配
	return currentJTI == jti, nil
}

// RemoveFromWhitelist 移除指定端的 JWT（登出时使用）
func (w *WhitelistManager) RemoveFromWhitelist(loginID string, source string) error {
	ctx := context.Background()
	key := w.buildKey(loginID, source)

	return w.rdb.Del(ctx, key).Err()
}

// RemoveAllForUser 移除用户所有端的 JWT（Role/Department 变化时使用）
func (w *WhitelistManager) RemoveAllForUser(loginID string) error {
	ctx := context.Background()

	// 移除 H5 端和 Mobile 端的白名单
	h5Key := w.buildKey(loginID, "h5")
	mobileKey := w.buildKey(loginID, "mobile")

	pipe := w.rdb.Pipeline()
	pipe.Del(ctx, h5Key)
	pipe.Del(ctx, mobileKey)

	_, err := pipe.Exec(ctx)
	return err
}

// buildKey 构建 Redis key
// 格式：jwt:whitelist:{loginId}:{source}
func (w *WhitelistManager) buildKey(loginID string, source string) string {
	return fmt.Sprintf("%s:%s:%s", whitelistKeyPrefix, loginID, source)
}
