package kuda

import (
	"os"
	"os/signal"
	"reflect"

	"github.com/gocraft/work"
)

type ProcessorContext struct {
	args string
}

type KudaWorkerPool struct {
	pool *work.WorkerPool
}

func CreateKudaProcessor(kudaProcessor *KudaProcessor) *work.WorkerPool {
	pool := work.NewWorkerPool(ProcessorContext{}, kudaProcessor.concurrency, kudaProcessor.qName, kudaProcessor.redisPool)

	registerJobs(pool, kudaProcessor.jobMap)

	registerMiddleware(pool, kudaProcessor.middlewares)

	return pool
}

func SimpleProcessor(kudaProcessor *KudaProcessor) *work.WorkerPool {
	pool := work.NewWorkerPool(ProcessorContext{}, kudaProcessor.concurrency, kudaProcessor.qName, kudaProcessor.redisPool)

	ProcessorMiddleware(pool)

	registerJobs(pool, kudaProcessor.jobMap)

	customizeOptions(pool)

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

func registerMiddleware(pool *work.WorkerPool, fn interface{}) error {
	value := reflect.ValueOf(fn)
	in := make([]reflect.Value, 1)
	in[0] = reflect.ValueOf(pool)
	value.Call(in)
	return nil
}

func customizeOptions(pool *work.WorkerPool) error {
	pool.JobWithOptions("export", work.JobOptions{Priority: 10, MaxFails: 5}, (*ProcessorContext).Export)
	return nil
}

func (c *ProcessorContext) Export(job *work.Job) error {
	return nil
}
