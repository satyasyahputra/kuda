package kuda

import "github.com/gomodule/redigo/redis"

type KudaRedis struct {
	MaxActive int    `env:"MAX_ACTIVE" envDefault:"5"`
	MaxIdle   int    `env:"MAX_IDLE" envDefault:"5"`
	Wait      bool   `env:"WAIT" envDefault:"true"`
	Endpoint  string `env:"ENDPOINT" envDefault:"localhost:6379"`
}

func (kr *KudaRedis) NewRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxActive: kr.MaxActive,
		MaxIdle:   kr.MaxIdle,
		Wait:      kr.Wait,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", kr.Endpoint)
		},
	}
}
