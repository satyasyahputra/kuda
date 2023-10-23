package my_worker

import (
	"log"

	"github.com/gocraft/work"
)

func Alias() string {
	return "my_worker"
}

func Run(job *work.Job) error {

	log.Printf("MyWorkerArgs: %v", job.Args)
	return nil
}
