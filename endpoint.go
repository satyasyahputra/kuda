package kuda

import (
	"net/http"

	"github.com/gorilla/mux"
)

type KudaHttpContext struct {
	router *mux.Router
	kec    *KudaEnqueuerContext
}

func NewKudaHttp(router *mux.Router, kec *KudaEnqueuerContext) *KudaHttpContext {
	router.Path("/enqueue").Methods("GET").HandlerFunc(enqueue(kec))
	return &KudaHttpContext{
		router: router,
		kec:    kec,
	}
}

func (khc *KudaHttpContext) StartHttp() {
	http.ListenAndServe("localhost:8080", khc.router)
}

func enqueue(kec *KudaEnqueuerContext) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		kec.Enqueuers["my_queue"].Enqueue("my_worker", map[string]interface{}{})
		w.Write([]byte{})
	}
}
