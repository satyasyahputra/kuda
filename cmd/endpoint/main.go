package main

import (
	"strings"

	"github.com/caarlos0/env/v9"
	"github.com/gorilla/mux"
	"github.com/satyasyahputra/kuda"
)

func main() {
	appConfig := loadEnv()
	redisPool := appConfig.kr.NewRedisPool()
	queues := strings.Split(appConfig.kq.Queues, ",")
	r := mux.NewRouter()
	kec, _ := kuda.NewKudaEnqueuerContext(redisPool, queues)
	khc := &appConfig.khc

	khc.Router(r).
		Enqueuer(kec).
		DefaultRoutes().
		StartHttp()
}

type appConfig struct {
	kr  kuda.KudaRedis
	khc kuda.KudaHttpContext
	kq  kuda.KudaQueue
}

func loadEnv() appConfig {
	kr := kuda.KudaRedis{}
	khc := kuda.KudaHttpContext{}
	kq := kuda.KudaQueue{}

	if err := env.ParseWithOptions(&kr, env.Options{Prefix: "KUDA_REDIS_"}); err != nil {
		panic(err)
	}
	if err := env.ParseWithOptions(&khc, env.Options{Prefix: "KUDA_API_"}); err != nil {
		panic(err)
	}
	if err := env.ParseWithOptions(&kq, env.Options{Prefix: "KUDA_"}); err != nil {
		panic(err)
	}
	return appConfig{
		kr:  kr,
		khc: khc,
		kq:  kq,
	}
}
