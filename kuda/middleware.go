package kuda

import (
	"fmt"

	"github.com/gocraft/work"
)

func ProcessorMiddleware(pool *work.WorkerPool) error {
	pool.Middleware((*ProcessorContext).log)
	return nil
}

func (c *ProcessorContext) log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Println("Starting job: ", job.Name)
	return next()
}
