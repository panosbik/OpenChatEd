package redis

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	ctx    = context.Background()
)

// Assuming you have a Redis client already set up
func SaveRefreshToken(token string, userID uint) error {
	err := Client.Set(ctx, token, userID, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetUserIDByRefreshToken(token string) (userID uint, err error) {
	userIDString, err := Client.Get(ctx, token).Result()
	if err == redis.Nil {
		// No user ID found for this refresh token
		return
	} else if err != nil {
		return
	}
	// Convert the user ID string to uint
	u64, _ := strconv.ParseUint(userIDString, 10, 32)
	userID = uint(u64)

	// Delete the refresh token
	err = Client.Del(ctx, token).Err()
	if err != nil {
		return
	}

	return
}
