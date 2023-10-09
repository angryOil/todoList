package main

import (
	"net/http"
	"todoList/cmd/app/handler"
)

func main() {
	todoHandler := handler.NewHandler()
	http.ListenAndServe(":8080", todoHandler)
}
