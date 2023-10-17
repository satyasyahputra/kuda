package kuda

import (
	"os"
	"os/signal"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

type KudaProcessor struct {
	qName       string
	concurrency uint
	redisPool   *redis.Pool
	jobMap      map[string]interface{}
	middlewares func(pool *work.WorkerPool)
}

type ProcessorContext struct {
	args string
}

func NewKudaProcessor(qName string, concurrency uint, redisPool *redis.Pool, jobMap map[string]interface{}, middlewares func(pool *work.WorkerPool)) *KudaProcessor {
	processor := &KudaProcessor{
		qName:       qName,
		concurrency: concurrency,
		redisPool:   redisPool,
		jobMap:      jobMap,
		middlewares: middlewares,
	}
	return processor
}

func CreateKudaProcessor(kudaProcessor *KudaProcessor) *work.WorkerPool {
	pool := work.NewWorkerPool(ProcessorContext{}, kudaProcessor.concurrency, kudaProcessor.qName, kudaProcessor.redisPool)

	registerJobs(pool, kudaProcessor.jobMap)

	kudaProcessor.middlewares(pool)

	return pool
}

func RunProcessors(pools []*work.WorkerPool) {
	for _, wp := range pools {
		go RunProcessor(wp)
	}
}

func RunProcessor(pool *work.WorkerPool) {
	pool.Start()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	pool.Stop()
}

func registerJobs(pool *work.WorkerPool, jobMap map[string]interface{}) error {
	for jobName, fn := range jobMap {
		pool.Job(jobName, fn)
	}
	return nil
}

func customizeOptions(pool *work.WorkerPool) error {
	pool.JobWithOptions("export", work.JobOptions{Priority: 10, MaxFails: 5}, (*ProcessorContext).Export)
	return nil
}

func (c *ProcessorContext) Export(job *work.Job) error {
	return nil
}
