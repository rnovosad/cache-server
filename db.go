package main

import (
	"cassius/env"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"strconv"
	"strings"
	"time"
)

type RedisDB struct {
	*redis.Client
}

type CacheScore struct {
	Score float64
	Member string
}

func NewRedisDB(options redis.Options) *RedisDB {
	client := redis.NewClient(&options)
	return &RedisDB{client}
}

func (r *RedisDB) HasCache(k string) (bool, string) {
	cache, _ := r.Get(k).Result()
	return cache != "", cache
}

func (r *RedisDB) GetNumKeys() int64 {
	size := r.ZCard("hits").Val()
	return size
}

func (r *RedisDB) GetDbSize() int64 {
	const emptyDbSize = 1050496
	memory := make(map[string]string)
	size := r.Info("memory").Val()
	list1 := strings.Split(size, "\r\n")
	for _, val := range list1 {
		item := strings.Split(val, ":")
		if len(item) == 2 {
			memory[item[0]] = item[1]
		}
	}
	actSize, _ := strconv.ParseInt(memory["used_memory"], 10, 64)
	fmt.Println("Requested BD size", actSize - emptyDbSize)
	return actSize - emptyDbSize
}

func (r *RedisDB) GetAllCache() []byte {
	res, _ := r.ZRevRangeWithScores("hits", 0, -1).Result()
	fmt.Println(res)
	data, _ := json.Marshal(res)
	if err := json.Unmarshal(data, &[]CacheScore{}); err != nil {
		log.Println(err)
	}
	return data
}

func (r *RedisDB) SetCache(k string, v []byte) {
	d, _ := time.ParseDuration(env.GetConfig().Cache.Duration)
	err := r.Set(k, v, d).Err()
	if err != nil {
		log.Println("arghh")
	}
}

func (r *RedisDB) RemoveCache(k string) {
	if err := r.Del(k).Err(); err != nil {
		log.Println(err)
	}
}

func (r *RedisDB) SetLastAccess(k string) {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	if err := r.ZAdd( "last_access", redis.Z{Score: float64(timestamp), Member: k}).Err(); err != nil {
		log.Println(err)
	}
}

func (r *RedisDB) IncreaseHit(k string) {
	if err := r.ZIncr("hits", redis.Z{Score: 1, Member: k}).Err(); err != nil {
		log.Println(err)
	}
}

func (r *RedisDB) PopOldest() string {
	result, err := r.ZPopMin("last_access").Result()
	if err != nil {
		log.Println(err)
	}
	member := fmt.Sprintf("%v", result[0].Member)
	r.RemoveCache(member)
	r.ZRem("hits", member)
	return  member
}
