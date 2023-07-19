package kuda

import (
	"errors"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var enqueuers = map[string]*work.Enqueuer{}

func RegisterEnqueuer(name string, redisPool *redis.Pool) (*work.Enqueuer, error) {
	enqueuers[name] = work.NewEnqueuer(name, redisPool)
	return enqueuers[name], nil
}

func GetEnqueuer(name string) (*work.Enqueuer, error) {
	nQR, exist := enqueuers[name]
	if !exist {
		return nil, errors.New("Enqueuer not found: " + name)
	}
	return nQR, nil
}
