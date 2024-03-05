package kuda

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

type KudaProcessor struct {
	redisPool   *redis.Pool
	middlewares func(pool *work.WorkerPool)
}

type ProcessorContext struct {
	Args   string `json:"args,omitempty"`
	Queue  string `json:"queue,omitempty"`
	Worker string `json:"worker,omitempty"`
}

func DecodeJSON(pc *ProcessorContext, input io.Reader, w http.ResponseWriter) error {
	if err := json.NewDecoder(input).Decode(&pc); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}
	return nil
}

func NewKudaProcessor(rp *redis.Pool, middlewares func(pool *work.WorkerPool)) *KudaProcessor {
	processor := &KudaProcessor{
		redisPool:   rp,
		middlewares: middlewares,
	}
	return processor
}

func CreateKudaProcessor(kp *KudaProcessor, queue string, jobMap map[string]func(job *work.Job) error) (*work.WorkerPool, error) {
	q, c, err := ExtractQueue(queue)
	if err != nil {
		return nil, err
	}

	pool := work.NewWorkerPool(ProcessorContext{}, c, q, kp.redisPool)
	registerJobs(pool, jobMap)
	kp.middlewares(pool)

	return pool, nil
}

func RunProcessors(pools []*work.WorkerPool) {
	for _, wp := range pools {
		log.Println("start")
		wp.Start()
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	for _, wp := range pools {
		log.Println("stop")
		wp.Stop()
	}
}

func registerJobs(pool *work.WorkerPool, jobMap map[string]func(job *work.Job) error) {
	for jobName, fn := range jobMap {
		pool.Job(jobName, fn)
	}
}

// func customizeOptions(pool *work.WorkerPool) {
// 	pool.JobWithOptions("export", work.JobOptions{Priority: 10, MaxFails: 5}, (*ProcessorContext).Export)
// }

func (c *ProcessorContext) Export(job *work.Job) error {
	return nil
}
