package kuda

import (
	"log"
	"os"
	"os/signal"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

type KudaProcessor struct {
	concurrency uint
	redisPool   *redis.Pool
	middlewares func(pool *work.WorkerPool)
}

type ProcessorContext struct {
	args string
}

func NewKudaProcessor(concurrency uint, redisPool *redis.Pool, middlewares func(pool *work.WorkerPool)) *KudaProcessor {
	processor := &KudaProcessor{
		concurrency: concurrency,
		redisPool:   redisPool,
		middlewares: middlewares,
	}
	return processor
}

func CreateKudaProcessor(kudaProcessor *KudaProcessor, qName string, jobMap map[string]func(pc *ProcessorContext, job *work.Job) error) *work.WorkerPool {
	pool := work.NewWorkerPool(ProcessorContext{}, kudaProcessor.concurrency, qName, kudaProcessor.redisPool)

	registerJobs(pool, jobMap)

	kudaProcessor.middlewares(pool)

	return pool
}

func RunProcessors(pools []*work.WorkerPool) {
	for _, wp := range pools {
		wp.Start()
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	for _, wp := range pools {
		log.Println("stop")
		wp.Stop()
	}
}

func RunProcessor(pool *work.WorkerPool) {
	pool.Start()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	pool.Stop()
}

func registerJobs(pool *work.WorkerPool, jobMap map[string]func(pc *ProcessorContext, job *work.Job) error) {
	for jobName, fn := range jobMap {
		pool.Job(jobName, fn)
	}
}

func customizeOptions(pool *work.WorkerPool) {
	pool.JobWithOptions("export", work.JobOptions{Priority: 10, MaxFails: 5}, (*ProcessorContext).Export)
}

func (c *ProcessorContext) Export(job *work.Job) error {
	return nil
}
