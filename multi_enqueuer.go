package kuda

import (
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var enqueuers = map[string]*work.Enqueuer{}

type KudaEnqueuerContext struct {
	RedisPool *redis.Pool
	Enqueuers map[string]*work.Enqueuer
}

func NewKudaEnqueuerContext(rp *redis.Pool, queues []string) *KudaEnqueuerContext {
	enqueuersMap := map[string]*work.Enqueuer{}
	for _, name := range queues {
		enqueuersMap[name] = work.NewEnqueuer(name, redisPool)
	}

	return &KudaEnqueuerContext{
		RedisPool: rp,
		Enqueuers: enqueuersMap,
	}
}
