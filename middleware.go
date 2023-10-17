package kuda

import (
	"fmt"

	"github.com/gocraft/work"
)

func ProcessorMiddleware(pool *work.WorkerPool) {
	pool.Middleware((*ProcessorContext).log)
}

func (c *ProcessorContext) log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Println("Starting job: ", job.Name)
	return next()
}
