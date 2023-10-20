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
	"todoList/page"
)

type TodoHandler struct {
	c controller.TodoController
}

func NewHandler(c controller.TodoController) http.Handler {
	h := TodoHandler{c: c}
	m := mux.NewRouter()
	m.HandleFunc("/todos", h.getTodos).Methods(http.MethodGet)
	m.HandleFunc("/todos/{id:[0-9]+}", h.getTodoDetail).Methods(http.MethodGet)
	m.HandleFunc("/todos/{id:[0-9]+}", h.deleteTodo).Methods(http.MethodDelete)
	m.HandleFunc("/todos", h.createTodo).Methods(http.MethodPost)
	m.HandleFunc("/todos", h.updateTodo).Methods(http.MethodPut)
	return m
}

// updateTodo godoc
// @Summary todo를 수정합니다.
// @Description 수정할내용을 요청합니다.
// @Tags todos
// @Param updateTodoDto body req.UpdateTodoDto true "updateTodo"
// @Accept  json
// @Produce  json
// @Success 200 {object} bool
// @Failure 400 "잘못된요청시"
// @Router /todos [put]
// @Security ApiKeyAuth
func (th TodoHandler) updateTodo(w http.ResponseWriter, r *http.Request) {
	t := &req.UpdateTodoDto{}
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = th.c.UpdateTodo(r.Context(), *t)
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
	w.Write([]byte(fmt.Sprintf("%t", true)))
}

// createTodo godoc
// @Summary  todo를 생성합니다.
// @Description 생성할 내용을 요청합니다.
// @Tags todos
// @Param create body req.CreateTodoDto true "createTodo"
// @Accept  json
// @Produce  json
// @Success 201 {object} bool
// @Failure 400 "잘못된요청시"
// @Failure 500 "알수없는 에러시"
// @Router /todos [post]
// @Security ApiKeyAuth
func (th TodoHandler) createTodo(w http.ResponseWriter, r *http.Request) {
	t := &req.CreateTodoDto{}
	err := json.NewDecoder(r.Body).Decode(t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = th.c.CreateTodo(r.Context(), *t)
	if err != nil {
		if strings.Contains(err.Error(), "is empty") || strings.Contains(err.Error(), "is not valid") {
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

// deleteTodo godoc
// @Summary  삭제합니다.
// @Description  삭제합니다
// @Tags todos
// @Param id path int true "todoID로 삭제합니다"
// @Accept  json
// @Produce  json
// @Success 200 {object} string
// @Failure 400 "토큰이 없을경우"
// @Failure 500 "알수없는 에러시"
// @Router /todos/{id} [delete]
// @Security ApiKeyAuth
func (th TodoHandler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = th.c.DeleteTodo(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(fmt.Sprintf("success delete id:%d", id)))
}

// getTodos godoc
// @Summary  todo를 상세조회합니다.
// @Description  todo를 상세 조회합니다
// @Tags todos
// @Param id path int true "todoID 로 검색합니다"
// @Accept  json
// @Produce  json
// @Success 200 {object} res.DetailDto
// @Failure 404 "없거나 해당 userID의 todo가 아닐경우"
// @Failure 500 "알수없는 에러시"
// @Router /todos/{id} [get]
// @Security ApiKeyAuth
func (th TodoHandler) getTodoDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	todo, err := th.c.GetDetail(r.Context(), id)
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

// getTodos godoc
// @Summary  todoList를 조회합니다.
// @Description list로 조회합니다.
// @Tags todos
// @Param page query int false "페이지입니다 0부터 시작입니다"
// @Param size query int false "한페이지 사이즈입니다 최솟값은10입니다"
// @Accept  json
// @Produce  json
// @Success 200 {array} res.ListDto
// @Failure 400 "잘못된요청시"
// @Failure 500 "알수없는 에러시"
// @Router /todos [get]
// @Security ApiKeyAuth
func (th TodoHandler) getTodos(w http.ResponseWriter, r *http.Request) {
	pageNum, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		pageNum = 0
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		pageSize = 0
	}
	reqPage := page.NewReqPage(pageNum, pageSize)

	todoDtos, count, err := th.c.GetTodos(r.Context(), reqPage)
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
