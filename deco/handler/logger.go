package handler

import (
	"log"
	"net/http"
	"time"
)

func Logger(w http.ResponseWriter, r *http.Request, h http.Handler) {
	start := time.Now()
	log.Printf("[%s][logger] %s start", r.Method, r.URL)
	h.ServeHTTP(w, r)
	log.Printf("[%s][logger] %s elapsed time:%s\n", r.Method, r.URL, time.Now().Sub(start))
}
