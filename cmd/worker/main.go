package main

import (
	"log"
	"strings"

	"github.com/caarlos0/env/v9"
	"github.com/gocraft/work"
	"github.com/satyasyahputra/kuda"
	"github.com/satyasyahputra/kuda/workers/http_worker"
	"github.com/satyasyahputra/kuda/workers/my_worker"
	"github.com/satyasyahputra/kuda/workers/your_worker"
)

func main() {
	config := loadEnv()
	redisPool := config.kr.NewRedisPool()
	queues := strings.Split(config.q.Queues, ",")

	jobMap := map[string]func(job *work.Job) error{
		my_worker.Alias():   my_worker.Run,
		http_worker.Alias(): http_worker.Run,
		your_worker.Alias(): your_worker.Run,
	}

	kudaProcessor := kuda.NewKudaProcessor(redisPool, kuda.ProcessorMiddleware)
	processors := []*work.WorkerPool{}

	for _, queue := range queues {
		pool, err := kuda.CreateKudaProcessor(kudaProcessor, queue, jobMap)
		if err != nil {
			log.Panicf("error occurred when create processor: %v", err)
			return
		}
		processors = append(processors, pool)
	}

	kuda.RunProcessors(processors)
}

type config struct {
	kr *kuda.KudaRedis
	q  *kuda.KudaQueue
}

func loadEnv() *config {
	kr := kuda.KudaRedis{}
	if err := env.ParseWithOptions(&kr, env.Options{Prefix: "KUDA_REDIS_"}); err != nil {
		panic(err)
	}
	q := kuda.KudaQueue{}
	if err := env.ParseWithOptions(&q, env.Options{Prefix: "KUDA_"}); err != nil {
		panic(err)
	}
	return &config{
		kr: &kr,
		q:  &q,
	}
}
