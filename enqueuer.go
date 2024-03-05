package kuda

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

type KudaEnqueuerContext struct {
	RedisPool *redis.Pool
	Enqueuers map[string]*work.Enqueuer
}

type KudaQueue struct {
	Queues string `env:"QUEUES"`
}

func NewKudaEnqueuerContext(rp *redis.Pool, queues []string) (*KudaEnqueuerContext, error) {
	enqueuersMap := map[string]*work.Enqueuer{}
	for _, queue := range queues {
		q, _, err := ExtractQueue(queue)
		if err != nil {
			return nil, err
		}
		enqueuersMap[q] = work.NewEnqueuer(q, rp)
	}

	return &KudaEnqueuerContext{
		RedisPool: rp,
		Enqueuers: enqueuersMap,
	}, nil
}

func ExtractQueue(queue string) (qName string, concurrency uint, err error) {
	arr := strings.Split(queue, ":")
	if len(arr) != 2 {
		return "", 0, errors.New("invalid queue format")
	}

	qName = arr[0]

	u64, err := strconv.ParseUint(arr[1], 10, 32)
	if err != nil {
		return qName, concurrency, err
	}
	concurrency = uint(u64)

	return qName, concurrency, err
}
