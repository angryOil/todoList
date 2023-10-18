package main

import (
	"net/http"
	"todoList/cmd/app/handler"
	"todoList/controller"
	handler2 "todoList/deco/handler"
	"todoList/repository"
	"todoList/repository/infla"
	"todoList/service"
)

func main() {
	c := getController()
	todoHandler := handler.NewHandler(c)
	// middleware 래핑
	middlewareWrapped := handler2.NewDecoHandler(todoHandler, handler2.AuthMiddleware)
	// logger 래핑
	wrappedHandler := handler2.NewDecoHandler(middlewareWrapped, handler2.Logger)

	http.ListenAndServe(":8080", wrappedHandler)
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
