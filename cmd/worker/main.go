package main

import (
	"github.com/caarlos0/env/v9"
	"github.com/satyasyahputra/kuda"
	"github.com/satyasyahputra/kuda/workers/example_worker"
)

func main() {
	kr := loadEnv()
	redisPool := kr.NewRedisPool()
	queue := "example"

	enqeueuer, _ := kuda.RegisterEnqueuer(queue, redisPool)
	enqeueuer.Enqueue("example_worker", nil)

	jobMap := map[string]interface{}{
		"example_worker": example_worker.Run,
	}

	customKudaProcessor := kuda.NewKudaProcessor(queue, 10, redisPool, jobMap, kuda.ProcessorMiddleware)
	custom_pool := kuda.CreateKudaProcessor(customKudaProcessor)

	kuda.RunProcessor(custom_pool)
}

func loadEnv() kuda.KudaRedis {
	kr := kuda.KudaRedis{}
	if err := env.ParseWithOptions(&kr, env.Options{Prefix: "KUDA_REDIS_"}); err != nil {
		panic(err)
	}
	return kr
}
