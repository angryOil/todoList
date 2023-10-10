package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"todoList/controller"
	"todoList/controller/req"
	"todoList/deco/handler"
	"todoList/page"
	"todoList/repository"
	"todoList/repository/infla"
	"todoList/service"
)

func NewHandler() http.Handler {
	m := mux.NewRouter()
	m.HandleFunc("/todos", getTodos).Methods(http.MethodGet)
	m.HandleFunc("/todos/{id:[0-9]+}", getTodoDetail).Methods(http.MethodGet)
	m.HandleFunc("/todos/{id:[0-9]+}", deleteTodo).Methods(http.MethodDelete)
	m.HandleFunc("/todos", createTodo).Methods(http.MethodPost)
	m.HandleFunc("/todos", updateTodo).Methods(http.MethodPut)
	// logger 래핑
	wrappedHandler := handler.NewDecoHandler(m, handler.Logger)
	return wrappedHandler
}

func getController() controller.TodoController {
	return controller.NewController(
		service.NewService(
			repository.NewRepository(
				infla.NewDB(),
			),
		),
	)
}

var c = getController()

func updateTodo(w http.ResponseWriter, r *http.Request) {
	t := &req.UpdateTodoDto{}
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = c.UpdateTodo(r.Context(), *t)
	if err != nil {
		if strings.Contains(err.Error(), "is empty") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte("success"))
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	t := &req.CreateTodoDto{}
	err := json.NewDecoder(r.Body).Decode(t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = c.CreateTodo(r.Context(), *t)
	if err != nil {
		if strings.Contains(err.Error(), "is empty") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("success"))
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = c.DeleteTodo(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(fmt.Sprintf("success delete id:%d", id)))
}

func getTodoDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	todo, err := c.GetDetail(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		fmt.Println("internal error is ", err.Error())
		return
	}

	data, err := json.Marshal(todo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(data))
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	pageNum, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		pageNum = 0
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		pageSize = 0
	}
	reqPage := page.NewReqPage(pageNum, pageSize)

	todoDtos, count, err := c.GetTodos(r.Context(), reqPage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	result := page.GetPagination(todoDtos, reqPage, count)
	data, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}
