package dbs

import (
	"context"
	"fmt"
	"time"
)

const USER_SERVER_MAPPING_KEY = "user_server_mapping:%s"

// SetUserServer 设置用户-服务器映射
func SetUserServer(userID, serverIP string) error {
	ctx := context.Background()
	key := fmt.Sprintf(USER_SERVER_MAPPING_KEY, userID)

	// 设置映射，并设置过期时间（例如30分钟无活动则自动清除）
	return RedisEngine.Set(ctx, key, serverIP, 30*time.Minute).Err()
}

// GetUserServer 获取用户所在服务器
func GetUserServer(userID string) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf(USER_SERVER_MAPPING_KEY, userID)

	return RedisEngine.Get(ctx, key).Result()
}

// RemoveUserServer 删除用户-服务器映射（用户断开连接时）
func RemoveUserServer(userID string) error {
	ctx := context.Background()
	key := fmt.Sprintf(USER_SERVER_MAPPING_KEY, userID)
	return RedisEngine.Del(ctx, key).Err()
}
