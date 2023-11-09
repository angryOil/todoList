package controller

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"todoList/controller/req"
	"todoList/controller/res"
	"todoList/domain"
	"todoList/domain/vo"
	"todoList/page"
	"todoList/repository"
	"todoList/repository/infla"
	"todoList/service"
)

var c = NewController(
	service.NewService(
		repository.NewRepository(
			infla.NewDB(),
		),
	),
)

type TodoControllerSuite struct {
	suite.Suite
	controller TodoController
	rollback   func() error
}

func TestTodoControllerTestSuite(t *testing.T) {
	suite.Run(t, &TodoControllerSuite{})
}

var testRepo repository.ITodoRepository

func (s *TodoControllerSuite) SetupTest() {
	log.Println("setup controller testSuite ... ")
	db := infla.NewDB()
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	testRepo = repository.NewRepository(tx)

	s.rollback = tx.Rollback
	s.controller = NewController(service.NewService(testRepo))
}

// 테스트시 기준이 될 데이터
var baseTestDomain = vo.Detail{
	UserId:    9999,
	Title:     "base test model title",
	Content:   "base test model content",
	OrderNum:  1,
	IsDeleted: false,
}

// 병합된 데이터 테스트 이므로 그대로
func getTestMergeTodo(todo domain.Todo) (vo.Save, error) {
	return todo.ToSave(), nil
}

// test 에선 멱등성을 위해 save(update) 메소드를 사용할거고
// update 에서는 getValidFunc 에서 나온 결과로(repo 에서 조회후 넣을 예정)
func validFunMaker(todoId int) func([]domain.Todo) (domain.Todo, error) {
	return func(todos []domain.Todo) (domain.Todo, error) {
		return domain.NewTodoBuilder().
			UserId(baseTestDomain.UserId).
			Id(todoId).
			Title(baseTestDomain.Title + strconv.Itoa(todoId)).
			Content(baseTestDomain.Content + strconv.Itoa(todoId)).
			OrderNum(baseTestDomain.OrderNum + todoId).
			IsDeleted(baseTestDomain.IsDeleted).
			Build(), nil
	}
}

// todoId,userId 는 0이어서는 안됨
func getTestSaveValid(todo domain.Todo) error {
	t := todo.ToInfo()
	if t.UserId == 0 {
		return errors.New("userId is zero")
	}
	if t.Id == 0 {
		return errors.New("todoId is zero")
	}
	return nil
}

var baseTestPageReq = page.NewReqPage(0, 0)

func getInitDomainArr() []domain.Todo {
	base := baseTestDomain
	result := make([]domain.Todo, 10)
	for i, _ := range result {
		id := i + 1
		result[i] = domain.NewTodoBuilder().
			Id(id).
			UserId(base.UserId).
			Title(base.Title + strconv.Itoa(id)).
			Content(base.Content + strconv.Itoa(id)).
			OrderNum(base.OrderNum + id).
			IsDeleted(base.IsDeleted).
			Build()
	}
	return result
}

func (s *TodoControllerSuite) BeforeTest(suiteName, testName string) {
	log.Printf("이것은 BeforeTest 입니다. 기본 값들을 생성하죠. %s %s", suiteName, testName)
	getTestValidFunc := validFunMaker
	for _, t := range getInitDomainArr() {
		//err := testRepo.Save(context.Background(), todoDomain.UserId, todoDomain.Id, getTestValidFunc(todoDomain.Id), getTestMergeTodo, getTestSaveValid)
		v := t.ToDetail()
		err := testRepo.Save(context.Background(), v.UserId, v.Id, getTestValidFunc(v.Id), getTestMergeTodo)
		if err != nil {
			panic(err)
		}
	}
}
func (s *TodoControllerSuite) AfterTest(suiteName, testName string) {
	log.Printf("이것은 AfterTest 입니다. 롤백처리를 하죠. %s %s", suiteName, testName)
	err := s.rollback()
	if err != nil {
		panic(err)
	}
}

func getTestUserValueCtx(userId int) context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, "userId", userId)
}

func (s *TodoControllerSuite) TestCreateTodo() {
	s.Run("정상적인 값을 넘겼을경우 error 는 nil 을반환한다", func() {
		ctx := getTestUserValueCtx(baseTestDomain.UserId)
		err := s.controller.CreateTodo(ctx, req.CreateTodoDto{
			Title:    "test title ",
			Content:  "test content",
			OrderNum: 22,
		})
		assert.NoError(s.T(), err)
	})
	s.Run("userId 가 없이 요청했을경우 error를 반환", func() {
		err := s.controller.CreateTodo(context.Background(), req.CreateTodoDto{
			Title:    "userId 가 없데요",
			Content:  "정말요? userId가없다구요? handler에서 거부를 안했나요?",
			OrderNum: 2222,
		})
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "user id is not valid")
	})
	s.Run("userId 를 0으로 요청했을경우 error를 반환", func() {
		ctx := getTestUserValueCtx(0)
		err := s.controller.CreateTodo(ctx, req.CreateTodoDto{
			Title:    "userId 가 0입니다",
			Content:  "정말요? userId가0 이라구요? handler 에서 거부를 안했나요?",
			OrderNum: 2222,
		})
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "invalid")
	})
	s.Run("제목이 없을경우 error를 반환", func() {
		ctx := getTestUserValueCtx(baseTestDomain.UserId)
		err := s.controller.CreateTodo(ctx, req.CreateTodoDto{
			Title:    "",
			Content:  "제목이 없는 글입니다.",
			OrderNum: 2222,
		})
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "invalid title")
	})
	s.Run("내용이 없을경우 error를 반환", func() {
		ctx := getTestUserValueCtx(baseTestDomain.UserId)
		err := s.controller.CreateTodo(ctx, req.CreateTodoDto{
			Title:    "이것은 내용이없어요",
			Content:  "",
			OrderNum: 2222,
		})
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "invalid")
	})
	s.Run("orderNum 이 없을경우 error를 반환", func() {
		ctx := getTestUserValueCtx(baseTestDomain.UserId)
		err := s.controller.CreateTodo(ctx, req.CreateTodoDto{
			Title:    "있어요",
			Content:  "제목이 있습니다 글입니다.",
			OrderNum: 0,
		})
		assert.NoError(s.T(), err)
	})
}

func (s *TodoControllerSuite) TestUpdateTodo() {
	s.Run("정상적인 업데이트 일경우", func() {
		target := getInitDomainArr()[rand.Intn(len(getInitDomainArr())-1)]
		t := target.ToDetail()
		ctx := getTestUserValueCtx(t.UserId)
		err := s.controller.UpdateTodo(ctx, req.UpdateTodoDto{
			Id:        t.Id,
			Title:     "mod!!!",
			Content:   "Con mod",
			OrderNum:  33,
			IsDeleted: false,
		})
		assert.NoError(s.T(), err)
	})

	s.Run("todoId없이 수정을 시도할경우 error를 리턴", func() {
		target := getInitDomainArr()[rand.Intn(len(getInitDomainArr())-1)]
		t := target.ToDetail()
		ctx := getTestUserValueCtx(t.UserId)
		err := s.controller.UpdateTodo(ctx, req.UpdateTodoDto{
			Id:        0,
			Title:     "todo id 가 없네요",
			Content:   "Todo Id 가 없어요!!",
			OrderNum:  22,
			IsDeleted: false,
		})
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "no")
	})
	s.Run("userId없이 수정을 시도할경우 error를 리턴", func() {
		target := getInitDomainArr()[rand.Intn(len(getInitDomainArr())-1)]
		t := target.ToDetail()
		err := s.controller.UpdateTodo(context.Background(), req.UpdateTodoDto{
			Id:        t.Id,
			Title:     "userID 없다는데 정말인가요?",
			Content:   "userId가  없어요...",
			OrderNum:  22,
			IsDeleted: false,
		})
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "user id is not valid")
	})
	s.Run("userId를 0으로 수정을 시도할경우 error를 리턴", func() {
		target := getInitDomainArr()[rand.Intn(len(getInitDomainArr())-1)]
		ctx := getTestUserValueCtx(0)
		t := target.ToDetail()
		err := s.controller.UpdateTodo(ctx, req.UpdateTodoDto{
			Id:        t.Id,
			Title:     "userID 없다는데 정말인가요?",
			Content:   "userId가  없어요...",
			OrderNum:  22,
			IsDeleted: false,
		})
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "no row")
	})
	s.Run("제목이없이 수정을 시도할경우 error를 리턴", func() {
		target := getInitDomainArr()[rand.Intn(len(getInitDomainArr())-1)].ToDetail()
		ctx := getTestUserValueCtx(target.UserId)
		err := s.controller.UpdateTodo(ctx, req.UpdateTodoDto{
			Id:        target.Id,
			Title:     "",
			Content:   "제목이 없어요...",
			OrderNum:  22,
			IsDeleted: false,
		})
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "invalid")
	})
	s.Run("내용없이 수정을 시도할경우 error를 리턴", func() {
		target := getInitDomainArr()[rand.Intn(len(getInitDomainArr())-1)].ToDetail()
		ctx := getTestUserValueCtx(target.UserId)
		err := s.controller.UpdateTodo(ctx, req.UpdateTodoDto{
			Id:        target.Id,
			Title:     "내용이 없어요 ",
			Content:   "",
			OrderNum:  22,
			IsDeleted: false,
		})
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "invalid")
	})
	s.Run("orderNum 없이 수정을 시도 할경우 error를 리턴", func() {
		target := getInitDomainArr()[rand.Intn(len(getInitDomainArr())-1)].ToDetail()
		ctx := getTestUserValueCtx(target.UserId)
		err := s.controller.UpdateTodo(ctx, req.UpdateTodoDto{
			Id:        target.Id,
			Title:     "mod title!@",
			Content:   "mod content?",
			OrderNum:  0,
			IsDeleted: false,
		})
		assert.NoError(s.T(), err)
	})
}

func (s *TodoControllerSuite) TestDeleteTodo() {
	s.Run("실제 있는 값을 삭제했을경우 error nil을 반환", func() {
		target := getInitDomainArr()[rand.Intn(len(getInitDomainArr())-1)].ToDetail()
		ctx := getTestUserValueCtx(target.UserId)
		dto, err := s.controller.GetDetail(ctx, target.Id)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), target.Id, dto.Id)
	})
	s.Run("없는 데이터를 삭제했을때 error nil을 반환", func() {
		// todoId가 0인 데이터는 존재할수없음
		ctx := getTestUserValueCtx(baseTestDomain.UserId)
		dto, err := s.controller.GetDetail(ctx, 0)
		assert.Error(s.T(), err)
		assert.Equal(s.T(), res.DetailDto{}, dto)

		err = s.controller.DeleteTodo(ctx, 0)
		assert.NoError(s.T(), err)
	})
	s.Run("userId없이 요청을 했을경우 error 를 반환", func() {
		err := s.controller.DeleteTodo(context.Background(), getInitDomainArr()[0].ToDetail().Id)
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "user id is not valid")
	})
	s.Run("userId를 0으로 요청을 했을경우 error 를 반환", func() {
		ctx := getTestUserValueCtx(0)
		err := s.controller.DeleteTodo(ctx, getInitDomainArr()[0].ToDetail().Id)
		assert.NoError(s.T(), err)
	})
}

func (s *TodoControllerSuite) TestGetTodos() {
	s.Run("특정 userId의 todoList전체 갯수를 구한후 전체 list를 요청한다", func() {
		ctx := getTestUserValueCtx(baseTestDomain.UserId)
		_, totalCnt, err := s.controller.GetTodos(ctx, page.NewReqPage(0, 0))
		assert.NoError(s.T(), err)
		assert.NotZero(s.T(), totalCnt)

		results, totalCnt, err := s.controller.GetTodos(ctx, page.NewReqPage(0, totalCnt))
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), totalCnt, len(results))
	})
	s.Run("userId 가 없을경우 error를 반환한다", func() {
		_, totalCnt, err := s.controller.GetTodos(context.Background(), page.NewReqPage(0, 0))
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "user id is not valid")
		assert.Equal(s.T(), 0, totalCnt)
	})
	s.Run("userId 0 일경우 error를 반환한다", func() {
		ctx := getTestUserValueCtx(0)
		_, totalCnt, err := s.controller.GetTodos(ctx, page.NewReqPage(0, 0))
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "invalid")
		assert.Equal(s.T(), 0, totalCnt)
	})
}
