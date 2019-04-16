package main

import (
	"cassius/env"
	"github.com/go-redis/redis"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUrlHandler(t *testing.T) {
	h := &Handler{cache: NewRedisDB(redis.Options{}), config: env.GetConfig()}
	req, err := http.NewRequest("GET", "/webpage/https%3A%2F%2Fmeyerweb.com%2Feric%2Ftools%2Fdencoder%2F1 ", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetUrlHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "There is no cache for page \n"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
