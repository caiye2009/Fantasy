package auth

import (
    "context"
    "errors"
    "fmt"
    "time"
    
    "github.com/redis/go-redis/v9"
)

// WhiteListWang Redis白名单管理器
type WhiteListWang struct {
    redis *redis.Client
}

// NewWhiteListWang 创建白名单管理器
func NewWhiteListWang(redis *redis.Client) *WhiteListWang {
    return &WhiteListWang{
        redis: redis,
    }
}

// SaveJTI 保存JTI到白名单
// deviceType: "mobile" (H5不使用此方法)
func (w *WhiteListWang) SaveJTI(
    ctx context.Context,
    userID int64,
    deviceType string,
    accessJTI string,
    refreshJTI string,
) error {
    accessKey := fmt.Sprintf("user:%d:%s:access_jti", userID, deviceType)
    refreshKey := fmt.Sprintf("user:%d:%s:refresh_jti", userID, deviceType)
    
    pipe := w.redis.Pipeline()
    pipe.Set(ctx, accessKey, accessJTI, 16*time.Minute)       // AccessToken TTL: 16分钟
    pipe.Set(ctx, refreshKey, refreshJTI, 31*24*time.Hour)    // RefreshToken TTL: 31天
    
    if _, err := pipe.Exec(ctx); err != nil {
        return fmt.Errorf("failed to save jti: %w", err)
    }
    
    return nil
}

// UpdateAccessJTI 只更新AccessToken的JTI(刷新Token时使用)
func (w *WhiteListWang) UpdateAccessJTI(
    ctx context.Context,
    userID int64,
    deviceType string,
    accessJTI string,
) error {
    key := fmt.Sprintf("user:%d:%s:access_jti", userID, deviceType)
    
    if err := w.redis.Set(ctx, key, accessJTI, 16*time.Minute).Err(); err != nil {
        return fmt.Errorf("failed to update access jti: %w", err)
    }
    
    return nil
}

// ValidateAccessJTI 验证AccessToken的JTI
func (w *WhiteListWang) ValidateAccessJTI(
    ctx context.Context,
    userID int64,
    deviceType string,
    jti string,
) error {
    key := fmt.Sprintf("user:%d:%s:access_jti", userID, deviceType)
    
    storedJTI, err := w.redis.Get(ctx, key).Result()
    if err == redis.Nil {
        return errors.New("session expired")
    }
    if err != nil {
        return fmt.Errorf("failed to get access jti: %w", err)
    }
    
    if storedJTI != jti {
        return errors.New("session replaced, your account has been logged in on another device")
    }
    
    return nil
}

// ValidateRefreshJTI 验证RefreshToken的JTI
func (w *WhiteListWang) ValidateRefreshJTI(
    ctx context.Context,
    userID int64,
    deviceType string,
    jti string,
) error {
    key := fmt.Sprintf("user:%d:%s:refresh_jti", userID, deviceType)
    
    storedJTI, err := w.redis.Get(ctx, key).Result()
    if err == redis.Nil {
        return errors.New("refresh token expired")
    }
    if err != nil {
        return fmt.Errorf("failed to get refresh jti: %w", err)
    }
    
    if storedJTI != jti {
        return errors.New("refresh token replaced")
    }
    
    return nil
}

// DeleteJTI 删除JTI(登出时使用)
func (w *WhiteListWang) DeleteJTI(
    ctx context.Context,
    userID int64,
    deviceType string,
) error {
    keys := []string{
        fmt.Sprintf("user:%d:%s:access_jti", userID, deviceType),
        fmt.Sprintf("user:%d:%s:refresh_jti", userID, deviceType),
    }
    
    if err := w.redis.Del(ctx, keys...).Err(); err != nil {
        return fmt.Errorf("failed to delete jti: %w", err)
    }
    
    return nil
}