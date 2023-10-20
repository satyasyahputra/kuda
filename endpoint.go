package kuda

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type KudaHttpContext struct {
	Host   string `env:"HOST" envDefault:"localhost"`
	Port   string `env:"PORT" envDefault:"8080"`
	router *mux.Router
	kec    *KudaEnqueuerContext
}

func (khc *KudaHttpContext) DefaultRoutes() *KudaHttpContext {
	khc.router.Path("/enqueue").Methods("GET").HandlerFunc(enqueue(khc.kec))
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
		kec.Enqueuers["my_queue"].Enqueue("my_worker", map[string]interface{}{})
		w.Write([]byte{})
	}
}
