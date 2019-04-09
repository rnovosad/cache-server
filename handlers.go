package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"
)

type handler struct {
	cache Storage
}

func (h *handler) GetUrlHandler(w http.ResponseWriter, r *http.Request) {
	rawUrl := mux.Vars(r)["id"]
	urlKey, err := url.QueryUnescape(rawUrl)
	if err != nil {
		log.Println("Can't parse URL", err)
	}
	hasCache, content := h.cache.HasCache(urlKey)
	if hasCache {
		if _, err:= w.Write([]byte(content)); err != nil {
			log.Println("can't get URL", err)
		}
	} else {
		content := fmt.Sprintf("There is no cache for page %v\n", urlKey)
		if _, err:= w.Write([]byte(content)); err != nil {
			log.Println("cant write get URL", err)
		}
	}
}

func (h *handler) PutUrlHandler(w http.ResponseWriter, r *http.Request) {
	rawUrl := mux.Vars(r)["id"]
	urlKey, err := url.QueryUnescape(rawUrl)
	if err != nil {
		log.Println("Can't parse URL", err)
	}
	fmt.Println(urlKey)
	hasCache, content := h.cache.HasCache(urlKey)
	if hasCache {
		if _, err:= w.Write([]byte(content)); err != nil {
			log.Println("can't get URL", err)
		}
	} else {
		resp, err := http.Get(urlKey)
		if err != nil {
			log.Println("Can't get URL", err)
		}
		buf := new(bytes.Buffer)
		size, err := buf.ReadFrom(resp.Body)
		fmt.Println("size size", size)
		if err != nil {
			log.Println("Can't read response", err)
		}
		h.cache.SetCache(urlKey, buf.Bytes())
		if _, err:= w.Write(buf.Bytes()); err != nil {
			log.Println("Can't write response", err)
		}
	}
}

func (h *handler) TopHandler(w http.ResponseWriter, r *http.Request) {
	content := h.cache.GetAllCache()
	data, _ := json.Marshal(content)
	if _, err:= w.Write(data); err != nil {
		log.Println("top", err)
	}
}

func (h *handler) DelUrlHandler(w http.ResponseWriter, r *http.Request) {
	rawUrl := mux.Vars(r)["id"]
	urlKey, _ := url.QueryUnescape(rawUrl)
	if hasCache, _ := h.cache.HasCache(urlKey); hasCache {
		h.cache.RemoveCache(urlKey)
	} else {
		content := fmt.Sprintf("There is no cache for page %v\n", urlKey)
		if _, err:= w.Write([]byte(content)); err != nil {
			log.Println("can't del URL", err)
		}
	}
}