package main

import (
	"cassius/env"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"time"
)

type RedisDB struct {
	*redis.Client
}

type Storage interface {
	HasCache(key string) (bool, string)
	GetAllCache() []redis.Z
	SetCache(k string, v []byte)
	RemoveCache(k string)
	IncreaseHit(k string)
}

type Item struct {
	content []byte
	hits int
}

func NewRedisDB(options redis.Options) *RedisDB {
	client := redis.NewClient(&options)
	return &RedisDB{client}
}

func (r *RedisDB) HasCache(k string) (bool, string) {
	cache, _ := r.HGet(k, "content").Result()
	return cache != "", cache
}

func (r *RedisDB) GetAllCache() []redis.Z {
	res, _ := r.ZRevRangeWithScores("keySet", 0, 99).Result()
	fmt.Println(res)
	return res
}

func (r *RedisDB) SetCache(k string, v []byte) {
	value := Item{v, 1}
	d, _ := time.ParseDuration(env.GetConfig().Cache.Duration)
	err := r.HMSet(k, map[string]interface{}{
		"content": value.content,
		"hits": value.hits,
	}).Err()
	r.Expire(k, d)
	if err != nil {
		log.Println()
	}
}

func (r *RedisDB) RemoveCache(k string) {
	err := r.Del(k).Err()
	if err != nil {
		log.Println(err)
	}
}

func (r *RedisDB) IncreaseHit(k string) {
	if err := r.ZIncr("keySet", redis.Z{Score: 1, Member: k}); err != nil {
		log.Println(err)
	}
	err := r.HIncrBy(k, "hits", 1).Err()
	if err != nil {
		log.Println(err)
	}
}