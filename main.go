package main

import (
	"cassius/env"
	"github.com/go-redis/redis"
	"log"
	"net/http"
)

// our main function

func main() {
	config := env.GetConfig()
	h := NewHandler(config, NewRedisDB(redis.Options{}))
	router := NewRouter(h)

	log.Printf("cache ttl %v", config.Cache.Duration)
	log.Printf("max items count %v", config.Cache.MaxItemsCount)
	log.Printf("max cache size %v", config.Cache.MaxDbSize)
	log.Printf("max item size %v", config.Cache.MaxItemSize)
	log.Printf("port for this application is %d", config.Server.Port)

	log.Fatal(http.ListenAndServe(":8000", router))
}