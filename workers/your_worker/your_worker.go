package your_worker

import (
	"log"

	"github.com/gocraft/work"
)

func Alias() string {
	return "your_worker"
}

func Run(job *work.Job) error {
	log.Println("common notify working")
	return nil
}
