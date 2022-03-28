package inits

import (
	"hitokoto-go/global"
	"os"
	"strings"
)

func Config() error {
	// From ENV

	global.Config.PGConnString = os.Getenv("POSTGRES_CONNECTION_STRING")
	global.Config.RedisConnString = os.Getenv("REDIS_CONNECTION_STRING")
	mode := os.Getenv("MODE")
	global.Config.IsProdMode = strings.HasPrefix(strings.ToLower(mode), "prod")

	return nil
}
