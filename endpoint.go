package kuda

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gocraft/work"
	"github.com/gorilla/mux"
)

type KudaHttpContext struct {
	Host   string `env:"HOST" envDefault:"localhost"`
	Port   string `env:"PORT" envDefault:"8080"`
	router *mux.Router
	kec    *KudaEnqueuerContext
}

func (khc *KudaHttpContext) DefaultRoutes() *KudaHttpContext {
	khc.router.Path("/enqueue").Methods("POST").HandlerFunc(enqueue(khc.kec))
	return khc
}

func (khc *KudaHttpContext) Router(router *mux.Router) *KudaHttpContext {
	khc.router = router
	return khc
}

func (khc *KudaHttpContext) Enqueuer(kec *KudaEnqueuerContext) *KudaHttpContext {
	khc.kec = kec
	return khc
}

func (khc *KudaHttpContext) StartHttp() {
	addr := khc.Host + ":" + khc.Port
	log.Printf("starting http server on `%s`", addr)
	http.ListenAndServe(addr, khc.router)
}

func enqueue(kec *KudaEnqueuerContext) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pc := ProcessorContext{}
		err := DecodeJSON(&pc, r.Body, w)
		if err != nil {
			return
		}

		args := work.Q{}
		err = json.NewDecoder(strings.NewReader(pc.Args)).Decode(&args)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		job, err := kec.Enqueuers[pc.Queue].Enqueue(pc.Worker, args)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "job": job})
	}
}
