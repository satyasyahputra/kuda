package main

import (
	"github.com/caarlos0/env/v9"
	"github.com/gocraft/work"
	"github.com/satyasyahputra/kuda"
	"github.com/satyasyahputra/kuda/workers/my_worker"
)

func main() {
	kr := loadEnv()
	redisPool := kr.NewRedisPool()
	queues := []string{"my_queue"}

	jobMap := map[string]func(c *kuda.ProcessorContext, job *work.Job) error{
		"my_worker": my_worker.Run,
	}

	kudaProcessor := kuda.NewKudaProcessor(10, redisPool, kuda.ProcessorMiddleware)
	processors := []*work.WorkerPool{}

	for _, queue := range queues {
		custom_pool := kuda.CreateKudaProcessor(kudaProcessor, queue, jobMap)
		processors = append(processors, custom_pool)
	}

	kuda.RunProcessors(processors)
}

func loadEnv() kuda.KudaRedis {
	kr := kuda.KudaRedis{}
	if err := env.ParseWithOptions(&kr, env.Options{Prefix: "KUDA_REDIS_"}); err != nil {
		panic(err)
	}
	return kr
}
