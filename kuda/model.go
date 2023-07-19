package kuda

import (
	"github.com/gomodule/redigo/redis"
)

type KudaProcessor struct {
	qName       string
	concurrency uint
	redisPool   *redis.Pool
	jobMap      map[string]interface{}
	middlewares interface{}
}

func NewKudaProcessor(qName string, concurrency uint, redisPool *redis.Pool, jobMap map[string]interface{}, middlewares interface{}) *KudaProcessor {
	processor := &KudaProcessor{
		qName:       qName,
		concurrency: concurrency,
		redisPool:   redisPool,
		jobMap:      jobMap,
		middlewares: middlewares,
	}
	return processor
}
