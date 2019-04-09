package main
//
//import (
//	"fmt"
//	"github.com/gorilla/mux"
//	"log"
//	"net/http"
//	"net/http/httptest"
//	"net/url"
//)
//
//func Cached(red *RedisDB, handler func(w http.ResponseWriter, r *http.Request, url string, redis *RedisDB)) http.HandlerFunc {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		rawUrl := mux.Vars(r)["id"]
//		urlKey, _ := url.QueryUnescape(rawUrl)
//		if red.HasCache(urlKey) {
//			red.IncreaseHit(urlKey)
//			fmt.Println("Cache Hit!")
//			content := red.GetContent(urlKey)
//			if _, err := w.Write([]byte(content)); err != nil {
//				log.Fatal(err)
//			}
//		} else {
//			c := httptest.NewRecorder()
//			handler(c, r, urlKey, red)
//			fmt.Print("Cache not hit!\n")
//			content := c.Body.Bytes()
//			if _, err := w.Write(content); err != nil {
//				log.Fatal(err)
//			}
//		}
//	})
//}