package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gogf/gf/database/gredis"
)

var (
	goRedisClientCache = make(map[gredis.Config]*redis.Client)
)

func GetGoRedisClient(config *gredis.Config) *redis.Client {
	if goRedisClientCache[*config] == nil {
		addr := fmt.Sprintf("%v:%v", config.Host, config.Port)
		if goRedisClientCache[*config] == nil {
			client := redis.NewClient(&redis.Options{
				Addr:        addr,
				DB:          config.Db,
				Password:    config.Pass,
				PoolSize:    config.MaxActive,
				MaxConnAge:  config.MaxConnLifetime,
				IdleTimeout: config.IdleTimeout,
			})
			client.AddHook(&PrefixHook{})
			goRedisClientCache[*config] = client
		}
	}
	return goRedisClientCache[*config]
}

func GetGRedisConfig(group string) *gredis.Config {
	return nil
}
