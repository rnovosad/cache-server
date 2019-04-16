package main

import (
	"bytes"
	"cassius/env"
	"fmt"
	"github.com/docker/go-units"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"
)

type storage interface {
	HasCache(key string) (bool, string)
	GetAllCache() []byte
	SetCache(k string, v []byte)
	RemoveCache(k string)
	IncreaseHit(k string)
	GetDbSize() int64
	GetNumKeys() int64
	SetLastAccess(k string)
	PopOldest() string
}

type Handler struct {
	cache storage
	config env.Configuration
}

func NewHandler(cfg env.Configuration, storage storage) *Handler {
	return &Handler{
		cache:  storage,
		config: cfg,
	}
}

func (h *Handler) GetUrlHandler(w http.ResponseWriter, r *http.Request) {
	rawUrl := mux.Vars(r)["id"]
	urlKey, err := url.QueryUnescape(rawUrl)
	if err != nil {
		log.Println("Can't parse URL", err)
	}
	hasCache, content := h.cache.HasCache(urlKey)
	if hasCache {
		h.cache.IncreaseHit(urlKey)
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

func (h *Handler) PutUrlHandler(w http.ResponseWriter, r *http.Request) {
	maxContentSize, _ := units.FromHumanSize(h.config.Cache.MaxItemSize)
	maxDbSize, _ := units.FromHumanSize(h.config.Cache.MaxDbSize)
	maxKeysCount := h.config.Cache.MaxItemsCount
	dbSize := h.cache.GetDbSize()
	keysCount := h.cache.GetNumKeys()
	rawUrl := mux.Vars(r)["id"]

	urlKey, err := url.QueryUnescape(rawUrl)

	if err != nil {
		log.Println("Can't parse URL", err)
	}

	resp, err := http.Get(urlKey)

	if err != nil {
		log.Println("Can't get URL", err)
	}

	buf := new(bytes.Buffer)
	bodySize, err := buf.ReadFrom(resp.Body)

	if err != nil {
		log.Println("Can't read response", err)
	}

	if bodySize > maxContentSize {
		content := fmt.Sprintf("Cache item bodySize is too big, %v", bodySize)
		if _, err:= w.Write([]byte(content)); err != nil {
			log.Println("cant write URL", err)
		}
		return
	}

	hasCache, _ := h.cache.HasCache(urlKey)
	for dbSize + bodySize > maxDbSize || (!hasCache && h.cache.GetNumKeys() >= maxKeysCount) {
		h.cache.PopOldest()
	}
	h.cache.SetCache(urlKey, buf.Bytes())
	h.cache.IncreaseHit(urlKey)
	h.cache.SetLastAccess(urlKey)
	log.Println("keys: ", keysCount)
	if _, err:= w.Write(buf.Bytes()); err != nil {
		log.Println("Can't write response", err)
	}
}

func (h *Handler) TopHandler(w http.ResponseWriter, r *http.Request) {
	content := h.cache.GetAllCache()
	if _, err:= w.Write(content); err != nil {
		log.Println("top", err)
	}
}

func (h *Handler) DelUrlHandler(w http.ResponseWriter, r *http.Request) {
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