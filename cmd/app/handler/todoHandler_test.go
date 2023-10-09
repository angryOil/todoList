package handler

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"todoList/controller/req"
	"todoList/controller/res"
)

func TestCreateTodo(t *testing.T) {
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	reqTodo := req.CreateTodoDto{Title: "제목", Content: "이것은?", OrderNum: 1}
	data, _ := json.Marshal(reqTodo)

	req, err := http.NewRequest("POST", ts.URL+"/todos", strings.NewReader(string(data)))
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	//resp, err := http.Get(ts.URL+"/todos")
}

func TestUpdateTodo(t *testing.T) {
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()
	reqTodo := req.UpdateTodoDto{Id: 3, Title: "title3", Content: "con2", OrderNum: 11, IsDeleted: false}
	data, err := json.Marshal(reqTodo)
	assert.NoError(t, err)
	newRequest, err := http.NewRequest("PUT", ts.URL+"/todos", strings.NewReader(string(data)))
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(newRequest)
	defer resp.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 업데이트후 확인 upsert
	newRequest, err = http.NewRequest("GET", ts.URL+"/todos/"+strconv.Itoa(reqTodo.Id), nil)
	assert.NoError(t, err)

	resp, err = http.DefaultClient.Do(newRequest)
	defer resp.Body.Close()
	assert.NoError(t, err)

	resTodo := &res.DetailDto{}
	err = json.NewDecoder(resp.Body).Decode(resTodo)

	assert.Equal(t, reqTodo.Id, resTodo.Id)
	assert.Equal(t, reqTodo.Title, resTodo.Title)
	assert.Equal(t, reqTodo.Content, resTodo.Content)
	assert.Equal(t, reqTodo.OrderNum, resTodo.OrderNum)
	assert.Equal(t, reqTodo.IsDeleted, resTodo.IsDeleted)

	// 수정후 확인 update
	reqTodo = req.UpdateTodoDto{Id: 3, Title: "mod title", Content: "modi content!", OrderNum: 21, IsDeleted: true}
	data, err = json.Marshal(reqTodo)
	assert.NoError(t, err)

	newRequest, err = http.NewRequest("PUT", ts.URL+"/todos", strings.NewReader(string(data)))
	assert.NoError(t, err)

	resp, err = http.DefaultClient.Do(newRequest)
	resp.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	newRequest, err = http.NewRequest("GET", ts.URL+"/todos/"+strconv.Itoa(reqTodo.Id), nil)
	assert.NoError(t, err)

	resp, err = http.DefaultClient.Do(newRequest)
	defer resp.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resTodo = &res.DetailDto{}
	err = json.NewDecoder(resp.Body).Decode(resTodo)
	assert.NoError(t, err)
	assert.Equal(t, reqTodo.Id, resTodo.Id)
	assert.Equal(t, reqTodo.Title, resTodo.Title)
	assert.Equal(t, reqTodo.Content, resTodo.Content)
	assert.Equal(t, reqTodo.OrderNum, resTodo.OrderNum)
	assert.Equal(t, reqTodo.IsDeleted, resTodo.IsDeleted)

}

func TestDeleteTodo(t *testing.T) {
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	// 생성
	reqTodo := req.UpdateTodoDto{Id: 3, Title: "delete test", Content: "this is delete test", OrderNum: 21, IsDeleted: true}
	data, err := json.Marshal(reqTodo)
	assert.NoError(t, err)

	newRequest, err := http.NewRequest("PUT", ts.URL+"/todos", strings.NewReader(string(data)))
	assert.NoError(t, err)

	resultTodo := &res.DetailDto{}

	resp, err := http.DefaultClient.Do(newRequest)
	defer resp.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 생성확인
	newRequest, err = http.NewRequest("GET", ts.URL+"/todos/"+strconv.Itoa(reqTodo.Id), nil)
	assert.NoError(t, err)
	resp, err = http.DefaultClient.Do(newRequest)
	defer resp.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(resultTodo)
	assert.NoError(t, err)

	assert.Equal(t, reqTodo.Id, resultTodo.Id)
	assert.Equal(t, reqTodo.Title, resultTodo.Title)

	// 삭제

	newRequest, err = http.NewRequest("DELETE", ts.URL+"/todos/"+strconv.Itoa(reqTodo.Id), nil)
	assert.NoError(t, err)

	resp, err = http.DefaultClient.Do(newRequest)
	defer resp.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 삭제후 데이터 확인
	resultTodo = &res.DetailDto{}
	newRequest, err = http.NewRequest("GET", ts.URL+"/todos/"+strconv.Itoa(reqTodo.Id), nil)
	assert.NoError(t, err)

	resp, err = http.DefaultClient.Do(newRequest)
	defer resp.Body.Close()
	assert.NoError(t, err)

	err = json.NewDecoder(resp.Body).Decode(resultTodo)
	assert.Error(t, err)
	assert.ErrorAs(t, io.EOF, &err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

}
