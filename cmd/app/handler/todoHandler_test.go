package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"todoList/controller"
	"todoList/controller/req"
	"todoList/controller/res"
	"todoList/deco/handler"
	"todoList/domain"
	"todoList/page"
	"todoList/repository"
	"todoList/repository/infla"
	"todoList/service"
)

type TodoHandlerTestSuite struct {
	suite.Suite
	handler  http.Handler
	rollback func() error
}

func TestTodoHandlerTestSuite(t *testing.T) {
	suite.Run(t, &TodoHandlerTestSuite{})
}

// queryString 으로 userId를 받음
func testMiddleware(w http.ResponseWriter, r *http.Request, h http.Handler) {
	ctx := context.WithValue(r.Context(), "userId", baseTestDomain.UserId)
	h.ServeHTTP(w, r.WithContext(ctx))
}

var testRepo repository.ITodoRepository

var c controller.TodoController

// 각 테스트실행전 1번씩 실행됨 전체 테스트 실해전 실행은 SetupSuite() 메소드 사용
func (s *TodoHandlerTestSuite) SetupTest() {
	log.Printf("setup testSuite ...")
	var db = infla.NewDB()
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		log.Panicf("tx setup fail err: %e", err)
	}
	s.rollback = tx.Rollback
	if err != nil {
		log.Panicf("location setup fail err: %e", err)
	}
	testRepo = repository.NewRepository(tx)
	h := NewHandler(controller.NewController(service.NewService(testRepo)))
	h = handler.NewDecoHandler(h, testMiddleware)
	s.handler = h
}

// userId 가9999 인 테스트 데이터를 10개 가지고 시작할거임

// 테스트시 기준이 될 데이터
var baseTestDomain = domain.Todo{
	UserId:    9999,
	Title:     "base test model title",
	Content:   "base test model content",
	OrderNum:  1,
	IsDeleted: false,
}

var domainArr []domain.Todo

func getInitDomainArr() []domain.Todo {
	base := baseTestDomain
	result := make([]domain.Todo, 10)
	for i, _ := range result {
		result[i] = domain.Todo{
			Id:        i + 1, //id 가 0이 되면안되므로
			UserId:    base.UserId,
			Title:     base.Title + strconv.Itoa(i),
			Content:   base.Content + strconv.Itoa(i),
			OrderNum:  base.OrderNum + i,
			IsDeleted: base.IsDeleted,
		}
	}
	return result
}

// test 에선 멱등성을 위해 save(update) 메소드를 사용할거고
// update 에서는 getValidFunc 에서 나온 결과로(repo 에서 조회후 넣을 예정)
func validFunMaker(todoId int) func([]domain.Todo) (domain.Todo, error) {
	return func(todos []domain.Todo) (domain.Todo, error) {
		return domain.Todo{
			UserId:    baseTestDomain.UserId,
			Id:        todoId + 1,
			Title:     baseTestDomain.Title + strconv.Itoa(todoId),
			Content:   baseTestDomain.Content + strconv.Itoa(todoId),
			OrderNum:  baseTestDomain.OrderNum + todoId,
			IsDeleted: baseTestDomain.IsDeleted,
		}, nil
	}
}

// 병합된 데이터 테스트 이므로 그대로
func getTestMergeTodo(todo domain.Todo) domain.Todo {
	return todo
}

// todoId,userId 는 0이어서는 안됨
func getTestSaveValid(todo domain.Todo) error {
	if todo.UserId == 0 {
		return errors.New("userId is zero")
	}
	if todo.Id == 0 {
		return errors.New("todoId is zero")
	}
	return nil
}

var (
	Repo = "repo"
)

var RepoTestFun = func() {
	log.Println("test")
}

// 각테스트전에 실행
func (s *TodoHandlerTestSuite) BeforeTest(suiteName, testName string) {
	log.Printf("이것은 BeforeTest 입니다. %s %s", suiteName, testName)
	getTestValidFunc := validFunMaker
	for _, todoDomain := range getInitDomainArr() {
		err := testRepo.Save(context.Background(), todoDomain.UserId, todoDomain.Id, getTestValidFunc(todoDomain.Id), getTestMergeTodo, getTestSaveValid)
		if err != nil {
			panic(err)
		}
	}
}

// 물론 BeforeTest 도 존재함 (suiteName,testName string)을 인자로 받아서 실행
func (s *TodoHandlerTestSuite) AfterTest(suiteName, testName string) {
	log.Printf("roll back / suiteName: %s, testName: %s", suiteName, testName)
	//err := s.commit()
	err := s.rollback()
	if err != nil {
		panic(err)
	}
}

func (s *TodoHandlerTestSuite) TestCreateTodo() {
	s.Run("todo 성공", func() {
		ts := httptest.NewServer(s.handler)
		defer ts.Close()
		reqTodo := req.CreateTodoDto{Title: "제목", Content: "이것은?", OrderNum: 1}
		data, _ := json.Marshal(reqTodo)

		re, err := http.NewRequest("POST", ts.URL+"/todos", strings.NewReader(string(data)))
		assert.NoError(s.T(), err)

		resp, err := http.DefaultClient.Do(re)
		defer resp.Body.Close()
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), http.StatusCreated, resp.StatusCode)
	})
	s.Run("잘못된 json 일경우", func() {
		ts := httptest.NewServer(s.handler)
		defer ts.Close()
		re, err := http.NewRequest("POST", ts.URL+"/todos", strings.NewReader("잘못된json 입니다"))
		assert.NoError(s.T(), err)

		resp, err := http.DefaultClient.Do(re)
		defer resp.Body.Close()
		readBody, err := io.ReadAll(resp.Body)
		assert.NoError(s.T(), err)
		assert.Contains(s.T(), string(readBody), "invalid")
		assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
	})
	s.Run("todo 제목이 없을때", func() {
		ts := httptest.NewServer(s.handler)
		defer ts.Close()
		// 제목이 없을경우
		reqTodo := req.CreateTodoDto{Content: "제목이 없네요??", OrderNum: 1}
		data, _ := json.Marshal(reqTodo)

		re, err := http.NewRequest("POST", ts.URL+"/todos", strings.NewReader(string(data)))
		assert.NoError(s.T(), err)

		resp, err := http.DefaultClient.Do(re)
		defer resp.Body.Close()
		readBody, err := io.ReadAll(resp.Body)
		assert.NoError(s.T(), err)
		assert.Contains(s.T(), string(readBody), "title is empty")
		assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
	})

}

func (s *TodoHandlerTestSuite) TestUpdateTodo() {
	s.Run("업데이트 테스트", func() {
		ts := httptest.NewServer(s.handler)
		defer ts.Close()
		reqTodo := req.UpdateTodoDto{Id: 3, Title: "title3", Content: "con2", OrderNum: 11, IsDeleted: false}
		data, err := json.Marshal(reqTodo)
		assert.NoError(s.T(), err)
		newRequest, err := http.NewRequest("PUT", ts.URL+"/todos", strings.NewReader(string(data)))
		assert.NoError(s.T(), err)

		resp, err := http.DefaultClient.Do(newRequest)
		defer resp.Body.Close()
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

		// 업데이트후 확인 upsert
		newRequest, err = http.NewRequest("GET", ts.URL+"/todos/"+strconv.Itoa(reqTodo.Id), nil)
		assert.NoError(s.T(), err)

		resp, err = http.DefaultClient.Do(newRequest)
		defer resp.Body.Close()
		assert.NoError(s.T(), err)

		resTodo := &res.DetailDto{}
		err = json.NewDecoder(resp.Body).Decode(resTodo)

		assert.Equal(s.T(), reqTodo.Id, resTodo.Id)
		assert.Equal(s.T(), reqTodo.Title, resTodo.Title)
		assert.Equal(s.T(), reqTodo.Content, resTodo.Content)
		assert.Equal(s.T(), reqTodo.OrderNum, resTodo.OrderNum)
		assert.Equal(s.T(), reqTodo.IsDeleted, resTodo.IsDeleted)

		// 수정후 확인 update
		reqTodo = req.UpdateTodoDto{Id: 3, Title: "mod title", Content: "modi content!", OrderNum: 21, IsDeleted: true}
		data, err = json.Marshal(reqTodo)
		assert.NoError(s.T(), err)

		newRequest, err = http.NewRequest("PUT", ts.URL+"/todos", strings.NewReader(string(data)))
		assert.NoError(s.T(), err)

		resp, err = http.DefaultClient.Do(newRequest)
		defer resp.Body.Close()
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

		newRequest, err = http.NewRequest("GET", ts.URL+"/todos/"+strconv.Itoa(reqTodo.Id), nil)
		assert.NoError(s.T(), err)

		resp, err = http.DefaultClient.Do(newRequest)
		defer resp.Body.Close()
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

		resTodo = &res.DetailDto{}
		err = json.NewDecoder(resp.Body).Decode(resTodo)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), reqTodo.Id, resTodo.Id)
		assert.Equal(s.T(), reqTodo.Title, resTodo.Title)
		assert.Equal(s.T(), reqTodo.Content, resTodo.Content)
		assert.Equal(s.T(), reqTodo.OrderNum, resTodo.OrderNum)
		assert.Equal(s.T(), reqTodo.IsDeleted, resTodo.IsDeleted)
	})
	s.Run("잘못된 json 값일때", func() {
		ts := httptest.NewServer(s.handler)
		defer ts.Close()
		newRequest, err := http.NewRequest("PUT", ts.URL+"/todos", strings.NewReader("string(data)"))
		assert.NoError(s.T(), err)

		resp, err := http.DefaultClient.Do(newRequest)
		defer resp.Body.Close()
		readBody, _ := io.ReadAll(resp.Body)
		assert.NoError(s.T(), err)
		assert.Contains(s.T(), string(readBody), "invalid")
		assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
	})
	s.Run("업데이트 제목이 없을경우 error를 반환 ", func() {
		ts := httptest.NewServer(s.handler)
		defer ts.Close()
		reqTodo := req.UpdateTodoDto{Id: 3, Title: "", Content: "con2", OrderNum: 11, IsDeleted: false}
		data, err := json.Marshal(reqTodo)
		assert.NoError(s.T(), err)
		newRequest, err := http.NewRequest("PUT", ts.URL+"/todos", strings.NewReader(string(data)))
		assert.NoError(s.T(), err)

		resp, err := http.DefaultClient.Do(newRequest)
		defer resp.Body.Close()
		readBody, _ := io.ReadAll(resp.Body)
		assert.NoError(s.T(), err)
		assert.Contains(s.T(), string(readBody), "is empty")
		assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
	})
	s.Run("업데이트 id가 0일경우 error를 반환 ", func() {
		ts := httptest.NewServer(s.handler)
		defer ts.Close()
		reqTodo := req.UpdateTodoDto{Id: 0, Title: "ㅁㅇㄹㄴ", Content: "con2", OrderNum: 11, IsDeleted: false}
		data, err := json.Marshal(reqTodo)
		assert.NoError(s.T(), err)
		newRequest, err := http.NewRequest("PUT", ts.URL+"/todos", strings.NewReader(string(data)))
		assert.NoError(s.T(), err)

		resp, err := http.DefaultClient.Do(newRequest)
		defer resp.Body.Close()
		readBody, _ := io.ReadAll(resp.Body)
		assert.NoError(s.T(), err)
		assert.Contains(s.T(), string(readBody), "is zero")
		assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
	})
}

func (s *TodoHandlerTestSuite) TestDeleteTodo() {
	s.Run("삭제 테스트", func() {
		ts := httptest.NewServer(s.handler)
		defer ts.Close()

		target := getInitDomainArr()[rand.Intn(len(getInitDomainArr())-1)]
		// 존재 여부확인및 타겟이 맞는지 비교
		re, err := http.NewRequest("GET", ts.URL+"/todos/"+strconv.Itoa(target.Id), nil)
		assert.NoError(s.T(), err)
		resp, err := http.DefaultClient.Do(re)
		assert.NoError(s.T(), err)
		defer resp.Body.Close()
		var result res.DetailDto
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(s.T(), err)

		assert.Equal(s.T(), target.Id, result.Id)
		assert.Equal(s.T(), target.UserId, result.UserId)
		assert.Equal(s.T(), target.Title, result.Title)
		assert.Equal(s.T(), target.Content, result.Content)

		// 삭제
		re, err = http.NewRequest("DELETE", ts.URL+"/todos/"+strconv.Itoa(target.Id), nil)
		assert.NoError(s.T(), err)
		resp, err = http.DefaultClient.Do(re)
		assert.NoError(s.T(), err)
		defer resp.Body.Close()
		assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

		//데이터 삭제 됐는지 확인
		re, err = http.NewRequest("GET", ts.URL+"/todos/"+strconv.Itoa(target.Id), nil)
		assert.NoError(s.T(), err)
		resp, err = http.DefaultClient.Do(re)
		assert.NoError(s.T(), err)
		defer resp.Body.Close()
		var result2 res.DetailDto
		err = json.NewDecoder(resp.Body).Decode(&result2)
		assert.Error(s.T(), err)
		assert.Equal(s.T(), res.DetailDto{}, result2)
		assert.Equal(s.T(), http.StatusNotFound, resp.StatusCode)
	})
	s.Run("삭제 테스트 id 형태를 다른걸로 줄경우", func() {
		ts := httptest.NewServer(s.handler)
		defer ts.Close()

		//target := getInitDomainArr()[rand.Intn(len(getInitDomainArr())-1)]
		// 존재 여부확인
		re, err := http.NewRequest("DELETE", ts.URL+"/todos/ads", nil)
		assert.NoError(s.T(), err)
		resp, err := http.DefaultClient.Do(re)
		assert.NoError(s.T(), err)
		defer resp.Body.Close()
		var result res.DetailDto
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.Error(s.T(), err)
	})
}

func (s *TodoHandlerTestSuite) TestGetTodos() {
	s.Run("페이지 조회 ", func() {
		ts := httptest.NewServer(s.handler)
		defer ts.Close()
		re, err := http.NewRequest("GET", ts.URL+"/todos", nil)
		assert.NoError(s.T(), err)
		resp, err := http.DefaultClient.Do(re)
		assert.NoError(s.T(), err)
		defer resp.Body.Close()
		assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

		var results page.Pagination[res.ListDto]
		err = json.NewDecoder(resp.Body).Decode(&results)
		assert.NoError(s.T(), err)
		assert.NotZero(s.T(), len(results.Contents))
		assert.NotZero(s.T(), results.Total)
	})
}
