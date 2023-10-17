package example_worker

import (
	"log"

	"github.com/gocraft/work"
	"github.com/satyasyahputra/kuda"
)

func Run(c *kuda.ProcessorContext, job *work.Job) error {
	log.Println("common notify working")
	return nil
}
