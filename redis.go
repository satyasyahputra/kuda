package kuda

import "github.com/gomodule/redigo/redis"

var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", "localhost:6379")
	},
}

type KudaRedis struct {
	MaxActive int    `env:"MAX_ACTIVE" envDefault:"5"`
	MaxIdle   int    `env:"MAX_IDLE" envDefault:"5"`
	Wait      bool   `env:"WAIT" envDefault:"true"`
	Host      string `env:"HOST" envDefault:"localhost"`
	Port      string `env:"PORT" envDefault:"6379"`
}

func (kr *KudaRedis) NewRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxActive: kr.MaxActive,
		MaxIdle:   kr.MaxIdle,
		Wait:      kr.Wait,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", kr.Host+":"+kr.Port)
		},
	}
}
