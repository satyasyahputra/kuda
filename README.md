# Kuda
Hard work like `Kuda`!

This project is simply enable multiple worker by templating of gocraft/work library.

# Example
```
package main

import (
	"log"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/kuda/kuda"
)

var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", "localhost:6379")
	},
}

func main() {

	enqeueuer, _ := kuda.RegisterEnqueuer("common", redisPool)
	enqeueuer.Enqueue("email_notify", nil)

	jobMap := map[string]interface{}{
		"email_notify": EmailNotify,
	}

	simpleKudaProcessor := kuda.NewKudaProcessor("common", 10, redisPool, jobMap, nil)
	simple_pool := kuda.SimpleProcessor(simpleKudaProcessor)

	go kuda.RunProcessor(simple_pool)

	customEnqeueuer, _ := kuda.RegisterEnqueuer("custom", redisPool)
	customEnqeueuer.Enqueue("custom_notify", nil)

	customJobMap := map[string]interface{}{
		"custom_notify": CustomNotify,
	}

	customKudaProcessor := kuda.NewKudaProcessor("custom", 10, redisPool, customJobMap, Middlewares)
	custom_pool := kuda.CreateKudaProcessor(customKudaProcessor)

	kuda.RunProcessor(custom_pool)
}

func EmailNotify(c *kuda.ProcessorContext, job *work.Job) error {
	log.Println("send email")
	return nil
}

func CustomNotify(c *kuda.ProcessorContext, job *work.Job) error {
	log.Println("custom notify")
	return nil
}

func Middlewares(pool *work.WorkerPool) {
	pool.Middleware(customLog)
}

func customLog(c *kuda.ProcessorContext, job *work.Job, next work.NextMiddlewareFunc) error {
	log.Println("custom log!")
	return next()
}

```