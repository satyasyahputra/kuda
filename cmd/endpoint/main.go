package main

import (
	"github.com/caarlos0/env/v9"
	"github.com/gorilla/mux"
	"github.com/satyasyahputra/kuda"
)

func main() {
	kr := loadEnv()
	redisPool := kr.NewRedisPool()
	queues := []string{"my_queue"}
	r := mux.NewRouter()
	kec := kuda.NewKudaEnqueuerContext(redisPool, queues)
	khc := kuda.NewKudaHttp(r, kec)
	khc.StartHttp()
}

func loadEnv() kuda.KudaRedis {
	kr := kuda.KudaRedis{}
	if err := env.ParseWithOptions(&kr, env.Options{Prefix: "KUDA_REDIS_"}); err != nil {
		panic(err)
	}
	return kr
}
