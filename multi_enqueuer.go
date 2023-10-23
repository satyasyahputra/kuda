package kuda

import (
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

type KudaEnqueuerContext struct {
	RedisPool *redis.Pool
	Enqueuers map[string]*work.Enqueuer
}

func NewKudaEnqueuerContext(rp *redis.Pool, queues []string) (*KudaEnqueuerContext, error) {
	enqueuersMap := map[string]*work.Enqueuer{}
	for _, queue := range queues {
		q, _, err := ExtractQueue(queue)
		if err != nil {
			return nil, err
		}
		enqueuersMap[q] = work.NewEnqueuer(q, redisPool)
	}

	return &KudaEnqueuerContext{
		RedisPool: rp,
		Enqueuers: enqueuersMap,
	}, nil
}
