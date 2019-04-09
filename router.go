package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Router struct {
	*mux.Router
}

func NewRouter(r Storage) *Router {
	h := &handler{r}
	router := mux.NewRouter()
	router.HandleFunc("/webpage/{id}", h.PutUrlHandler)
	router.UseEncodedPath().HandleFunc("/webpage/{id}", h.GetUrlHandler).Methods(http.MethodGet)
	router.UseEncodedPath().HandleFunc("/webpage/{id}", h.PutUrlHandler).Methods(http.MethodPut)
	router.UseEncodedPath().HandleFunc("/top", h.TopHandler).Methods(http.MethodGet)
	router.UseEncodedPath().HandleFunc("/webpage/{id}", h.DelUrlHandler).Methods(http.MethodDelete)
	return &Router{router}
}