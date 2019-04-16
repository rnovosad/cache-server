package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type handler interface {
	GetUrlHandler(w http.ResponseWriter, r *http.Request)
	PutUrlHandler(w http.ResponseWriter, r *http.Request)
	TopHandler(w http.ResponseWriter, r *http.Request)
	DelUrlHandler(w http.ResponseWriter, r *http.Request)
}

type Router struct {
	*mux.Router
	handler handler
}

func NewRouter(h handler) *Router {
	router := mux.NewRouter()
	router.UseEncodedPath().HandleFunc("/webpage/{id}", h.GetUrlHandler).Methods(http.MethodGet)
	router.UseEncodedPath().HandleFunc("/webpage/{id}", h.PutUrlHandler).Methods(http.MethodPut)
	router.UseEncodedPath().HandleFunc("/top", h.TopHandler).Methods(http.MethodGet)
	router.UseEncodedPath().HandleFunc("/webpage/{id}", h.DelUrlHandler).Methods(http.MethodDelete)
	return &Router{
		Router:  router,
		handler: h,
	}
}
