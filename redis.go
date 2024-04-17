// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package gonet

import (
	"github.com/redis/go-redis/v9"
	"reflect"
)

type RedisService interface {
	Set(key string, value any) error
	Get(key string) (any, error)
}

type RedisOptions = redis.Options

type redisService struct {
	options *RedisOptions
}

func newRedisService(options *RedisOptions) *redisService {
	return &redisService{}
}

func AddRedis(options *RedisOptions) {
	AddSingleton[RedisService](func() *redisService {
		return newRedisService(options)
	})
}

func AddRedisFactory(factory any) {
	factoryVal := reflect.ValueOf(factory)
	deps := resolveValueDeps(factoryVal)
	options := factoryVal.Call(deps)[0].Interface().(*RedisOptions)

	AddSingleton[RedisService](func() *redisService {
		return newRedisService(options)
	})
}
