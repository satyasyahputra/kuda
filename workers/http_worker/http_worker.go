package http_worker

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gocraft/work"
)

type HttpWorker struct {
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Body    string            `json:"body"`
	Headers map[string]string `json:"headers"`
}

func decode(s string) HttpWorker {
	hw := HttpWorker{}
	if err := json.NewDecoder(strings.NewReader(s)).Decode(&hw); err != nil {
		panic("Unable to decode")
	}
	return hw
}

func Alias() string {
	return "http_worker"
}

func Run(job *work.Job) error {
	log.Printf("run http worker: %v", job.Args)
	req := buildRequest(job)
	res := call(req)
	if res.StatusCode >= 400 {
		log.Printf("Error: %d", res.StatusCode)
	} else {
		log.Println("success")
	}
	return nil
}

func buildRequest(job *work.Job) *http.Request {
	req, err := http.NewRequest(job.ArgString("method"), job.ArgString("url"), strings.NewReader(job.ArgString("body")))
	if err != nil {
		panic("Failed request: " + job.ArgString("method") + " " + job.ArgString("url") + " " + job.ArgString("body"))
	}

	return req
}

func call(req *http.Request) *http.Response {
	client := &http.Client{
		CheckRedirect: http.DefaultClient.CheckRedirect,
	}

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	return res
}
