package main

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"todoList/cmd/app/handler"
	"todoList/controller"
	handler2 "todoList/deco/handler"
	_ "todoList/docs"
	"todoList/repository"
	"todoList/repository/infla"
	"todoList/service"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name token
func main() {
	r := mux.NewRouter()
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	c := getController()
	todoHandler := handler.NewHandler(c)
	// middleware 래핑
	middlewareWrapped := handler2.NewDecoHandler(todoHandler, handler2.AuthMiddleware)
	// logger 래핑
	wrappedHandler := handler2.NewDecoHandler(middlewareWrapped, handler2.Logger)
	r.PathPrefix("/todos").Handler(wrappedHandler)
	http.ListenAndServe(":8082", r)
}

// 주입 방식변경
func getController() controller.TodoController {
	return controller.NewController(
		service.NewService(
			repository.NewRepository(
				infla.NewDB(),
			),
		),
	)
}
