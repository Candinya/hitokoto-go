package inits

import (
	"context"
	"github.com/go-redis/redis/v8"
	"hitokoto-go/global"
	"time"
)

func Redis() error {
	ctx, _ := context.WithTimeout(context.TODO(), 3*time.Second)
	opt, err := redis.ParseURL(global.Config.RedisConnString)
	if err != nil {
		return err
	}
	global.Redis = redis.NewClient(opt)
	_, err = global.Redis.Ping(ctx).Result()
	return err
}
