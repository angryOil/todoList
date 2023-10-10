package handler

import (
	"log"
	"net/http"
	"time"
)

func Logger(w http.ResponseWriter, r *http.Request, h http.Handler) {
	start := time.Now()
	log.Printf("[logger] %s start", r.URL)
	h.ServeHTTP(w, r)
	log.Printf("[logger] %s elapsed time:%s\n", r.URL, time.Now().Sub(start))
}
